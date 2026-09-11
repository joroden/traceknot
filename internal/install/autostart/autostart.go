package autostart

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
)

func home() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return dir
}

func daemonLogPath() string {
	return filepath.Join(home(), ".traceknot", "daemon.log")
}

func runQuiet(ctx context.Context, name string, args ...string) {
	_ = exec.CommandContext(ctx, name, args...).Run()
}
