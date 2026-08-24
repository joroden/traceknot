package claude

import (
	"strings"

	"traceknot/internal/model"
	"traceknot/internal/normalize/shared"
	"traceknot/internal/ptr"
)

const rootOwner = ""

type turnEvents struct {
	promptID string
	events   []Event
}

func (builder *Builder) buildTurns(sessionID string, events []Event, links *agentLinks, approvals map[string]approvalInfo, cacheTierSplits map[string]cacheTierSplit) *model.SessionContent {
	content := &model.SessionContent{}
	session := newSessionState(cacheTierSplits)
	for _, turn := range groupByTurn(events) {
		builder.buildTurn(sessionID, turn, links, approvals, content, session)
	}
	return content
}

func groupByTurn(events []Event) []*turnEvents {
	current := &turnEvents{}
	turns := []*turnEvents{current}
	for _, event := range events {
		if event.Name == eventUserPrompt {
			current = &turnEvents{promptID: event.PromptID}
			turns = append(turns, current)
		}
		current.events = append(current.events, event)
	}
	return turns
}

type sessionState struct {
	entryNodeByOwner       map[string]string
	agentSeedByID          map[string]*model.AgentSeed
	matchedAgents          map[string]bool
	spawnOrder             []string
	cacheTierSplits        map[string]cacheTierSplit
	internalUsageByToolUse map[string]*model.ChatSeed
}

func newSessionState(cacheTierSplits map[string]cacheTierSplit) *sessionState {
	return &sessionState{
		entryNodeByOwner:       map[string]string{},
		agentSeedByID:          map[string]*model.AgentSeed{},
		matchedAgents:          map[string]bool{},
		cacheTierSplits:        cacheTierSplits,
		internalUsageByToolUse: map[string]*model.ChatSeed{},
	}
}

type turnState struct {
	session            *sessionState
	currentStepByOwner map[string]*model.ChatSeed
	stepByRequestID    map[string]*model.ChatSeed
}

func newTurnState(session *sessionState, rootNodeID string) *turnState {
	session.entryNodeByOwner[rootOwner] = rootNodeID
	return &turnState{
		session:            session,
		currentStepByOwner: map[string]*model.ChatSeed{},
		stepByRequestID:    map[string]*model.ChatSeed{},
	}
}

func (builder *Builder) buildTurn(sessionID string, turn *turnEvents, links *agentLinks, approvals map[string]approvalInfo, content *model.SessionContent, session *sessionState) {
	if len(turn.events) == 0 {
		return
	}

	rootNodeID := ""
	eventsToWalk := turn.events
	turnContent := &model.SessionContent{}
	if turn.promptID != "" {
		if turn.events[0].Name != eventUserPrompt {
			return
		}
		userEvent := turn.events[0]
		prompt, _ := attributeString(userEvent.Attributes, "prompt")
		if prompt == "" {
			return
		}
		turnChat := builder.userChat(sessionID, userEvent, prompt)
		turnContent.Chats = []*model.ChatSeed{turnChat}
		rootNodeID = turnChat.NodeID
		eventsToWalk = turn.events[1:]
	}

	state := newTurnState(session, rootNodeID)
	builder.collectSpawns(sessionID, eventsToWalk, links, turnContent, state)
	completions := builder.walkTurnEvents(sessionID, eventsToWalk, links, approvals, turnContent, state)

	minChats := 0
	if turn.promptID != "" {
		minChats = 1
	}
	if len(turnContent.Chats) <= minChats && len(turnContent.ToolCalls) == 0 && len(turnContent.Agents) == 0 {
		return
	}

	matchSubagentCompletions(session.spawnOrder, session.agentSeedByID, session.matchedAgents, completions)
	content.Chats = append(content.Chats, turnContent.Chats...)
	content.ToolCalls = append(content.ToolCalls, turnContent.ToolCalls...)
	content.Agents = append(content.Agents, turnContent.Agents...)
}

