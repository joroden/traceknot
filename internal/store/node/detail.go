package node

import (
	"context"
	"database/sql"
	"fmt"
)

type Detail struct {
	NodeID                string
	SessionID             string
	Kind                  string
	Name                  *string
	Model                 *string
	Status                *string
	StartedAtUnixMs       *int64
	EndedAtUnixMs         *int64
	DurationMs            *float64
	InputTokens           int64
	CachedInputTokens     int64
	CacheWriteTokens      int64
	OutputTokens          int64
	ReasoningTokens       int64
	Cost                  float64
	EstimatedInputTokens  int64
	EstimatedOutputTokens int64
	TokenEstimateMethod   *string
	OwningAgentID         *string
	OwningAgentName       *string
	MetadataJSON          string

	Chat     *ChatDetail
	ToolCall *ToolCallDetail
	Agent    *AgentDetail
}

type ChatDetail struct {
	SystemText    *string
	PromptText    *string
	OutputText    *string
	ReasoningText *string
}

type ToolCallDetail struct {
	ToolName         *string
	ToolCallID       *string
	ArgumentsJSON    *string
	ResultText       *string
	ErrorDetailsJSON *string
	ApprovalDecision *string
	ApprovalSource   *string
}

type AgentDetail struct {
	AgentName       *string
	AgentType       *string
	SpawnPrompt     *string
	SpawnToolCallID *string
	ResultSummary   *string
}

func LoadDetail(ctx context.Context, db Querier, nodeID string) (*Detail, error) {
	var detail Detail
	err := db.QueryRowContext(ctx, `
		SELECT node_id, session_id, kind, name,
			model, status, started_at_unix_ms, ended_at_unix_ms, duration_ms,
			input_tokens, cached_input_tokens, cache_write_tokens,
			output_tokens, reasoning_tokens, cost,
			estimated_input_tokens, estimated_output_tokens, token_estimate_method,
			owning_agent_id, owning_agent_name, metadata_json
		FROM nodes
		WHERE node_id = ?`,
		nodeID,
	).Scan(
		&detail.NodeID, &detail.SessionID, &detail.Kind, &detail.Name,
		&detail.Model, &detail.Status,
		&detail.StartedAtUnixMs, &detail.EndedAtUnixMs, &detail.DurationMs,
		&detail.InputTokens, &detail.CachedInputTokens, &detail.CacheWriteTokens,
		&detail.OutputTokens, &detail.ReasoningTokens, &detail.Cost,
		&detail.EstimatedInputTokens, &detail.EstimatedOutputTokens, &detail.TokenEstimateMethod,
		&detail.OwningAgentID, &detail.OwningAgentName, &detail.MetadataJSON,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("load node detail: %w", err)
	}
	switch detail.Kind {
	case "chat":
		detail.Chat, err = loadChatDetail(ctx, db, nodeID)
	case "tool_call":
		detail.ToolCall, err = loadToolCallDetail(ctx, db, nodeID)
	case "agent":
		detail.Agent, err = loadAgentDetail(ctx, db, nodeID)
	}
	if err != nil {
		return nil, err
	}
	return &detail, nil
}

func loadChatDetail(ctx context.Context, db Querier, nodeID string) (*ChatDetail, error) {
	var detail ChatDetail
	if err := db.QueryRowContext(ctx, `
		SELECT system_text, prompt_text, output_text, reasoning_text
		FROM chat_nodes WHERE node_id = ?`,
		nodeID,
	).Scan(&detail.SystemText, &detail.PromptText, &detail.OutputText, &detail.ReasoningText); err != nil {
		if err == sql.ErrNoRows {
			return &ChatDetail{}, nil
		}
		return nil, fmt.Errorf("load chat node detail: %w", err)
	}
	return &detail, nil
}

func loadToolCallDetail(ctx context.Context, db Querier, nodeID string) (*ToolCallDetail, error) {
	var detail ToolCallDetail
	if err := db.QueryRowContext(ctx, `
		SELECT tool_name, tool_call_id, arguments_json, result_text, error_details_json, approval_decision, approval_source
		FROM tool_call_nodes WHERE node_id = ?`,
		nodeID,
	).Scan(
		&detail.ToolName, &detail.ToolCallID, &detail.ArgumentsJSON,
		&detail.ResultText, &detail.ErrorDetailsJSON,
		&detail.ApprovalDecision, &detail.ApprovalSource,
	); err != nil {
		if err == sql.ErrNoRows {
			return &ToolCallDetail{}, nil
		}
		return nil, fmt.Errorf("load tool call node detail: %w", err)
	}
	return &detail, nil
}

func loadAgentDetail(ctx context.Context, db Querier, nodeID string) (*AgentDetail, error) {
	var detail AgentDetail
	if err := db.QueryRowContext(ctx, `
		SELECT agent_name, agent_type, spawn_prompt, spawn_tool_call_id, result_summary
		FROM agent_nodes WHERE node_id = ?`,
		nodeID,
	).Scan(
		&detail.AgentName, &detail.AgentType, &detail.SpawnPrompt, &detail.SpawnToolCallID, &detail.ResultSummary,
	); err != nil {
		if err == sql.ErrNoRows {
			return &AgentDetail{}, nil
		}
		return nil, fmt.Errorf("load agent node detail: %w", err)
	}
	return &detail, nil
}
