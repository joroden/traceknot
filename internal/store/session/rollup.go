package session

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
)

type Rollup struct {
	StartedAtUnixMs      *int64
	EndedAtUnixMs        *int64
	DurationMs           *float64
	Status               string
	TurnCount            int64
	NodeCount            int64
	ToolCallCount        int64
	AgentRunCount        int64
	InputTokens          int64
	CachedInputTokens    int64
	NonCachedInputTokens int64
	OutputTokens         int64
	ReasoningTokens      int64
	CacheWriteTokens     int64
	Cost                 float64
	Models               []string
}

func ComputeRollup(ctx context.Context, db Querier, sessionID string) (Rollup, error) {
	rollup, err := queryRollupTotals(ctx, db, sessionID)
	if err != nil {
		return Rollup{}, err
	}
	models, err := queryRollupModels(ctx, db, sessionID)
	if err != nil {
		return Rollup{}, err
	}
	rollup.Models = models
	return rollup, nil
}

func queryRollupTotals(ctx context.Context, db Querier, sessionID string) (Rollup, error) {
	var rollup Rollup
	err := db.QueryRowContext(ctx, `
		SELECT
			MIN(started_at_unix_ms),
			MAX(ended_at_unix_ms),
			COALESCE(SUM(CASE WHEN kind = 'chat' AND name = 'user' THEN 1 ELSE 0 END), 0),
			COUNT(*),
			COALESCE(SUM(CASE WHEN kind = 'tool_call' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN kind = 'agent' AND parent_node_id IS NOT NULL THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(input_tokens), 0),
			COALESCE(SUM(cached_input_tokens), 0),
			COALESCE(SUM(cache_write_tokens), 0),
			COALESCE(SUM(output_tokens), 0),
			COALESCE(SUM(reasoning_tokens), 0),
			COALESCE(SUM(cost), 0)
		FROM nodes
		WHERE session_id = ?`,
		sessionID,
	).Scan(
		&rollup.StartedAtUnixMs,
		&rollup.EndedAtUnixMs,
		&rollup.TurnCount,
		&rollup.NodeCount,
		&rollup.ToolCallCount,
		&rollup.AgentRunCount,
		&rollup.InputTokens,
		&rollup.CachedInputTokens,
		&rollup.CacheWriteTokens,
		&rollup.OutputTokens,
		&rollup.ReasoningTokens,
		&rollup.Cost,
	)
	if err != nil {
		return Rollup{}, fmt.Errorf("compute session rollup: %w", err)
	}

	if rollup.NodeCount > 0 {
		rollup.Status = "ok"
	} else {
		rollup.Status = "empty"
	}
	if rollup.StartedAtUnixMs != nil && rollup.EndedAtUnixMs != nil {
		duration := float64(*rollup.EndedAtUnixMs - *rollup.StartedAtUnixMs)
		rollup.DurationMs = &duration
	}
	rollup.NonCachedInputTokens = max(rollup.InputTokens-rollup.CachedInputTokens-rollup.CacheWriteTokens, 0)
	return rollup, nil
}

func queryRollupModels(ctx context.Context, db Querier, sessionID string) ([]string, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT DISTINCT model FROM nodes
		WHERE session_id = ? AND model IS NOT NULL AND model != ''
		ORDER BY model`,
		sessionID,
	)
	if err != nil {
		return nil, fmt.Errorf("list session models: %w", err)
	}
	defer rows.Close()

	models := []string{}
	for rows.Next() {
		var model string
		if err := rows.Scan(&model); err != nil {
			return nil, fmt.Errorf("scan session model: %w", err)
		}
		models = append(models, model)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate session models: %w", err)
	}
	sort.Strings(models)
	return models, nil
}

func UpdateRollup(ctx context.Context, db Querier, sessionID string, rollup Rollup) error {
	modelsJSON, err := json.Marshal(rollup.Models)
	if err != nil {
		return fmt.Errorf("marshal session models: %w", err)
	}

	_, err = db.ExecContext(ctx, `
		UPDATE sessions SET
			started_at_unix_ms = ?,
			ended_at_unix_ms = ?,
			duration_ms = ?,
			status = ?,
			turn_count = ?,
			node_count = ?,
			tool_call_count = ?,
			agent_run_count = ?,
			input_tokens = ?,
			cached_input_tokens = ?,
			non_cached_input_tokens = ?,
			output_tokens = ?,
			reasoning_tokens = ?,
			cache_write_tokens = ?,
			cost = ?,
			models_json = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE session_id = ?`,
		rollup.StartedAtUnixMs,
		rollup.EndedAtUnixMs,
		rollup.DurationMs,
		rollup.Status,
		rollup.TurnCount,
		rollup.NodeCount,
		rollup.ToolCallCount,
		rollup.AgentRunCount,
		rollup.InputTokens,
		rollup.CachedInputTokens,
		rollup.NonCachedInputTokens,
		rollup.OutputTokens,
		rollup.ReasoningTokens,
		rollup.CacheWriteTokens,
		rollup.Cost,
		string(modelsJSON),
		sessionID,
	)
	if err != nil {
		return fmt.Errorf("update session rollup: %w", err)
	}
	return nil
}
