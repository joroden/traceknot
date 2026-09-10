package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"traceknot/internal/install/agentenv"
	"traceknot/internal/install/autostart"
	"traceknot/internal/settings"
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

	wasRunning := daemonRunning(ctx)
	fresh := isFreshInstall(ctx, wasRunning)

	if wasRunning {
		stopDaemon(ctx, defaultServerURL)
	}
	if err := startDaemonNow(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "traceknot: start:", err)
		return 1
	}

	if !fresh {
		fmt.Println("Existing installation detected — run `traceknot` anytime to reconfigure.")
		return 0
	}

	setupURL := defaultServerURL + "/setup"
	if !canOpenBrowser() {
		fmt.Println("Open " + setupURL + " in your browser to finish setup (hooks, work items, skills).")
		return 0
	}
	if err := openBrowser(setupURL); err != nil {
		fmt.Println("Open " + setupURL + " in your browser to finish setup (hooks, work items, skills).")
		return 0
	}
	fmt.Println("Finish setup in the browser window that just opened.")
	return 0
}

func isFreshInstall(ctx context.Context, wasRunning bool) bool {
	if wasRunning || autostart.Enabled(ctx) {
		return false
	}
	state, err := settings.Current(ctx)
	if err != nil {
		return true
	}
	for _, hook := range state.Hooks {
		if hook.Enabled {
			return false
		}
	}
	return true
}
