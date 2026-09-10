package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"traceknot/internal/install/agentenv"
)

func RunPostInstall() int {
	ctx := context.Background()
	exe := resolveExe()

	if err := agentenv.ApplyEnv(filepath.Dir(exe)); err != nil {
		fmt.Fprintln(os.Stderr, "traceknot: env:", err)
		return 1
	}
	if err := agentenv.ApplyCodex(); err != nil {
		fmt.Fprintln(os.Stderr, "traceknot: codex config:", err)
	}

	if daemonRunning(ctx) {
		stopDaemon(ctx, defaultServerURL)
	}
	if err := startDaemonNow(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "traceknot: start:", err)
		return 1
	}

	setupURL := defaultServerURL + "/setup"
	if !canOpenBrowser() {
		fmt.Println("Open " + setupURL + " in your browser to review settings (hooks, work items, skills).")
		return 0
	}
	if err := openBrowser(setupURL); err != nil {
		fmt.Println("Open " + setupURL + " in your browser to review settings (hooks, work items, skills).")
		return 0
	}
	fmt.Println("Review your settings in the browser window that just opened.")
	return 0
}
