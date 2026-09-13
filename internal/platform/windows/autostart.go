//go:build windows

package windows

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

const startupScriptName = "traceknot.vbs"

func AutostartEnabled(ctx context.Context) bool {
	path, err := startupScriptPath()
	if err != nil {
		return false
	}
	_, err = os.Stat(path)
	return err == nil
}

func AutostartEnable(ctx context.Context, exe string) error {
	path, err := startupScriptPath()
	if err != nil {
		return fmt.Errorf("locate startup folder: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create startup folder: %w", err)
	}
	unblock(exe)
	if err := os.WriteFile(path, []byte(launcherScript(exe)), 0o644); err != nil {
		return fmt.Errorf("write startup script: %w", err)
	}
	return nil
}

func AutostartDisable(ctx context.Context, exe string) error {
	path, err := startupScriptPath()
	if err != nil {
		return fmt.Errorf("locate startup folder: %w", err)
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove startup script: %w", err)
	}
	return nil
}

func startupScriptPath() (string, error) {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		return "", fmt.Errorf("APPDATA is not set")
	}
	return filepath.Join(appData, "Microsoft", "Windows", "Start Menu", "Programs", "Startup", startupScriptName), nil
}

func launcherScript(exe string) string {
	return "CreateObject(\"WScript.Shell\").Run Chr(34) & \"" + exe + "\" & Chr(34) & \" daemon\", 0, False\n"
}

// unblock strips the Zone.Identifier alternate data stream (Windows' Mark of
// the Web) from exe, if present, so launching it doesn't trigger the "Open
// File - Security Warning" prompt. Absence of the stream is not an error.
func unblock(exe string) {
	_ = os.Remove(exe + ":Zone.Identifier")
}
