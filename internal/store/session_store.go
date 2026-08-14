package store

import (
	"context"
	"database/sql"
	"time"

	"traceknot/internal/model"
	"traceknot/internal/store/claim"
	"traceknot/internal/store/session"
)

func (store *Store) ResolveSessionID(ctx context.Context, rawID string) (*string, error) {
	return session.ResolveID(ctx, store.db, rawID)
}

func (store *Store) ComputeSessionRollup(ctx context.Context, tx *sql.Tx, sessionID string) (session.Rollup, error) {
	return session.ComputeRollup(ctx, tx, sessionID)
}

func (store *Store) UpdateSessionRollup(ctx context.Context, tx *sql.Tx, sessionID string, rollup session.Rollup) error {
	return session.UpdateRollup(ctx, tx, sessionID, rollup)
}

func (store *Store) UpsertSession(ctx context.Context, tx *sql.Tx, seed *model.SessionSeed) error {
	if err := session.Upsert(ctx, tx, seed); err != nil {
		return err
	}
	rawIDs := []string{}
	for _, value := range []*string{seed.ExternalConversationID, seed.NativeSessionID} {
		if value != nil && *value != "" {
			rawIDs = append(rawIDs, *value)
		}
	}
	if len(rawIDs) > 0 {
		_ = claim.AttachByExternalID(ctx, tx, seed.SessionID, rawIDs, time.Now().UnixMilli())
	}
	return nil
}

func (store *Store) LoadSessionMeta(ctx context.Context, sessionID string) (*SessionMeta, error) {
	return session.LoadMeta(ctx, store.db, sessionID)
}

func (store *Store) ListSessions(ctx context.Context, filter ListFilter) (ListResult, error) {
	return session.List(ctx, store.db, filter)
}

func (store *Store) ListWorkItemGroups(ctx context.Context, filter GroupListFilter) (GroupListResult, error) {
	return session.ListGroups(ctx, store.db, filter)
}

func (store *Store) LoadSessionTree(ctx context.Context, sessionID string) ([]TreeRow, error) {
	return session.LoadTree(ctx, store.db, sessionID)
}
