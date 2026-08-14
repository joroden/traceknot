package claude

import (
	"encoding/json"
	"sort"

	"traceknot/internal/model"
	"traceknot/internal/normalize/shared"
	"traceknot/internal/pricing"
	"traceknot/internal/ptr"
	"traceknot/internal/tokenize"
)

type Session struct {
	SessionID string
	Events    []Event
	Spans     []Span
}

type Builder struct {
	catalog   *pricing.Catalog
	estimator *tokenize.Estimator
}

func NewBuilder(catalog *pricing.Catalog, estimator *tokenize.Estimator) *Builder {
	return &Builder{catalog: catalog, estimator: estimator}
}

func (builder *Builder) buildSession(session *Session) (*model.SessionSeed, *model.SessionContent) {
	events := sortedEvents(session.Events)
	if len(events) == 0 {
		return nil, nil
	}
	cacheTierSplits, events := joinUsageSupplements(events)
	if len(events) == 0 {
		return nil, nil
	}
	links := buildAgentLinks(session.Spans)
	approvals := approvalsByToolUseID(events)

	seed := builder.sessionSeed(session.SessionID, events)
	content := builder.buildTurns(seed.SessionID, events, links, approvals, cacheTierSplits)

	if len(content.Chats) == 0 && len(content.ToolCalls) == 0 && len(content.Agents) == 0 {
		return nil, nil
	}
	return seed, content
}

func (builder *Builder) sessionSeed(sessionID string, events []Event) *model.SessionSeed {
	started := events[0].TimestampMs
	ended := started
	firstPrompt := ""
	title := ""
	systemPrompt := ""
	for _, event := range events {
		if event.TimestampMs > ended {
			ended = event.TimestampMs
		}
		if firstPrompt == "" && event.Name == eventUserPrompt {
			firstPrompt, _ = attributeString(event.Attributes, "prompt")
		}
		if title == "" && event.Name == eventAssistantResp {
			if source, _ := attributeString(event.Attributes, "query_source"); source == "generate_session_title" {
				if response, ok := attributeString(event.Attributes, "response"); ok {
					title = titleFromResponse(response)
				}
			}
		}
		if event.Name == eventAPIRequestBody {

			if candidate, ok := attributeString(event.Attributes, "system_prompt"); ok && len(candidate) > len(systemPrompt) {
				systemPrompt = candidate
			}
		}
	}
	if title == "" {
		title = shared.Title(firstPrompt)
	}

	dbSessionID := shared.SessionID("claude", sessionID)
	return &model.SessionSeed{
		SessionID:              dbSessionID,
		ExternalConversationID: ptr.String(sessionID),
		NativeSessionID:        ptr.String(sessionID),
		SessionIDSource:        "external_conversation_id",
		Provider:               "claude",
		Title:                  title,
		ServiceName:            ptr.String("claude-code"),
		StartedAtUnixMs:        ptr.Int64(started),
		EndedAtUnixMs:          ptr.Int64(ended),
		Metadata: map[string]any{
			"signal":        "claude",
			"system_prompt": systemPrompt,
		},
	}
}

func titleFromResponse(response string) string {
	var parsed struct {
		Title string `json:"title"`
	}
	if err := json.Unmarshal([]byte(response), &parsed); err != nil {
		return ""
	}
	return parsed.Title
}

func sortedEvents(events []Event) []Event {
	sorted := make([]Event, len(events))
	copy(sorted, events)
	sort.SliceStable(sorted, func(left, right int) bool {
		if sorted[left].TimestampMs != sorted[right].TimestampMs {
			return sorted[left].TimestampMs < sorted[right].TimestampMs
		}
		return sorted[left].Sequence < sorted[right].Sequence
	})
	return sorted
}
