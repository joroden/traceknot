package dashboard

import (
	"context"
	"fmt"
)

type PeriodSummary struct {
	Cost                 float64
	Tokens               int64
	InputTokens          int64
	OutputTokens         int64
	SessionCount         int64
	ClaimedCount         int64
	UnattributedCost     float64
	UnattributedSessions int64
}

func loadPeriodSummary(ctx context.Context, db Querier, startMs, endMs int64) (PeriodSummary, error) {
	var summary PeriodSummary
	err := db.QueryRowContext(ctx, `
		SELECT
			COALESCE(SUM(cost), 0),
			COALESCE(SUM(input_tokens + output_tokens), 0),
			COALESCE(SUM(input_tokens), 0),
			COALESCE(SUM(output_tokens), 0),
			COUNT(*)
		FROM sessions
		WHERE started_at_unix_ms >= ? AND started_at_unix_ms < ? AND node_count > 0`,
		startMs,
		endMs,
	).Scan(&summary.Cost, &summary.Tokens, &summary.InputTokens, &summary.OutputTokens, &summary.SessionCount)
	if err != nil {
		return summary, fmt.Errorf("period summary: %w", err)
	}

	err = db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM sessions s
		JOIN claims c ON c.session_id = s.session_id
		WHERE s.started_at_unix_ms >= ? AND s.started_at_unix_ms < ?
			AND s.node_count > 0 AND c.status = 'claimed'`,
		startMs,
		endMs,
	).Scan(&summary.ClaimedCount)
	if err != nil {
		return summary, fmt.Errorf("period claimed count: %w", err)
	}

	err = db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(s.cost), 0), COUNT(*)
		FROM sessions s
		LEFT JOIN claims c ON c.session_id = s.session_id AND c.status = 'claimed'
		WHERE s.started_at_unix_ms >= ? AND s.started_at_unix_ms < ?
			AND s.node_count > 0 AND c.session_id IS NULL`,
		startMs,
		endMs,
	).Scan(&summary.UnattributedCost, &summary.UnattributedSessions)
	if err != nil {
		return summary, fmt.Errorf("period unattributed: %w", err)
	}
	return summary, nil
}
