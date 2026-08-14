package claudeusage

import (
	"bufio"
	"io"
	"os"

	"traceknot/internal/normalize/claude"
	"traceknot/internal/normalize/shared"
)

func (watcher *Watcher) readNewLines(file, sessionID string) []shared.RawRecord {
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
			if record, ok := watcher.parseLine(sessionID, line); ok {
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

func (watcher *Watcher) parseLine(sessionID string, line []byte) (shared.RawRecord, bool) {
	event, ok := claude.ParseUsageSupplementLine(sessionID, line)
	if !ok {
		return shared.RawRecord{}, false
	}
	return claude.UsageSupplementRecord(event)
}
