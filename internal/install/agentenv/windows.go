//go:build windows

package agentenv

import (
	"fmt"
	"os"
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
	return nil
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
