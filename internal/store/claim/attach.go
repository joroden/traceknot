package claim

import (
	"context"
	"database/sql"
	"fmt"

	"traceknot/internal/sqlutil"
)

func AttachByExternalID(
	ctx context.Context,
	db Querier,
	sessionID string,
	rawIDs []string,
	nowUnixMs int64,
) error {
	if len(rawIDs) == 0 {
		return nil
	}

	var claimID string
	err := db.QueryRowContext(ctx, `
		SELECT claim_id FROM claims
		WHERE session_id IN (`+sqlutil.Placeholders(len(rawIDs))+`)
		ORDER BY claimed_at_unix_ms DESC
		LIMIT 1`,
		sqlutil.ToAnySlice(rawIDs)...,
	).Scan(&claimID)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return fmt.Errorf("find claim by external id: %w", err)
	}

	_, err = db.ExecContext(ctx, `
		UPDATE claims SET session_id = ?, updated_at_unix_ms = ?
		WHERE claim_id = ?`,
		sessionID,
		nowUnixMs,
		claimID,
	)
	if err != nil {
		return fmt.Errorf("attach claim to session: %w", err)
	}
	return nil
}

type DetachedClaim struct {
	ClaimID         string
	SessionID       string
	ClaimedAtUnixMs int64
}

func ListDetached(ctx context.Context, db Querier) ([]DetachedClaim, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT claim_id, session_id, claimed_at_unix_ms
		FROM claims
		WHERE session_id IS NOT NULL
			AND session_id NOT IN (SELECT session_id FROM sessions)
		ORDER BY claimed_at_unix_ms`)
	if err != nil {
		return nil, fmt.Errorf("list detached claims: %w", err)
	}
	defer rows.Close()

	claims := []DetachedClaim{}
	for rows.Next() {
		var item DetachedClaim
		if err := rows.Scan(&item.ClaimID, &item.SessionID, &item.ClaimedAtUnixMs); err != nil {
			return nil, fmt.Errorf("scan detached claim: %w", err)
		}
		claims = append(claims, item)
	}
	return claims, rows.Err()
}
