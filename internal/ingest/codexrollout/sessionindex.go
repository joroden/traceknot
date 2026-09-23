package codexrollout

import (
	"bufio"
	"os"
	"path/filepath"

	"traceknot/internal/normalize/codex"
	"traceknot/internal/normalize/shared"
)

func (watcher *Watcher) readSessionIndex() []shared.RawRecord {
	file := filepath.Join(filepath.Dir(watcher.sessionsDir), "session_index.jsonl")
	handle, err := os.Open(file)
	if err != nil {
		return nil
	}
	defer handle.Close()

	latest := make(map[string]shared.RawRecord)
	scanner := bufio.NewScanner(handle)
	for scanner.Scan() {
		if record, ok := codex.SessionIndexRecord(scanner.Bytes()); ok {
			previous, exists := latest[record.NativeID]
			if !exists || record.TimestampMs >= previous.TimestampMs {
				latest[record.NativeID] = record
			}
		}
	}
	if err := scanner.Err(); err != nil {
		watcher.logger.Warn("codex session index read failed", "error", err)
		return nil
	}

	var records []shared.RawRecord
	for id, record := range latest {
		if watcher.indexEntries[id] == record.DedupKey {
			continue
		}
		watcher.indexEntries[id] = record.DedupKey
		records = append(records, record)
	}
	return records
}
