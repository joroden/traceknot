//go:build linux

package autostart

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

func Enabled(_ context.Context) bool {
	return systemdUnitExists()
}

func Enable(ctx context.Context) error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("resolve binary: %w", err)
	}
	return writeSystemdUnit(ctx, exe)
}

func Disable(ctx context.Context) error {
	if !systemdUnitExists() {
		return nil
	}
	runQuiet(ctx, "systemctl", "--user", "disable", "traceknot.service")
	if err := os.Remove(systemdUnitPath()); err != nil {
		return fmt.Errorf("remove unit: %w", err)
	}
	runQuiet(ctx, "systemctl", "--user", "daemon-reload")
	return nil
}

func SystemdUnitExists() bool { return systemdUnitExists() }

func LaunchAgentExists() bool { return false }

func LaunchAgentPath() string { return "" }

func RepairLaunchAgent(_ string) error { return nil }

func systemdUnitExists() bool {
	_, err := os.Stat(systemdUnitPath())
	return err == nil
}

func systemdUnitPath() string {
	return filepath.Join(home(), ".config", "systemd", "user", "traceknot.service")
}

func writeSystemdUnit(ctx context.Context, exe string) error {
	path := systemdUnitPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("mkdir unit dir: %w", err)
	}
	unit := "[Unit]\n" +
		"Description=traceknot telemetry daemon\n" +
		"After=network.target\n\n" +
		"[Service]\n" +
		"Type=simple\n" +
		"ExecStart=" + exe + "\n" +
		"Restart=always\n" +
		"RestartSec=2\n\n" +
		"[Install]\n" +
		"WantedBy=default.target\n"
	if err := os.WriteFile(path, []byte(unit), 0o644); err != nil {
		return fmt.Errorf("write unit: %w", err)
	}
	runQuiet(ctx, "systemctl", "--user", "daemon-reload")
	runQuiet(ctx, "systemctl", "--user", "enable", "traceknot.service")
	return nil
}
