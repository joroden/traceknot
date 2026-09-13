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

const PromptDuplicateWindowMs = 10000

func IsDuplicateOffer(fingerprint string, promptFingerprint string, outcome *Outcome, nowUnixMs int64) bool {
	if fingerprint != "" && outcome.OfferFingerprint != "" {
		return fingerprint == outcome.OfferFingerprint
	}
	return promptFingerprint != "" &&
		promptFingerprint == outcome.OfferPromptFingerprint &&
		nowUnixMs-outcome.ClaimedAtUnixMs <= PromptDuplicateWindowMs
}

type Outcome struct {
	Status                 string
	WorkItemKey            string
	ClaimedAtUnixMs        int64
	OfferFingerprint       string
	OfferPromptFingerprint string
}

func OutcomeForSession(ctx context.Context, db Querier, sessionID string) (*Outcome, error) {
	var outcome Outcome
	err := db.QueryRowContext(ctx, `
		SELECT status, work_item_key, claimed_at_unix_ms, offer_fingerprint, offer_prompt_fingerprint
		FROM claims
		WHERE session_id = ?
		LIMIT 1`,
		sessionID,
	).Scan(&outcome.Status, &outcome.WorkItemKey, &outcome.ClaimedAtUnixMs, &outcome.OfferFingerprint, &outcome.OfferPromptFingerprint)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("outcome for session: %w", err)
	}
	return &outcome, nil
}

func Offer(ctx context.Context, db Querier, sessionID string, fingerprint string, promptFingerprint string, nowUnixMs int64) (*Outcome, bool, error) {
	result, err := db.ExecContext(ctx, `
		INSERT INTO claims (
			claim_id, session_id, work_item_key, work_item_title,
			provider, project, source, status, claimed_at_unix_ms, updated_at_unix_ms,
			offer_fingerprint, offer_prompt_fingerprint
		) VALUES (?, ?, '', '', '', '', 'offer', 'pending', ?, ?, ?, ?)
		ON CONFLICT(session_id) WHERE session_id IS NOT NULL DO NOTHING`,
		stableid.From("offer", sessionID),
		sessionID,
		nowUnixMs,
		nowUnixMs,
		fingerprint,
		promptFingerprint,
	)
	if err != nil {
		return nil, false, fmt.Errorf("record offer: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return nil, false, fmt.Errorf("offer rows affected: %w", err)
	}
	if affected == 1 {
		return &Outcome{
			Status:                 StatusPending,
			ClaimedAtUnixMs:        nowUnixMs,
			OfferFingerprint:       fingerprint,
			OfferPromptFingerprint: promptFingerprint,
		}, true, nil
	}
	outcome, err := OutcomeForSession(ctx, db, sessionID)
	if err != nil {
		return nil, false, err
	}
	return outcome, false, nil
}

func ResetPending(ctx context.Context, db Querier, sessionID string, fingerprint string, promptFingerprint string, nowUnixMs int64) error {
	if _, err := db.ExecContext(ctx, `
		UPDATE claims SET status = 'pending', claimed_at_unix_ms = ?, updated_at_unix_ms = ?,
			offer_fingerprint = ?, offer_prompt_fingerprint = ?
		WHERE session_id = ?`,
		nowUnixMs,
		nowUnixMs,
		fingerprint,
		promptFingerprint,
		sessionID,
	); err != nil {
		return fmt.Errorf("reset pending: %w", err)
	}
	return nil
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
