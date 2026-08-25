package codex

import (
	"traceknot/internal/model"
	"traceknot/internal/normalize/shared"
	"traceknot/internal/pricing"
	"traceknot/internal/ptr"
	"traceknot/internal/tokenize"
)

const spawnToolName = "multi_agent_v1spawn_agent"

const execWrapperTool = "exec"

type Builder struct {
	catalog   *pricing.Catalog
	estimator *tokenize.Estimator
}

func NewBuilder(catalog *pricing.Catalog, estimator *tokenize.Estimator) *Builder {
	return &Builder{
		catalog:   catalog,
		estimator: estimator,
	}
}

func (builder *Builder) Build(all map[string][]Event, touchedIDs []string) []shared.BuildResult {
	rootID := rootConversation(all)
	if rootID == "" {
		return nil
	}
	subagents := subagentMap(all)
	approvals := approvalsByCallID(all)
	waits := waitResultsByAgentID(all)
	affected := affectedRoots(subagents, touchedIDs)
	var results []shared.BuildResult
	for id, events := range all {

		if id == "" || subagents[id] != nil || !affected[id] {
			continue
		}
		seed, content := builder.buildSession(id, events, all, subagents, approvals, waits)
		if seed == nil {
			continue
		}
		results = append(results, shared.BuildResult{Seed: seed, Content: content})
	}
	return results
}

func (builder *Builder) buildSession(
	conversationID string,
	events []Event,
	all map[string][]Event,
	subagents map[string]*subagentLink,
	approvals map[string]approvalInfo,
	waits map[string]string,
) (*model.SessionSeed, *model.SessionContent) {
	sorted := sortedEvents(events)
	if len(sorted) == 0 {
		return nil, nil
	}

	if firstEvent(sorted, eventUserPrompt) == nil {
		return nil, nil
	}

	starts := firstEvent(sorted, eventConversationStarts)
	if starts == nil {
		starts = &sorted[0]
	}

	session := builder.sessionSeed(conversationID, sorted, starts)
	content := &model.SessionContent{}
	rootAgent := builder.rootAgent(starts, session.SessionID)
	content.Agents = append(content.Agents, rootAgent)
	tree := builder.buildTree(sorted, subagents, all, approvals, waits, session.SessionID, ptr.String(rootAgent.NodeID))
	content.Chats = tree.chats
	content.ToolCalls = tree.tools
	content.Agents = append(content.Agents, tree.agents...)
	return session, content
}

func (builder *Builder) rootAgent(starts *Event, sessionID string) *model.AgentSeed {
	return &model.AgentSeed{
		NodeSeed: model.NodeSeed{
			NodeID:          shared.AgentID(sessionID, "invoke"),
			Kind:            model.NodeKindAgent,
			Name:            "main agent",
			Model:           builder.model(*starts),
			Status:          ptr.String("ok"),
			StartedAtUnixMs: ptr.Int64(starts.TimestampMs),
			PreviewText:     "main agent",
			MetadataJSON:    shared.EventMetadata("codex.invoke_agent"),
		},
		AgentName: "main agent",
	}
}
