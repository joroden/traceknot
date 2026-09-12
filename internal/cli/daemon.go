package cli

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"traceknot/internal/platform"
	"traceknot/internal/server"
)

func RunDaemon() int {
	ctx := context.Background()
	platform.Current.FreePort(ctx, portFromURL(defaultServerURL))

	if err := os.MkdirAll(daemonLogDir(), 0o755); err != nil {
		return 1
	}
	logFile, err := os.OpenFile(daemonLogPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return 1
	}
	logger := slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: slog.LevelInfo}))
	defer WriteDaemonPID()()
	if err := server.Run(defaultDBPath(), "127.0.0.1:4318", logger); err != nil {
		logger.Error("server failed", "error", err)
		return 1
	}
	return 0
}

func defaultDBPath() string {
	return filepath.Join(mustHome(), ".traceknot", "telemetry.sqlite")
}
