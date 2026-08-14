package claude

import (
	"encoding/json"
	"time"

	"traceknot/internal/normalize/shared"
)

const eventUsageSupplement = "usage_supplement"

type transcriptLine struct {
	Timestamp string             `json:"timestamp"`
	Message   *transcriptMessage `json:"message"`
}

type transcriptMessage struct {
	ID    string           `json:"id"`
	Model string           `json:"model"`
	Usage *transcriptUsage `json:"usage"`
}

type transcriptUsage struct {
	CacheCreationInputTokens int64                    `json:"cache_creation_input_tokens"`
	CacheCreation            *transcriptCacheCreation `json:"cache_creation"`
}

type transcriptCacheCreation struct {
	Ephemeral1hInputTokens int64 `json:"ephemeral_1h_input_tokens"`
	Ephemeral5mInputTokens int64 `json:"ephemeral_5m_input_tokens"`
}

func ParseUsageSupplementLine(sessionID string, line []byte) (Event, bool) {
	var parsed transcriptLine
	if err := json.Unmarshal(line, &parsed); err != nil {
		return Event{}, false
	}
	if parsed.Message == nil || parsed.Message.Usage == nil || parsed.Message.ID == "" {
		return Event{}, false
	}
	timestampMs, ok := parseTranscriptTimestamp(parsed.Timestamp)
	if !ok {
		return Event{}, false
	}
	var cache1h, cache5m int64
	if creation := parsed.Message.Usage.CacheCreation; creation != nil {
		cache1h = creation.Ephemeral1hInputTokens
		cache5m = creation.Ephemeral5mInputTokens
	}
	return Event{
		Name:        eventUsageSupplement,
		SessionID:   sessionID,
		TimestampMs: timestampMs,
		Attributes: map[string]any{
			"message_id":                  parsed.Message.ID,
			"model":                       parsed.Message.Model,
			"cache_creation_input_tokens": parsed.Message.Usage.CacheCreationInputTokens,
			"cache_creation_1h_tokens":    cache1h,
			"cache_creation_5m_tokens":    cache5m,
		},
	}, true
}

func parseTranscriptTimestamp(value string) (int64, bool) {
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return 0, false
	}
	return parsed.UnixMilli(), true
}

func UsageSupplementRecord(event Event) (shared.RawRecord, bool) {
	messageID, ok := attributeString(event.Attributes, "message_id")
	if !ok {
		return shared.RawRecord{}, false
	}
	payload, err := json.Marshal(event)
	if err != nil {
		return shared.RawRecord{}, false
	}
	return shared.RawRecord{
		NativeID:    event.SessionID,
		Signal:      "log",
		DedupKey:    contentHash(eventUsageSupplement, event.SessionID, messageID),
		TimestampMs: event.TimestampMs,
		PayloadJSON: string(payload),
	}, true
}
