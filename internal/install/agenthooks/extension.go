package agenthooks

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const copilotExtensionName = "user:traceknot"

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
	return nil
}

func extensionRemove(path string, _ string) error {
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
	escapedExe, err := json.Marshal(exe)
	if err != nil {
		return false
	}
	return strings.Contains(string(content), string(escapedExe))
}

func SeedCopilotExtensionPermission(location string) error {
	settingsPath, err := copilotPermissionsPath()
	if err != nil {
		return err
	}
	document, err := readJSONMap(settingsPath)
	if err != nil {
		return err
	}
	locations, _ := document["locations"].(map[string]any)
	if locations == nil {
		locations = map[string]any{}
		document["locations"] = locations
	}
	entry, _ := locations[location].(map[string]any)
	if entry == nil {
		entry = map[string]any{}
		locations[location] = entry
	}
	approvals, _ := entry["tool_approvals"].([]any)
	for _, approval := range approvals {
		fields, ok := approval.(map[string]any)
		if !ok {
			continue
		}
		if fields["kind"] == "extension-permission-access" && fields["extensionName"] == copilotExtensionName {
			return nil
		}
	}
	entry["tool_approvals"] = append(approvals, map[string]any{
		"kind":          "extension-permission-access",
		"extensionName": copilotExtensionName,
	})
	return writeJSON(settingsPath, document)
}

func removeCopilotExtensionPermissions() error {
	settingsPath, err := copilotPermissionsPath()
	if err != nil {
		return err
	}
	if !fileExists(settingsPath) {
		return nil
	}
	document, err := readJSONMap(settingsPath)
	if err != nil {
		return err
	}
	locations, _ := document["locations"].(map[string]any)
	changed := false
	for _, value := range locations {
		entry, ok := value.(map[string]any)
		if !ok {
			continue
		}
		approvals, ok := entry["tool_approvals"].([]any)
		if !ok {
			continue
		}
		kept := make([]any, 0, len(approvals))
		for _, approval := range approvals {
			fields, ok := approval.(map[string]any)
			if ok && fields["kind"] == "extension-permission-access" && fields["extensionName"] == copilotExtensionName {
				changed = true
				continue
			}
			kept = append(kept, approval)
		}
		entry["tool_approvals"] = kept
	}
	if !changed {
		return nil
	}
	return writeJSON(settingsPath, document)
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
