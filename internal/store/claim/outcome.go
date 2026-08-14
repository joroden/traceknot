package claim

import (
	"context"
	"database/sql"
	"fmt"

	"traceknot/internal/stableid"
)

const (
	StatusClaimed = "claimed"
	StatusSkipped = "skipped"
	StatusPending = "pending"
)

type Outcome struct {
	Status      string
	WorkItemKey string
}

func OutcomeForSession(ctx context.Context, db Querier, sessionID string) (*Outcome, error) {
	var outcome Outcome
	err := db.QueryRowContext(ctx, `
		SELECT status, work_item_key
		FROM claims
		WHERE session_id = ?
		LIMIT 1`,
		sessionID,
	).Scan(&outcome.Status, &outcome.WorkItemKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("outcome for session: %w", err)
	}
	return &outcome, nil
}

func Offer(ctx context.Context, db Querier, sessionID string, nowUnixMs int64) (string, bool, error) {
	result, err := db.ExecContext(ctx, `
		INSERT INTO claims (
			claim_id, session_id, work_item_key, work_item_title,
			provider, project, source, status, claimed_at_unix_ms, updated_at_unix_ms
		) VALUES (?, ?, '', '', '', '', 'offer', 'pending', ?, ?)
		ON CONFLICT(session_id) WHERE session_id IS NOT NULL DO NOTHING`,
		stableid.From("offer", sessionID),
		sessionID,
		nowUnixMs,
		nowUnixMs,
	)
	if err != nil {
		return "", false, fmt.Errorf("record offer: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return "", false, fmt.Errorf("offer rows affected: %w", err)
	}
	if affected == 1 {
		return StatusPending, true, nil
	}
	outcome, err := OutcomeForSession(ctx, db, sessionID)
	if err != nil {
		return "", false, err
	}
	if outcome == nil {
		return "", false, nil
	}
	return outcome.Status, false, nil
}

func RecordSkip(ctx context.Context, db Querier, sessionID string, nowUnixMs int64) error {
	if _, err := db.ExecContext(ctx, `
		UPDATE claims SET status = 'skipped', updated_at_unix_ms = ?
		WHERE session_id = ?`,
		nowUnixMs,
		sessionID,
	); err != nil {
		return fmt.Errorf("record skip: %w", err)
	}
	return nil
}
