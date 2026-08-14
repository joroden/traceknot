package cli

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"traceknot/internal/install/autostart"
)

const defaultServerURL = "http://127.0.0.1:4318"

func daemonRunning(ctx context.Context) bool {
	return daemonHealthy(ctx, defaultServerURL)
}

func daemonStatusLine(ctx context.Context) string {
	if !daemonHealthy(ctx, defaultServerURL) {
		return "Server: not running"
	}
	mode := "background"
	switch {
	case autostart.SystemdUnitExists():
		mode = "systemd user service"
	case autostart.LaunchAgentExists():
		mode = "LaunchAgent"
	}
	line := "Server: running at " + defaultServerURL + " (" + mode
	if pid := readPidFile(); pid != "" {
		line += ", pid " + pid
	}
	return line + ")"
}

func startDaemonNow(ctx context.Context) error {
	freePort(ctx, portFromURL(defaultServerURL))
	switch {
	case autostart.SystemdUnitExists():
		runQuiet(ctx, "systemctl", "--user", "start", "traceknot.service")
		if waitHealthy(ctx, defaultServerURL) {
			fmt.Println("daemon started (systemd user service)")
			return nil
		}
	case autostart.LaunchAgentExists():
		plist := autostart.LaunchAgentPath()
		runQuiet(ctx, "launchctl", "bootstrap", "gui/"+fmt.Sprint(os.Getuid()), plist)
		if waitHealthy(ctx, defaultServerURL) {
			fmt.Println("daemon started (LaunchAgent)")
			return nil
		}
	}
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot resolve binary: %w", err)
	}
	if err := startDaemonBackground([]string{exe}); err != nil {
		return err
	}
	if waitHealthy(ctx, defaultServerURL) {
		fmt.Println("daemon started (background, log at " + daemonLogPath() + ")")
		return nil
	}
	return fmt.Errorf("daemon failed to start; check %s", daemonLogPath())
}

func stopDaemon(ctx context.Context, server string) string {
	if !daemonHealthy(ctx, server) {
		if pid := readPidFile(); pid != "" {
			killProcess(ctx, pid)
			_ = os.Remove(daemonPidPath())
			return "stopped"
		}
		return "not_running"
	}
	switch {
	case autostart.SystemdUnitExists():
		runQuiet(ctx, "systemctl", "--user", "stop", "traceknot.service")
		return afterStopStatus(ctx, server)
	case autostart.LaunchAgentExists():
		runQuiet(ctx, "launchctl", "bootout", "gui/"+fmt.Sprint(os.Getuid()), autostart.LaunchAgentPath())
		return afterStopStatus(ctx, server)
	default:
		if pid := readPidFile(); pid != "" {
			killProcess(ctx, pid)
			_ = os.Remove(daemonPidPath())
			return "stopped"
		}
		return "manual"
	}
}

func afterStopStatus(ctx context.Context, server string) string {
	if !daemonHealthy(ctx, server) {
		return "stopped"
	}
	if pid := readPidFile(); pid != "" {
		killProcess(ctx, pid)
		_ = os.Remove(daemonPidPath())
		if !daemonHealthy(ctx, server) {
			return "stopped"
		}
	}
	return "manual"
}

func killProcess(ctx context.Context, pid string) {
	if runtime.GOOS == "windows" {
		runQuiet(ctx, "taskkill", "/PID", pid, "/F")
		return
	}
	runQuiet(ctx, "kill", pid)
}

func daemonLogDir() string {
	return filepath.Join(mustHome(), ".traceknot")
}

func daemonPidPath() string {
	return filepath.Join(daemonLogDir(), "daemon.pid")
}

func daemonLogPath() string {
	return filepath.Join(daemonLogDir(), "daemon.log")
}

func readPidFile() string {
	content, err := os.ReadFile(daemonPidPath())
	if err != nil {
		return ""
	}
	return string(content)
}

func waitHealthy(ctx context.Context, server string) bool {
	for range 10 {
		if daemonHealthy(ctx, server) {
			return true
		}
		time.Sleep(500 * time.Millisecond)
	}
	return false
}

func portFromURL(server string) string {
	parsed, err := url.Parse(server)
	if err != nil {
		return "4318"
	}
	if parsed.Port() != "" {
		return parsed.Port()
	}
	return "4318"
}

func freePort(ctx context.Context, port string) {
	switch runtime.GOOS {
	case "linux":
		runQuiet(ctx, "fuser", "-k", port+"/tcp")
	case "windows":
		script := "Get-NetTCPConnection -LocalPort " + port + " -State Listen -ErrorAction SilentlyContinue | ForEach-Object { Stop-Process -Id $_.OwningProcess -Force -ErrorAction SilentlyContinue }"
		_ = exec.CommandContext(ctx, "powershell", "-NoProfile", "-Command", script).Run()
	}
}
