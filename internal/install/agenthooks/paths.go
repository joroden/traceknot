package agenthooks

import (
	"fmt"
	"os"
	"path/filepath"
)

func claudeSettingsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home: %w", err)
	}
	return filepath.Join(home, ".claude", "settings.json"), nil
}

func codexHooksPath() (string, error) {
	codexHome := os.Getenv("CODEX_HOME")
	if codexHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("resolve home: %w", err)
		}
		codexHome = filepath.Join(home, ".codex")
	}
	return filepath.Join(codexHome, "hooks.json"), nil
}

func copilotHooksPath() (string, error) {
	copilotHome := os.Getenv("COPILOT_HOME")
	if copilotHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("resolve home: %w", err)
		}
		copilotHome = filepath.Join(home, ".copilot")
	}
	return filepath.Join(copilotHome, "hooks", "traceknot.json"), nil
}
