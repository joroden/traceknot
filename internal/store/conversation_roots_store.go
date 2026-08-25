package store

import (
	"context"

	"traceknot/internal/store/conversationroots"
)

func (store *Store) ConversationRoots(ctx context.Context, provider string, nativeIDs []string) (map[string]string, error) {
	return conversationroots.RootsFor(ctx, store.db, provider, nativeIDs)
}

func (store *Store) ConversationFamily(ctx context.Context, provider string, roots []string) ([]string, error) {
	return conversationroots.FamilyMembers(ctx, store.db, provider, roots)
}

func (store *Store) SaveConversationRoots(ctx context.Context, provider string, roots map[string]string) error {
	return conversationroots.SaveRoots(ctx, store.db, provider, roots)
}
