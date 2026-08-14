package api

import (
	"net/http"
	"strconv"
	"time"

	"traceknot/internal/httputil"
	"traceknot/internal/store"
)

type Dashboard struct {
	store *store.Store
}

func NewDashboard(storeHandle *store.Store) *Dashboard {
	return &Dashboard{store: storeHandle}
}

func (dashboard *Dashboard) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/dashboard", dashboard.handleQuery)
	return mux
}

type dashboardResponse struct {
	FirstRun           bool                `json:"first_run"`
	Period             periodResponse      `json:"period"`
	ByWorkItem         []workItemCostJSON  `json:"by_work_item"`
	ByModel            []namedCostJSON     `json:"by_model"`
	ByAgent            []namedCostJSON     `json:"by_agent"`
	OverTime           []trendBucketJSON   `json:"over_time"`
	TrendGranularityMs int64               `json:"trend_granularity_ms"`
	RecentSessions     []recentSessionJSON `json:"recent_sessions"`
}

type periodResponse struct {
	StartUnixMs              int64    `json:"start_unix_ms"`
	EndUnixMs                int64    `json:"end_unix_ms"`
	Cost                     float64  `json:"cost"`
	CostDeltaPct             *float64 `json:"cost_delta_pct"`
	Tokens                   int64    `json:"tokens"`
	TokensDeltaPct           *float64 `json:"tokens_delta_pct"`
	InputTokens              int64    `json:"input_tokens"`
	OutputTokens             int64    `json:"output_tokens"`
	SessionCount             int64    `json:"session_count"`
	CoveragePct              *float64 `json:"coverage_pct"`
	CoverageDeltaPct         *float64 `json:"coverage_delta_pct"`
	UnattributedCost         float64  `json:"unattributed_cost"`
	UnattributedSessionCount int64    `json:"unattributed_session_count"`
}

type workItemCostJSON struct {
	Key          string  `json:"key"`
	Title        string  `json:"title"`
	Provider     string  `json:"provider"`
	Project      string  `json:"project"`
	Cost         float64 `json:"cost"`
	SessionCount int64   `json:"session_count"`
}

type namedCostJSON struct {
	Name string  `json:"name"`
	Cost float64 `json:"cost"`
}

type trendBucketJSON struct {
	StartUnixMs int64           `json:"start_unix_ms"`
	TotalCost   float64         `json:"total_cost"`
	Models      []namedCostJSON `json:"models"`
}

type recentSessionJSON struct {
	SessionID       string   `json:"session_id"`
	Provider        string   `json:"provider"`
	StartedAtUnixMs *int64   `json:"started_at_unix_ms"`
	Cost            float64  `json:"cost"`
	Tokens          int64    `json:"tokens"`
	Status          string   `json:"status"`
	NodeCount       int64    `json:"node_count"`
	Models          []string `json:"models"`
	Title           string   `json:"title"`
}

func (dashboard *Dashboard) handleQuery(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	startUnixMs, _ := strconv.ParseInt(query.Get("start_unix_ms"), 10, 64)
	endUnixMs, _ := strconv.ParseInt(query.Get("end_unix_ms"), 10, 64)

	period, err := store.ResolveDashboardPeriod(time.Now(), store.DashboardRequest{
		Range:       query.Get("range"),
		StartUnixMs: startUnixMs,
		EndUnixMs:   endUnixMs,
	})
	if err != nil {
		httputil.WriteError(writer, http.StatusBadRequest, "invalid_range", err.Error())
		return
	}

	snapshot, err := dashboard.store.DashboardSnapshot(request.Context(), period)
	if err != nil {
		httputil.WriteError(writer, http.StatusInternalServerError, "dashboard_failed", err.Error())
		return
	}

	httputil.WriteJSON(writer, http.StatusOK, toDashboardResponse(period, snapshot))
}

