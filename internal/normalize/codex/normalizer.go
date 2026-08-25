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

func (normalizer *Normalizer) RebuildScope() shared.RebuildScope {
	return shared.RebuildScopeProvider
}

func (normalizer *Normalizer) ExtractLogs(data *logspb.LogsData) []shared.RawRecord {
	return eventsToRecords(ExtractEvents(data))
}

func (normalizer *Normalizer) ExtractTraces(*tracepb.TracesData) []shared.RawRecord {
	return nil
}

func (normalizer *Normalizer) Rebuild(byNativeID map[string][]shared.RawRecord, touchedIDs []string) []shared.BuildResult {
	all := recordsToEventsByConversation(byNativeID)
	return normalizer.builder.Build(all, touchedIDs)
}

func (normalizer *Normalizer) ResolveRoots(byNativeID map[string][]shared.RawRecord) map[string]string {
	all := recordsToEventsByConversation(byNativeID)
	subagents := subagentMap(all)
	roots := make(map[string]string, len(all))
	for nativeID := range all {
		roots[nativeID] = rootFor(nativeID, subagents)
	}
	return roots
}
