package codex

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"maps"

	"traceknot/internal/jsonutil"
	"traceknot/internal/normalize/shared"
)

var piiAttributeKeys = []string{
	"user.email", "user.account_id", "terminal.type",
	"app.version", "auth_mode", "slug", "originator",
}

func stripPII(attributes map[string]any) map[string]any {
	stripped := make(map[string]any, len(attributes))
	maps.Copy(stripped, attributes)
	for _, key := range piiAttributeKeys {
		delete(stripped, key)
	}
	return stripped
}

var relevantEventNames = map[string]bool{
	eventConversationStarts: true,
	eventUserPrompt:         true,
	eventToolResult:         true,
	eventToolDecision:       true,
	eventSSEEvent:           true,
}

func eventsToRecords(events []Event) []shared.RawRecord {
	records := make([]shared.RawRecord, 0, len(events))
	for _, event := range events {
		if !relevantEventNames[event.Name] || event.ConversationID == "" {
			continue
		}
		narrowed := event
		narrowed.Attributes = stripPII(event.Attributes)
		payload, err := json.Marshal(narrowed)
		if err != nil {
			continue
		}
		records = append(records, shared.RawRecord{
			NativeID:    narrowed.ConversationID,
			Signal:      "log",
			DedupKey:    contentHash(narrowed.Name, narrowed.ConversationID, narrowed.TimestampMs, narrowed.Attributes),
			TimestampMs: narrowed.TimestampMs,
			PayloadJSON: string(payload),
		})
	}
	return records
}

func recordsToEventsByConversation(byNativeID map[string][]shared.RawRecord) map[string][]Event {
	all := make(map[string][]Event, len(byNativeID))
	for nativeID, records := range byNativeID {
		events := make([]Event, 0, len(records))
		for _, record := range records {
			if record.Signal != "log" && record.Signal != "rollout" {
				continue
			}
			var event Event
			if err := json.Unmarshal([]byte(record.PayloadJSON), &event); err == nil {
				events = append(events, event)
			}
		}
		all[nativeID] = events
	}
	return all
}

func contentHash(parts ...any) string {
	sum := sha256.Sum256([]byte(jsonutil.ToCanonicalJSON(parts)))
	return hex.EncodeToString(sum[:])
}
