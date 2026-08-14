package agentenv

import (
	"context"
	"os/exec"
	"runtime"
	"strings"
)

func applyUnixEnv(binDir string) error {
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
	if IsWSL() {
		vscodeServerEnvSetup := vscodeServerEnvSetupPath()
		if err := upsertBlock(vscodeServerEnvSetup, string(envBlock)); err != nil {
			return err
		}
		if err := upsertBlock(vscodeServerEnvSetup, renderedPath); err != nil {
			return err
		}
	}
	if runtime.GOOS == "darwin" {
		setLaunchctlEnv(envBlock)
	}
	return nil
}

func removeUnixEnv() error {
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
	if IsWSL() {
		vscodeServerEnvSetup := vscodeServerEnvSetupPath()
		if err := removeBlock(vscodeServerEnvSetup, envStart, envEnd); err != nil {
			return err
		}
		if err := removeBlock(vscodeServerEnvSetup, pathStart, pathEnd); err != nil {
			return err
		}
	}
	if runtime.GOOS == "darwin" {
		unsetLaunchctlEnv(envBlock)
	}
	return nil
}

func setLaunchctlEnv(block []byte) {
	ctx := context.Background()
	for _, pair := range unixEnvPairs(block) {
		_ = exec.CommandContext(ctx, "launchctl", "setenv", pair.name, pair.value).Run()
	}
}

func unsetLaunchctlEnv(block []byte) {
	ctx := context.Background()
	for _, pair := range unixEnvPairs(block) {
		_ = exec.CommandContext(ctx, "launchctl", "unsetenv", pair.name).Run()
	}
}