func (builder *Builder) collectSpawns(sessionID string, events []Event, links *agentLinks, turnContent *model.SessionContent, state *turnState) {
	for _, event := range events {
		if event.Name != eventToolResult {
			continue
		}
		toolUseID, _ := attributeString(event.Attributes, "tool_use_id")
		toolName, _ := attributeString(event.Attributes, "tool_name")
		if toolUseID == "" || toolName != toolNameAgent {
			continue
		}
		agentID := links.agentIDBySpawnToolUse[toolUseID]
		if agentID == "" {
			continue
		}
		spawnNodeID := shared.NodeID(sessionID, "tool", toolUseID)
		agent := builder.agentSeedFromToolResult(sessionID, event, toolUseID, spawnNodeID)
		turnContent.Agents = append(turnContent.Agents, agent)
		state.session.entryNodeByOwner[agentID] = agent.NodeID
		state.session.agentSeedByID[agentID] = agent
		state.session.spawnOrder = append(state.session.spawnOrder, agentID)
	}
}

func (builder *Builder) walkTurnEvents(
	sessionID string,
	events []Event,
	links *agentLinks,
	approvals map[string]approvalInfo,
	turnContent *model.SessionContent,
	state *turnState,
) []Event {
	var completions []Event
	for _, event := range events {
		switch event.Name {
		case eventAPIRequest:
			requestID, ok := attributeString(event.Attributes, "request_id")
			if !ok {
				continue
			}
			querySource, _ := attributeString(event.Attributes, "query_source")
			owner := links.agentIDByRequestID[requestID]
			parentNodeID, ok := state.session.entryNodeByOwner[owner]
			if !ok {
				continue
			}
			toolUseID := ""
			if owner == rootOwner {
				toolUseID = links.parentToolUseIDByRequest[requestID]
			}
			step := builder.assistantStep(sessionID, event, requestID, parentNodeID, state.session.cacheTierSplits[requestID])
			state.stepByRequestID[requestID] = step
			switch {
			case toolUseID != "":
				state.session.internalUsageByToolUse[toolUseID] = step
			case owner == rootOwner && !isVisibleQuerySource(querySource):
				step.Name = chatNameMeta
				turnContent.Chats = append(turnContent.Chats, step)
			default:
				turnContent.Chats = append(turnContent.Chats, step)
				if isVisibleQuerySource(querySource) {
					state.currentStepByOwner[owner] = step
				}
			}

		case eventAssistantResp:
			requestID, _ := attributeString(event.Attributes, "request_id")
			step, ok := state.stepByRequestID[requestID]
			if !ok {
				continue
			}
			response, _ := attributeString(event.Attributes, "response")
			step.OutputText = response
			step.PreviewText = shared.Preview("assistant", response)
			if owner := links.agentIDByRequestID[requestID]; owner != rootOwner {
				if agent, ok := state.session.agentSeedByID[owner]; ok {
					agent.ResultSummary = response
				}
			}

		case eventToolResult:
			toolUseID, _ := attributeString(event.Attributes, "tool_use_id")
			toolName, _ := attributeString(event.Attributes, "tool_name")
			if toolUseID == "" || toolName == "" {
				continue
			}
			owner := links.ownerAgentIDByToolUse[toolUseID]
			parentNodeID, ok := currentStepOrEntryNode(state, owner)
			if !ok {
				continue
			}
			internalUsage := state.session.internalUsageByToolUse[toolUseID]
			tool := builder.toolCall(sessionID, event, toolUseID, toolName, parentNodeID, links, approvals, internalUsage)
			turnContent.ToolCalls = append(turnContent.ToolCalls, tool)

		case eventSubagentComplete:
			completions = append(completions, event)
		}
	}
	return completions
}

func currentStepOrEntryNode(state *turnState, owner string) (string, bool) {
	if step, ok := state.currentStepByOwner[owner]; ok {
		return step.NodeID, true
	}
	if nodeID, ok := state.session.entryNodeByOwner[owner]; ok {
		return nodeID, true
	}
	return "", false
}

func isVisibleQuerySource(source string) bool {
	return source == querySourceMain || strings.HasPrefix(source, querySourceAgentPrefix)
}

func (builder *Builder) userChat(sessionID string, event Event, prompt string) *model.ChatSeed {
	return &model.ChatSeed{
		NodeSeed: model.NodeSeed{
			NodeID:          shared.NodeID(sessionID, "user", event.PromptID),
			Kind:            model.NodeKindChat,
			Name:            "user",
			Status:          ptr.String("ok"),
			StartedAtUnixMs: ptr.Int64(event.TimestampMs),
			PreviewText:     shared.Preview("user", prompt),
			MetadataJSON:    shared.EventMetadata(eventUserPrompt),
		},
		PromptText: prompt,
	}
}

