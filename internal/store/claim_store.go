package store

import (
	"context"
	"time"

	"traceknot/internal/store/claim"
	"traceknot/internal/store/session"
)

func (store *Store) UpsertClaim(ctx context.Context, c Claim) error {
	return claim.Upsert(ctx, store.db, c)
}

func (store *Store) AttachClaimToSession(
	ctx context.Context,
	claimID string,
	sessionID string,
	nowUnixMs int64,
) error {
	return claim.AttachToSession(ctx, store.db, claimID, sessionID, nowUnixMs)
}

func (store *Store) OutcomeForSession(ctx context.Context, sessionID string) (*Outcome, error) {
	return claim.OutcomeForSession(ctx, store.db, sessionID)
}

func (store *Store) OfferPicker(ctx context.Context, sessionID string, nowUnixMs int64) (string, bool, error) {
	return claim.Offer(ctx, store.db, sessionID, nowUnixMs)
}

func (store *Store) RecordSkip(ctx context.Context, sessionID string, nowUnixMs int64) error {
	return claim.RecordSkip(ctx, store.db, sessionID, nowUnixMs)
}

func (store *Store) UpsertRecentWorkItem(ctx context.Context, item RecentWorkItem) error {
	return claim.UpsertRecent(ctx, store.db, item)
}

func (store *Store) ListRecentWorkItems(
	ctx context.Context,
	provider string,
	limit int,
) ([]RecentWorkItem, error) {
	return claim.ListRecent(ctx, store.db, provider, limit)
}

func (store *Store) ReconcileClaimAttachments(ctx context.Context) error {
	detached, err := claim.ListDetached(ctx, store.db)
	if err != nil {
		return err
	}
	for _, item := range detached {
		resolved, err := session.ResolveID(ctx, store.db, item.SessionID)
		if err != nil {
			return err
		}
		if resolved == nil {
			continue
		}
		if err := claim.AttachToSession(ctx, store.db, item.ClaimID, *resolved, time.Now().UnixMilli()); err != nil {
			return err
		}
	}
	return nil
}
