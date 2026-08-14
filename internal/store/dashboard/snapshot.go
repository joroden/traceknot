package dashboard

import (
	"context"
	"database/sql"
	"fmt"
)

type Querier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type Snapshot struct {
	FirstRun           bool
	Period             PeriodSummary
	PreviousPeriod     PeriodSummary
	ByWorkItem         []WorkItemCost
	ByModel            []NamedCost
	ByAgent            []NamedCost
	OverTime           []TrendBucket
	TrendGranularityMs int64
	RecentSessions     []RecentSession
}

func minSessionStart(ctx context.Context, db Querier) (int64, error) {
	var minStart sql.NullInt64
	if err := db.QueryRowContext(ctx, `
		SELECT MIN(started_at_unix_ms)
		FROM sessions
		WHERE started_at_unix_ms IS NOT NULL`).Scan(&minStart); err != nil {
		return 0, fmt.Errorf("min session start: %w", err)
	}
	return minStart.Int64, nil
}

func Load(ctx context.Context, db Querier, period Range) (Snapshot, error) {
	snapshot := Snapshot{}

	totalSessions, err := countSessions(ctx, db)
	if err != nil {
		return snapshot, err
	}
	snapshot.FirstRun = totalSessions == 0

	if period.Key == "all" {
		if minStart, err := minSessionStart(ctx, db); err == nil && minStart > 0 {
			period.StartMs = minStart
		}
	}

	if snapshot.Period, err = loadPeriodSummary(ctx, db, period.StartMs, period.EndMs); err != nil {
		return snapshot, err
	}
	if snapshot.PreviousPeriod, err = loadPeriodSummary(ctx, db, period.PrevStartMs, period.PrevEndMs); err != nil {
		return snapshot, err
	}
	if snapshot.ByWorkItem, err = loadByWorkItem(ctx, db, period.StartMs, period.EndMs, 20); err != nil {
		return snapshot, err
	}
	if snapshot.ByModel, err = loadByNodeColumn(ctx, db, period.StartMs, period.EndMs, "model", 20); err != nil {
		return snapshot, err
	}
	if snapshot.ByAgent, err = loadByAgent(ctx, db, period.StartMs, period.EndMs, 20); err != nil {
		return snapshot, err
	}
	if snapshot.OverTime, snapshot.TrendGranularityMs, err = loadTrend(ctx, db, period); err != nil {
		return snapshot, err
	}
	if snapshot.RecentSessions, err = loadRecent(ctx, db, 5); err != nil {
		return snapshot, err
	}
	return snapshot, nil
}

func countSessions(ctx context.Context, db Querier) (int64, error) {
	var count int64
	err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sessions`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count sessions: %w", err)
	}
	return count, nil
}
