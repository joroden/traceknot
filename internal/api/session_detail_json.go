package api

import (
	"encoding/json"

	"traceknot/internal/store"
)

type sessionMetaJSON struct {
	SessionID         string   `json:"session_id"`
	Provider          string   `json:"provider"`
	ServiceName       *string  `json:"service_name"`
	Status            *string  `json:"status"`
	StartedAtUnixMs   *int64   `json:"started_at_unix_ms"`
	EndedAtUnixMs     *int64   `json:"ended_at_unix_ms"`
	DurationMs        *float64 `json:"duration_ms"`
	TurnCount         int64    `json:"turn_count"`
	NodeCount         int64    `json:"node_count"`
	ToolCallCount     int64    `json:"tool_call_count"`
	AgentRunCount     int64    `json:"agent_run_count"`
	InputTokens       int64    `json:"input_tokens"`
	CachedInputTokens int64    `json:"cached_input_tokens"`
	OutputTokens      int64    `json:"output_tokens"`
	ReasoningTokens   int64    `json:"reasoning_tokens"`
	CacheWriteTokens  int64    `json:"cache_write_tokens"`
	Cost              float64  `json:"cost"`
	Models            []string `json:"models"`
}

type treeRowJSON struct {
	NodeID                string   `json:"node_id"`
	Kind                  string   `json:"kind"`
	Name                  *string  `json:"name"`
	AgentName             *string  `json:"agent_name"`
	ToolName              *string  `json:"tool_name"`
	ToolCallID            *string  `json:"tool_call_id"`
	Model                 *string  `json:"model"`
	Status                *string  `json:"status"`
	StartedAtUnixMs       *int64   `json:"started_at_unix_ms"`
	DurationMs            *float64 `json:"duration_ms"`
	InputTokens           int64    `json:"input_tokens"`
	CachedInputTokens     int64    `json:"cached_input_tokens"`
	CacheWriteTokens      int64    `json:"cache_write_tokens"`
	OutputTokens          int64    `json:"output_tokens"`
	ReasoningTokens       int64    `json:"reasoning_tokens"`
	Cost                  float64  `json:"cost"`
	EstimatedInputTokens  int64    `json:"estimated_input_tokens"`
	EstimatedOutputTokens int64    `json:"estimated_output_tokens"`
	OwningAgentID         *string  `json:"owning_agent_id"`
	OwningAgentName       *string  `json:"owning_agent_name"`
	ParentNodeID          *string  `json:"parent_node_id"`
	HasContent            bool     `json:"has_content"`
	AggInputTokens        int64    `json:"agg_input_tokens"`
	AggCachedInputTokens  int64    `json:"agg_cached_input_tokens"`
	AggCacheWriteTokens   int64    `json:"agg_cache_write_tokens"`
	AggOutputTokens       int64    `json:"agg_output_tokens"`
	AggReasoningTokens    int64    `json:"agg_reasoning_tokens"`
	AggCost               float64  `json:"agg_cost"`
	DescendantCount       int64    `json:"descendant_count"`
	SubagentCount         int64    `json:"subagent_count"`
}

func metaJSON(meta *store.SessionMeta) sessionMetaJSON {
	models := []string{}
	_ = json.Unmarshal([]byte(meta.ModelsJSON), &models)
	return sessionMetaJSON{
		SessionID:         meta.SessionID,
		Provider:          meta.Provider,
		ServiceName:       meta.ServiceName,
		Status:            meta.Status,
		StartedAtUnixMs:   meta.StartedAtUnixMs,
		EndedAtUnixMs:     meta.EndedAtUnixMs,
		DurationMs:        meta.DurationMs,
		TurnCount:         meta.TurnCount,
		NodeCount:         meta.NodeCount,
		ToolCallCount:     meta.ToolCallCount,
		AgentRunCount:     meta.AgentRunCount,
		InputTokens:       meta.InputTokens,
		CachedInputTokens: meta.CachedInputTokens,
		OutputTokens:      meta.OutputTokens,
		ReasoningTokens:   meta.ReasoningTokens,
		CacheWriteTokens:  meta.CacheWriteTokens,
		Cost:              meta.Cost,
		Models:            models,
	}
}

