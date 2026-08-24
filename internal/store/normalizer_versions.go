package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func (store *Store) NormalizerVersion(ctx context.Context, provider string) (int, error) {
	var version int
	err := store.db.QueryRowContext(ctx,
		`SELECT version FROM normalizer_versions WHERE provider = ?`, provider,
	).Scan(&version)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("load normalizer version for %s: %w", provider, err)
	}
	return version, nil
}

func (store *Store) SetNormalizerVersion(ctx context.Context, provider string, version int) error {
	_, err := store.db.ExecContext(ctx, `
		INSERT INTO normalizer_versions (provider, version, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(provider) DO UPDATE SET
			version = excluded.version,
			updated_at = excluded.updated_at`,
		provider, version,
	)
	if err != nil {
		return fmt.Errorf("set normalizer version for %s: %w", provider, err)
	}
	return nil
}

func (store *Store) ListSessionIDsForProvider(ctx context.Context, provider string) ([]string, error) {
	rows, err := store.db.QueryContext(ctx,
		`SELECT session_id FROM sessions WHERE provider = ?`, provider)
	if err != nil {
		return nil, fmt.Errorf("list session ids for %s: %w", provider, err)
	}
	defer rows.Close()

	var sessionIDs []string
	for rows.Next() {
		var sessionID string
		if err := rows.Scan(&sessionID); err != nil {
			return nil, fmt.Errorf("scan session id: %w", err)
		}
		sessionIDs = append(sessionIDs, sessionID)
	}
	return sessionIDs, rows.Err()
}

func (store *Store) UpdateRawSignalNativeID(ctx context.Context, provider, oldNativeID, dedupKey, newNativeID string) error {
	if oldNativeID == newNativeID {
		return nil
	}
	_, err := store.db.ExecContext(ctx, `
		UPDATE raw_signal SET native_id = ?
		WHERE provider = ? AND native_id = ? AND dedup_key = ?`,
		newNativeID, provider, oldNativeID, dedupKey,
	)
	if err != nil {
		return fmt.Errorf("update raw_signal native_id: %w", err)
	}
	return nil
}

func (store *Store) CleanupOrphanedSessions(ctx context.Context, sessionIDs []string) error {
	if len(sessionIDs) == 0 {
		return nil
	}
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

	now := time.Now().UnixMilli()
	for _, sessionID := range sessionIDs {
		var replacement sql.NullString
		err := tx.QueryRowContext(ctx, `
			SELECT COALESCE(external_conversation_id, native_session_id)
			FROM sessions WHERE session_id = ?`, sessionID,
		).Scan(&replacement)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("look up raw id for orphaned session %s: %w", sessionID, err)
		}
		if !replacement.Valid || replacement.String == "" {
			continue
		}
		if _, err := tx.ExecContext(ctx, `
			UPDATE OR IGNORE claims SET session_id = ?, updated_at_unix_ms = ?
			WHERE session_id = ?`,
			replacement.String, now, sessionID,
		); err != nil {
			return fmt.Errorf("detach claim from orphaned session %s: %w", sessionID, err)
		}
	}

	if err := store.DeleteSessions(ctx, tx, sessionIDs...); err != nil {
		return fmt.Errorf("delete orphaned sessions: %w", err)
	}

	if err := store.CommitBatch(ctx, tx); err != nil {
		return err
	}
	committed = true
	return nil
}
