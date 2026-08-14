package claude

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"maps"

	"traceknot/internal/jsonutil"
	"traceknot/internal/normalize/shared"
)

var piiAttributeKeys = []string{
	"user.id", "user.email", "user.account_uuid", "user.account_id",
	"organization.id", "terminal.type",
}

func stripPII(attributes map[string]any) map[string]any {
	stripped := make(map[string]any, len(attributes))
	maps.Copy(stripped, attributes)
	for _, key := range piiAttributeKeys {
		delete(stripped, key)
	}
	return stripped
}

func narrowEvent(event Event) (Event, bool) {
	switch event.Name {
	case eventAPIResponseBody:
		return Event{}, false
	case eventAPIRequestBody:
		body, _ := attributeString(event.Attributes, "body")
		narrowed := event
		narrowed.Attributes = map[string]any{"system_prompt": systemPromptFromBody(body)}
		return narrowed, true
	default:
		narrowed := event
		narrowed.Attributes = stripPII(event.Attributes)
		return narrowed, true
	}
}

func eventsToRecords(events []Event) []shared.RawRecord {
	records := make([]shared.RawRecord, 0, len(events))
	for _, event := range events {
		if event.SessionID == "" {
			continue
		}
		narrowed, ok := narrowEvent(event)
		if !ok {
			continue
		}
		payload, err := json.Marshal(narrowed)
		if err != nil {
			continue
		}
		records = append(records, shared.RawRecord{
			NativeID:    narrowed.SessionID,
			Signal:      "log",
			DedupKey:    contentHash(narrowed.Name, narrowed.SessionID, narrowed.PromptID, narrowed.Sequence, narrowed.Attributes),
			TimestampMs: narrowed.TimestampMs,
			PayloadJSON: string(payload),
		})
	}
	return records
}

func spansToRecords(spans []Span) []shared.RawRecord {
	records := make([]shared.RawRecord, 0, len(spans))
	for _, span := range spans {
		if span.SessionID == "" {
			continue
		}
		narrowed := span
		narrowed.Attributes = stripPII(span.Attributes)
		payload, err := json.Marshal(narrowed)
		if err != nil {
			continue
		}
		records = append(records, shared.RawRecord{
			NativeID:    narrowed.SessionID,
			Signal:      "span",
			DedupKey:    narrowed.SpanID,
			TimestampMs: narrowed.TimestampMs,
			PayloadJSON: string(payload),
		})
	}
	return records
}

func recordsToSession(nativeID string, records []shared.RawRecord) *Session {
	session := &Session{SessionID: nativeID}
	for _, record := range records {
		switch record.Signal {
		case "log":
			var event Event
			if err := json.Unmarshal([]byte(record.PayloadJSON), &event); err == nil {
				session.Events = append(session.Events, event)
			}
		case "span":
			var span Span
			if err := json.Unmarshal([]byte(record.PayloadJSON), &span); err == nil {
				session.Spans = append(session.Spans, span)
			}
		}
	}
	return session
}

func contentHash(parts ...any) string {
	sum := sha256.Sum256([]byte(jsonutil.ToCanonicalJSON(parts)))
	return hex.EncodeToString(sum[:])
}