func toTreeRowJSON(row store.TreeRow) treeRowJSON {
	return treeRowJSON{
		NodeID:                row.NodeID,
		Kind:                  row.Kind,
		Name:                  row.Name,
		AgentName:             row.AgentName,
		ToolName:              row.ToolName,
		ToolCallID:            row.ToolCallID,
		Model:                 row.Model,
		Status:                row.Status,
		StartedAtUnixMs:       row.StartedAtUnixMs,
		DurationMs:            row.DurationMs,
		InputTokens:           row.InputTokens,
		CachedInputTokens:     row.CachedInputTokens,
		CacheWriteTokens:      row.CacheWriteTokens,
		OutputTokens:          row.OutputTokens,
		ReasoningTokens:       row.ReasoningTokens,
		Cost:                  row.Cost,
		EstimatedInputTokens:  row.EstimatedInputTokens,
		EstimatedOutputTokens: row.EstimatedOutputTokens,
		OwningAgentID:         row.OwningAgentID,
		OwningAgentName:       row.OwningAgentName,
		ParentNodeID:          row.ParentNodeID,
		HasContent:            row.HasContent,
		AggInputTokens:        row.AggInputTokens,
		AggCachedInputTokens:  row.AggCachedInputTokens,
		AggCacheWriteTokens:   row.AggCacheWriteTokens,
		AggOutputTokens:       row.AggOutputTokens,
		AggReasoningTokens:    row.AggReasoningTokens,
		AggCost:               row.AggCost,
		DescendantCount:       row.DescendantCount,
		SubagentCount:         row.SubagentCount,
	}
}

func nodeDetailJSON(detail *store.NodeDetail) map[string]any {
	promptText, outputText, reasoningText, argumentsJSON, resultText, errorDetailsJSON := "", "", "", "", "", ""
	agentName, toolName, toolCallID, spawnPrompt := "", "", "", ""
	if detail.Chat != nil {
		promptText = derefStr(detail.Chat.PromptText)
		outputText = derefStr(detail.Chat.OutputText)
		reasoningText = derefStr(detail.Chat.ReasoningText)
	}
	if detail.ToolCall != nil {
		toolName = derefStr(detail.ToolCall.ToolName)
		toolCallID = derefStr(detail.ToolCall.ToolCallID)
		argumentsJSON = derefStr(detail.ToolCall.ArgumentsJSON)
		resultText = derefStr(detail.ToolCall.ResultText)
		errorDetailsJSON = derefStr(detail.ToolCall.ErrorDetailsJSON)
	}
	if detail.Agent != nil {
		agentName = derefStr(detail.Agent.AgentName)
		spawnPrompt = derefStr(detail.Agent.SpawnPrompt)
	}
	return map[string]any{
		"node_id":                 detail.NodeID,
		"session_id":              detail.SessionID,
		"kind":                    detail.Kind,
		"name":                    detail.Name,
		"agent_name":              nullableString(agentName),
		"tool_name":               nullableString(toolName),
		"tool_call_id":            nullableString(toolCallID),
		"spawn_prompt":            nullableString(spawnPrompt),
		"model":                   detail.Model,
		"status":                  detail.Status,
		"started_at_unix_ms":      detail.StartedAtUnixMs,
		"ended_at_unix_ms":        detail.EndedAtUnixMs,
		"duration_ms":             detail.DurationMs,
		"input_tokens":            detail.InputTokens,
		"cached_input_tokens":     detail.CachedInputTokens,
		"cache_write_tokens":      detail.CacheWriteTokens,
		"output_tokens":           detail.OutputTokens,
		"reasoning_tokens":        detail.ReasoningTokens,
		"cost":                    detail.Cost,
		"estimated_input_tokens":  detail.EstimatedInputTokens,
		"estimated_output_tokens": detail.EstimatedOutputTokens,
		"token_estimate_method":   detail.TokenEstimateMethod,
		"prompt_text":             nullableString(promptText),
		"output_text":             nullableString(outputText),
		"reasoning_text":          nullableString(reasoningText),
		"arguments_json":          nullableString(argumentsJSON),
		"result_text":             nullableString(resultText),
		"error_details_json":      nullableString(errorDetailsJSON),
		"owning_agent_id":         detail.OwningAgentID,
		"owning_agent_name":       detail.OwningAgentName,
		"metadata_json":           detail.MetadataJSON,
	}
}

func derefStr(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func nullableString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func syntheticRootDetail(rows []store.TreeRow, sessionID string) *store.NodeDetail {
	prefix := "synthetic:root:" + sessionID
	for _, row := range rows {
		if row.NodeID != prefix {
			continue
		}
		return &store.NodeDetail{
			NodeID:            row.NodeID,
			SessionID:         sessionID,
			Kind:              "agent",
			Name:              row.Name,
			Status:            row.Status,
			StartedAtUnixMs:   row.StartedAtUnixMs,
			DurationMs:        row.DurationMs,
			InputTokens:       row.AggInputTokens,
			CachedInputTokens: row.AggCachedInputTokens,
			CacheWriteTokens:  row.AggCacheWriteTokens,
			OutputTokens:      row.AggOutputTokens,
			ReasoningTokens:   row.AggReasoningTokens,
			Cost:              row.AggCost,
		}
	}
	return nil
}
