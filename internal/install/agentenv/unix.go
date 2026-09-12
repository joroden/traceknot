//go:build !windows

package agentenv

import (
	"runtime"
	"strings"

	"traceknot/internal/platform"
)

func ApplyEnv(binDir string) error {
	envBlock, err := loadSnippet("unix/env.sh")
	if err != nil {
		return err
	}
	pathBlock, err := loadSnippet("unix/path.sh")
	if err != nil {
		return err
	}
	renderedPath := strings.ReplaceAll(string(pathBlock), "{{BIN_DIR}}", binDir)
	for _, profile := range profileCandidates() {
		if err := upsertBlock(profile, string(envBlock)); err != nil {
			return err
		}
		if err := upsertBlock(profile, renderedPath); err != nil {
			return err
		}
	}
	if platform.Current.IsWSL() {
		vscodeServerEnvSetup := vscodeServerEnvSetupPath()
		if err := upsertBlock(vscodeServerEnvSetup, string(envBlock)); err != nil {
			return err
		}
		if err := upsertBlock(vscodeServerEnvSetup, renderedPath); err != nil {
			return err
		}
	}
	if runtime.GOOS == "darwin" {
		platform.Current.SetLaunchctlEnv(pairMap(unixEnvPairs(envBlock)))
	}
	return nil
}

func RemoveEnv() error {
	envBlock, err := loadSnippet("unix/env.sh")
	if err != nil {
		return err
	}
	pathBlock, err := loadSnippet("unix/path.sh")
	if err != nil {
		return err
	}
	envStart, envEnd := blockMarkers(envBlock)
	pathStart, pathEnd := blockMarkers(pathBlock)
	for _, profile := range profileCandidates() {
		if err := removeBlock(profile, envStart, envEnd); err != nil {
			return err
		}
		if err := removeBlock(profile, pathStart, pathEnd); err != nil {
			return err
		}
	}
	if platform.Current.IsWSL() {
		vscodeServerEnvSetup := vscodeServerEnvSetupPath()
		if err := removeBlock(vscodeServerEnvSetup, envStart, envEnd); err != nil {
			return err
		}
		if err := removeBlock(vscodeServerEnvSetup, pathStart, pathEnd); err != nil {
			return err
		}
	}
	if runtime.GOOS == "darwin" {
		platform.Current.UnsetLaunchctlEnv(pairNames(unixEnvPairs(envBlock)))
	}
	return nil
}
