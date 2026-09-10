package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"traceknot/internal/install/agentenv"
)

func RunSettings() int {
	ctx := context.Background()
	exe := resolveExe()

	if err := agentenv.ApplyEnv(filepath.Dir(exe)); err != nil {
		fmt.Fprintln(os.Stderr, "traceknot: env:", err)
		return 1
	}
	if err := agentenv.ApplyCodex(); err != nil {
		fmt.Fprintln(os.Stderr, "traceknot: codex config:", err)
	}

	if !daemonRunning(ctx) {
		if err := startDaemonNow(ctx); err != nil {
			fmt.Fprintln(os.Stderr, "traceknot: start:", err)
			return 1
		}
	}

	settingsURL := defaultServerURL + "/settings"
	if !canOpenBrowser() {
		fmt.Println("Open " + settingsURL + " in your browser to reconfigure traceknot.")
		return 0
	}
	if err := openWithDefaultHandler(settingsURL); err != nil {
		fmt.Println("Open " + settingsURL + " in your browser to reconfigure traceknot.")
		return 0
	}
	fmt.Println("Opened traceknot settings in your browser.")
	return 0
}
