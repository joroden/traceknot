package store

import (
	"context"
	"database/sql"

	"traceknot/internal/model"
	"traceknot/internal/store/node"
)

func (store *Store) UpsertNodes(ctx context.Context, tx *sql.Tx, sessionID string, content *model.SessionContent) error {
	for _, seed := range content.Chats {
		if err := node.UpsertBase(ctx, tx, sessionID, &seed.NodeSeed, chatHasContent(seed)); err != nil {
			return err
		}
		if err := node.UpsertChat(ctx, tx, seed); err != nil {
			return err
		}
	}
	for _, seed := range content.ToolCalls {
		if err := node.UpsertBase(ctx, tx, sessionID, &seed.NodeSeed, toolCallHasContent(seed)); err != nil {
			return err
		}
		if err := node.UpsertToolCall(ctx, tx, seed); err != nil {
			return err
		}
	}
	for _, seed := range content.Agents {
		if err := node.UpsertBase(ctx, tx, sessionID, &seed.NodeSeed, agentHasContent(seed)); err != nil {
			return err
		}
		if err := node.UpsertAgent(ctx, tx, seed); err != nil {
			return err
		}
	}
	return nil
}

func (store *Store) ReassignOwners(ctx context.Context, tx *sql.Tx, sessionID string) error {
	return node.ReassignOwners(ctx, tx, sessionID)
}

func (store *Store) LoadNodeDetail(ctx context.Context, nodeID string) (*NodeDetail, error) {
	return node.LoadDetail(ctx, store.db, nodeID)
}

func chatHasContent(seed *model.ChatSeed) bool {
	return seed.SystemText != "" || seed.PromptText != "" || seed.OutputText != "" || seed.ReasoningText != ""
}

func toolCallHasContent(seed *model.ToolCallSeed) bool {
	return seed.ArgumentsJSON != "" || seed.ResultText != "" || seed.ErrorDetailsJSON != ""
}

func agentHasContent(seed *model.AgentSeed) bool {
	return seed.SpawnPrompt != ""
}
