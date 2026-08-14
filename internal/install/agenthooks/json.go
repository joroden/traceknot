package agenthooks

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func applyTemplate(path string, template []byte, exe string) error {
	escapedExe, marshalErr := json.Marshal(exe)
	if marshalErr != nil {
		return fmt.Errorf("escape binary path: %w", marshalErr)
	}
	placeholder := string(escapedExe[1 : len(escapedExe)-1])
	rendered := strings.ReplaceAll(string(template), "{{TRACEKNOT_BIN}}", placeholder)
	var templateDoc map[string]any
	if err := json.Unmarshal([]byte(rendered), &templateDoc); err != nil {
		return fmt.Errorf("parse embedded template: %w", err)
	}

	document, err := readJSONMap(path)
	if err != nil {
		return err
	}
	for key, value := range templateDoc {
		if key == "hooks" {
			continue
		}
		document[key] = value
	}

	templateHooks, _ := templateDoc["hooks"].(map[string]any)
	documentHooks, _ := document["hooks"].(map[string]any)
	if documentHooks == nil {
		documentHooks = map[string]any{}
		document["hooks"] = documentHooks
	}

	for event, groups := range documentHooks {
		eventGroups, _ := groups.([]any)
		documentHooks[event] = removeTraceknotGroups(eventGroups, exe)
	}
	for event, templateGroups := range templateHooks {
		groups, _ := documentHooks[event].([]any)
		templateGroupsValue, _ := templateGroups.([]any)
		documentHooks[event] = append(groups, templateGroupsValue...)
	}
	return writeJSON(path, document)
}

func removeTraceknotGroups(groups []any, exe string) []any {
	kept := groups[:0]
	for _, group := range groups {
		entry, ok := group.(map[string]any)
		if !ok {
			kept = append(kept, group)
			continue
		}
		if containsTraceknotCommand(entry, exe) {
			continue
		}
		kept = append(kept, group)
	}
	return kept
}

func containsTraceknotCommand(group map[string]any, exe string) bool {
	for _, field := range []string{"command", "bash", "powershell"} {
		if command, ok := group[field].(string); ok && isTraceknotCommand(command, exe) {
			return true
		}
	}
	handlers, _ := group["hooks"].([]any)
	for _, handler := range handlers {
		item, ok := handler.(map[string]any)
		if !ok {
			continue
		}
		for _, field := range []string{"command", "bash", "powershell"} {
			if command, ok := item[field].(string); ok && isTraceknotCommand(command, exe) {
				return true
			}
		}
	}
	return false
}

func isTraceknotCommand(command string, exe string) bool {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return false
	}
	return fields[0] == exe || filepath.Base(fields[0]) == "traceknot"
}

func readJSONMap(path string) (map[string]any, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]any{}, nil
		}
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var document map[string]any
	if err := json.Unmarshal(content, &document); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return document, nil
}

func writeJSON(path string, document map[string]any) error {
	content, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		return fmt.Errorf("encode %s: %w", path, err)
	}
	content = append(content, '\n')
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, content, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}
