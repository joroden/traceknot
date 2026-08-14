package session

import (
	"context"
	"database/sql"
	"fmt"
)

type Meta struct {
	SessionID         string
	Provider          string
	Title             string
	ServiceName       *string
	Status            *string
	StartedAtUnixMs   *int64
	EndedAtUnixMs     *int64
	DurationMs        *float64
	TurnCount         int64
	NodeCount         int64
	ToolCallCount     int64
	AgentRunCount     int64
	InputTokens       int64
	CachedInputTokens int64
	OutputTokens      int64
	ReasoningTokens   int64
	CacheWriteTokens  int64
	Cost              float64
	ModelsJSON        string
	MetadataJSON      string
}

func LoadMeta(ctx context.Context, db Querier, sessionID string) (*Meta, error) {
	var meta Meta
	err := db.QueryRowContext(ctx, `
		SELECT session_id, provider, title, service_name, status,
			started_at_unix_ms, ended_at_unix_ms, duration_ms,
			turn_count, node_count, tool_call_count, agent_run_count,
			input_tokens, cached_input_tokens, output_tokens, reasoning_tokens,
			cache_write_tokens, cost, models_json, metadata_json
		FROM sessions
		WHERE session_id = ?`,
		sessionID,
	).Scan(
		&meta.SessionID, &meta.Provider, &meta.Title, &meta.ServiceName, &meta.Status,
		&meta.StartedAtUnixMs, &meta.EndedAtUnixMs, &meta.DurationMs,
		&meta.TurnCount, &meta.NodeCount, &meta.ToolCallCount, &meta.AgentRunCount,
		&meta.InputTokens, &meta.CachedInputTokens, &meta.OutputTokens, &meta.ReasoningTokens,
		&meta.CacheWriteTokens, &meta.Cost, &meta.ModelsJSON, &meta.MetadataJSON,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("load session meta: %w", err)
	}
	return &meta, nil
}
