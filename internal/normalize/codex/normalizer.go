package codex

import (
	"traceknot/internal/normalize/shared"
	"traceknot/internal/pricing"
	"traceknot/internal/tokenize"
	logspb "go.opentelemetry.io/proto/otlp/logs/v1"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
)

type Normalizer struct {
	builder *Builder
}

func NewNormalizer(catalog *pricing.Catalog, estimator *tokenize.Estimator) *Normalizer {
	return &Normalizer{builder: NewBuilder(catalog, estimator)}
}

func (normalizer *Normalizer) Provider() string {
	return "codex"
}

func (normalizer *Normalizer) ExtractLogs(data *logspb.LogsData) []shared.RawRecord {
	return eventsToRecords(ExtractEvents(data))
}

func (normalizer *Normalizer) ExtractTraces(*tracepb.TracesData) []shared.RawRecord {
	return nil
}

func (normalizer *Normalizer) Rebuild(byNativeID map[string][]shared.RawRecord, _ []string) []shared.BuildResult {
	all := recordsToEventsByConversation(byNativeID)
	return normalizer.builder.Build(all)
}
