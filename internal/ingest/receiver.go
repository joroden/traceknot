package ingest

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"traceknot/internal/normalize/shared"
	"traceknot/internal/store"
	"google.golang.org/protobuf/encoding/protojson"
)

type Receiver struct {
	store       *store.Store
	normalizers map[string]shared.Normalizer
	logger      *slog.Logger
	captureDir  string
	captureFile string
	marshal     protojson.MarshalOptions
}

func NewReceiver(storeHandle *store.Store, normalizers map[string]shared.Normalizer, logger *slog.Logger) *Receiver {
	return &Receiver{
		store:       storeHandle,
		normalizers: normalizers,
		logger:      logger,
		captureDir:  os.Getenv("TRACEKNOT_CAPTURE_DIR"),
		captureFile: os.Getenv("TRACEKNOT_CAPTURE_FILE"),
		marshal: protojson.MarshalOptions{
			EmitUnpopulated: false,
			UseProtoNames:   true,
		},
	}
}

func (receiver *Receiver) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/logs", func(writer http.ResponseWriter, request *http.Request) {
		receiver.handleLogs(writer, request)
	})
	mux.HandleFunc("POST /v1/traces", func(writer http.ResponseWriter, request *http.Request) {
		receiver.handleTraces(writer, request)
	})
	return mux
}