func toDashboardResponse(period store.DashboardRange, snapshot store.DashboardSnapshot) dashboardResponse {
	return dashboardResponse{
		FirstRun:           snapshot.FirstRun,
		Period:             toPeriodResponse(period, snapshot),
		ByWorkItem:         toWorkItemCostJSONs(snapshot.ByWorkItem),
		ByModel:            toNamedCostJSONs(snapshot.ByModel),
		ByAgent:            toNamedCostJSONs(snapshot.ByAgent),
		OverTime:           toTrendBucketJSONs(snapshot.OverTime),
		TrendGranularityMs: snapshot.TrendGranularityMs,
		RecentSessions:     toRecentSessionJSONs(snapshot.RecentSessions),
	}
}

func toPeriodResponse(period store.DashboardRange, snapshot store.DashboardSnapshot) periodResponse {
	return periodResponse{
		StartUnixMs:              period.StartMs,
		EndUnixMs:                period.EndMs,
		Cost:                     snapshot.Period.Cost,
		CostDeltaPct:             deltaPct(snapshot.Period.Cost, snapshot.PreviousPeriod.Cost),
		Tokens:                   snapshot.Period.Tokens,
		TokensDeltaPct:           deltaPct(float64(snapshot.Period.Tokens), float64(snapshot.PreviousPeriod.Tokens)),
		InputTokens:              snapshot.Period.InputTokens,
		OutputTokens:             snapshot.Period.OutputTokens,
		SessionCount:             snapshot.Period.SessionCount,
		CoveragePct:              coveragePct(snapshot.Period.ClaimedCount, snapshot.Period.SessionCount),
		CoverageDeltaPct:         coverageDeltaPct(snapshot.Period, snapshot.PreviousPeriod),
		UnattributedCost:         snapshot.Period.UnattributedCost,
		UnattributedSessionCount: snapshot.Period.UnattributedSessions,
	}
}

func toWorkItemCostJSONs(items []store.WorkItemCost) []workItemCostJSON {
	result := make([]workItemCostJSON, 0, len(items))
	for _, item := range items {
		result = append(result, workItemCostJSON{
			Key:          item.Key,
			Title:        item.Title,
			Provider:     item.Provider,
			Project:      item.Project,
			Cost:         item.Cost,
			SessionCount: item.SessionCount,
		})
	}
	return result
}

func toNamedCostJSONs(items []store.NamedCost) []namedCostJSON {
	result := make([]namedCostJSON, 0, len(items))
	for _, item := range items {
		result = append(result, namedCostJSON{Name: item.Name, Cost: item.Cost})
	}
	return result
}

func toTrendBucketJSONs(buckets []store.TrendBucket) []trendBucketJSON {
	result := make([]trendBucketJSON, 0, len(buckets))
	for _, bucket := range buckets {
		result = append(result, trendBucketJSON{
			StartUnixMs: bucket.StartUnixMs,
			TotalCost:   bucket.TotalCost,
			Models:      toNamedCostJSONs(bucket.Models),
		})
	}
	return result
}

func toRecentSessionJSONs(sessions []store.RecentSession) []recentSessionJSON {
	result := make([]recentSessionJSON, 0, len(sessions))
	for _, session := range sessions {
		result = append(result, recentSessionJSON{
			SessionID:       session.SessionID,
			Provider:        session.Provider,
			StartedAtUnixMs: session.StartedAtUnixMs,
			Cost:            session.Cost,
			Tokens:          session.Tokens,
			Status:          session.Status,
			NodeCount:       session.NodeCount,
			Models:          session.Models,
			Title:           session.Title,
		})
	}
	return result
}

func deltaPct(current, previous float64) *float64 {
	if previous == 0 {
		return nil
	}
	value := (current - previous) / previous * 100
	return &value
}

func coveragePct(claimed, total int64) *float64 {
	if total == 0 {
		return nil
	}
	value := float64(claimed) / float64(total) * 100
	return &value
}

func coverageDeltaPct(current, previous store.PeriodSummary) *float64 {
	currentPct := coveragePct(current.ClaimedCount, current.SessionCount)
	previousPct := coveragePct(previous.ClaimedCount, previous.SessionCount)
	if currentPct == nil || previousPct == nil {
		return nil
	}
	value := *currentPct - *previousPct
	return &value
}
