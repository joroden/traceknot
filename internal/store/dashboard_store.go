package store

import (
	"context"

	"traceknot/internal/store/dashboard"
)

func (store *Store) DashboardSnapshot(ctx context.Context, period DashboardRange) (DashboardSnapshot, error) {
	return dashboard.Load(ctx, store.db, period)
}
