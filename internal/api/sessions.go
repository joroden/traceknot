package api

import (
	"net/http"
	"strconv"
	"strings"

	"traceknot/internal/httputil"
	"traceknot/internal/store"
)

type Sessions struct {
	store *store.Store
}

func NewSessions(storeHandle *store.Store) *Sessions {
	return &Sessions{store: storeHandle}
}

func (sessions *Sessions) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/sessions", sessions.handleList)
	return mux
}

type tokenBreakdownJSON struct {
	Total     int64 `json:"total"`
	Raw       int64 `json:"raw,omitempty"`
	Cached    int64 `json:"cached,omitempty"`
	Write     int64 `json:"write,omitempty"`
	Reasoning int64 `json:"reasoning,omitempty"`
}

type claimStateJSON struct {
	Status        string `json:"status"`
	WorkItemKey   string `json:"work_item_key,omitempty"`
	WorkItemTitle string `json:"work_item_title,omitempty"`
}

type sessionRowJSON struct {
	SessionID       string             `json:"session_id"`
	Provider        string             `json:"provider"`
	Title           string             `json:"title"`
	StartedAtUnixMs *int64             `json:"started_at_unix_ms"`
	EndedAtUnixMs   *int64             `json:"ended_at_unix_ms"`
	DurationMs      *float64           `json:"duration_ms"`
	Models          []string           `json:"models"`
	Cost            float64            `json:"cost"`
	InputTokens     tokenBreakdownJSON `json:"input_tokens"`
	OutputTokens    tokenBreakdownJSON `json:"output_tokens"`
	Claim           claimStateJSON     `json:"claim"`
}

func toSessionRowJSON(row store.ListRow) sessionRowJSON {
	return sessionRowJSON{
		SessionID:       row.SessionID,
		Provider:        row.Provider,
		Title:           row.Title,
		StartedAtUnixMs: row.StartedAtUnixMs,
		EndedAtUnixMs:   row.EndedAtUnixMs,
		DurationMs:      row.DurationMs,
		Models:          row.Models,
		Cost:            row.Cost,
		InputTokens: tokenBreakdownJSON{
			Total:  row.InputTokens,
			Raw:    row.NonCachedInputTokens,
			Cached: row.CachedInputTokens,
			Write:  row.CacheWriteTokens,
		},
		OutputTokens: tokenBreakdownJSON{
			Total:     row.OutputTokens,
			Reasoning: row.ReasoningTokens,
		},
		Claim: claimStateJSON{
			Status:        row.ClaimStatus,
			WorkItemKey:   row.WorkItemKey,
			WorkItemTitle: row.WorkItemTitle,
		},
	}
}

func queryInt64(request *http.Request, key string) int64 {
	value, err := strconv.ParseInt(request.URL.Query().Get(key), 10, 64)
	if err != nil {
		return 0
	}
	return value
}

func queryInt(request *http.Request, key string) int {
	value, err := strconv.Atoi(request.URL.Query().Get(key))
	if err != nil {
		return 0
	}
	return value
}

func parseSorts(raw string) []store.SortSpec {
	if raw == "" {
		return nil
	}
	segments := strings.Split(raw, ",")
	sorts := make([]store.SortSpec, 0, len(segments))
	for _, segment := range segments {
		column, dir, found := strings.Cut(segment, ":")
		if column == "" {
			continue
		}
		if !found {
			dir = "desc"
		}
		sorts = append(sorts, store.SortSpec{Column: column, Dir: dir})
	}
	return sorts
}

func (sessions *Sessions) handleList(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	filter := store.ListFilter{
		Provider:         query.Get("provider"),
		Model:            query.Get("model"),
		Query:            query.Get("q"),
		StartUnixMs:      queryInt64(request, "start_unix_ms"),
		EndUnixMs:        queryInt64(request, "end_unix_ms"),
		WorkItemKey:      query.Get("work_item_key"),
		WorkItemProvider: query.Get("work_item_provider"),
		Unclaimed:        query.Get("unclaimed") == "true",
		Sorts:            parseSorts(query.Get("sort")),
		Offset:           queryInt(request, "offset"),
		Limit:            queryInt(request, "limit"),
	}

	result, err := sessions.store.ListSessions(request.Context(), filter)
	if err != nil {
		httputil.WriteError(writer, http.StatusInternalServerError, "sessions_failed", err.Error())
		return
	}

	rows := make([]sessionRowJSON, 0, len(result.Rows))
	for _, row := range result.Rows {
		rows = append(rows, toSessionRowJSON(row))
	}
	httputil.WriteJSON(writer, http.StatusOK, map[string]any{
		"sessions":    rows,
		"total_count": result.TotalCount,
	})
}
