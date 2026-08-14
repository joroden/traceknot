package store

import (
	"time"

	"traceknot/internal/store/claim"
	"traceknot/internal/store/dashboard"
	"traceknot/internal/store/node"
	"traceknot/internal/store/rawsignal"
	"traceknot/internal/store/session"
)

type (
	ListFilter      = session.ListFilter
	ListResult      = session.ListResult
	ListRow         = session.ListRow
	SortSpec        = session.SortSpec
	GroupListFilter = session.GroupListFilter
	GroupListResult = session.GroupListResult
	GroupSortSpec   = session.GroupSortSpec
	SessionMeta     = session.Meta
	TreeRow         = session.TreeRow

	NodeDetail = node.Detail

	Claim          = claim.Claim
	RecentWorkItem = claim.RecentWorkItem
	Outcome        = claim.Outcome

	DashboardRange    = dashboard.Range
	DashboardRequest  = dashboard.Request
	DashboardSnapshot = dashboard.Snapshot
	PeriodSummary     = dashboard.PeriodSummary
	WorkItemCost      = dashboard.WorkItemCost
	NamedCost         = dashboard.NamedCost
	TrendBucket       = dashboard.TrendBucket
	RecentSession     = dashboard.RecentSession

	RawSignalRecord = rawsignal.Record
)

const (
	ClaimStatusClaimed = claim.StatusClaimed
	ClaimStatusSkipped = claim.StatusSkipped
)

func ResolveDashboardPeriod(now time.Time, request DashboardRequest) (DashboardRange, error) {
	return dashboard.Resolve(now, request)
}
