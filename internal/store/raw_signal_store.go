package store

import (
	"context"

	"traceknot/internal/store/rawsignal"
)

func (store *Store) InsertRawSignal(ctx context.Context, records []RawSignalRecord) error {
	return rawsignal.Insert(ctx, store.db, records)
}

func (store *Store) LoadRawSignalByProvider(ctx context.Context, provider string) (map[string][]RawSignalRecord, error) {
	return rawsignal.LoadByProvider(ctx, store.db, provider)
}
