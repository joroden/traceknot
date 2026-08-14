package codex

import (
	"traceknot/internal/model"
	"traceknot/internal/normalize/shared"
	"traceknot/internal/ptr"
)

func (builder *Builder) sessionSeed(conversationID string, events []Event, starts *Event) *model.SessionSeed {
	started := starts.TimestampMs
	ended := started
	firstPrompt := ""
	for _, event := range events {
		if event.TimestampMs > ended {
			ended = event.TimestampMs
		}
		if firstPrompt == "" && event.Name == eventUserPrompt {
			firstPrompt, _ = attributeString(event.Attributes, "prompt")
		}
	}

	metadata := map[string]any{
		"signal": "codex",
	}
	for _, key := range []string{
		"model", "reasoning_effort", "reasoning_summary", "approval_policy",
		"sandbox_policy", "provider_name", "originator", "terminal.type",
	} {
		if value, ok := attributeString(starts.Attributes, key); ok {
			metadata[key] = value
		}
	}

	return &model.SessionSeed{
		SessionID:              shared.SessionID("conversation", conversationID),
		ExternalConversationID: ptr.String(conversationID),
		NativeSessionID:        ptr.String(conversationID),
		SessionIDSource:        "external_conversation_id",
		Provider:               "codex",
		Title:                  shared.Title(firstPrompt),
		ServiceName:            ptr.String("codex"),
		StartedAtUnixMs:        ptr.Int64(started),
		EndedAtUnixMs:          ptr.Int64(ended),
		Metadata:               metadata,
	}
}
