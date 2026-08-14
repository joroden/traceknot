package dashboard

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"traceknot/internal/ptr"
)

type RecentSession struct {
	SessionID       string
	Provider        string
	StartedAtUnixMs *int64
	Cost            float64
	Tokens          int64
	Status          string
	NodeCount       int64
	Models          []string
	Title           string
}

func loadRecent(ctx context.Context, db Querier, limit int) ([]RecentSession, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT session_id, provider, started_at_unix_ms, cost,
			input_tokens + output_tokens, status, node_count, models_json, title
		FROM sessions
		WHERE node_count > 0
		ORDER BY started_at_unix_ms DESC
		LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("recent sessions: %w", err)
	}
	defer rows.Close()

	sessions := []RecentSession{}
	for rows.Next() {
		var session RecentSession
		var startedAt sql.NullInt64
		var modelsJSON string
		if err := rows.Scan(
			&session.SessionID,
			&session.Provider,
			&startedAt,
			&session.Cost,
			&session.Tokens,
			&session.Status,
			&session.NodeCount,
			&modelsJSON,
			&session.Title,
		); err != nil {
			return nil, fmt.Errorf("scan recent session: %w", err)
		}
		if startedAt.Valid {
			session.StartedAtUnixMs = ptr.Int64(startedAt.Int64)
		}
		if err := json.Unmarshal([]byte(modelsJSON), &session.Models); err != nil || session.Models == nil {
			session.Models = []string{}
		}
		sessions = append(sessions, session)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return sessions, nil
}
