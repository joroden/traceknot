package claude

import "traceknot/internal/jsonutil"

const (
	spanLLMRequest  = "claude_code.llm_request"
	spanTool        = "claude_code.tool"
	spawnToolName   = "Agent"
	toolOutputEvent = "tool.output"
)

type agentLinks struct {
	agentIDByRequestID    map[string]string
	spawnToolUseIDByAgent map[string]string
	agentIDBySpawnToolUse map[string]string
	ownerAgentIDByToolUse map[string]string
	toolOutputByToolUse   map[string]string
}

func buildAgentLinks(spans []Span) *agentLinks {
	byID := make(map[string]*Span, len(spans))
	for index := range spans {
		byID[spans[index].SpanID] = &spans[index]
	}

	links := &agentLinks{
		agentIDByRequestID:    make(map[string]string),
		spawnToolUseIDByAgent: make(map[string]string),
		agentIDBySpawnToolUse: make(map[string]string),
		ownerAgentIDByToolUse: make(map[string]string),
		toolOutputByToolUse:   make(map[string]string),
	}

	for index := range spans {
		span := &spans[index]
		switch span.Name {
		case spanLLMRequest:
			requestID, ok := attributeString(span.Attributes, "request_id")
			if !ok {
				continue
			}
			agentID, _ := attributeString(span.Attributes, "agent_id")
			links.agentIDByRequestID[requestID] = agentID
			if agentID != "" {
				spawnToolUseID := nearestAncestorToolUseID(span.ParentSpanID, byID)
				links.spawnToolUseIDByAgent[agentID] = spawnToolUseID
				if spawnToolUseID != "" {
					links.agentIDBySpawnToolUse[spawnToolUseID] = agentID
				}
			}
		case spanTool:
			toolUseID, ok := attributeString(span.Attributes, "tool_use_id")
			if !ok {
				continue
			}

			agentID, _ := attributeString(span.Attributes, "agent_id")
			links.ownerAgentIDByToolUse[toolUseID] = agentID
			if output := toolOutputFromEvents(span.Events); output != "" {
				links.toolOutputByToolUse[toolUseID] = output
			}
		}
	}
	return links
}

func nearestAncestorToolUseID(parentSpanID string, byID map[string]*Span) string {
	seen := map[string]bool{}
	for parentSpanID != "" && byID[parentSpanID] != nil && !seen[parentSpanID] {
		seen[parentSpanID] = true
		span := byID[parentSpanID]
		if span.Name == spanTool {
			if name, _ := attributeString(span.Attributes, "tool_name"); name == spawnToolName {
				toolUseID, _ := attributeString(span.Attributes, "tool_use_id")
				return toolUseID
			}
		}
		parentSpanID = span.ParentSpanID
	}
	return ""
}

var toolOutputContentKeys = []string{"output", "content", "diff"}

func toolOutputFromEvents(events []SpanEvent) string {
	for _, event := range events {
		if event.Name != toolOutputEvent {
			continue
		}
		for _, key := range toolOutputContentKeys {
			if value, ok := attributeString(event.Attributes, key); ok {
				return value
			}
		}
		if len(event.Attributes) > 0 {
			return jsonutil.ToCanonicalJSON(event.Attributes)
		}
	}
	return ""
}

type approvalInfo struct {
	decision string
	source   string
}

func approvalsByToolUseID(events []Event) map[string]approvalInfo {
	output := make(map[string]approvalInfo)
	for _, event := range events {
		if event.Name != eventToolDecision {
			continue
		}
		toolUseID, ok := attributeString(event.Attributes, "tool_use_id")
		if !ok || toolUseID == "" {
			continue
		}
		decision, _ := attributeString(event.Attributes, "decision")
		source, _ := attributeString(event.Attributes, "source")
		output[toolUseID] = approvalInfo{decision: decision, source: source}
	}
	return output
}
