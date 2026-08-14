package claude

type cacheTierSplit struct {
	cache1h int64
	cache5m int64
}

func joinUsageSupplements(events []Event) (map[string]cacheTierSplit, []Event) {
	var supplements []Event
	rest := make([]Event, 0, len(events))
	for _, event := range events {
		if event.Name == eventUsageSupplement {
			supplements = append(supplements, event)
			continue
		}
		rest = append(rest, event)
	}
	if len(supplements) == 0 {
		return nil, rest
	}

	byModelOTel := map[string][]Event{}
	for _, event := range rest {
		if event.Name != eventAPIRequest {
			continue
		}
		querySource, _ := attributeString(event.Attributes, "query_source")
		if !isVisibleQuerySource(querySource) {
			continue
		}
		modelName, ok := attributeString(event.Attributes, "model")
		if !ok {
			continue
		}
		byModelOTel[modelName] = append(byModelOTel[modelName], event)
	}

	byModelSupplement := map[string][]Event{}
	for _, event := range supplements {
		modelName, ok := attributeString(event.Attributes, "model")
		if !ok {
			continue
		}
		byModelSupplement[modelName] = append(byModelSupplement[modelName], event)
	}

	const matchWindow = 5

	splits := map[string]cacheTierSplit{}
	for modelName, otelEvents := range byModelOTel {
		supplementEvents := byModelSupplement[modelName]
		used := make([]bool, len(supplementEvents))
		for index, otelEvent := range otelEvents {
			otelCacheCreate, _ := attributeInt(otelEvent.Attributes, "cache_creation_tokens")
			lo := max(index-matchWindow, 0)
			hi := min(index+matchWindow, len(supplementEvents)-1)
			matchIndex := -1
			for candidate := lo; candidate <= hi; candidate++ {
				if used[candidate] {
					continue
				}
				supplementCacheCreate, _ := attributeInt(supplementEvents[candidate].Attributes, "cache_creation_input_tokens")
				if supplementCacheCreate == otelCacheCreate {
					matchIndex = candidate
					break
				}
			}
			if matchIndex == -1 {
				continue
			}
			used[matchIndex] = true
			requestID, ok := attributeString(otelEvent.Attributes, "request_id")
			if !ok {
				continue
			}
			cache1h, _ := attributeInt(supplementEvents[matchIndex].Attributes, "cache_creation_1h_tokens")
			cache5m, _ := attributeInt(supplementEvents[matchIndex].Attributes, "cache_creation_5m_tokens")
			splits[requestID] = cacheTierSplit{cache1h: cache1h, cache5m: cache5m}
		}
	}
	return splits, rest
}
