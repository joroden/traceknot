package dashboard

import (
	"context"
	"fmt"
)

const (
	hourDurationMs  = int64(60 * 60 * 1000)
	dayDurationMs   = int64(24 * 60 * 60 * 1000)
	weekDurationMs  = int64(7 * 24 * 60 * 60 * 1000)
	monthDurationMs = int64(30 * 24 * 60 * 60 * 1000)
)

type TrendBucket struct {
	StartUnixMs int64
	TotalCost   float64
	Models      []NamedCost
}

func bucketMsFor(durationMs int64) int64 {
	switch {
	case durationMs <= 2*dayDurationMs:
		return hourDurationMs
	case durationMs <= 90*dayDurationMs:
		return dayDurationMs
	case durationMs <= 730*dayDurationMs:
		return weekDurationMs
	default:
		return monthDurationMs
	}
}

func loadTrend(ctx context.Context, db Querier, period Range) ([]TrendBucket, int64, error) {
	durationMs := period.EndMs - period.StartMs
	bucketMs := bucketMsFor(durationMs)
	bucketCount := int((durationMs + bucketMs - 1) / bucketMs)

	buckets := make([]TrendBucket, bucketCount)
	for index := range buckets {
		buckets[index].StartUnixMs = period.StartMs + int64(index)*bucketMs
		buckets[index].Models = []NamedCost{}
	}

	rows, err := db.QueryContext(ctx, `
		SELECT CAST((n.started_at_unix_ms - ?) / ? AS INTEGER), COALESCE(n.model, 'unknown'), SUM(n.cost)
		FROM nodes n
		JOIN sessions s ON s.session_id = n.session_id
		WHERE n.started_at_unix_ms >= ? AND n.started_at_unix_ms < ?
			AND n.model IS NOT NULL AND n.model != ''
			AND n.model NOT IN ('auto', 'AUTO', 'Auto')
		GROUP BY 1, n.model`,
		period.StartMs,
		bucketMs,
		period.StartMs,
		period.EndMs,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("trend: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var bucketIndex int
		var model string
		var cost float64
		if err := rows.Scan(&bucketIndex, &model, &cost); err != nil {
			return nil, 0, fmt.Errorf("scan trend bucket: %w", err)
		}
		if bucketIndex < 0 || bucketIndex >= bucketCount {
			continue
		}
		buckets[bucketIndex].Models = append(buckets[bucketIndex].Models, NamedCost{Name: model, Cost: cost})
		buckets[bucketIndex].TotalCost += cost
	}
	return buckets, bucketMs, rows.Err()
}
