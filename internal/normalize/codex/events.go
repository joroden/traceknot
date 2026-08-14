package codex

import (
	"sort"

	logspb "go.opentelemetry.io/proto/otlp/logs/v1"
)

const (
	eventConversationStarts = "codex.conversation_starts"
	eventUserPrompt         = "codex.user_prompt"
	eventToolResult         = "codex.tool_result"
	eventToolDecision       = "codex.tool_decision"
	eventSSEEvent           = "codex.sse_event"

	eventRolloutMessage = "codex.rollout_message"
	eventRolloutCall    = "codex.rollout_call"
)

type Event struct {
	Name           string
	ConversationID string
	TimestampMs    int64
	Attributes     map[string]any
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
				conversationID, _ := attributeString(attributes, "conversation.id")
				events = append(events, Event{
					Name:           name,
					ConversationID: conversationID,
					TimestampMs:    timestampMillis(record),
					Attributes:     attributes,
				})
			}
		}
	}
	return events
}

func toolName(event Event) string {
	name, _ := attributeString(event.Attributes, "tool_name")
	return name
}

func sortedEvents(events []Event) []Event {
	sorted := make([]Event, len(events))
	copy(sorted, events)
	sort.SliceStable(sorted, func(left, right int) bool {
		if sorted[left].TimestampMs != sorted[right].TimestampMs {
			return sorted[left].TimestampMs < sorted[right].TimestampMs
		}
		return sorted[left].Name < sorted[right].Name
	})
	return sorted
}

func firstEvent(events []Event, name string) *Event {
	for index := range events {
		if events[index].Name == name {
			return &events[index]
		}
	}
	return nil
}
