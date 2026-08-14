package api

import (
	"net/http"
	"strings"

	"traceknot/internal/httputil"
	"traceknot/internal/store"
)

type WorkItems struct {
	store *store.Store
}

func NewWorkItems(storeHandle *store.Store) *WorkItems {
	return &WorkItems{store: storeHandle}
}

func (workItems *WorkItems) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/work-items", workItems.handleList)
	return mux
}

type workItemGroupJSON struct {
	WorkItemKey      string   `json:"work_item_key"`
	WorkItemProvider string   `json:"work_item_provider"`
	Title            string   `json:"title"`
	IsUnclaimed      bool     `json:"is_unclaimed"`
	SessionCount     int64    `json:"session_count"`
	Cost             float64  `json:"cost"`
	DurationMs       *float64 `json:"duration_ms"`
	InputTokens      int64    `json:"input_tokens"`
	OutputTokens     int64    `json:"output_tokens"`
}

func parseGroupSort(raw string) store.GroupSortSpec {
	if raw == "" {
		return store.GroupSortSpec{}
	}
	column, dir, found := strings.Cut(raw, ":")
	if !found {
		dir = "desc"
	}
	return store.GroupSortSpec{Column: column, Dir: dir}
}

func (workItems *WorkItems) handleList(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	filter := store.GroupListFilter{
		Provider:         query.Get("provider"),
		Model:            query.Get("model"),
		Query:            query.Get("q"),
		StartUnixMs:      queryInt64(request, "start_unix_ms"),
		EndUnixMs:        queryInt64(request, "end_unix_ms"),
		WorkItemKey:      query.Get("work_item_key"),
		WorkItemProvider: query.Get("work_item_provider"),
		Sort:             parseGroupSort(query.Get("sort")),
		Offset:           queryInt(request, "offset"),
		Limit:            queryInt(request, "limit"),
	}

	result, err := workItems.store.ListWorkItemGroups(request.Context(), filter)
	if err != nil {
		httputil.WriteError(writer, http.StatusInternalServerError, "work_items_failed", err.Error())
		return
	}

	groups := make([]workItemGroupJSON, 0, len(result.Groups))
	for _, row := range result.Groups {
		groups = append(groups, workItemGroupJSON{
			WorkItemKey:      row.WorkItemKey,
			WorkItemProvider: row.WorkItemProvider,
			Title:            row.Title,
			IsUnclaimed:      row.IsUnclaimed,
			SessionCount:     row.SessionCount,
			Cost:             row.Cost,
			DurationMs:       row.DurationMs,
			InputTokens:      row.InputTokens,
			OutputTokens:     row.OutputTokens,
		})
	}
	httputil.WriteJSON(writer, http.StatusOK, map[string]any{
		"groups":      groups,
		"total_count": result.TotalCount,
	})
}
