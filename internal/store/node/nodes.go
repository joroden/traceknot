package node

import (
	"context"
	"database/sql"
	"fmt"

	"traceknot/internal/model"
)

type Querier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

func UpsertBase(ctx context.Context, db Querier, sessionID string, seed *model.NodeSeed, hasContent bool) error {
	_, err := db.ExecContext(ctx, `
		INSERT INTO nodes (
			node_id, session_id, parent_node_id, kind, name, status, model,
			started_at_unix_ms, ended_at_unix_ms, duration_ms, preview_text,
			input_tokens, cached_input_tokens, cache_write_tokens,
			output_tokens, reasoning_tokens, cost,
			estimated_input_tokens, estimated_output_tokens, token_estimate_method,
			has_content, metadata_json
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(node_id) DO UPDATE SET
			parent_node_id = excluded.parent_node_id,
			kind = excluded.kind,
			name = excluded.name,
			status = excluded.status,
			model = excluded.model,
			started_at_unix_ms = excluded.started_at_unix_ms,
			ended_at_unix_ms = excluded.ended_at_unix_ms,
			duration_ms = excluded.duration_ms,
			preview_text = excluded.preview_text,
			input_tokens = excluded.input_tokens,
			cached_input_tokens = excluded.cached_input_tokens,
			cache_write_tokens = excluded.cache_write_tokens,
			output_tokens = excluded.output_tokens,
			reasoning_tokens = excluded.reasoning_tokens,
			cost = excluded.cost,
			estimated_input_tokens = excluded.estimated_input_tokens,
			estimated_output_tokens = excluded.estimated_output_tokens,
			token_estimate_method = excluded.token_estimate_method,
			has_content = excluded.has_content,
			metadata_json = excluded.metadata_json,
			updated_at = CURRENT_TIMESTAMP`,
		seed.NodeID,
		sessionID,
		seed.ParentNodeID,
		string(seed.Kind),
		seed.Name,
		seed.Status,
		seed.Model,
		seed.StartedAtUnixMs,
		seed.EndedAtUnixMs,
		seed.DurationMs,
		seed.PreviewText,
		seed.InputTokens,
		seed.CachedInputTokens,
		seed.CacheWriteTokens,
		seed.OutputTokens,
		seed.ReasoningTokens,
		seed.Cost,
		seed.EstimatedInputTokens,
		seed.EstimatedOutputTokens,
		seed.TokenEstimateMethod,
		boolInt(hasContent),
		seed.MetadataJSON,
	)
	if err != nil {
		return fmt.Errorf("upsert node %s: %w", seed.NodeID, err)
	}
	return nil
}

func UpsertChat(ctx context.Context, db Querier, seed *model.ChatSeed) error {
	if _, err := db.ExecContext(ctx, `
		INSERT INTO chat_nodes (node_id, system_text, prompt_text, output_text, reasoning_text)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(node_id) DO UPDATE SET
			system_text = excluded.system_text,
			prompt_text = excluded.prompt_text,
			output_text = excluded.output_text,
			reasoning_text = excluded.reasoning_text`,
		seed.NodeID,
		nullIfEmpty(seed.SystemText),
		nullIfEmpty(seed.PromptText),
		nullIfEmpty(seed.OutputText),
		nullIfEmpty(seed.ReasoningText),
	); err != nil {
		return fmt.Errorf("upsert chat node %s: %w", seed.NodeID, err)
	}
	return nil
}

func UpsertToolCall(ctx context.Context, db Querier, seed *model.ToolCallSeed) error {
	if _, err := db.ExecContext(ctx, `
		INSERT INTO tool_call_nodes (node_id, tool_name, tool_call_id, arguments_json, result_text, error_details_json, approval_decision, approval_source)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(node_id) DO UPDATE SET
			tool_name = excluded.tool_name,
			tool_call_id = excluded.tool_call_id,
			arguments_json = excluded.arguments_json,
			result_text = excluded.result_text,
			error_details_json = excluded.error_details_json,
			approval_decision = excluded.approval_decision,
			approval_source = excluded.approval_source`,
		seed.NodeID,
		nullIfEmpty(seed.ToolName),
		nullIfEmpty(seed.ToolCallID),
		nullIfEmpty(seed.ArgumentsJSON),
		nullIfEmpty(seed.ResultText),
		nullIfEmpty(seed.ErrorDetailsJSON),
		seed.ApprovalDecision,
		seed.ApprovalSource,
	); err != nil {
		return fmt.Errorf("upsert tool call node %s: %w", seed.NodeID, err)
	}
	return nil
}

func UpsertAgent(ctx context.Context, db Querier, seed *model.AgentSeed) error {
	if _, err := db.ExecContext(ctx, `
		INSERT INTO agent_nodes (node_id, agent_name, agent_type, spawn_prompt, spawn_tool_call_id, result_summary)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(node_id) DO UPDATE SET
			agent_name = excluded.agent_name,
			agent_type = excluded.agent_type,
			spawn_prompt = excluded.spawn_prompt,
			spawn_tool_call_id = excluded.spawn_tool_call_id,
			result_summary = excluded.result_summary`,
		seed.NodeID,
		nullIfEmpty(seed.AgentName),
		nullIfEmpty(seed.AgentType),
		nullIfEmpty(seed.SpawnPrompt),
		nullIfEmpty(seed.SpawnToolCallID),
		nullIfEmpty(seed.ResultSummary),
	); err != nil {
		return fmt.Errorf("upsert agent node %s: %w", seed.NodeID, err)
	}
	return nil
}

func DeleteForSessions(ctx context.Context, db Querier, placeholders string, args ...any) error {
	queries := []string{
		"DELETE FROM chat_nodes WHERE node_id IN (SELECT node_id FROM nodes WHERE session_id IN (" + placeholders + "))",
		"DELETE FROM tool_call_nodes WHERE node_id IN (SELECT node_id FROM nodes WHERE session_id IN (" + placeholders + "))",
		"DELETE FROM agent_nodes WHERE node_id IN (SELECT node_id FROM nodes WHERE session_id IN (" + placeholders + "))",
		"DELETE FROM nodes WHERE session_id IN (" + placeholders + ")",
	}
	for _, query := range queries {
		if _, err := db.ExecContext(ctx, query, args...); err != nil {
			return fmt.Errorf("delete session nodes: %w", err)
		}
	}
	return nil
}

func nullIfEmpty(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
