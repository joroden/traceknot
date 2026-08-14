package claude

import (
	"encoding/json"

	"traceknot/internal/jsonutil"
	"traceknot/internal/model"
	"traceknot/internal/normalize/shared"
	"traceknot/internal/ptr"
)

func (builder *Builder) agentSeedFromToolResult(sessionID string, event Event, toolUseID, spawnNodeID string) *model.AgentSeed {
	toolInput, _ := attributeString(event.Attributes, "tool_input")
	var parsed struct {
		Description  string `json:"description"`
		Prompt       string `json:"prompt"`
		SubagentType string `json:"subagent_type"`
	}
	_ = json.Unmarshal([]byte(toolInput), &parsed)

	name := parsed.SubagentType
	if name == "" {
		name = "subagent"
	}
	preview := name
	if parsed.Description != "" {
		preview = parsed.Description
	}

	return &model.AgentSeed{
		NodeSeed: model.NodeSeed{
			NodeID:          shared.AgentID(sessionID, toolUseID),
			ParentNodeID:    ptr.String(spawnNodeID),
			Kind:            model.NodeKindAgent,
			Name:            "subagent",
			Status:          ptr.String("ok"),
			StartedAtUnixMs: ptr.Int64(event.TimestampMs),
			PreviewText:     "subagent · " + preview,
			MetadataJSON:    shared.EventMetadata("tool_result:Agent"),
		},
		AgentName:       name,
		AgentType:       parsed.SubagentType,
		SpawnPrompt:     parsed.Prompt,
		SpawnToolCallID: toolUseID,
	}
}

func matchSubagentCompletions(spawnOrder []string, agentSeedByID map[string]*model.AgentSeed, matched map[string]bool, completions []Event) {
	for _, completion := range completions {
		agentType, _ := attributeString(completion.Attributes, "agent_type")
		for _, agentID := range spawnOrder {
			if matched[agentID] {
				continue
			}
			agent, ok := agentSeedByID[agentID]
			if !ok || agent.AgentType != agentType {
				continue
			}
			applySubagentCompletion(agent, completion)
			matched[agentID] = true
			break
		}
	}
}

func applySubagentCompletion(agent *model.AgentSeed, event Event) {
	if duration, ok := attributeInt(event.Attributes, "duration_ms"); ok && duration > 0 {
		agent.EndedAtUnixMs = ptr.Int64(event.TimestampMs)
		agent.DurationMs = ptr.Float64(float64(duration))
	}
	totalTokens, _ := attributeInt(event.Attributes, "total_tokens")
	totalToolUses, _ := attributeInt(event.Attributes, "total_tool_uses")
	finalModel, _ := attributeString(event.Attributes, "final_model")
	modelSwapped, _ := attributeBool(event.Attributes, "model_swapped")
	agent.MetadataJSON = jsonutil.ToCanonicalJSON(map[string]any{
		"eventName":       "subagent_completed",
		"total_tokens":    totalTokens,
		"total_tool_uses": totalToolUses,
		"final_model":     finalModel,
		"model_swapped":   modelSwapped,
	})
}
