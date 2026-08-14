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

	byProvider, err := receiver.store.LoadRawSignalByProvider(ctx, provider)
	if err != nil {
		receiver.logger.Error("load raw_signal failed", "provider", provider, "error", err)
		return
	}
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

	touchedIDs := make([]string, 0, len(touched))
	for nativeID := range touched {
		touchedIDs = append(touchedIDs, nativeID)
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
