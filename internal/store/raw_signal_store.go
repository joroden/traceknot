package store

import (
	"context"
	"fmt"

	"traceknot/internal/store/rawsignal"
)

func (store *Store) InsertRawSignal(ctx context.Context, records []RawSignalRecord) error {
	if len(records) == 0 {
		return nil
	}
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin raw_signal insert: %w", err)
	}
	if err := rawsignal.Insert(ctx, tx, records); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit raw_signal insert: %w", err)
	}
	return nil
}

func (store *Store) LoadRawSignalByProvider(ctx context.Context, provider string) (map[string][]RawSignalRecord, error) {
	return rawsignal.LoadByProvider(ctx, store.db, provider)
}

func (store *Store) LoadRawSignalByNativeIDs(ctx context.Context, provider string, nativeIDs []string) (map[string][]RawSignalRecord, error) {
	return rawsignal.LoadByNativeIDs(ctx, store.db, provider, nativeIDs)
}
