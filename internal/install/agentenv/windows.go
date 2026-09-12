//go:build windows

package agentenv

import (
	"fmt"
	"os"
	"path/filepath"

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
