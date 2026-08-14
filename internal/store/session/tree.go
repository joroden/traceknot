package session

import (
	"context"
	"database/sql"
	"fmt"
	"math"

	"traceknot/internal/ptr"
)

type TreeRow struct {
	NodeID                string
	Kind                  string
	Name                  *string
	AgentName             *string
	ToolName              *string
	ToolCallID            *string
	Model                 *string
	Status                *string
	StartedAtUnixMs       *int64
	DurationMs            *float64
	InputTokens           int64
	CachedInputTokens     int64
	CacheWriteTokens      int64
	OutputTokens          int64
	ReasoningTokens       int64
	Cost                  float64
	EstimatedInputTokens  int64
	EstimatedOutputTokens int64
	OwningAgentID         *string
	OwningAgentName       *string
	ParentNodeID          *string
	HasContent            bool
	AggInputTokens        int64
	AggCachedInputTokens  int64
	AggCacheWriteTokens   int64
	AggOutputTokens       int64
	AggReasoningTokens    int64
	AggCost               float64
	DescendantCount       int64
	SubagentCount         int64
}

func LoadTree(ctx context.Context, db Querier, sessionID string) ([]TreeRow, error) {
	all, err := queryTreeRows(ctx, db, sessionID)
	if err != nil {
		return nil, err
	}

	children, roots := buildChildIndex(all)
	all, children, roots = unifyRoots(all, children, roots, sessionID)
	aggregateTree(all, children, roots)
	return all, nil
}

func queryTreeRows(ctx context.Context, db Querier, sessionID string) ([]TreeRow, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT n.node_id, n.kind, n.name, a.agent_name, t.tool_name, t.tool_call_id,
			n.model, n.status, n.started_at_unix_ms, n.duration_ms,
			n.input_tokens, n.cached_input_tokens, n.cache_write_tokens,
			n.output_tokens, n.reasoning_tokens, n.cost,
			n.estimated_input_tokens, n.estimated_output_tokens,
			n.owning_agent_id, n.owning_agent_name, n.parent_node_id, n.has_content
		FROM nodes n
		LEFT JOIN agent_nodes a ON a.node_id = n.node_id
		LEFT JOIN tool_call_nodes t ON t.node_id = n.node_id
		WHERE n.session_id = ?
		ORDER BY COALESCE(n.started_at_unix_ms, 0), n.rowid`,
		sessionID,
	)
	if err != nil {
		return nil, fmt.Errorf("load session tree rows: %w", err)
	}
	defer rows.Close()

	var all []TreeRow
	for rows.Next() {
		var row TreeRow
		var parentNodeID sql.NullString
		if err := rows.Scan(
			&row.NodeID, &row.Kind, &row.Name, &row.AgentName, &row.ToolName, &row.ToolCallID,
			&row.Model, &row.Status, &row.StartedAtUnixMs, &row.DurationMs,
			&row.InputTokens, &row.CachedInputTokens, &row.CacheWriteTokens,
			&row.OutputTokens, &row.ReasoningTokens, &row.Cost,
			&row.EstimatedInputTokens, &row.EstimatedOutputTokens,
			&row.OwningAgentID, &row.OwningAgentName, &parentNodeID, &row.HasContent,
		); err != nil {
			return nil, fmt.Errorf("scan session tree row: %w", err)
		}
		if parentNodeID.Valid {
			row.ParentNodeID = &parentNodeID.String
		}
		all = append(all, row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate session tree rows: %w", err)
	}
	return all, nil
}

func buildChildIndex(all []TreeRow) ([][]int, []int) {
	children := make([][]int, len(all))
	byNodeID := make(map[string]int, len(all))
	for index, row := range all {
		byNodeID[row.NodeID] = index
	}
	var roots []int
	for index, row := range all {
		if row.ParentNodeID == nil {
			roots = append(roots, index)
			continue
		}
		parentIndex, ok := byNodeID[*row.ParentNodeID]
		if !ok {
			roots = append(roots, index)
			continue
		}
		children[parentIndex] = append(children[parentIndex], index)
	}
	return children, roots
}

func aggregateTree(all []TreeRow, children [][]int, roots []int) {
	visited := make([]bool, len(all))
	var walk func(int) treeAgg
	walk = func(index int) treeAgg {
		visited[index] = true
		agg := treeAgg{
			input:     all[index].InputTokens,
			cachedIn:  all[index].CachedInputTokens,
			cacheWr:   all[index].CacheWriteTokens,
			output:    all[index].OutputTokens,
			reasoning: all[index].ReasoningTokens,
			cost:      all[index].Cost,
		}
		for _, child := range children[index] {
			if visited[child] {
				continue
			}
			sub := walk(child)
			agg.input += sub.input
			agg.cachedIn += sub.cachedIn
			agg.cacheWr += sub.cacheWr
			agg.output += sub.output
			agg.reasoning += sub.reasoning
			agg.cost += sub.cost
			agg.descendants += sub.descendants + 1
			agg.subagents += sub.subagents
			if all[child].Kind == "agent" {
				agg.subagents++
			}
		}
		all[index].AggInputTokens = agg.input
		all[index].AggCachedInputTokens = agg.cachedIn
		all[index].AggCacheWriteTokens = agg.cacheWr
		all[index].AggOutputTokens = agg.output
		all[index].AggReasoningTokens = agg.reasoning
		all[index].AggCost = agg.cost
		all[index].DescendantCount = agg.descendants
		all[index].SubagentCount = agg.subagents
		return agg
	}
	for _, root := range roots {
		walk(root)
	}
}

type treeAgg struct {
	input       int64
	cachedIn    int64
	cacheWr     int64
	output      int64
	reasoning   int64
	cost        float64
	descendants int64
	subagents   int64
}

func unifyRoots(all []TreeRow, children [][]int, roots []int, sessionID string) ([]TreeRow, [][]int, []int) {
	if len(roots) <= 1 {
		return all, children, roots
	}
	synth := syntheticRootRow(roots, all, sessionID)
	index := len(all)
	all = append(all, synth)
	childSet := make([]int, 0, len(roots))
	for _, root := range roots {
		all[root].ParentNodeID = &synth.NodeID
		childSet = append(childSet, root)
	}
	children = append(children, childSet)
	return all, children, []int{index}
}

func syntheticRootRow(roots []int, all []TreeRow, sessionID string) TreeRow {
	var minStart int64 = math.MaxInt64
	var maxEnd float64
	valid := false
	for _, index := range roots {
		if all[index].StartedAtUnixMs == nil {
			continue
		}
		start := *all[index].StartedAtUnixMs
		if start < minStart {
			minStart = start
		}
		end := float64(start)
		if all[index].DurationMs != nil {
			end += *all[index].DurationMs
		}
		if end > maxEnd {
			maxEnd = end
		}
		valid = true
	}
	row := TreeRow{
		NodeID: "synthetic:root:" + sessionID,
		Kind:   "agent",
		Name:   ptr.String("main agent"),
		Status: ptr.String("ok"),
	}
	if valid {
		row.StartedAtUnixMs = &minStart
		duration := maxEnd - float64(minStart)
		row.DurationMs = &duration
	}
	return row
}
