package export

import (
	"context"
	"fmt"

	"traceknot/internal/store"
)

type WorkItem struct {
	Key      string
	Provider string
	Title    string
	Sessions []*Session
}

func LoadWorkItem(ctx context.Context, st *store.Store, workItemKey, workItemProvider string) (*WorkItem, error) {
	result, err := st.ListSessions(ctx, store.ListFilter{
		WorkItemKey:      workItemKey,
		WorkItemProvider: workItemProvider,
		Limit:            500,
	})
	if err != nil {
		return nil, fmt.Errorf("list work item sessions: %w", err)
	}
	if len(result.Rows) == 0 {
		return nil, fmt.Errorf("no sessions found for work item %s/%s", workItemProvider, workItemKey)
	}

	sessions := make([]*Session, 0, len(result.Rows))
	for _, row := range result.Rows {
		sess, err := LoadSession(ctx, st, row.SessionID)
		if err != nil {
			return nil, fmt.Errorf("load session %s: %w", row.SessionID, err)
		}
		sessions = append(sessions, sess)
	}

	return &WorkItem{
		Key:      workItemKey,
		Provider: workItemProvider,
		Title:    result.Rows[0].WorkItemTitle,
		Sessions: sessions,
	}, nil
}
