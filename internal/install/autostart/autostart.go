package autostart

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func Enabled(ctx context.Context) bool {
	switch runtime.GOOS {
	case "linux":
		return systemdUnitExists()
	case "darwin":
		return launchAgentExists()
	case "windows":
		return runKeyExists(ctx)
	}
	return false
}

func Enable(ctx context.Context) error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("resolve binary: %w", err)
	}
	switch runtime.GOOS {
	case "linux":
		return writeSystemdUnit(ctx, exe)
	case "darwin":
		return writeLaunchAgent(exe)
	case "windows":
		return setRunKey(ctx, exe)
	}
	return nil
}

func Disable(ctx context.Context) error {
	switch runtime.GOOS {
	case "linux":
		if !systemdUnitExists() {
			return nil
		}
		runQuiet(ctx, "systemctl", "--user", "disable", "traceknot.service")
		if err := os.Remove(systemdUnitPath()); err != nil {
			return fmt.Errorf("remove unit: %w", err)
		}
		runQuiet(ctx, "systemctl", "--user", "daemon-reload")
	case "darwin":
		if !launchAgentExists() {
			return nil
		}
		if err := os.Remove(launchAgentPath()); err != nil {
			return fmt.Errorf("remove plist: %w", err)
		}
	case "windows":
		if !runKeyExists(ctx) {
			return nil
		}
		_ = exec.CommandContext(ctx, "reg", "delete", `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`, "/v", "traceknot", "/f").Run()
	}
	return nil
}

func SystemdUnitExists() bool {
	return runtime.GOOS == "linux" && systemdUnitExists()
}

func LaunchAgentExists() bool {
	return runtime.GOOS == "darwin" && launchAgentExists()
}

func LaunchAgentPath() string {
	return launchAgentPath()
}

func RepairLaunchAgent(exe string) error {
	if runtime.GOOS != "darwin" || !launchAgentExists() {
		return nil
	}
	return writeLaunchAgent(exe)
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

func runKeyPath() string {
	return `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`
}

func runKeyExists(ctx context.Context) bool {
	command := exec.CommandContext(ctx, "reg", "query", runKeyPath(), "/v", "traceknot")
	return command.Run() == nil
}

func setRunKey(ctx context.Context, exe string) error {
	command := exec.CommandContext(ctx, "reg", "add", runKeyPath(), "/v", "traceknot", "/t", "REG_SZ", "/d", "\""+exe+"\"", "/f")
	if err := command.Run(); err != nil {
		return fmt.Errorf("set run key: %w", err)
	}
	return nil
}

func runQuiet(ctx context.Context, name string, args ...string) {
	_ = exec.CommandContext(ctx, name, args...).Run()
}

func home() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return dir
}

func daemonLogPath() string {
	return filepath.Join(home(), ".traceknot", "daemon.log")
}
