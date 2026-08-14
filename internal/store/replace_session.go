package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"traceknot/internal/model"
	"traceknot/internal/store/node"
)

func (store *Store) DeleteSessions(ctx context.Context, tx *sql.Tx, sessionIDs ...string) error {
	if len(sessionIDs) == 0 {
		return nil
	}
	placeholders := make([]string, len(sessionIDs))
	args := make([]any, len(sessionIDs))
	for index, sessionID := range sessionIDs {
		placeholders[index] = "?"
		args[index] = sessionID
	}
	inClause := strings.Join(placeholders, ", ")
	if err := node.DeleteForSessions(ctx, tx, inClause, args...); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx,
		"DELETE FROM sessions WHERE session_id IN ("+inClause+")", args...); err != nil {
		return fmt.Errorf("delete sessions: %w", err)
	}
	return nil
}

func (store *Store) ReplaceSession(ctx context.Context, seed *model.SessionSeed, content *model.SessionContent) error {
	tx, err := store.BeginBatch(ctx)
	if err != nil {
		return err
	}
	committed := false
	defer func() {
		if !committed {
			_ = store.RollbackBatch(ctx, tx)
		}
	}()

	if err := store.DeleteSessions(ctx, tx, seed.SessionID); err != nil {
		return fmt.Errorf("delete session %s: %w", seed.SessionID, err)
	}
	if err := store.UpsertSession(ctx, tx, seed); err != nil {
		return fmt.Errorf("upsert session %s: %w", seed.SessionID, err)
	}
	if err := store.UpsertNodes(ctx, tx, seed.SessionID, content); err != nil {
		return fmt.Errorf("upsert nodes for session %s: %w", seed.SessionID, err)
	}
	if err := store.ReassignOwners(ctx, tx, seed.SessionID); err != nil {
		return fmt.Errorf("reassign owners for session %s: %w", seed.SessionID, err)
	}
	rollup, err := store.ComputeSessionRollup(ctx, tx, seed.SessionID)
	if err != nil {
		return fmt.Errorf("compute rollup for session %s: %w", seed.SessionID, err)
	}
	if err := store.UpdateSessionRollup(ctx, tx, seed.SessionID, rollup); err != nil {
		return fmt.Errorf("update rollup for session %s: %w", seed.SessionID, err)
	}
	if err := store.CommitBatch(ctx, tx); err != nil {
		return err
	}
	committed = true
	return nil
}
