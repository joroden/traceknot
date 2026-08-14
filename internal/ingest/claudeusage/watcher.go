package claudeusage

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"traceknot/internal/ingest"
	"traceknot/internal/normalize/shared"
)

const pollInterval = 3 * time.Second

type Watcher struct {
	receiver    *ingest.Receiver
	normalizer  shared.Normalizer
	logger      *slog.Logger
	projectsDir string

	offsets map[string]int64
}

func New(receiver *ingest.Receiver, normalizer shared.Normalizer, logger *slog.Logger) *Watcher {
	return &Watcher{
		receiver:    receiver,
		normalizer:  normalizer,
		logger:      logger,
		projectsDir: projectsDir(),
		offsets:     make(map[string]int64),
	}
}

func (watcher *Watcher) Run(ctx context.Context) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	for {
		watcher.poll(ctx)
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (watcher *Watcher) poll(ctx context.Context) {
	var records []shared.RawRecord
	records = append(records, watcher.pollGlob(filepath.Join(watcher.projectsDir, "*", "*.jsonl"), mainSessionID)...)
	records = append(records, watcher.pollGlob(filepath.Join(watcher.projectsDir, "*", "*", "subagents", "*.jsonl"), subagentSessionID)...)
	if len(records) == 0 {
		return
	}
	watcher.receiver.Ingest(ctx, watcher.normalizer, records)
}

func (watcher *Watcher) pollGlob(pattern string, sessionIDFromPath func(string) string) []shared.RawRecord {
	files, err := filepath.Glob(pattern)
	if err != nil {
		watcher.logger.Warn("claude usage glob failed", "error", err)
		return nil
	}
	var records []shared.RawRecord
	for _, file := range files {
		sessionID := sessionIDFromPath(file)
		if sessionID == "" {
			continue
		}
		records = append(records, watcher.readNewLines(file, sessionID)...)
	}
	return records
}

func mainSessionID(path string) string {
	base := filepath.Base(path)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

func subagentSessionID(path string) string {
	subagentsDir := filepath.Dir(path)
	return filepath.Base(filepath.Dir(subagentsDir))
}

func projectsDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".claude", "projects")
}
