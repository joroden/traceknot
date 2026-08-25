package store

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

//go:embed migrations/001-core.sql
var coreMigrationSQL string

//go:embed migrations/002-picker.sql
var pickerMigrationSQL string

//go:embed migrations/003-normalizer-versions.sql
var normalizerVersionsMigrationSQL string

//go:embed migrations/004-scaling-indexes.sql
var scalingIndexesMigrationSQL string

type Store struct {
	db        *sql.DB
	batchLock sync.Mutex
}

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	db.SetMaxOpenConns(1)

	store := &Store{db: db}
	if _, err := store.db.Exec("PRAGMA journal_mode = WAL"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("enable WAL: %w", err)
	}
	if _, err := store.db.Exec("PRAGMA synchronous = NORMAL"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("set synchronous NORMAL: %w", err)
	}
	if _, err := store.db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}
	if _, err := store.db.ExecContext(context.Background(), coreMigrationSQL); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("apply core migration: %w", err)
	}
	if _, err := store.db.ExecContext(context.Background(), pickerMigrationSQL); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("apply picker migration: %w", err)
	}
	if _, err := store.db.ExecContext(context.Background(), normalizerVersionsMigrationSQL); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("apply normalizer versions migration: %w", err)
	}
	if _, err := store.db.ExecContext(context.Background(), scalingIndexesMigrationSQL); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("apply scaling indexes migration: %w", err)
	}
	if err := store.seedInstallTime(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func (store *Store) seedInstallTime() error {
	var hasRow bool
	err := store.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM app_meta WHERE id = 1)`).Scan(&hasRow)
	if err != nil {
		return fmt.Errorf("check install_ts: %w", err)
	}
	if hasRow {
		return nil
	}
	if _, err := store.db.Exec(`INSERT INTO app_meta (id, install_ts) VALUES (1, ?)`, time.Now().UnixMilli()); err != nil {
		return fmt.Errorf("seed install_ts: %w", err)
	}
	return nil
}

func (store *Store) Close() error {
	return store.db.Close()
}

func (store *Store) DB() *sql.DB {
	return store.db
}

func (store *Store) BeginBatch(ctx context.Context) (*sql.Tx, error) {
	store.batchLock.Lock()
	if _, err := store.db.ExecContext(ctx, "PRAGMA foreign_keys = OFF"); err != nil {
		store.batchLock.Unlock()
		return nil, fmt.Errorf("disable foreign keys for batch: %w", err)
	}
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		store.batchLock.Unlock()
		return nil, fmt.Errorf("begin batch: %w", err)
	}
	return tx, nil
}

func (store *Store) CommitBatch(ctx context.Context, tx *sql.Tx) error {
	defer store.batchLock.Unlock()
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit batch: %w", err)
	}
	if _, err := store.db.ExecContext(ctx, "PRAGMA foreign_keys = ON"); err != nil {
		return fmt.Errorf("re-enable foreign keys: %w", err)
	}
	return nil
}

func (store *Store) RollbackBatch(ctx context.Context, tx *sql.Tx) error {
	defer store.batchLock.Unlock()
	if err := tx.Rollback(); err != nil {
		return fmt.Errorf("rollback batch: %w", err)
	}
	if _, err := store.db.ExecContext(ctx, "PRAGMA foreign_keys = ON"); err != nil {
		return fmt.Errorf("re-enable foreign keys: %w", err)
	}
	return nil
}
