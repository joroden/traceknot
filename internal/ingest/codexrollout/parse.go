package codexrollout

import (
	"bufio"
	"io"
	"os"

	"traceknot/internal/normalize/codex"
	"traceknot/internal/normalize/shared"
)

func (watcher *Watcher) readNewLines(file string) []shared.RawRecord {
	handle, err := os.Open(file)
	if err != nil {
		return nil
	}
	defer handle.Close()

	offset := watcher.offsets[file]
	if _, err := handle.Seek(offset, io.SeekStart); err != nil {
		return nil
	}

	var records []shared.RawRecord
	reader := bufio.NewReader(handle)
	read := offset
	for {
		line, err := reader.ReadBytes('\n')
		if len(line) > 0 && err == nil {
			read += int64(len(line))
			if record, ok := watcher.parseLine(file, line); ok {
				records = append(records, record)
			}
		}
		if err != nil {
			break
		}
	}
	watcher.offsets[file] = read
	return records
}

func (watcher *Watcher) parseLine(file string, line []byte) (shared.RawRecord, bool) {
	conversationID, resolved := watcher.conversationIDs[file]
	if !resolved {
		if id, ok := codex.RolloutConversationID(line); ok {
			watcher.conversationIDs[file] = id
		}
		return shared.RawRecord{}, false
	}
	event, ok := codex.ParseRolloutLine(conversationID, line)
	if !ok {
		return shared.RawRecord{}, false
	}
	return codex.RolloutRecord(event)
}
