package claim

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Querier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type Claim struct {
	ClaimID         string
	SessionID       *string
	WorkItemKey     string
	WorkItemTitle   string
	Provider        string
	Project         string
	Source          string
	ClaimedAtUnixMs int64
	UpdatedAtUnixMs int64
}

func Upsert(ctx context.Context, db Querier, claim Claim) error {
	_, err := db.ExecContext(ctx, `
		INSERT INTO claims (
			claim_id, session_id, work_item_key, work_item_title,
			provider, project, source, status, claimed_at_unix_ms, updated_at_unix_ms
		) VALUES (?, ?, ?, ?, ?, ?, ?, 'claimed', ?, ?)
		ON CONFLICT(session_id) WHERE session_id IS NOT NULL DO UPDATE SET
			claim_id = excluded.claim_id,
			work_item_key = excluded.work_item_key,
			work_item_title = excluded.work_item_title,
			provider = excluded.provider,
			project = excluded.project,
			source = excluded.source,
			status = 'claimed',
			claimed_at_unix_ms = excluded.claimed_at_unix_ms,
			updated_at_unix_ms = excluded.updated_at_unix_ms`,
		claim.ClaimID,
		claim.SessionID,
		claim.WorkItemKey,
		claim.WorkItemTitle,
		claim.Provider,
		claim.Project,
		claim.Source,
		claim.ClaimedAtUnixMs,
		claim.UpdatedAtUnixMs,
	)
	if err != nil {
		return fmt.Errorf("upsert claim: %w", err)
	}
	return nil
}

var ErrClaimNotFound = errors.New("claim not found")

func AttachToSession(
	ctx context.Context,
	db Querier,
	claimID string,
	sessionID string,
	nowUnixMs int64,
) error {
	result, err := db.ExecContext(ctx, `
		UPDATE claims SET session_id = ?, updated_at_unix_ms = ?
		WHERE claim_id = ?`,
		sessionID,
		nowUnixMs,
		claimID,
	)
	if err != nil {
		return fmt.Errorf("attach claim to session: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("attach claim to session: %w", err)
	}
	if affected == 0 {
		return ErrClaimNotFound
	}
	return nil
}
