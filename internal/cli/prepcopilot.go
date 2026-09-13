package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"traceknot/internal/install/agenthooks"
)

// RunPrepCopilot seeds Copilot's permissions-config.json so the traceknot
// extension can load for the current repo, without requiring an interactive
// approval. It is invoked by the copilot shell function before every
// `copilot` launch and must never block that launch, so it always returns 0.
func RunPrepCopilot() int {
	if err := agenthooks.SeedCopilotExtensionPermission(copilotLocation()); err != nil {
		fmt.Fprintln(os.Stderr, "traceknot: prep-copilot:", err)
	}
	return 0
}

func copilotLocation() string {
	output, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err == nil {
		if root := strings.TrimSpace(string(output)); root != "" {
			return filepath.Clean(root)
		}
	}
	if cwd, err := os.Getwd(); err == nil {
		return filepath.Clean(cwd)
	}
	return "."
}
