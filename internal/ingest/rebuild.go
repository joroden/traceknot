package ingest

import (
	"context"
	"fmt"

	"traceknot/internal/normalize/shared"
)

type RegroupFunc func(record shared.RawRecord) (nativeID string, ok bool)

func (receiver *Receiver) RebuildProvider(ctx context.Context, normalizer shared.Normalizer, regroup RegroupFunc) error {
	provider := normalizer.Provider()

	byNativeID, err := receiver.store.LoadRawSignalByProvider(ctx, provider)
	if err != nil {
		return fmt.Errorf("load raw_signal for rebuild: %w", err)
	}
	if len(byNativeID) == 0 {
		return nil
	}

	newByNativeID := make(map[string][]shared.RawRecord, len(byNativeID))
	for oldNativeID, rows := range byNativeID {
		for _, row := range rows {
			record := shared.RawRecord{
				NativeID:    row.NativeID,
				Signal:      row.Signal,
				DedupKey:    row.DedupKey,
				TimestampMs: row.TimestampMs,
				PayloadJSON: row.PayloadJSON,
			}
			if regroup != nil {
				if newNativeID, ok := regroup(record); ok && newNativeID != oldNativeID {
					if err := receiver.store.UpdateRawSignalNativeID(ctx, provider, oldNativeID, record.DedupKey, newNativeID); err != nil {
						return fmt.Errorf("persist regrouped native id: %w", err)
					}
					record.NativeID = newNativeID
				}
			}
			newByNativeID[record.NativeID] = append(newByNativeID[record.NativeID], record)
		}
	}

	touchedIDs := make([]string, 0, len(newByNativeID))
	for nativeID := range newByNativeID {
		touchedIDs = append(touchedIDs, nativeID)
	}

	results := normalizer.Rebuild(newByNativeID, touchedIDs)

	if linked, ok := normalizer.(shared.LinkedNormalizer); ok {
		if err := receiver.store.SaveConversationRoots(ctx, provider, linked.ResolveRoots(newByNativeID)); err != nil {
			return fmt.Errorf("save conversation roots during rebuild: %w", err)
		}
	}

	oldSessionIDs, err := receiver.store.ListSessionIDsForProvider(ctx, provider)
	if err != nil {
		return fmt.Errorf("list existing sessions for rebuild: %w", err)
	}
	newSessionIDs := make(map[string]bool, len(results))
	for _, result := range results {
		newSessionIDs[result.Seed.SessionID] = true
	}
	var orphaned []string
	for _, sessionID := range oldSessionIDs {
		if !newSessionIDs[sessionID] {
			orphaned = append(orphaned, sessionID)
		}
	}
	if len(orphaned) > 0 {
		if err := receiver.store.CleanupOrphanedSessions(ctx, orphaned); err != nil {
			return fmt.Errorf("clean up orphaned sessions: %w", err)
		}
	}

	for _, result := range results {
		if err := receiver.store.ReplaceSession(ctx, result.Seed, result.Content); err != nil {
			return fmt.Errorf("replace session %s during rebuild: %w", result.Seed.SessionID, err)
		}
	}
	return nil
}
