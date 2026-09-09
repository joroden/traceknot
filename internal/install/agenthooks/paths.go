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

func copilotHome() (string, error) {
	copilotHome := os.Getenv("COPILOT_HOME")
	if copilotHome != "" {
		return copilotHome, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home: %w", err)
	}
	return filepath.Join(home, ".copilot"), nil
}

func copilotExtensionPath() (string, error) {
	home, err := copilotHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "extensions", "traceknot", "extension.mjs"), nil
}

func copilotSettingsPath() (string, error) {
	home, err := copilotHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "settings.json"), nil
}
