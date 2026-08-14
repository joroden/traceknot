package node

import (
	"context"
	"database/sql"
	"fmt"
)

type ownerRow struct {
	nodeID    string
	parent    *string
	kind      string
	agentName *string
}

func ReassignOwners(ctx context.Context, db Querier, sessionID string) error {
	rows, err := db.QueryContext(ctx, `
		SELECT n.node_id, n.parent_node_id, n.kind, a.agent_name
		FROM nodes n
		LEFT JOIN agent_nodes a ON a.node_id = n.node_id
		WHERE n.session_id = ?`, sessionID)
	if err != nil {
		return fmt.Errorf("select session nodes: %w", err)
	}
	defer rows.Close()

	var all []ownerRow
	byNodeID := map[string]int{}
	for rows.Next() {
		var entry ownerRow
		var parent sql.NullString
		if err := rows.Scan(&entry.nodeID, &parent, &entry.kind, &entry.agentName); err != nil {
			return fmt.Errorf("scan session node: %w", err)
		}
		if parent.Valid {
			entry.parent = &parent.String
		}
		index := len(all)
		all = append(all, entry)
		byNodeID[entry.nodeID] = index
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate session nodes: %w", err)
	}

	update, err := db.PrepareContext(ctx, `
		UPDATE nodes SET owning_agent_id = ?, owning_agent_name = ? WHERE node_id = ?`)
	if err != nil {
		return fmt.Errorf("prepare owner update: %w", err)
	}
	defer update.Close()

	for index := range all {
		owner := ownerIndex(all, index, byNodeID)
		var ownerID, ownerName any
		if owner >= 0 {
			ownerID = all[owner].nodeID
			ownerName = all[owner].agentName
		}
		if _, err := update.ExecContext(ctx, ownerID, ownerName, all[index].nodeID); err != nil {
			return fmt.Errorf("update owner for %s: %w", all[index].nodeID, err)
		}
	}
	return nil
}

func ownerIndex(all []ownerRow, start int, byNodeID map[string]int) int {
	current := start
	visited := map[int]bool{start: true}
	for {
		entry := all[current]
		if entry.kind == "agent" && current != start {
			return current
		}
		if entry.parent == nil {
			return -1
		}
		next, ok := byNodeID[*entry.parent]
		if !ok || visited[next] {
			return -1
		}
		visited[next] = true
		current = next
	}
}
