package agenthooks

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func extensionInstall(path string, exe string) error {
	template, err := ExtensionTemplate()
	if err != nil {
		return err
	}
	escapedExe, err := json.Marshal(exe)
	if err != nil {
		return fmt.Errorf("escape binary path: %w", err)
	}
	rendered := strings.ReplaceAll(string(template), `"TRACEKNOT_BIN_PLACEHOLDER"`, string(escapedExe))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(rendered), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return enableCopilotExperimental()
}

func enableCopilotExperimental() error {
	settingsPath, err := copilotSettingsPath()
	if err != nil {
		return err
	}
	document, err := readJSONMap(settingsPath)
	if err != nil {
		return err
	}
	if enabled, _ := document["experimental"].(bool); enabled {
		return nil
	}
	document["experimental"] = true
	return writeJSON(settingsPath, document)
}

func extensionRemove(path string, exe string) error {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove %s: %w", path, err)
	}
	_ = os.Remove(filepath.Dir(path))
	return nil
}

func extensionInstalled(path string, exe string) bool {
	content, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return strings.Contains(string(content), exe)
}
