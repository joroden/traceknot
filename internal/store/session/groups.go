package session

import (
	"context"
	"fmt"
	"strings"
)

type GroupRow struct {
	WorkItemKey      string
	WorkItemProvider string
	Title            string
	IsUnclaimed      bool
	SessionCount     int64
	Cost             float64

	DurationMs   *float64
	InputTokens  int64
	OutputTokens int64
}

type GroupSortSpec struct {
	Column string
	Dir    string
}

type GroupListFilter struct {
	Provider    string
	Model       string
	Query       string
	StartUnixMs int64
	EndUnixMs   int64

	WorkItemKey      string
	WorkItemProvider string
	Sort             GroupSortSpec
	Offset           int
	Limit            int
}

type GroupListResult struct {
	Groups     []GroupRow
	TotalCount int64
}

var groupSortColumns = map[string]string{
	"cost":          "cost",
	"session_count": "session_count",
	"name":          "title",
	"duration":      "duration_ms",
	"input_tokens":  "input_tokens",
	"output_tokens": "output_tokens",
}

func buildFilteredSessionsWhere(filter GroupListFilter) (string, []any) {
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
	return strings.Join(where, " AND "), args
}

func buildGroupsCTE(filter GroupListFilter) (string, []any) {
	whereSQL, args := buildFilteredSessionsWhere(filter)
	cte := `
		filtered AS (
			SELECT s.cost, s.duration_ms, s.input_tokens, s.output_tokens,
				c.status, c.work_item_key, c.work_item_title, c.provider AS work_item_provider
			FROM sessions s
			LEFT JOIN claims c ON c.session_id = s.session_id
			WHERE ` + whereSQL + `
		),
		grouped AS (
			SELECT
				work_item_key, work_item_provider, work_item_title AS title, 0 AS is_unclaimed,
				COUNT(*) AS session_count, SUM(cost) AS cost,
				SUM(duration_ms) AS duration_ms, SUM(input_tokens) AS input_tokens, SUM(output_tokens) AS output_tokens
			FROM filtered
			WHERE status = 'claimed'
			GROUP BY work_item_provider, work_item_key
			UNION ALL
			SELECT
				'' AS work_item_key, '' AS work_item_provider, 'Unclaimed' AS title, 1 AS is_unclaimed,
				COUNT(*) AS session_count, SUM(cost) AS cost,
				SUM(duration_ms) AS duration_ms, SUM(input_tokens) AS input_tokens, SUM(output_tokens) AS output_tokens
			FROM filtered
			WHERE status IS NULL OR status != 'claimed'
			HAVING COUNT(*) > 0
		)`
	return cte, args
}

func buildGroupOrderBy(sort GroupSortSpec) string {
	column, ok := groupSortColumns[sort.Column]
	if !ok {
		column = "cost"
	}
	dir := "DESC"
	if sort.Dir == "asc" {
		dir = "ASC"
	}
	if column == "title" {
		return "is_unclaimed DESC, " + column + " " + dir
	}
	return "is_unclaimed DESC, " + column + " " + dir + ", title ASC"
}

func ListGroups(ctx context.Context, db Querier, filter GroupListFilter) (GroupListResult, error) {
	cte, args := buildGroupsCTE(filter)

	scopeSQL := ""
	scopeArgs := []any{}
	if filter.WorkItemKey != "" {
		scopeSQL = " WHERE work_item_key = ? AND work_item_provider = ?"
		scopeArgs = append(scopeArgs, filter.WorkItemKey, filter.WorkItemProvider)
	}

	limit := filter.Limit
	if limit <= 0 || limit > 500 {
		limit = 100
	}

	var result GroupListResult
	result.Groups = []GroupRow{}

	listQuery := "WITH " + cte + `
		SELECT work_item_key, work_item_provider, title, is_unclaimed, session_count, cost,
			duration_ms, input_tokens, output_tokens, COUNT(*) OVER () AS total_count
		FROM grouped` + scopeSQL + `
		ORDER BY ` + buildGroupOrderBy(filter.Sort) + `
		LIMIT ? OFFSET ?`
	listArgs := append(append([]any{}, args...), scopeArgs...)
	listArgs = append(listArgs, limit, filter.Offset)

	rows, err := db.QueryContext(ctx, listQuery, listArgs...)
	if err != nil {
		return GroupListResult{}, fmt.Errorf("list work item groups: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var row GroupRow
		var isUnclaimed int
		if err := rows.Scan(
			&row.WorkItemKey, &row.WorkItemProvider, &row.Title, &isUnclaimed,
			&row.SessionCount, &row.Cost,
			&row.DurationMs, &row.InputTokens, &row.OutputTokens,
			&result.TotalCount,
		); err != nil {
			return GroupListResult{}, fmt.Errorf("scan work item group: %w", err)
		}
		row.IsUnclaimed = isUnclaimed != 0
		result.Groups = append(result.Groups, row)
	}
	if err := rows.Err(); err != nil {
		return GroupListResult{}, fmt.Errorf("iterate work item groups: %w", err)
	}

	if len(result.Groups) == 0 && filter.Offset > 0 {
		countQuery := "WITH " + cte + " SELECT COUNT(*) FROM grouped" + scopeSQL
		countArgs := append(append([]any{}, args...), scopeArgs...)
		if err := db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&result.TotalCount); err != nil {
			return GroupListResult{}, fmt.Errorf("count work item groups: %w", err)
		}
	}
	return result, nil
}
