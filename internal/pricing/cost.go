package pricing

const webSearchFeeUSD = 0.01

func clampTokenCount(value int64) int64 {
	if value > 0 {
		return value
	}
	return 0
}

func (catalog *Catalog) CalculateCost(input CostInput) float64 {
	normalizedInput := clampTokenCount(input.InputTokens)
	normalizedCached := clampTokenCount(input.CachedInputTokens)
	if normalizedCached > normalizedInput {
		normalizedCached = normalizedInput
	}
	normalizedNonCached := clampTokenCount(input.NonCachedInputTokens)
	if normalizedNonCached > normalizedInput-normalizedCached {
		normalizedNonCached = normalizedInput - normalizedCached
	}
	normalizedOutput := clampTokenCount(input.OutputTokens)
	normalizedCacheWrite := clampTokenCount(input.CacheWriteTokens)
	normalizedCacheWrite1h := clampTokenCount(input.CacheWrite1hTokens)
	if normalizedCacheWrite1h > normalizedCacheWrite {
		normalizedCacheWrite1h = normalizedCacheWrite
	}
	normalizedCacheWrite5m := normalizedCacheWrite - normalizedCacheWrite1h

	catalogEntry := catalog.ResolveEntry(input.Model, input.UsageTimestampUnix)
	if catalogEntry == nil {
		return 0
	}

	cachedRate := catalogEntry.Input
	if catalogEntry.InputCached != nil {
		cachedRate = *catalogEntry.InputCached
	}

	billedNonCached := normalizedNonCached + normalizedCacheWrite
	billedCacheWrite := float64(0)
	if catalogEntry.CacheWrite != nil {
		billedNonCached = normalizedNonCached
		rate1h := *catalogEntry.CacheWrite
		if catalogEntry.CacheWrite1h != nil {
			rate1h = *catalogEntry.CacheWrite1h
		}
		billedCacheWrite = float64(normalizedCacheWrite5m)**catalogEntry.CacheWrite + float64(normalizedCacheWrite1h)*rate1h
	}

	return toMillionTokens(float64(billedNonCached))*catalogEntry.Input +
		toMillionTokens(float64(normalizedCached))*cachedRate +
		toMillionTokens(float64(normalizedOutput))*catalogEntry.Output +
		toMillionTokens(billedCacheWrite) +
		float64(clampTokenCount(input.WebSearchQueries))*webSearchFeeUSD
}

func toMillionTokens(value float64) float64 {
	return value / 1_000_000
}

type CostInput struct {
	Model                string
	InputTokens          int64
	CachedInputTokens    int64
	NonCachedInputTokens int64
	OutputTokens         int64
	CacheWriteTokens     int64
	CacheWrite1hTokens   int64
	WebSearchQueries     int64
	UsageTimestampUnix   *int64
}
