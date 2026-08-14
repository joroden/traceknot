package shared

import (
	"traceknot/internal/model"
	logspb "go.opentelemetry.io/proto/otlp/logs/v1"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
)

type RawRecord struct {
	NativeID    string
	Signal      string
	DedupKey    string
	TimestampMs int64
	PayloadJSON string
}

type BuildResult struct {
	Seed    *model.SessionSeed
	Content *model.SessionContent
}

type Normalizer interface {
	Provider() string
	ExtractLogs(*logspb.LogsData) []RawRecord
	ExtractTraces(*tracepb.TracesData) []RawRecord

	Rebuild(byNativeID map[string][]RawRecord, touchedIDs []string) []BuildResult
}
