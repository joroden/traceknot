package session

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

type ListRow struct {
	SessionID            string
	Provider             string
	Title                string
	StartedAtUnixMs      *int64
	EndedAtUnixMs        *int64
	DurationMs           *float64
	InputTokens          int64
	CachedInputTokens    int64
	NonCachedInputTokens int64
	CacheWriteTokens     int64
	OutputTokens         int64
	ReasoningTokens      int64
	Cost                 float64
	Models               []string
	ClaimStatus          string
	WorkItemKey          string
	WorkItemTitle        string
}

type SortSpec struct {
	Column string
	Dir    string
}

type ListFilter struct {
	Provider         string
	Model            string
	Query            string
	StartUnixMs      int64
	EndUnixMs        int64
	WorkItemKey      string
	WorkItemProvider string

	Unclaimed bool
	Sorts     []SortSpec
	Offset    int
	Limit     int
}

type ListResult struct {
	Rows       []ListRow
	TotalCount int64
}

var listSortColumns = map[string]string{
	"cost":          "s.cost",
	"input_tokens":  "s.input_tokens",
	"output_tokens": "s.output_tokens",
	"started":       "s.started_at_unix_ms",
	"duration":      "s.duration_ms",
	"last_active":   "s.ended_at_unix_ms",
}

func buildListWhere(filter ListFilter) (string, []any) {
	where := []string{"s.node_count > 0"}
	args := []any{}

	if filter.Provider != "" {
		where = append(where, "s.provider = ?")
		args = append(args, filter.Provider)
	}
	if filter.Model != "" {
		where = append(where, "s.models_json LIKE ?")
		args = append(args, "%\""+filter.Model+"\"%")
	}
	if filter.Query != "" {
		where = append(where, "s.title LIKE ?")
		args = append(args, "%"+filter.Query+"%")
	}
	if filter.StartUnixMs > 0 {
		where = append(where, "s.started_at_unix_ms >= ?")
		args = append(args, filter.StartUnixMs)
	}
	if filter.EndUnixMs > 0 {
		where = append(where, "s.started_at_unix_ms < ?")
		args = append(args, filter.EndUnixMs)
	}
	if filter.Unclaimed {
		where = append(where, "(c.status IS NULL OR c.status != 'claimed')")
	} else if filter.WorkItemKey != "" {
		where = append(where, "c.status = 'claimed' AND c.work_item_key = ? AND c.provider = ?")
		args = append(args, filter.WorkItemKey, filter.WorkItemProvider)
	}
	return strings.Join(where, " AND "), args
}

func buildOrderBy(sorts []SortSpec) string {
	terms := make([]string, 0, len(sorts)+1)
	for _, sort := range sorts {
		column, ok := listSortColumns[sort.Column]
		if !ok {
			continue
		}
		dir := "DESC"
		if sort.Dir == "asc" {
			dir = "ASC"
		}
		terms = append(terms, column+" "+dir)
	}
	if len(terms) == 0 {
		terms = append(terms, "s.ended_at_unix_ms DESC")
	}
	terms = append(terms, "s.session_id")
	return strings.Join(terms, ", ")
}

func List(ctx context.Context, db Querier, filter ListFilter) (ListResult, error) {
	whereSQL, args := buildListWhere(filter)
	limit := clampListLimit(filter.Limit)

	var result ListResult
	if err := countListRows(ctx, db, whereSQL, args, &result.TotalCount); err != nil {
		return ListResult{}, err
	}

	rows, err := queryListRows(ctx, db, whereSQL, buildOrderBy(filter.Sorts), args, limit, filter.Offset)
	if err != nil {
		return ListResult{}, err
	}
	result.Rows = rows
	return result, nil
}

func clampListLimit(limit int) int {
	if limit <= 0 || limit > 500 {
		return 100
	}
	return limit
}

func countListRows(ctx context.Context, db Querier, whereSQL string, args []any, out *int64) error {
	countQuery := `
		SELECT COUNT(*)
		FROM sessions s
		LEFT JOIN claims c ON c.session_id = s.session_id
		WHERE ` + whereSQL
	if err := db.QueryRowContext(ctx, countQuery, args...).Scan(out); err != nil {
		return fmt.Errorf("count sessions: %w", err)
	}
	return nil
}

func queryListRows(
	ctx context.Context,
	db Querier,
	whereSQL string,
	orderBy string,
	args []any,
	limit int,
	offset int,
) ([]ListRow, error) {
	listQuery := `
		SELECT
			s.session_id, s.provider, s.title,
			s.started_at_unix_ms, s.ended_at_unix_ms, s.duration_ms,
			s.input_tokens, s.cached_input_tokens, s.non_cached_input_tokens,
			s.cache_write_tokens, s.output_tokens, s.reasoning_tokens,
			s.cost, s.models_json,
			c.status, c.work_item_key, c.work_item_title
		FROM sessions s
		LEFT JOIN claims c ON c.session_id = s.session_id
		WHERE ` + whereSQL + `
		ORDER BY ` + orderBy + `
		LIMIT ? OFFSET ?`
	listArgs := append(append([]any{}, args...), limit, offset)

	rows, err := db.QueryContext(ctx, listQuery, listArgs...)
	if err != nil {
		return nil, fmt.Errorf("list sessions: %w", err)
	}
	defer rows.Close()

	result := []ListRow{}
	for rows.Next() {
		var row ListRow
		var modelsJSON string
		var claimStatus, workItemKey, workItemTitle sql.NullString
		if err := rows.Scan(
			&row.SessionID, &row.Provider, &row.Title,
			&row.StartedAtUnixMs, &row.EndedAtUnixMs, &row.DurationMs,
			&row.InputTokens, &row.CachedInputTokens, &row.NonCachedInputTokens,
			&row.CacheWriteTokens, &row.OutputTokens, &row.ReasoningTokens,
			&row.Cost, &modelsJSON,
			&claimStatus, &workItemKey, &workItemTitle,
		); err != nil {
			return nil, fmt.Errorf("scan session row: %w", err)
		}
		if err := json.Unmarshal([]byte(modelsJSON), &row.Models); err != nil || row.Models == nil {
			row.Models = []string{}
		}
		if claimStatus.Valid {
			row.ClaimStatus = claimStatus.String
		} else {
			row.ClaimStatus = "unclaimed"
		}
		if row.ClaimStatus == "claimed" {
			row.WorkItemKey = workItemKey.String
			row.WorkItemTitle = workItemTitle.String
		}
		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate session rows: %w", err)
	}
	return result, nil
}
