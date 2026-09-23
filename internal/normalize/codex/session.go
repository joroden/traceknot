package codex

import (
	"strings"

	"traceknot/internal/model"
	"traceknot/internal/normalize/shared"
	"traceknot/internal/ptr"
)

func (builder *Builder) sessionSeed(conversationID string, events []Event, starts *Event) *model.SessionSeed {
	started := starts.TimestampMs
	ended := started
	firstPrompt := ""
	title := ""
	for _, event := range events {
		if event.Name != eventSessionTitle && event.TimestampMs > ended {
			ended = event.TimestampMs
		}
		if firstPrompt == "" && event.Name == eventUserPrompt {
			firstPrompt, _ = attributeString(event.Attributes, "prompt")
		}
		if event.Name == eventSessionTitle {
			title, _ = attributeString(event.Attributes, "title")
		}
	}
	if title == "" {
		title = shared.Title(firstPrompt)
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
		Title:                  title,
		ServiceName:            ptr.String("codex"),
		StartedAtUnixMs:        ptr.Int64(started),
		EndedAtUnixMs:          ptr.Int64(ended),
		Metadata:               metadata,
	}
}

func isTitleGenerationPrompt(event Event) bool {
	prompt, _ := attributeString(event.Attributes, "prompt")
	return strings.HasPrefix(prompt, "Generate a concise, single-line task title") &&
		strings.Contains(prompt, "\n\nUser prompt:\n")
}
