package rawsignal

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Querier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type Record struct {
	NativeID    string
	Provider    string
	Signal      string
	DedupKey    string
	TimestampMs int64
	PayloadJSON string
}

const selectColumns = "native_id, provider, signal, dedup_key, timestamp_unix_ms, payload_json"

func scanRecord(rows *sql.Rows) (Record, error) {
	var record Record
	err := rows.Scan(
		&record.NativeID, &record.Provider, &record.Signal, &record.DedupKey,
		&record.TimestampMs, &record.PayloadJSON,
	)
	return record, err
}

func Insert(ctx context.Context, db Querier, records []Record) error {
	for _, record := range records {
		_, err := db.ExecContext(ctx, `
			INSERT INTO raw_signal (native_id, provider, signal, dedup_key, timestamp_unix_ms, payload_json)
			VALUES (?, ?, ?, ?, ?, ?)
			ON CONFLICT (native_id, dedup_key) DO NOTHING`,
			record.NativeID,
			record.Provider,
			record.Signal,
			record.DedupKey,
			record.TimestampMs,
			record.PayloadJSON,
		)
		if err != nil {
			return fmt.Errorf("insert raw_signal for native_id %s: %w", record.NativeID, err)
		}
	}
	return nil
}

func LoadByNativeIDs(ctx context.Context, db Querier, provider string, nativeIDs []string) (map[string][]Record, error) {
	byNativeID := make(map[string][]Record, len(nativeIDs))
	if len(nativeIDs) == 0 {
		return byNativeID, nil
	}

	placeholders := make([]string, len(nativeIDs))
	args := make([]any, 0, len(nativeIDs)+1)
	args = append(args, provider)
	for index, nativeID := range nativeIDs {
		placeholders[index] = "?"
		args = append(args, nativeID)
	}

	rows, err := db.QueryContext(ctx, `
		SELECT `+selectColumns+`
		FROM raw_signal
		WHERE provider = ? AND native_id IN (`+strings.Join(placeholders, ", ")+`)
		ORDER BY timestamp_unix_ms, id`, args...)
	if err != nil {
		return nil, fmt.Errorf("load raw_signal by native id: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		record, err := scanRecord(rows)
		if err != nil {
			return nil, fmt.Errorf("scan raw_signal row: %w", err)
		}
		byNativeID[record.NativeID] = append(byNativeID[record.NativeID], record)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate raw_signal rows: %w", err)
	}
	return byNativeID, nil
}

func LoadByProvider(ctx context.Context, db Querier, provider string) (map[string][]Record, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT `+selectColumns+`
		FROM raw_signal
		WHERE provider = ?
		ORDER BY timestamp_unix_ms, id`, provider)
	if err != nil {
		return nil, fmt.Errorf("load raw_signal by provider: %w", err)
	}
	defer rows.Close()

	byNativeID := make(map[string][]Record)
	for rows.Next() {
		record, err := scanRecord(rows)
		if err != nil {
			return nil, fmt.Errorf("scan raw_signal row: %w", err)
		}
		byNativeID[record.NativeID] = append(byNativeID[record.NativeID], record)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate raw_signal rows: %w", err)
	}
	return byNativeID, nil
}