func (builder *Builder) assistantStep(sessionID string, event Event, requestID string, parentNodeID string, tierSplit cacheTierSplit) *model.ChatSeed {

	nonCached, _ := attributeInt(event.Attributes, "input_tokens")
	output, _ := attributeInt(event.Attributes, "output_tokens")
	cacheRead, _ := attributeInt(event.Attributes, "cache_read_tokens")
	cacheCreate, _ := attributeInt(event.Attributes, "cache_creation_tokens")
	input := nonCached + cacheRead + cacheCreate
	var modelPtr *string
	if modelName, ok := attributeString(event.Attributes, "model"); ok {
		modelPtr = &modelName
	}

	startedAt := event.TimestampMs
	duration, hasDuration := attributeInt(event.Attributes, "duration_ms")
	if hasDuration && duration > 0 {
		startedAt = event.TimestampMs - duration
	}

	var parentNodePtr *string
	if parentNodeID != "" {
		parentNodePtr = ptr.String(parentNodeID)
	}
	seed := &model.ChatSeed{
		NodeSeed: model.NodeSeed{
			NodeID:            shared.NodeID(sessionID, "assistant", requestID),
			ParentNodeID:      parentNodePtr,
			Kind:              model.NodeKindChat,
			Name:              "assistant",
			Model:             modelPtr,
			Status:            ptr.String("ok"),
			StartedAtUnixMs:   ptr.Int64(startedAt),
			PreviewText:       shared.Preview("assistant", ""),
			InputTokens:       input,
			CachedInputTokens: cacheRead,
			CacheWriteTokens:  cacheCreate,
			OutputTokens:      output,
			MetadataJSON:      shared.EventMetadata(eventAPIRequest),
		},
	}
	if hasDuration && duration > 0 {
		seed.EndedAtUnixMs = ptr.Int64(event.TimestampMs)
		seed.DurationMs = ptr.Float64(float64(duration))
	}
	var webSearchQueries int64
	if querySource, _ := attributeString(event.Attributes, "query_source"); querySource == querySourceWebSearch {
		webSearchQueries = 1
	}
	seed.Cost = shared.NodeCost(seed.Model, input, cacheRead, cacheCreate, tierSplit.cache1h, output, webSearchQueries, seed.StartedAtUnixMs, builder.catalog)
	return seed
}

func (builder *Builder) toolCall(
	sessionID string,
	event Event,
	toolUseID, toolName, parentNodeID string,
	links *agentLinks,
	approvals map[string]approvalInfo,
	internalUsage *model.ChatSeed,
) *model.ToolCallSeed {
	arguments, _ := attributeString(event.Attributes, "tool_input")
	success, _ := attributeBool(event.Attributes, "success")
	status := "ok"
	if !success {
		status = "error"
	}
	output := links.toolOutputByToolUse[toolUseID]
	if output == "" && internalUsage != nil {
		output = internalUsage.OutputText
	}

	startedAt := event.TimestampMs
	duration, hasDuration := attributeInt(event.Attributes, "duration_ms")
	if hasDuration && duration > 0 {
		startedAt = event.TimestampMs - duration
	}

	var parentNodePtr *string
	if parentNodeID != "" {
		parentNodePtr = ptr.String(parentNodeID)
	}
	seed := &model.ToolCallSeed{
		NodeSeed: model.NodeSeed{
			NodeID:          shared.NodeID(sessionID, "tool", toolUseID),
			ParentNodeID:    parentNodePtr,
			Kind:            model.NodeKindToolCall,
			Name:            toolName,
			Status:          ptr.String(status),
			StartedAtUnixMs: ptr.Int64(startedAt),
			PreviewText:     shared.Preview(toolName, arguments),
			MetadataJSON:    shared.EventMetadata(eventToolResult),
		},
		ToolName:      toolName,
		ToolCallID:    toolUseID,
		ArgumentsJSON: arguments,
		ResultText:    output,
	}
	if hasDuration && duration > 0 {
		seed.EndedAtUnixMs = ptr.Int64(event.TimestampMs)
		seed.DurationMs = ptr.Float64(float64(duration))
	}
	if internalUsage != nil {
		seed.Model = internalUsage.Model
		seed.InputTokens = internalUsage.InputTokens
		seed.CachedInputTokens = internalUsage.CachedInputTokens
		seed.CacheWriteTokens = internalUsage.CacheWriteTokens
		seed.OutputTokens = internalUsage.OutputTokens
		seed.Cost = internalUsage.Cost
	}
	if info, ok := approvals[toolUseID]; ok {
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
