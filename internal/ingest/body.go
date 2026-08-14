package ingest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	logspb "go.opentelemetry.io/proto/otlp/logs/v1"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
	"google.golang.org/protobuf/proto"
)

const defaultBodyLimitBytes = 64 * 1024 * 1024

func (receiver *Receiver) readBody(request *http.Request) ([]byte, string, error) {
	rawBody, err := io.ReadAll(io.LimitReader(request.Body, defaultBodyLimitBytes+1))
	if err != nil {
		return nil, "", err
	}
	if int64(len(rawBody)) > defaultBodyLimitBytes {
		return nil, "", fmt.Errorf("request body exceeded the %d byte limit", defaultBodyLimitBytes)
	}
	decoded, err := decodeBody(rawBody, request.Header.Get("content-encoding"))
	if err != nil {
		return nil, "", err
	}
	return decoded, request.Header.Get("content-type"), nil
}

func (receiver *Receiver) capture(message proto.Message, signalType string, request *http.Request) {
	if receiver.captureDir == "" {
		return
	}
	if err := os.MkdirAll(receiver.captureDir, 0o755); err != nil {
		return
	}
	fileName := receiver.captureFile
	if fileName == "" {
		fileName = "ingest.jsonl"
	}
	path := filepath.Join(receiver.captureDir, fileName)

	var spanCount, logCount int
	switch typed := message.(type) {
	case *tracepb.TracesData:
		spanCount = tracesCount(typed)
	case *logspb.LogsData:
		logCount = logsCount(typed)
	}

	record := map[string]any{
		"signal":      signalType,
		"user_agent":  request.Header.Get("user-agent"),
		"received_at": time.Now().UTC().Format(time.RFC3339Nano),
		"span_count":  spanCount,
		"log_count":   logCount,
	}
	rawJSON, err := receiver.marshal.Marshal(message)
	if err != nil {
		rawJSON = []byte(fmt.Sprintf(`{"marshal_error":%q}`, err.Error()))
	}
	record["raw_json"] = string(rawJSON)

	encoded, err := json.Marshal(record)
	if err != nil {
		return
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer file.Close()
	_, _ = file.Write(append(encoded, '\n'))
}

func tracesCount(data *tracepb.TracesData) int {
	count := 0
	for _, resourceSpans := range data.ResourceSpans {
		for _, scopeSpans := range resourceSpans.ScopeSpans {
			count += len(scopeSpans.Spans)
		}
	}
	return count
}

func logsCount(data *logspb.LogsData) int {
	count := 0
	for _, resourceLogs := range data.ResourceLogs {
		for _, scopeLogs := range resourceLogs.ScopeLogs {
			count += len(scopeLogs.LogRecords)
		}
	}
	return count
}
