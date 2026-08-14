package dashboard

import (
	"context"
	"fmt"
)

type WorkItemCost struct {
	Provider     string
	Key          string
	Title        string
	Project      string
	Cost         float64
	SessionCount int64
}

type NamedCost struct {
	Name string
	Cost float64
}

func loadByWorkItem(ctx context.Context, db Querier, startMs, endMs int64, limit int) ([]WorkItemCost, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT c.provider, c.work_item_key, c.work_item_title, c.project,
			SUM(s.cost), COUNT(s.session_id)
		FROM claims c
		JOIN sessions s ON s.session_id = c.session_id
		WHERE s.started_at_unix_ms >= ? AND s.started_at_unix_ms < ?
			AND c.status = 'claimed'
		GROUP BY c.provider, c.work_item_key
		ORDER BY SUM(s.cost) DESC
		LIMIT ?`,
		startMs,
		endMs,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("cost by work item: %w", err)
	}
	defer rows.Close()

	items := []WorkItemCost{}
	for rows.Next() {
		var item WorkItemCost
		if err := rows.Scan(&item.Provider, &item.Key, &item.Title, &item.Project, &item.Cost, &item.SessionCount); err != nil {
			return nil, fmt.Errorf("scan work item cost: %w", err)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func loadByNodeColumn(ctx context.Context, db Querier, startMs, endMs int64, column string, limit int) ([]NamedCost, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT COALESCE(n.`+column+`, 'unknown'), SUM(n.cost)
		FROM nodes n
		JOIN sessions s ON s.session_id = n.session_id
		WHERE s.started_at_unix_ms >= ? AND s.started_at_unix_ms < ?
			AND n.`+column+` IS NOT NULL AND n.`+column+` != ''
			AND n.`+column+` NOT IN ('auto', 'AUTO', 'Auto')
		GROUP BY n.`+column+`
		ORDER BY SUM(n.cost) DESC
		LIMIT ?`,
		startMs,
		endMs,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("cost by %s: %w", column, err)
	}
	defer rows.Close()

	items := []NamedCost{}
	for rows.Next() {
		var item NamedCost
		if err := rows.Scan(&item.Name, &item.Cost); err != nil {
			return nil, fmt.Errorf("scan %s cost: %w", column, err)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func loadByAgent(ctx context.Context, db Querier, startMs, endMs int64, limit int) ([]NamedCost, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT s.provider, SUM(s.cost)
		FROM sessions s
		WHERE s.started_at_unix_ms >= ? AND s.started_at_unix_ms < ?
		GROUP BY s.provider
		ORDER BY SUM(s.cost) DESC
		LIMIT ?`,
		startMs,
		endMs,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("cost by agent: %w", err)
	}
	defer rows.Close()

	items := []NamedCost{}
	for rows.Next() {
		var item NamedCost
		if err := rows.Scan(&item.Name, &item.Cost); err != nil {
			return nil, fmt.Errorf("scan agent cost: %w", err)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
