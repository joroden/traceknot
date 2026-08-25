package copilot

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
	return "copilot"
}

func (normalizer *Normalizer) RebuildScope() shared.RebuildScope {
	return shared.RebuildScopeTouched
}

func (normalizer *Normalizer) ExtractLogs(*logspb.LogsData) []shared.RawRecord {
	return nil
}

func (normalizer *Normalizer) ExtractTraces(data *tracepb.TracesData) []shared.RawRecord {
	return spansToRecords(ExtractSpans(data))
}

func (normalizer *Normalizer) Rebuild(byNativeID map[string][]shared.RawRecord, touchedIDs []string) []shared.BuildResult {
	var results []shared.BuildResult
	seen := make(map[string]bool, len(touchedIDs))
	for _, nativeID := range touchedIDs {
		if seen[nativeID] {
			continue
		}
		seen[nativeID] = true
		session := recordsToSession(nativeID, byNativeID[nativeID])
		seed, content := normalizer.builder.buildSession(session)
		if seed == nil {
			continue
		}
		results = append(results, shared.BuildResult{Seed: seed, Content: content})
	}
	return results
}
