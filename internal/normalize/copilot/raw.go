package copilot

import (
	"encoding/json"
	"maps"

	"traceknot/internal/normalize/shared"
)

var piiAttributeKeys = []string{"enduser.pseudo.id"}

func stripPII(attributes map[string]any) map[string]any {
	stripped := make(map[string]any, len(attributes))
	maps.Copy(stripped, attributes)
	for _, key := range piiAttributeKeys {
		delete(stripped, key)
	}
	return stripped
}

func spansToRecords(spans []Span) []shared.RawRecord {
	records := make([]shared.RawRecord, 0, len(spans))
	for _, span := range spans {
		if span.TraceID == "" {
			continue
		}
		narrowed := span
		narrowed.Attributes = stripPII(span.Attributes)
		payload, err := json.Marshal(narrowed)
		if err != nil {
			continue
		}
		records = append(records, shared.RawRecord{
			NativeID:    narrowed.TraceID,
			Signal:      "span",
			DedupKey:    narrowed.SpanID,
			TimestampMs: narrowed.TimestampMs,
			PayloadJSON: string(payload),
		})
	}
	return records
}

func recordsToSession(nativeID string, records []shared.RawRecord) *Session {
	session := &Session{TraceID: nativeID}
	for _, record := range records {
		if record.Signal != "span" {
			continue
		}
		var span Span
		if err := json.Unmarshal([]byte(record.PayloadJSON), &span); err == nil {
			session.Spans = append(session.Spans, span)
		}
	}
	return session
}
