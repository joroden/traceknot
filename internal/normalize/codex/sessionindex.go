package codex

import (
	"encoding/json"
	"time"

	"traceknot/internal/normalize/shared"
)

func SessionIndexRecord(line []byte) (shared.RawRecord, bool) {
	var entry struct {
		ID         string `json:"id"`
		ThreadName string `json:"thread_name"`
		UpdatedAt  string `json:"updated_at"`
	}
	if err := json.Unmarshal(line, &entry); err != nil || entry.ID == "" || entry.ThreadName == "" {
		return shared.RawRecord{}, false
	}
	updatedAt, err := time.Parse(time.RFC3339Nano, entry.UpdatedAt)
	if err != nil {
		return shared.RawRecord{}, false
	}
	event := Event{
		Name:           eventSessionTitle,
		ConversationID: entry.ID,
		TimestampMs:    updatedAt.UnixMilli(),
		Attributes:     map[string]any{"title": entry.ThreadName},
	}
	payload, err := json.Marshal(event)
	if err != nil {
		return shared.RawRecord{}, false
	}
	return shared.RawRecord{
		NativeID:    entry.ID,
		Signal:      "session_index",
		DedupKey:    contentHash(event.Name, entry.ID, entry.UpdatedAt, entry.ThreadName),
		TimestampMs: event.TimestampMs,
		PayloadJSON: string(payload),
	}, true
}
