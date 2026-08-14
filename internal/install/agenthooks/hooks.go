package agenthooks

import (
	"os"
)

type Provider struct {
	Binary string
	path   func() (string, error)
}

var Providers = []Provider{
	{Binary: "claude", path: claudeSettingsPath},
	{Binary: "codex", path: codexHooksPath},
	{Binary: "copilot", path: copilotHooksPath},
}

func (p Provider) ConfigPath() (string, error) {
	return p.path()
}

func (p Provider) Install(exe string) error {
	path, err := p.ConfigPath()
	if err != nil {
		return err
	}
	template, err := HookTemplate(p.Binary)
	if err != nil {
		return err
	}
	return applyTemplate(path, template, exe)
}

func (p Provider) Remove(exe string) error {
	path, err := p.ConfigPath()
	if err != nil {
		return err
	}
	if !fileExists(path) {
		return nil
	}
	document, err := readJSONMap(path)
	if err != nil {
		return err
	}
	documentHooks, _ := document["hooks"].(map[string]any)
	if documentHooks == nil {
		return nil
	}
	changed := false
	for event, groups := range documentHooks {
		eventGroups, _ := groups.([]any)
		kept := removeTraceknotGroups(eventGroups, exe)
		if len(kept) != len(eventGroups) {
			documentHooks[event] = kept
			changed = true
		}
	}
	if !changed {
		return nil
	}
	return writeJSON(path, document)
}

func (p Provider) Installed(exe string) bool {
	path, err := p.ConfigPath()
	if err != nil {
		return false
	}
	document, err := readJSONMap(path)
	if err != nil {
		return false
	}
	hooks, _ := document["hooks"].(map[string]any)
	if hooks == nil {
		return false
	}
	for _, groupsValue := range hooks {
		groups, _ := groupsValue.([]any)
		if len(removeTraceknotGroups(groups, exe)) != len(groups) {
			return true
		}
	}
	return false
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
