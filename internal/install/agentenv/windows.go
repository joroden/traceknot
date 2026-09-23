//go:build windows

package agentenv

import (
	"encoding/json"
	"errors"
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
			fmt.Println(host.label + ": enabled profile scripts for the Copilot integration; the previous execution policy will be restored on uninstall.")
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
	effective, err := getExecutionPolicy(shell, "Get-ExecutionPolicy")
	if err != nil {
		return false, fmt.Errorf("check policy: %w", err)
	}
	if effective != "Restricted" {
		return false, nil
	}
	previous, err := getExecutionPolicy(shell, "Get-ExecutionPolicy -Scope CurrentUser")
	if err != nil {
		return false, fmt.Errorf("check current-user policy: %w", err)
	}
	if !validExecutionPolicy(previous) || previous == "RemoteSigned" {
		return false, fmt.Errorf("cannot safely change current-user policy %q", previous)
	}
	state, err := loadPolicyState()
	if err != nil {
		return false, err
	}
	if _, recorded := state[shell]; recorded {
		return false, fmt.Errorf("previous policy is already recorded; leaving current policy unchanged")
	}
	state[shell] = previous
	if err := savePolicyState(state); err != nil {
		return false, err
	}
	if err := setExecutionPolicy(shell, "RemoteSigned"); err != nil {
		return false, err
	}
	current, err := getExecutionPolicy(shell, "Get-ExecutionPolicy -Scope CurrentUser")
	if err != nil {
		return false, fmt.Errorf("verify current-user policy: %w", err)
	}
	if current != "RemoteSigned" {
		return false, fmt.Errorf("current-user policy remained %q after update", current)
	}
	return true, nil
}

func policyStatePath() string {
	return filepath.Join(homeDir(), ".traceknot", "powershell-execution-policy.json")
}

func loadPolicyState() (map[string]string, error) {
	data, err := os.ReadFile(policyStatePath())
	if os.IsNotExist(err) {
		return map[string]string{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read saved execution policy: %w", err)
	}
	state := map[string]string{}
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("decode saved execution policy: %w", err)
	}
	if state == nil {
		return nil, fmt.Errorf("saved execution policy has no entries")
	}
	return state, nil
}

func savePolicyState(state map[string]string) error {
	path := policyStatePath()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create policy state directory: %w", err)
	}
	data, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("encode saved execution policy: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("save execution policy: %w", err)
	}
	return nil
}

func getExecutionPolicy(shell string, command string) (string, error) {
	output, err := exec.Command(shell, "-NoProfile", "-NonInteractive", "-Command", command).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s: %s", command, strings.TrimSpace(string(output)))
	}
	return strings.TrimSpace(string(output)), nil
}

func validExecutionPolicy(value string) bool {
	switch value {
	case "Undefined", "Restricted", "AllSigned", "RemoteSigned", "Unrestricted", "Bypass":
		return true
	default:
		return false
	}
}

func setExecutionPolicy(shell string, value string) error {
	if !validExecutionPolicy(value) {
		return fmt.Errorf("invalid execution policy %q", value)
	}
	_, err := getExecutionPolicy(shell, "Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy "+value+" -Force")
	return err
}

func restoreProfileScriptPolicies() error {
	state, err := loadPolicyState()
	if err != nil {
		return err
	}
	var failures []error
	for shell, previous := range state {
		if !validExecutionPolicy(previous) {
			failures = append(failures, fmt.Errorf("%s: invalid saved policy %q", shell, previous))
			continue
		}
		if _, err := exec.LookPath(shell); err != nil {
			failures = append(failures, fmt.Errorf("%s: shell unavailable for policy restoration", shell))
			continue
		}
		current, err := getExecutionPolicy(shell, "Get-ExecutionPolicy -Scope CurrentUser")
		if err != nil {
			failures = append(failures, fmt.Errorf("%s: %w", shell, err))
			continue
		}
		if current == "RemoteSigned" {
			if err := setExecutionPolicy(shell, previous); err != nil {
				failures = append(failures, fmt.Errorf("%s: restore policy: %w", shell, err))
				continue
			}
		}
		delete(state, shell)
	}
	if len(state) == 0 {
		if err := os.Remove(policyStatePath()); err != nil && !os.IsNotExist(err) {
			failures = append(failures, fmt.Errorf("remove saved execution policy: %w", err))
		}
	} else if err := savePolicyState(state); err != nil {
		failures = append(failures, err)
	}
	return errors.Join(failures...)
}

func RemoveCopilotShim() error {
	shimBlock, err := loadSnippet("windows/copilot-shim.ps1")
	if err != nil {
		return err
	}
	start, end := blockMarkers(shimBlock)
	var failures []error
	for _, profile := range powershellProfileCandidates() {
		if err := removeBlock(profile, start, end); err != nil {
			failures = append(failures, err)
		}
	}
	if err := restoreProfileScriptPolicies(); err != nil {
		failures = append(failures, err)
	}
	return errors.Join(failures...)
}

func powershellProfileCandidates() []string {
	home := homeDir()
	return []string{
		filepath.Join(home, "Documents", "PowerShell", "profile.ps1"),
		filepath.Join(home, "Documents", "WindowsPowerShell", "profile.ps1"),
	}
}
