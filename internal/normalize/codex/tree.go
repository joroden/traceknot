package codex

import (
	"traceknot/internal/model"
	"traceknot/internal/normalize/shared"
	"traceknot/internal/ptr"
)

type treeResult struct {
	chats  []*model.ChatSeed
	tools  []*model.ToolCallSeed
	agents []*model.AgentSeed
}

type treeContext struct {
	subagents map[string]*subagentLink
	all       map[string][]Event
	approvals map[string]approvalInfo
	waits     map[string]string
	sessionID string
}

type treeBuildState struct {
	result            *treeResult
	containerByCallID map[string]*model.ChatSeed
	sawRealTool       bool
	pendingTools      []Event
}

func (builder *Builder) buildTree(
	events []Event,
	subagents map[string]*subagentLink,
	all map[string][]Event,
	approvals map[string]approvalInfo,
	waits map[string]string,
	sessionID string,
	agentRoot *string,
) *treeResult {
	ctx := &treeContext{subagents: subagents, all: all, approvals: approvals, waits: waits, sessionID: sessionID}
	state := &treeBuildState{result: &treeResult{}, containerByCallID: make(map[string]*model.ChatSeed)}
	sorted := sortedEvents(events)

	var currentUser *model.ChatSeed
	var currentStep *model.ChatSeed
	var stepIndex int

	var turns []*turnSpan
	var currentTurn *turnSpan

	for _, event := range sorted {
		switch event.Name {
		case eventUserPrompt:
			builder.flushPendingTools(currentUser, ctx, state)
			currentUser = builder.userChat(event, sessionID, agentRoot)
			currentStep = nil
			stepIndex = 0
			state.sawRealTool = false
			currentTurn = nil
			if currentUser != nil {
				state.result.chats = append(state.result.chats, currentUser)
				currentTurn = &turnSpan{userNode: currentUser}
				turns = append(turns, currentTurn)
			}
		case eventSSEEvent:
			if currentUser == nil {
				continue
			}
			step := builder.assistantStep(event, sessionID, currentUser, stepIndex)
			if step == nil {
				continue
			}
			stepIndex++
			currentStep = step
			state.sawRealTool = false
			state.result.chats = append(state.result.chats, step)
			if currentTurn != nil {
				currentTurn.lastStep = step
			}
			builder.flushPendingTools(step, ctx, state)
		case eventToolResult:
			if currentStep == nil {
				state.pendingTools = append(state.pendingTools, event)
				continue
			}
			builder.processTool(event, currentStep, ctx, state)
		}
	}
	builder.flushPendingTools(currentUser, ctx, state)
	resolveRolloutMessages(turns, sorted, state.containerByCallID)
	return state.result
}

func (builder *Builder) processTool(event Event, container *model.ChatSeed, ctx *treeContext, state *treeBuildState) {
	toolName, _ := attributeString(event.Attributes, "tool_name")
	callID, _ := attributeString(event.Attributes, "call_id")
	if callID != "" && container != nil {
		state.containerByCallID[callID] = container
	}
	if toolName == execWrapperTool && state.sawRealTool {
		state.sawRealTool = false
		return
	}
	tool := builder.toolCall(event, ctx.sessionID, container, ctx.approvals)
	if tool == nil {
		return
	}
	state.sawRealTool = toolName != execWrapperTool
	state.result.tools = append(state.result.tools, tool)
	if link := findSpawnLink(ctx.subagents, tool.ToolCallID); link != nil {
		agent := builder.agentNode(event, link, tool.NodeID, ctx.waits)
		state.result.agents = append(state.result.agents, agent)
		if subEvents, ok := ctx.all[link.conversationID]; ok {
			sub := builder.buildTree(subEvents, ctx.subagents, ctx.all, ctx.approvals, ctx.waits, ctx.sessionID, ptr.String(agent.NodeID))
			state.result.chats = append(state.result.chats, sub.chats...)
			state.result.tools = append(state.result.tools, sub.tools...)
			state.result.agents = append(state.result.agents, sub.agents...)
		}
	}
}

func (builder *Builder) flushPendingTools(container *model.ChatSeed, ctx *treeContext, state *treeBuildState) {
	for _, pending := range state.pendingTools {
		builder.processTool(pending, container, ctx, state)
	}
	state.pendingTools = nil
}

func (builder *Builder) userChat(event Event, sessionID string, agentRoot *string) *model.ChatSeed {
	prompt, _ := attributeString(event.Attributes, "prompt")
	if prompt == "" {
		return nil
	}
	return &model.ChatSeed{
		NodeSeed: model.NodeSeed{
			NodeID:          shared.NodeID(sessionID, "user", event.ConversationID, itoa(event.TimestampMs)),
			ParentNodeID:    agentRoot,
			Kind:            model.NodeKindChat,
			Name:            "user",
			Model:           builder.model(event),
			Status:          ptr.String("ok"),
			StartedAtUnixMs: ptr.Int64(event.TimestampMs),
			PreviewText:     shared.Preview("user", prompt),
			MetadataJSON:    shared.EventMetadata("codex.user_prompt"),
		},
		PromptText: prompt,
	}
}

