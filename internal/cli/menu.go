package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/charmbracelet/huh"

	"traceknot/internal/install/agentenv"
	"traceknot/internal/install/autostart"
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
	if wasRunning {
		stopDaemon(ctx, defaultServerURL)
		if err := startDaemonNow(ctx); err != nil {
			fmt.Fprintln(os.Stderr, "traceknot: restart:", err)
		}
	}

	autostartWasOn, _ := autostartStatus(ctx)
	choices := providerChoices(exe)
	if !isFreshInstall(wasRunning, autostartWasOn, choices) {
		fmt.Println("Existing installation detected — run `traceknot` anytime to reconfigure.")
		return 0
	}

	if !IsInteractive() {
		fmt.Println("Run `traceknot` to update your settings (server, autostart, hooks).")
		return 0
	}

	return RunMenu()
}

func RunMenu() int {
	ctx := context.Background()
	exe := resolveExe()

	if err := agentenv.ApplyEnv(filepath.Dir(exe)); err != nil {
		fmt.Fprintln(os.Stderr, "traceknot: env:", err)
		return 1
	}
	if err := agentenv.ApplyCodex(); err != nil {
		fmt.Fprintln(os.Stderr, "traceknot: codex config:", err)
	}

	serverWasRunning := daemonRunning(ctx)
	autostartWasOn, autostartLabel := autostartStatus(ctx)
	choices := providerChoices(exe)
	fresh := isFreshInstall(serverWasRunning, autostartWasOn, choices)

	serverOn := serverWasRunning || fresh
	autostartOn := autostartWasOn || fresh
	selected := selectedHooks(choices)
	if fresh {
		selected = allHookBinaries(choices)
	}

	form := buildMenuForm(daemonStatusLine(ctx), autostartLabel, hookOptions(choices), &serverOn, &autostartOn, &selected)
	if err := form.Run(); err != nil {
		return 1
	}

	applyServerToggle(ctx, serverWasRunning, serverOn)
	applyAutostartToggle(ctx, autostartWasOn, autostartOn)
	disableAutostartIfRequested(ctx)
	applyHookSelection(ctx, choices, selected, exe)
	return 0
}

func buildMenuForm(status string, autostartLabel string, options []huh.Option[string], serverOn *bool, autostartOn *bool, selected *[]string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Description(status+"\n"+autostartLabel),
			huh.NewConfirm().Title("Server").Affirmative("On").Negative("Off").Value(serverOn),
			huh.NewConfirm().Title("Start on login").Affirmative("On").Negative("Off").Value(autostartOn),
		).Title("Server"),
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Session hooks").
				Description("Space to toggle, enter to confirm.").
				Options(options...).
				Value(selected),
		).Title("Hooks"),
	)
}

func hookOptions(choices []providerChoice) []huh.Option[string] {
	options := make([]huh.Option[string], 0, len(choices)+1)
	for _, choice := range choices {
		options = append(options, huh.NewOption("session hook for "+choice.provider.Binary, choice.provider.Binary))
	}
	return append(options, huh.NewOption(restartSessionsLabel(), "restart_sessions"))
}

func selectedHooks(choices []providerChoice) []string {
	selected := make([]string, 0, len(choices))
	for _, choice := range choices {
		if choice.want {
			selected = append(selected, choice.provider.Binary)
		}
	}
	return selected
}

func allHookBinaries(choices []providerChoice) []string {
	all := make([]string, 0, len(choices))
	for _, choice := range choices {
		all = append(all, choice.provider.Binary)
	}
	return all
}

func isFreshInstall(serverWasRunning bool, autostartWasOn bool, choices []providerChoice) bool {
	if serverWasRunning || autostartWasOn {
		return false
	}
	for _, choice := range choices {
		if choice.want {
			return false
		}
	}
	return true
}

func applyHookSelection(ctx context.Context, choices []providerChoice, selected []string, exe string) {
	for i := range choices {
		choices[i].want = slices.Contains(selected, choices[i].provider.Binary)
	}
	applyChoices(choices, exe)
	if slices.Contains(selected, "restart_sessions") {
		reportSessionRestart(ctx)
	}
}

func applyServerToggle(ctx context.Context, was bool, want bool) {
	if was == want {
		return
	}
	if want {
		if err := startDaemonNow(ctx); err != nil {
			fmt.Fprintln(os.Stderr, "traceknot: start:", err)
		}
		return
	}
	switch stopDaemon(ctx, defaultServerURL) {
	case "stopped":
		fmt.Println("daemon stopped")
	case "manual":
		fmt.Println("daemon is running but was not started by traceknot; stop it manually")
	}
}

func applyAutostartToggle(ctx context.Context, was bool, want bool) {
	if was == want {
		return
	}
	if want {
		if err := autostart.Enable(ctx); err != nil {
			fmt.Fprintln(os.Stderr, "traceknot: autostart:", err)
		} else {
			fmt.Println("autostart: enabled")
		}
		return
	}
	if err := autostart.Disable(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "traceknot: autostart:", err)
	} else {
		fmt.Println("autostart: disabled")
	}
}

func restartSessionsLabel() string {
	if agentenv.IsWSL() {
		return "Disconnect VS Code / close Copilot CLI so they reconnect with the new environment (the VS Code window itself stays open and reconnects on its own)"
	}
	return "Close and reopen VS Code / Copilot CLI (needed for them to pick up the new environment)"
}

func reportSessionRestart(ctx context.Context) {
	closed, err := agentenv.CloseInteractiveSessions(ctx)
	switch {
	case err != nil:
		fmt.Fprintln(os.Stderr, "traceknot: close vscode/copilot:", err)
	case closed == 0:
		fmt.Println("vscode/copilot: nothing running to close")
	case agentenv.IsWSL():
		fmt.Println("vscode/copilot: disconnected", closed, "running session(s); VS Code will reconnect automatically with the new environment (Copilot CLI sessions need a manual restart)")
	default:
		fmt.Println("vscode/copilot: closed", closed, "running session(s) and reopened VS Code")
	}
}

func disableAutostartIfRequested(ctx context.Context) {
	if os.Getenv("TRACEKNOT_NO_AUTOSTART") == "" {
		return
	}
	if err := autostart.Disable(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "autostart: disable (TRACEKNOT_NO_AUTOSTART set):", err)
		return
	}
	fmt.Println("autostart: disabled (TRACEKNOT_NO_AUTOSTART set)")
}
