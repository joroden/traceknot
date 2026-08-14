package claim

import (
	"context"
	"fmt"
)

type RecentWorkItem struct {
	Key                    string `json:"key"`
	Provider               string `json:"provider"`
	Title                  string `json:"title"`
	Project                string `json:"project"`
	LastAttributedAtUnixMs int64  `json:"last_attributed_at_unix_ms"`
	AttributionCount       int64  `json:"attribution_count"`
}

func UpsertRecent(ctx context.Context, db Querier, item RecentWorkItem) error {
	_, err := db.ExecContext(ctx, `
		INSERT INTO work_item_recent (
			key, provider, title, project, last_attributed_at_unix_ms, attribution_count
		) VALUES (?, ?, ?, ?, ?, 1)
		ON CONFLICT(provider, key) DO UPDATE SET
			title = excluded.title,
			project = excluded.project,
			last_attributed_at_unix_ms = excluded.last_attributed_at_unix_ms,
			attribution_count = attribution_count + 1`,
		item.Key,
		item.Provider,
		item.Title,
		item.Project,
		item.LastAttributedAtUnixMs,
	)
	if err != nil {
		return fmt.Errorf("upsert recent work item: %w", err)
	}
	return nil
}

func ListRecent(
	ctx context.Context,
	db Querier,
	provider string,
	limit int,
) ([]RecentWorkItem, error) {
	if limit <= 0 {
		limit = 50
	}

	query := `
		SELECT key, provider, title, project, last_attributed_at_unix_ms, attribution_count
		FROM work_item_recent`
	var args []any
	if provider != "" {
		query += ` WHERE provider = ?`
		args = append(args, provider)
	}
	query += ` ORDER BY last_attributed_at_unix_ms DESC, attribution_count DESC LIMIT ?`
	args = append(args, limit)

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list recent work items: %w", err)
	}
	defer rows.Close()

	items := make([]RecentWorkItem, 0)
	for rows.Next() {
		var item RecentWorkItem
		if err := rows.Scan(
			&item.Key,
			&item.Provider,
			&item.Title,
			&item.Project,
			&item.LastAttributedAtUnixMs,
			&item.AttributionCount,
		); err != nil {
			return nil, fmt.Errorf("scan recent work item: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate recent work items: %w", err)
	}
	return items, nil
}
