package claude

import (
	logspb "go.opentelemetry.io/proto/otlp/logs/v1"
)

type Event struct {
	Name        string
	SessionID   string
	PromptID    string
	TimestampMs int64
	Sequence    int64
	Attributes  map[string]any
}

func ExtractEvents(request *logspb.LogsData) []Event {
	var events []Event
	for _, resourceLogs := range request.ResourceLogs {
		if resourceLogs == nil {
			continue
		}
		for _, scopeLogs := range resourceLogs.ScopeLogs {
			if scopeLogs == nil {
				continue
			}
			for _, record := range scopeLogs.LogRecords {
				if record == nil {
					continue
				}
				attributes := otlpAttributes(record.Attributes)
				name, ok := attributeString(attributes, "event.name")
				if !ok {
					continue
				}
				sessionID, _ := attributeString(attributes, "session.id")
				promptID, _ := attributeString(attributes, "prompt.id")
				sequence, _ := attributeInt(attributes, "event.sequence")
				events = append(events, Event{
					Name:        name,
					SessionID:   sessionID,
					PromptID:    promptID,
					TimestampMs: timestampMillis(record),
					Sequence:    sequence,
					Attributes:  attributes,
				})
			}
		}
	}
	return events
}
