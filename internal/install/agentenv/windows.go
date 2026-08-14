package agentenv

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func applyWindowsEnv(binDir string) error {
	envBlock, err := loadSnippet("windows/env.ps1")
	if err != nil {
		return err
	}
	var statements []string
	for _, pair := range windowsEnvPairs(envBlock) {
		statements = append(statements,
			"Set-ItemProperty -Path 'HKCU:\\Environment' -Name '"+psQuote(pair.name)+
				"' -Value '"+psQuote(pair.value)+"' -Type ExpandString;")
	}
	statements = append(statements,
		"$bin = '"+psQuote(binDir)+"';"+
			"$path = [Environment]::GetEnvironmentVariable('Path','User');"+
			"if (-not $path) { $path = '' };"+
			"if ($path -notlike \"*$bin*\") {"+
			"if ($path) { $path = $path.TrimEnd(';') + ';' };"+
			"$path += $bin;"+
			"Set-ItemProperty -Path 'HKCU:\\Environment' -Name 'Path' -Value $path -Type ExpandString;}")
	return runPowerShell(statements)
}

func removeWindowsEnv() error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("resolve binary: %w", err)
	}
	binDir := filepath.Dir(exe)
	envBlock, err := loadSnippet("windows/env.ps1")
	if err != nil {
		return err
	}
	var statements []string
	statements = append(statements,
		"$bin = '"+psQuote(binDir)+"';"+
			"$path = [Environment]::GetEnvironmentVariable('Path', 'User');"+
			"if ($path) {"+
			"$kept = @();"+
			"foreach ($entry in $path.Split(';')) {"+
			"if ($entry.Trim() -and -not [string]::Equals($entry.Trim(), $bin, [StringComparison]::OrdinalIgnoreCase)) { $kept += $entry }"+
			"};"+
			"$newPath = $kept -join ';';"+
			"if (-not [string]::Equals($newPath, $path.Trim(), [StringComparison]::OrdinalIgnoreCase)) {"+
			"Set-ItemProperty -Path 'HKCU:\\Environment' -Name 'Path' -Value $newPath -Type ExpandString"+
			"}}")
	for _, pair := range windowsEnvPairs(envBlock) {
		statements = append(statements,
			"Remove-ItemProperty -Path 'HKCU:\\Environment' -Name '"+psQuote(pair.name)+"' -ErrorAction SilentlyContinue;")
	}
	return runPowerShell(statements)
}

func psQuote(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}

func runPowerShell(statements []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return exec.CommandContext(ctx, "powershell", "-NoProfile", "-Command", strings.Join(statements, " ")).Run()
}
