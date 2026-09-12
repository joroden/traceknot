//go:build linux

package linux

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func AutostartEnabled(_ context.Context) bool {
	return systemdUnitExists()
}

func AutostartEnable(ctx context.Context, exe string) error {
	return writeSystemdUnit(ctx, exe)
}

func AutostartDisable(ctx context.Context) error {
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

func StartManaged(ctx context.Context) string {
	if !systemdUnitExists() {
		return ""
	}
	runQuiet(ctx, "systemctl", "--user", "start", "traceknot.service")
	return "systemd user service"
}

func StopManaged(ctx context.Context) bool {
	if !systemdUnitExists() {
		return false
	}
	runQuiet(ctx, "systemctl", "--user", "stop", "traceknot.service")
	return true
}

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
		"ExecStart=" + strings.ReplaceAll(exe, "%", "%%") + "\n" +
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

func runQuiet(ctx context.Context, name string, args ...string) {
	_ = exec.CommandContext(ctx, name, args...).Run()
}
