package conversationroots

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Querier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

func RootsFor(ctx context.Context, db Querier, provider string, nativeIDs []string) (map[string]string, error) {
	roots := make(map[string]string, len(nativeIDs))
	if len(nativeIDs) == 0 {
		return roots, nil
	}

	placeholders := make([]string, len(nativeIDs))
	args := make([]any, 0, len(nativeIDs)+1)
	args = append(args, provider)
	for index, nativeID := range nativeIDs {
		placeholders[index] = "?"
		args = append(args, nativeID)
	}

	rows, err := db.QueryContext(ctx, `
		SELECT native_id, root_native_id FROM conversation_roots
		WHERE provider = ? AND native_id IN (`+strings.Join(placeholders, ", ")+`)`, args...)
	if err != nil {
		return nil, fmt.Errorf("load conversation roots: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var nativeID, rootID string
		if err := rows.Scan(&nativeID, &rootID); err != nil {
			return nil, fmt.Errorf("scan conversation root: %w", err)
		}
		roots[nativeID] = rootID
	}
	return roots, rows.Err()
}

func FamilyMembers(ctx context.Context, db Querier, provider string, roots []string) ([]string, error) {
	members := []string{}
	if len(roots) == 0 {
		return members, nil
	}

	placeholders := make([]string, len(roots))
	args := make([]any, 0, len(roots)+1)
	args = append(args, provider)
	for index, root := range roots {
		placeholders[index] = "?"
		args = append(args, root)
	}

	rows, err := db.QueryContext(ctx, `
		SELECT native_id FROM conversation_roots
		WHERE provider = ? AND root_native_id IN (`+strings.Join(placeholders, ", ")+`)`, args...)
	if err != nil {
		return nil, fmt.Errorf("load conversation family: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var nativeID string
		if err := rows.Scan(&nativeID); err != nil {
			return nil, fmt.Errorf("scan conversation family member: %w", err)
		}
		members = append(members, nativeID)
	}
	return members, rows.Err()
}

func SaveRoots(ctx context.Context, db Querier, provider string, roots map[string]string) error {
	for nativeID, rootID := range roots {
		if nativeID == rootID {
			continue
		}
		if _, err := db.ExecContext(ctx, `
			INSERT INTO conversation_roots (provider, native_id, root_native_id)
			VALUES (?, ?, ?)
			ON CONFLICT (provider, native_id) DO UPDATE SET root_native_id = excluded.root_native_id`,
			provider, nativeID, rootID,
		); err != nil {
			return fmt.Errorf("save conversation root for %s: %w", nativeID, err)
		}
	}
	return nil
}
