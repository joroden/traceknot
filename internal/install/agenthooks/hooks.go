package agenthooks

import "os"

type Provider struct {
	Binary    string
	path      func() (string, error)
	install   func(path string, exe string) error
	remove    func(path string, exe string) error
	installed func(path string, exe string) bool
}

var Providers = []Provider{
	{Binary: "claude", path: claudeSettingsPath, install: claudeInstall, remove: claudeRemove, installed: jsonHookInstalled},
	{Binary: "codex", path: codexHooksPath, install: codexInstall, remove: codexRemove, installed: jsonHookInstalled},
	{Binary: "copilot", path: copilotExtensionPath, install: copilotInstall, remove: copilotRemove, installed: extensionInstalled},
}

func (p Provider) ConfigPath() (string, error) {
	return p.path()
}

func (p Provider) Install(exe string) error {
	path, err := p.ConfigPath()
	if err != nil {
		return err
	}
	return p.install(path, exe)
}

func (p Provider) Remove(exe string) error {
	path, err := p.ConfigPath()
	if err != nil {
		return err
	}
	if !fileExists(path) {
		return nil
	}
	return p.remove(path, exe)
}

func (p Provider) Installed(exe string) bool {
	path, err := p.ConfigPath()
	if err != nil {
		return false
	}
	return p.installed(path, exe)
}

func claudeInstall(path, exe string) error {
	return jsonHookInstall("claude")(path, exe)
}

func claudeRemove(path, exe string) error {
	return jsonHookRemove(path, exe)
}

func codexInstall(path, exe string) error {
	return jsonHookInstall("codex")(path, exe)
}

func codexRemove(path, exe string) error {
	return jsonHookRemove(path, exe)
}

func copilotInstall(path, exe string) error {
	if err := extensionInstall(path, exe); err != nil {
		return err
	}
	if err := enableCopilotExperimental(); err != nil {
		return err
	}
	hooksPath, err := copilotHooksPath()
	if err != nil {
		return err
	}
	return jsonHookInstall("copilot")(hooksPath, exe)
}

func copilotRemove(path, exe string) error {
	if err := extensionRemove(path, exe); err != nil {
		return err
	}
	hooksPath, err := copilotHooksPath()
	if err != nil {
		return err
	}
	if !fileExists(hooksPath) {
		return nil
	}
	return jsonHookRemove(hooksPath, exe)
}

func jsonHookInstall(vendor string) func(path string, exe string) error {
	return func(path string, exe string) error {
		template, err := HookTemplate(vendor)
		if err != nil {
			return err
		}
		return applyTemplate(path, template, exe)
	}
}

func jsonHookRemove(path string, exe string) error {
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

func jsonHookInstalled(path string, exe string) bool {
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
