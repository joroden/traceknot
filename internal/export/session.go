package export

import (
	"context"
	"fmt"

	"traceknot/internal/store"
)

type Node struct {
	store.TreeRow
	Detail   *store.NodeDetail
	Children []*Node
}

type Session struct {
	Meta  *store.SessionMeta
	Root  *Node
	Nodes []*Node
}

func LoadSession(ctx context.Context, st *store.Store, sessionID string) (*Session, error) {
	meta, err := st.LoadSessionMeta(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("load session meta: %w", err)
	}
	if meta == nil {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}

	rows, err := st.LoadSessionTree(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("load session tree: %w", err)
	}

	nodes := make([]*Node, len(rows))
	byID := make(map[string]*Node, len(rows))
	for i, row := range rows {
		node := &Node{TreeRow: row}
		nodes[i] = node
		byID[row.NodeID] = node
	}

	var root *Node
	for _, node := range nodes {
		if node.HasContent {
			detail, err := st.LoadNodeDetail(ctx, node.NodeID)
			if err != nil {
				return nil, fmt.Errorf("load node detail %s: %w", node.NodeID, err)
			}
			node.Detail = detail
		}
		if node.ParentNodeID == nil {
			root = node
			continue
		}
		if parent, ok := byID[*node.ParentNodeID]; ok {
			parent.Children = append(parent.Children, node)
		}
	}

	return &Session{Meta: meta, Root: root, Nodes: nodes}, nil
}
