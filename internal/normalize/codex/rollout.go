package codex

import (
	"encoding/json"
	"time"

	"traceknot/internal/normalize/shared"
)

type rolloutLine struct {
	Timestamp string          `json:"timestamp"`
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
}

type rolloutSessionMeta struct {
	ThreadID string `json:"id"`
}

type rolloutContentPart struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type rolloutResponseItem struct {
	Type    string               `json:"type"`
	Role    string               `json:"role"`
	Phase   string               `json:"phase"`
	Content []rolloutContentPart `json:"content"`
	CallID  string               `json:"call_id"`
}

func RolloutConversationID(line []byte) (string, bool) {
	var parsed rolloutLine
	if err := json.Unmarshal(line, &parsed); err != nil || parsed.Type != "session_meta" {
		return "", false
	}
	var meta rolloutSessionMeta
	if err := json.Unmarshal(parsed.Payload, &meta); err != nil || meta.ThreadID == "" {
		return "", false
	}
	return meta.ThreadID, true
}

func ParseRolloutLine(conversationID string, line []byte) (Event, bool) {
	var parsed rolloutLine
	if err := json.Unmarshal(line, &parsed); err != nil || parsed.Type != "response_item" {
		return Event{}, false
	}
	var item rolloutResponseItem
	if err := json.Unmarshal(parsed.Payload, &item); err != nil {
		return Event{}, false
	}
	timestampMs, ok := parseRolloutTimestamp(parsed.Timestamp)
	if !ok {
		return Event{}, false
	}
	switch {
	case item.Type == "message" && item.Role == "assistant":
		text := rolloutMessageText(item.Content)
		if text == "" {
			return Event{}, false
		}
		return RolloutMessageEvent(conversationID, timestampMs, item.Phase, text), true
	case item.Type == "custom_tool_call" && item.CallID != "":
		return RolloutCallEvent(conversationID, timestampMs, item.CallID), true
	default:
		return Event{}, false
	}
}

func rolloutMessageText(parts []rolloutContentPart) string {
	var text string
	for _, part := range parts {
		if part.Type != "output_text" || part.Text == "" {
			continue
		}
		if text != "" {
			text += "\n"
		}
		text += part.Text
	}
	return text
}

func parseRolloutTimestamp(value string) (int64, bool) {
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return 0, false
	}
	return parsed.UnixMilli(), true
}

func RolloutMessageEvent(conversationID string, timestampMs int64, phase, text string) Event {
	return Event{
		Name:           eventRolloutMessage,
		ConversationID: conversationID,
		TimestampMs:    timestampMs,
		Attributes: map[string]any{
			"text":  text,
			"phase": phase,
		},
	}
}

func RolloutCallEvent(conversationID string, timestampMs int64, callID string) Event {
	return Event{
		Name:           eventRolloutCall,
		ConversationID: conversationID,
		TimestampMs:    timestampMs,
		Attributes: map[string]any{
			"call_id": callID,
		},
	}
}

func RolloutRecord(event Event) (shared.RawRecord, bool) {
	payload, err := json.Marshal(event)
	if err != nil {
		return shared.RawRecord{}, false
	}
	return shared.RawRecord{
		NativeID:    event.ConversationID,
		Signal:      "rollout",
		DedupKey:    contentHash(event.Name, event.ConversationID, event.TimestampMs, event.Attributes),
		TimestampMs: event.TimestampMs,
		PayloadJSON: string(payload),
	}, true
}