func (receiver *Receiver) handleLogs(writer http.ResponseWriter, request *http.Request) {
	decoded, contentType, err := receiver.readBody(request)
	if err != nil {
		receiver.writeError(writer, http.StatusBadRequest, "body_read_failed", err.Error())
		return
	}
	logsData, _, err := parseLogsRequest(decoded, contentType)
	if err != nil {
		receiver.writeError(writer, http.StatusBadRequest, "invalid_otlp_payload", err.Error())
		return
	}
	receiver.capture(logsData, "logs", request)
	if normalizer := receiver.normalizers[logsServiceName(logsData)]; normalizer != nil {
		receiver.Ingest(request.Context(), normalizer, normalizer.ExtractLogs(logsData))
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	fmt.Fprint(writer, `{"ok":true}`)
}

func (receiver *Receiver) handleTraces(writer http.ResponseWriter, request *http.Request) {
	decoded, contentType, err := receiver.readBody(request)
	if err != nil {
		receiver.writeError(writer, http.StatusBadRequest, "body_read_failed", err.Error())
		return
	}
	tracesData, _, err := parseTraceRequest(decoded, contentType)
	if err != nil {
		receiver.writeError(writer, http.StatusBadRequest, "invalid_otlp_payload", err.Error())
		return
	}
	receiver.capture(tracesData, "traces", request)
	if normalizer := receiver.normalizers[tracesServiceName(tracesData)]; normalizer != nil {
		receiver.Ingest(request.Context(), normalizer, normalizer.ExtractTraces(tracesData))
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	fmt.Fprint(writer, `{"ok":true}`)
}

func (receiver *Receiver) writeError(writer http.ResponseWriter, status int, code string, message string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	body, _ := json.Marshal(map[string]any{"error": map[string]any{"code": code, "message": message}})
	_, _ = writer.Write(body)
}

func (receiver *Receiver) Ingest(ctx context.Context, normalizer shared.Normalizer, records []shared.RawRecord) {
	if len(records) == 0 {
		return
	}
	provider := normalizer.Provider()
	stored := make([]store.RawSignalRecord, 0, len(records))
	touched := make(map[string]struct{}, len(records))
	for _, record := range records {
		stored = append(stored, store.RawSignalRecord{
			NativeID:    record.NativeID,
			Provider:    provider,
			Signal:      record.Signal,
			DedupKey:    record.DedupKey,
			TimestampMs: record.TimestampMs,
			PayloadJSON: record.PayloadJSON,
		})
		touched[record.NativeID] = struct{}{}
	}
	if err := receiver.store.InsertRawSignal(ctx, stored); err != nil {
		receiver.logger.Error("insert raw_signal failed", "provider", provider, "error", err)
		return
	}

	touchedIDs := make([]string, 0, len(touched))
	for nativeID := range touched {
		touchedIDs = append(touchedIDs, nativeID)
	}

	byProvider, err := receiver.loadScoped(ctx, normalizer, provider, touchedIDs)
	if err != nil {
		receiver.logger.Error("load raw_signal failed", "provider", provider, "error", err)
		return
	}
	byNativeID := toRawRecords(byProvider)

	if linked, ok := normalizer.(shared.LinkedNormalizer); ok {
		if err := receiver.store.SaveConversationRoots(ctx, provider, linked.ResolveRoots(byNativeID)); err != nil {
			receiver.logger.Error("save conversation roots failed", "provider", provider, "error", err)
		}
	}

	for _, result := range normalizer.Rebuild(byNativeID, touchedIDs) {
		if err := receiver.store.ReplaceSession(ctx, result.Seed, result.Content); err != nil {
			receiver.logger.Error("ingest persist session failed",
				"session", result.Seed.SessionID,
				"provider", result.Seed.Provider,
				"error", err)
			continue
		}
		receiver.logger.Info("ingest session updated",
			"session", result.Seed.SessionID,
			"provider", result.Seed.Provider,
			"chats", len(result.Content.Chats),
			"tools", len(result.Content.ToolCalls),
			"agents", len(result.Content.Agents))
	}
}

func (receiver *Receiver) loadScoped(
	ctx context.Context,
	normalizer shared.Normalizer,
	provider string,
	touchedIDs []string,
) (map[string][]store.RawSignalRecord, error) {
	if _, ok := normalizer.(shared.LinkedNormalizer); ok {
		return receiver.loadFamilyScoped(ctx, provider, touchedIDs)
	}
	if normalizer.RebuildScope() == shared.RebuildScopeProvider {
		return receiver.store.LoadRawSignalByProvider(ctx, provider)
	}
	return receiver.store.LoadRawSignalByNativeIDs(ctx, provider, touchedIDs)
}

func (receiver *Receiver) loadFamilyScoped(
	ctx context.Context,
	provider string,
	touchedIDs []string,
) (map[string][]store.RawSignalRecord, error) {
	knownRoots, err := receiver.store.ConversationRoots(ctx, provider, touchedIDs)
	if err != nil {
		return nil, err
	}

	rootSet := make(map[string]struct{}, len(touchedIDs))
	for _, nativeID := range touchedIDs {
		if root, ok := knownRoots[nativeID]; ok {
			rootSet[root] = struct{}{}
		} else {
			rootSet[nativeID] = struct{}{}
		}
	}
	roots := make([]string, 0, len(rootSet))
	for root := range rootSet {
		roots = append(roots, root)
	}

	family, err := receiver.store.ConversationFamily(ctx, provider, roots)
	if err != nil {
		return nil, err
	}

	scope := make(map[string]struct{}, len(family)+len(touchedIDs)+len(roots))
	for _, nativeID := range family {
		scope[nativeID] = struct{}{}
	}
	for _, nativeID := range touchedIDs {
		scope[nativeID] = struct{}{}
	}
	for _, root := range roots {
		scope[root] = struct{}{}
	}
	scopedIDs := make([]string, 0, len(scope))
	for nativeID := range scope {
		scopedIDs = append(scopedIDs, nativeID)
	}
	return receiver.store.LoadRawSignalByNativeIDs(ctx, provider, scopedIDs)
}

func toRawRecords(byProvider map[string][]store.RawSignalRecord) map[string][]shared.RawRecord {
	byNativeID := make(map[string][]shared.RawRecord, len(byProvider))
	for nativeID, rows := range byProvider {
		converted := make([]shared.RawRecord, 0, len(rows))
		for _, row := range rows {
			converted = append(converted, shared.RawRecord{
				NativeID:    row.NativeID,
				Signal:      row.Signal,
				DedupKey:    row.DedupKey,
				TimestampMs: row.TimestampMs,
				PayloadJSON: row.PayloadJSON,
			})
		}
		byNativeID[nativeID] = converted
	}
	return byNativeID
}
