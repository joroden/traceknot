package shared

import "traceknot/internal/pricing"

func NodeCost(
	effectiveModel *string,
	input, cached, cacheWrite, cacheWrite1h, output, webSearchQueries int64,
	timestamp *int64,
	catalog *pricing.Catalog,
) float64 {
	if effectiveModel == nil {
		return 0
	}
	nonCached := max(input-cached-cacheWrite, 0)
	return catalog.CalculateCost(pricing.CostInput{
		Model:                *effectiveModel,
		InputTokens:          input,
		CachedInputTokens:    cached,
		NonCachedInputTokens: nonCached,
		OutputTokens:         output,
		CacheWriteTokens:     cacheWrite,
		CacheWrite1hTokens:   cacheWrite1h,
		WebSearchQueries:     webSearchQueries,
		UsageTimestampUnix:   timestamp,
	})
}
