//go:build windows

package agentenv

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"traceknot/internal/platform"
)

func ApplyEnv(binDir string) error {
	envBlock, err := loadSnippet("windows/env.ps1")
	if err != nil {
		return err
	}
	return platform.Current.ApplyEnv(pairMap(windowsEnvPairs(envBlock)), binDir)
}

func RemoveEnv() error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("resolve binary: %w", err)
	}
	envBlock, err := loadSnippet("windows/env.ps1")
	if err != nil {
		return err
	}
	return platform.Current.RemoveEnv(pairNames(windowsEnvPairs(envBlock)), filepath.Dir(exe))
}

func ApplyCopilotShim(exe string) error {
	shimBlock, err := loadSnippet("windows/copilot-shim.ps1")
	if err != nil {
		return err
	}
	rendered := strings.ReplaceAll(string(shimBlock), "{{TRACEKNOT_BIN}}", exe)
	for _, profile := range powershellProfileCandidates() {
		if err := upsertBlock(profile, rendered); err != nil {
			return err
		}
	}
	for _, host := range powershellHosts {
		changed, err := ensureProfileScriptsCanRun(host.exe)
		switch {
		case changed:
			fmt.Println(host.label + ": traceknot's Copilot integration needs permission to run a script on startup, which was turned off on this account by default; traceknot turned it on. To undo this, run in " + host.label + ": Set-ExecutionPolicy -Scope CurrentUser Restricted")
		case err != nil:
			fmt.Fprintln(os.Stderr, "traceknot: couldn't enable "+host.label+" scripts for the Copilot integration:", err)
		}
	}
	return nil
}

var powershellHosts = []struct {
	exe   string
	label string
}{
	{"powershell", "Windows PowerShell"},
	{"pwsh", "PowerShell 7"},
}

func ensureProfileScriptsCanRun(shell string) (bool, error) {
	if _, err := exec.LookPath(shell); err != nil {
		return false, nil
	}
	out, err := exec.Command(shell, "-NoProfile", "-NonInteractive", "-Command", "Get-ExecutionPolicy").Output()
	if err != nil {
		return false, fmt.Errorf("check policy: %w", err)
	}
	if strings.TrimSpace(string(out)) != "Restricted" {
		return false, nil
	}
	output, err := exec.Command(shell, "-NoProfile", "-NonInteractive", "-Command",
		"Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy RemoteSigned -Force").CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("set policy: %s", strings.TrimSpace(string(output)))
	}
	return true, nil
}

func RemoveCopilotShim() error {
	shimBlock, err := loadSnippet("windows/copilot-shim.ps1")
	if err != nil {
		return err
	}
	start, end := blockMarkers(shimBlock)
	for _, profile := range powershellProfileCandidates() {
		if err := removeBlock(profile, start, end); err != nil {
			return err
		}
	}
	return nil
}

func powershellProfileCandidates() []string {
	home := homeDir()
	return []string{
		filepath.Join(home, "Documents", "PowerShell", "profile.ps1"),
		filepath.Join(home, "Documents", "WindowsPowerShell", "profile.ps1"),
	}
}
