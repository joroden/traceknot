package copilot

import (
	"sort"

	"traceknot/internal/model"
	"traceknot/internal/normalize/shared"
	"traceknot/internal/pricing"
	"traceknot/internal/ptr"
	"traceknot/internal/tokenize"
)

const (
	opInvokeAgent = "invoke_agent"
	opChat        = "chat"
	opExecuteTool = "execute_tool"
	opExecuteHook = "execute_hook"
	toolNameTask  = "task"
)

type Session struct {
	TraceID string
	Spans   []Span
}

type Builder struct {
	catalog   *pricing.Catalog
	estimator *tokenize.Estimator
}

func NewBuilder(catalog *pricing.Catalog, estimator *tokenize.Estimator) *Builder {
	return &Builder{catalog: catalog, estimator: estimator}
}

func (builder *Builder) buildSession(session *Session) (*model.SessionSeed, *model.SessionContent) {

	spans := excludeHookSpans(session.Spans)
	if len(spans) == 0 {
		return nil, nil
	}

	rootSpanID := findRootSpanID(spans)
	if rootSpanID == "" {
		return nil, nil
	}
	byParent := groupByParent(spans)

	seed := builder.sessionSeed(session.TraceID, findSpanByID(spans, rootSpanID), spans)
	content := &model.SessionContent{}
	builder.walk(seed.SessionID, rootSpanID, nil, true, byParent, content)
	if len(content.Chats)+len(content.ToolCalls)+len(content.Agents) == 0 {

		return nil, nil
	}
	return seed, content
}

func findRootSpanID(spans []Span) string {
	known := make(map[string]bool, len(spans))
	for _, span := range spans {
		known[span.SpanID] = true
	}
	for _, span := range spans {
		if span.ParentSpanID != "" && !known[span.ParentSpanID] {
			return span.ParentSpanID
		}
	}

	for index := range spans {
		if spans[index].ParentSpanID == "" {
			return spans[index].SpanID
		}
	}
	return ""
}

func excludeHookSpans(spans []Span) []Span {
	filtered := make([]Span, 0, len(spans))
	for _, span := range spans {
		if op, _ := attributeString(span.Attributes, "gen_ai.operation.name"); op == opExecuteHook {
			continue
		}
		filtered = append(filtered, span)
	}
	return filtered
}

func findSpanByID(spans []Span, spanID string) *Span {
	for index := range spans {
		if spans[index].SpanID == spanID {
			return &spans[index]
		}
	}
	return nil
}

func groupByParent(spans []Span) map[string][]*Span {
	byParent := make(map[string][]*Span)
	for index := range spans {
		span := &spans[index]
		if span.ParentSpanID == "" {
			continue
		}
		byParent[span.ParentSpanID] = append(byParent[span.ParentSpanID], span)
	}
	for parent := range byParent {
		children := byParent[parent]
		sort.SliceStable(children, func(left, right int) bool {
			return children[left].TimestampMs < children[right].TimestampMs
		})
	}
	return byParent
}

func (builder *Builder) sessionSeed(traceID string, root *Span, spans []Span) *model.SessionSeed {
	started, ended := spans[0].TimestampMs, spans[0].TimestampMs
	conversationID, systemPrompt := "", ""
	for _, span := range spans {
		if span.TimestampMs < started {
			started = span.TimestampMs
		}
		if span.TimestampMs > ended {
			ended = span.TimestampMs
		}
		if conversationID == "" {
			conversationID, _ = attributeString(span.Attributes, "gen_ai.conversation.id")
		}
		if systemPrompt == "" {
			systemPrompt, _ = attributeString(span.Attributes, "gen_ai.system_instructions")
		}
	}

	title := ""
	if root != nil {
		if input, ok := attributeString(root.Attributes, "gen_ai.input.messages"); ok {
			if text := stripInjectedContext(firstRoleText(input, "user")); text != "" {
				title = shared.Title(text)
			}
		}
	}

	metadata := map[string]any{
		"signal":        "copilot",
		"system_prompt": systemPrompt,
	}
	if root != nil {
		for _, key := range []string{
			"github.copilot.git.repository", "github.copilot.git.branch",
			"github.copilot.git.commit_sha", "github.copilot.github.org",
			"github.copilot.turn_count", "gen_ai.agent.version",

			"copilot_chat.turn_count",
		} {
			if value, ok := attributeString(root.Attributes, key); ok {
				metadata[key] = value
			}
		}
	}

	dbSessionID := shared.SessionID("copilot", traceID)
	return &model.SessionSeed{
		SessionID:              dbSessionID,
		ExternalConversationID: ptr.String(conversationID),
		NativeSessionID:        ptr.String(traceID),
		SessionIDSource:        "trace_id",
		Provider:               "copilot",
		Title:                  title,
		ServiceName:            ptr.String("copilot-cli"),
		StartedAtUnixMs:        ptr.Int64(started),
		EndedAtUnixMs:          ptr.Int64(ended),
		Metadata:               metadata,
	}
}
