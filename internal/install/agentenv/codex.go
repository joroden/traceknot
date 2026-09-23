package agentenv

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func ApplyCodex() error {
	codexHome := filepath.Dir(codexConfigPath())
	if !codexInstalled(codexHome) {
		return nil
	}
	block, err := loadSnippet("codex/config.toml")
	if err != nil {
		return err
	}
	path := codexConfigPath()
	if err := recordCodexConfigPath(path); err != nil {
		return err
	}
	if err := upsertBlock(path, string(block)); err != nil {
		return fmt.Errorf("codex config: %w", err)
	}
	return nil
}

func RemoveCodex() error {
	block, err := loadSnippet("codex/config.toml")
	if err != nil {
		return err
	}
	start, end := blockMarkers(block)
	paths, err := recordedCodexConfigPaths()
	if err != nil {
		return err
	}
	paths = append(paths, codexConfigPath(), filepath.Join(homeDir(), ".codex", "config.toml"))
	seen := map[string]bool{}
	for _, path := range paths {
		if seen[path] || !fileExists(path) {
			continue
		}
		seen[path] = true
		if err := removeBlock(path, start, end); err != nil {
			return err
		}
	}
	if err := os.Remove(codexConfigPathsRecord()); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove saved Codex config paths: %w", err)
	}
	return nil
}

func codexConfigPath() string {
	codexHome := os.Getenv("CODEX_HOME")
	if codexHome == "" {
		codexHome = filepath.Join(homeDir(), ".codex")
	}
	if absolute, err := filepath.Abs(codexHome); err == nil {
		codexHome = absolute
	}
	return filepath.Join(codexHome, "config.toml")
}

func codexConfigPathsRecord() string {
	return filepath.Join(homeDir(), ".traceknot", "codex-config-paths.json")
}

func recordedCodexConfigPaths() ([]string, error) {
	data, err := os.ReadFile(codexConfigPathsRecord())
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read saved Codex config paths: %w", err)
	}
	var paths []string
	if err := json.Unmarshal(data, &paths); err != nil {
		return nil, fmt.Errorf("decode saved Codex config paths: %w", err)
	}
	return paths, nil
}

func recordCodexConfigPath(path string) error {
	paths, err := recordedCodexConfigPaths()
	if err != nil {
		return err
	}
	for _, existing := range paths {
		if existing == path {
			return nil
		}
	}
	paths = append(paths, path)
	data, err := json.Marshal(paths)
	if err != nil {
		return fmt.Errorf("encode Codex config paths: %w", err)
	}
	record := codexConfigPathsRecord()
	if err := os.MkdirAll(filepath.Dir(record), 0o700); err != nil {
		return fmt.Errorf("create Codex path record directory: %w", err)
	}
	if err := os.WriteFile(record, data, 0o600); err != nil {
		return fmt.Errorf("save Codex config paths: %w", err)
	}
	return nil
}

func codexInstalled(codexHome string) bool {
	if _, err := exec.LookPath("codex"); err == nil {
		return true
	}
	return dirExists(codexHome)
}
