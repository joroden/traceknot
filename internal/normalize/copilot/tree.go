package copilot

import (
	"encoding/json"

	"traceknot/internal/model"
	"traceknot/internal/normalize/shared"
	"traceknot/internal/ptr"
)

func (builder *Builder) walk(sessionID string, containerSpanID string, parentNodeID *string, isRoot bool, byParent map[string][]*Span, content *model.SessionContent) {
	var currentChatNodeID *string
	for _, span := range byParent[containerSpanID] {
		op, _ := attributeString(span.Attributes, "gen_ai.operation.name")
		switch op {
		case opChat:
			chat := builder.chatSeed(sessionID, span, parentNodeID, isRoot)
			content.Chats = append(content.Chats, chat)
			currentChatNodeID = ptr.String(chat.NodeID)

		case opExecuteTool:
			toolParentNodeID := parentNodeID
			if currentChatNodeID != nil {
				toolParentNodeID = currentChatNodeID
			}
			tool := builder.toolSeed(sessionID, span, toolParentNodeID)
			content.ToolCalls = append(content.ToolCalls, tool)

			toolName, _ := attributeString(span.Attributes, "gen_ai.tool.name")
			if toolName != toolNameTask {
				continue
			}
			for _, child := range byParent[span.SpanID] {
				childOp, _ := attributeString(child.Attributes, "gen_ai.operation.name")
				if childOp != opInvokeAgent {
					continue
				}
				agent := builder.agentSeed(sessionID, child, span, tool.NodeID)
				content.Agents = append(content.Agents, agent)
				builder.walk(sessionID, child.SpanID, ptr.String(agent.NodeID), false, byParent, content)
			}
		}
	}
}

func (builder *Builder) chatSeed(sessionID string, span *Span, parentNodeID *string, isRoot bool) *model.ChatSeed {
	name := "assistant"
	if isRoot {
		name = "user"
	}
	input, _ := attributeString(span.Attributes, "gen_ai.input.messages")
	output, _ := attributeString(span.Attributes, "gen_ai.output.messages")
	promptText := stripInjectedContext(lastRoleText(input, "user"))
	outputText := lastRoleText(output, "assistant")

	modelName, ok := attributeString(span.Attributes, "gen_ai.response.model")
	if !ok {
		modelName, ok = attributeString(span.Attributes, "gen_ai.request.model")
	}
	var modelPtr *string
	if ok {
		modelPtr = &modelName
	}

	status := "ok"
	if _, failed := attributeString(span.Attributes, "error.type"); failed {
		status = "error"
	}

	input64, _ := attributeInt(span.Attributes, "gen_ai.usage.input_tokens")
	output64, _ := attributeInt(span.Attributes, "gen_ai.usage.output_tokens")
	cacheRead, _ := attributeInt(span.Attributes, "gen_ai.usage.cache_read.input_tokens")

	cacheCreate, _ := attributeInt(span.Attributes, "gen_ai.usage.cache_creation.input_tokens")
	reasoning, _ := attributeInt(span.Attributes, "gen_ai.usage.reasoning.output_tokens")
	systemText, _ := attributeString(span.Attributes, "gen_ai.system_instructions")

	seed := &model.ChatSeed{
		NodeSeed: model.NodeSeed{
			NodeID:            shared.NodeID(sessionID, "span", span.SpanID),
			ParentNodeID:      parentNodeID,
			Kind:              model.NodeKindChat,
			Name:              name,
			Model:             modelPtr,
			Status:            ptr.String(status),
			StartedAtUnixMs:   ptr.Int64(span.TimestampMs),
			PreviewText:       shared.Preview("assistant", outputText),
			InputTokens:       input64,
			CachedInputTokens: cacheRead,
			CacheWriteTokens:  cacheCreate,
			OutputTokens:      output64,
			ReasoningTokens:   reasoning,
			MetadataJSON:      shared.EventMetadata("chat"),
		},
		SystemText: systemText,
		PromptText: promptText,
		OutputText: outputText,
	}
	seed.Cost = shared.NodeCost(seed.Model, input64, cacheRead, cacheCreate, 0, output64, 0, seed.StartedAtUnixMs, builder.catalog)
	return seed
}

func (builder *Builder) toolSeed(sessionID string, span *Span, parentNodeID *string) *model.ToolCallSeed {
	toolName, _ := attributeString(span.Attributes, "gen_ai.tool.name")
	toolCallID, _ := attributeString(span.Attributes, "gen_ai.tool.call.id")
	arguments, _ := attributeString(span.Attributes, "gen_ai.tool.call.arguments")
	result, _ := attributeString(span.Attributes, "gen_ai.tool.call.result")

	status := "ok"
	if _, failed := attributeString(span.Attributes, "error.type"); failed {
		status = "error"
	}

	seed := &model.ToolCallSeed{
		NodeSeed: model.NodeSeed{
			NodeID:          shared.NodeID(sessionID, "span", span.SpanID),
			ParentNodeID:    parentNodeID,
			Kind:            model.NodeKindToolCall,
			Name:            toolName,
			Status:          ptr.String(status),
			StartedAtUnixMs: ptr.Int64(span.TimestampMs),
			PreviewText:     shared.Preview(toolName, arguments),
			MetadataJSON:    shared.EventMetadata("execute_tool"),
		},
		ToolName:      toolName,
		ToolCallID:    toolCallID,
		ArgumentsJSON: arguments,
		ResultText:    result,
	}
	seed.EstimatedInputTokens, seed.EstimatedOutputTokens, seed.TokenEstimateMethod =
		estimateToolTokens(arguments, result, builder.estimator)
	return seed
}

func (builder *Builder) agentSeed(sessionID string, agentSpan *Span, spawnSpan *Span, spawnNodeID string) *model.AgentSeed {
	arguments, _ := attributeString(spawnSpan.Attributes, "gen_ai.tool.call.arguments")
	var parsed struct {
		Description string `json:"description"`
		Prompt      string `json:"prompt"`
		AgentType   string `json:"agent_type"`
		Name        string `json:"name"`
	}
	_ = json.Unmarshal([]byte(arguments), &parsed)

	name := parsed.Name
	if name == "" {
		name = parsed.AgentType
	}
	preview := name
	if parsed.Description != "" {
		preview = parsed.Description
	}

	status := "ok"
	if _, failed := attributeString(agentSpan.Attributes, "error.type"); failed {
		status = "error"
	}

	output, _ := attributeString(agentSpan.Attributes, "gen_ai.output.messages")

	return &model.AgentSeed{
		NodeSeed: model.NodeSeed{
			NodeID:          shared.AgentID(sessionID, spawnSpan.SpanID),
			ParentNodeID:    ptr.String(spawnNodeID),
			Kind:            model.NodeKindAgent,
			Name:            "subagent",
			Status:          ptr.String(status),
			StartedAtUnixMs: ptr.Int64(agentSpan.TimestampMs),
			PreviewText:     "subagent · " + preview,
			MetadataJSON:    shared.EventMetadata("invoke_agent"),
		},
		AgentName:       name,
		AgentType:       parsed.AgentType,
		SpawnPrompt:     parsed.Prompt,
		SpawnToolCallID: spawnSpan.SpanID,
		ResultSummary:   lastRoleText(output, "assistant"),
	}
}
