package cli

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"traceknot/internal/install/agentenv"
	"traceknot/internal/install/agenthooks"
	"traceknot/internal/platform"
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
	if err := platform.Current.AutostartDisable(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "uninstall: autostart:", err)
	}
	if err := platform.Current.RemoveBinary(exe); err != nil {
		fmt.Println("  - binary: could not remove: " + exe)
	}

	fmt.Println("Removed traceknot, local database preserved at " + filepath.Join(mustHome(), ".traceknot"))
	return 0
}
