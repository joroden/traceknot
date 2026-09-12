//go:build windows

package windows

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func RemoveBinary(exe string) error {
	stale := exe + ".tk-old"
	if err := os.Rename(exe, stale); err != nil {
		return fmt.Errorf("rename running binary: %w", err)
	}
	script := filepath.Join(os.TempDir(), "traceknot-uninstall.ps1")
	psQuote := func(value string) string { return strings.ReplaceAll(value, "'", "''") }
	content := "Start-Sleep -Milliseconds 500; Remove-Item -Force -LiteralPath '" + psQuote(stale) +
		"'; Remove-Item -Force -LiteralPath '" + psQuote(script) + "'"
	if err := os.WriteFile(script, []byte(content), 0o644); err != nil {
		return nil
	}
	command := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", script)
	hidden(command)
	_ = command.Start()
	return nil
}