func (builder *Builder) assistantStep(event Event, sessionID string, parent *model.ChatSeed, stepIndex int) *model.ChatSeed {
	kind, _ := attributeString(event.Attributes, "event.kind")
	if kind != "response.completed" {
		return nil
	}
	rawInput, _ := attributeInt(event.Attributes, "input_token_count")
	output, _ := attributeInt(event.Attributes, "output_token_count")
	cached, _ := attributeInt(event.Attributes, "cached_token_count")
	cacheWrite, _ := attributeInt(event.Attributes, "cache_write_token_count")
	reasoning, _ := attributeInt(event.Attributes, "reasoning_token_count")
	if rawInput == 0 && output == 0 && reasoning == 0 && cached == 0 && cacheWrite == 0 {
		return nil
	}

	input := rawInput

	effectiveModel := builder.model(event)
	if effectiveModel == nil {
		effectiveModel = parent.Model
	}
	seed := &model.ChatSeed{
		NodeSeed: model.NodeSeed{
			NodeID:            shared.NodeID(sessionID, "assistant", event.ConversationID, itoa(event.TimestampMs), itoa(int64(stepIndex))),
			ParentNodeID:      ptr.String(parent.NodeID),
			Kind:              model.NodeKindChat,
			Name:              "assistant",
			Model:             effectiveModel,
			Status:            ptr.String("ok"),
			StartedAtUnixMs:   ptr.Int64(event.TimestampMs),
			PreviewText:       shared.Preview("assistant", ""),
			InputTokens:       input,
			CachedInputTokens: cached,
			CacheWriteTokens:  cacheWrite,
			OutputTokens:      output,
			ReasoningTokens:   reasoning,
			MetadataJSON:      shared.EventMetadata("codex.sse_event"),
		},
	}
	seed.Cost = shared.NodeCost(seed.Model, input, cached, cacheWrite, 0, output, 0, seed.StartedAtUnixMs, builder.catalog)
	return seed
}

func (builder *Builder) toolCall(event Event, sessionID string, container *model.ChatSeed, approvals map[string]approvalInfo) *model.ToolCallSeed {
	toolName, _ := attributeString(event.Attributes, "tool_name")
	callID, _ := attributeString(event.Attributes, "call_id")
	if toolName == "" || callID == "" {
		return nil
	}
	arguments, _ := attributeString(event.Attributes, "arguments")
	output, _ := attributeString(event.Attributes, "output")
	success, _ := attributeBool(event.Attributes, "success")
	status := "ok"
	if !success {
		status = "error"
	}

	startedAt := event.TimestampMs
	duration, hasDuration := attributeInt(event.Attributes, "duration_ms")
	if hasDuration && duration > 0 {
		startedAt = event.TimestampMs - duration
	}

	seed := &model.ToolCallSeed{
		NodeSeed: model.NodeSeed{
			NodeID:          shared.NodeID(sessionID, "tool", toolName, callID),
			Kind:            model.NodeKindToolCall,
			Name:            toolName,
			Status:          ptr.String(status),
			StartedAtUnixMs: ptr.Int64(startedAt),
			PreviewText:     shared.Preview(toolName, arguments),
			MetadataJSON:    shared.EventMetadata("codex.tool_result"),
		},
		ToolName:      toolName,
		ToolCallID:    callID,
		ArgumentsJSON: arguments,
		ResultText:    output,
	}
	if container != nil {
		seed.ParentNodeID = ptr.String(container.NodeID)
	}
	if hasDuration && duration > 0 {
		seed.EndedAtUnixMs = ptr.Int64(event.TimestampMs)
		seed.DurationMs = ptr.Float64(float64(duration))
	}
	if info, ok := approvals[callID]; ok {
		if info.decision != "" {
			seed.ApprovalDecision = ptr.String(info.decision)
		}
		if info.source != "" {
			seed.ApprovalSource = ptr.String(info.source)
		}
	}
	seed.EstimatedInputTokens, seed.EstimatedOutputTokens, seed.TokenEstimateMethod =
		estimateToolTokens(arguments, output, builder.estimator)
	return seed
}

func (builder *Builder) agentNode(event Event, link *subagentLink, spawnNodeID string, waits map[string]string) *model.AgentSeed {
	nickname := link.nickname
	if nickname == "" {
		nickname = "subagent"
	}
	return &model.AgentSeed{
		NodeSeed: model.NodeSeed{
			NodeID:          shared.AgentID(link.conversationID),
			ParentNodeID:    ptr.String(spawnNodeID),
			Kind:            model.NodeKindAgent,
			Name:            "subagent",
			Model:           builder.model(event),
			Status:          ptr.String("ok"),
			StartedAtUnixMs: ptr.Int64(event.TimestampMs),
			PreviewText:     "subagent · " + nickname,
			MetadataJSON:    shared.EventMetadata("codex.subagent"),
		},
		AgentName:       nickname,
		SpawnPrompt:     link.spawnPrompt,
		SpawnToolCallID: link.spawnCallID,
		ResultSummary:   waits[link.conversationID],
	}
}
