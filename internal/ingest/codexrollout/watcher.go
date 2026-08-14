package codexrollout

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"traceknot/internal/ingest"
	"traceknot/internal/normalize/shared"
)

const pollInterval = 3 * time.Second

type Watcher struct {
	receiver    *ingest.Receiver
	normalizer  shared.Normalizer
	logger      *slog.Logger
	sessionsDir string

	offsets         map[string]int64
	conversationIDs map[string]string
}

func New(receiver *ingest.Receiver, normalizer shared.Normalizer, logger *slog.Logger) *Watcher {
	return &Watcher{
		receiver:        receiver,
		normalizer:      normalizer,
		logger:          logger,
		sessionsDir:     sessionsDir(),
		offsets:         make(map[string]int64),
		conversationIDs: make(map[string]string),
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
	files, err := filepath.Glob(filepath.Join(watcher.sessionsDir, "*", "*", "*", "rollout-*.jsonl"))
	if err != nil {
		watcher.logger.Warn("codex rollout glob failed", "error", err)
		return
	}
	var records []shared.RawRecord
	for _, file := range files {
		records = append(records, watcher.readNewLines(file)...)
	}
	if len(records) == 0 {
		return
	}
	watcher.receiver.Ingest(ctx, watcher.normalizer, records)
}

func sessionsDir() string {
	codexHome := os.Getenv("CODEX_HOME")
	if codexHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		codexHome = filepath.Join(home, ".codex")
	}
	return filepath.Join(codexHome, "sessions")
}
