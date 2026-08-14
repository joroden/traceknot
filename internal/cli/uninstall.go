package cli

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"traceknot/internal/install/agentenv"
	"traceknot/internal/install/agenthooks"
	"traceknot/internal/install/autostart"
)

func RunUninstall(args []string) int {
	_ = flag.NewFlagSet("uninstall", flag.ExitOnError).Parse(args)

	fmt.Println("Removing traceknot...")

	ctx := context.Background()
	exe := resolveExe()

	for _, provider := range agenthooks.Providers {
		if err := provider.Remove(exe); err != nil {
			fmt.Fprintln(os.Stderr, "uninstall: "+provider.Binary+":", err)
		}
	}
	if err := agentenv.RemoveEnv(); err != nil {
		fmt.Fprintln(os.Stderr, "uninstall: env:", err)
	}
	if err := agentenv.RemoveCodex(); err != nil {
		fmt.Fprintln(os.Stderr, "uninstall: codex config:", err)
	}
	if result := stopDaemon(ctx, defaultServerURL); result == "manual" {
		fmt.Fprintln(os.Stderr, "daemon is running but was not started by traceknot; stop it manually")
	}
	if err := autostart.Disable(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "uninstall: autostart:", err)
	}
	removeBinary(exe)

	fmt.Println("Removed traceknot, local database preserved at " + filepath.Join(mustHome(), ".traceknot"))
	return 0
}

func removeBinary(exe string) {
	if runtime.GOOS != "windows" {
		if err := os.Remove(exe); err == nil {
			return
		}
		fmt.Println("  - binary: could not remove: " + exe)
		return
	}
	stale := exe + ".tk-old"
	if err := os.Rename(exe, stale); err != nil {
		fmt.Println("  - binary: could not remove: " + exe + " (delete it manually)")
		return
	}
	batch := filepath.Join(os.TempDir(), "traceknot-uninstall.ps1")
	script := "Start-Sleep -Milliseconds 500; Remove-Item -Force -LiteralPath '" + stale + "'; Remove-Item -Force -LiteralPath '" + batch + "'"
	if err := os.WriteFile(batch, []byte(script), 0o644); err == nil {
		command := exec.Command("powershell", "-NoProfile", "-WindowStyle", "Hidden",
			"-ExecutionPolicy", "Bypass", "-File", batch)
		_ = command.Start()
	}
}
