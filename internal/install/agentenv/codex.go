package agentenv

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func ApplyCodex() error {
	codexHome := os.Getenv("CODEX_HOME")
	if codexHome == "" {
		codexHome = filepath.Join(homeDir(), ".codex")
	}
	if !codexInstalled(codexHome) {
		return nil
	}
	block, err := loadSnippet("codex/config.toml")
	if err != nil {
		return err
	}
	if err := upsertBlock(filepath.Join(codexHome, "config.toml"), string(block)); err != nil {
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
	path := filepath.Join(homeDir(), ".codex", "config.toml")
	if !fileExists(path) {
		return nil
	}
	return removeBlock(path, start, end)
}

func codexInstalled(codexHome string) bool {
	if _, err := exec.LookPath("codex"); err == nil {
		return true
	}
	return dirExists(codexHome)
}
