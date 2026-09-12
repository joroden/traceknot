//go:build darwin

package darwin

import (
	"context"
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func AutostartEnabled(_ context.Context) bool {
	return launchAgentExists()
}

func AutostartEnable(_ context.Context, exe string) error {
	return writeLaunchAgent(exe)
}

func AutostartDisable(_ context.Context) error {
	if !launchAgentExists() {
		return nil
	}
	if err := os.Remove(launchAgentPath()); err != nil {
		return fmt.Errorf("remove plist: %w", err)
	}
	return nil
}

func launchAgentExists() bool {
	_, err := os.Stat(launchAgentPath())
	return err == nil
}

func launchAgentPath() string {
	return filepath.Join(home(), "Library", "LaunchAgents", "dev.traceknot.daemon.plist")
}

func StartManaged(ctx context.Context, exe string) string {
	if !launchAgentExists() {
		return ""
	}
	_ = writeLaunchAgent(exe)
	_ = exec.CommandContext(ctx, "launchctl", "bootstrap", fmt.Sprintf("gui/%d", os.Getuid()), launchAgentPath()).Run()
	return "LaunchAgent"
}

func StopManaged(ctx context.Context) bool {
	if !launchAgentExists() {
		return false
	}
	_ = exec.CommandContext(ctx, "launchctl", "bootout", fmt.Sprintf("gui/%d", os.Getuid()), launchAgentPath()).Run()
	return true
}

func writeLaunchAgent(exe string) error {
	path := launchAgentPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("mkdir plist dir: %w", err)
	}
	logPath := daemonLogPath()
	if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err != nil {
		return fmt.Errorf("mkdir log dir: %w", err)
	}
	plist := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
		"<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">\n" +
		"<plist version=\"1.0\">\n" +
		"  <dict>\n" +
		"    <key>Label</key>\n" +
		"    <string>dev.traceknot.daemon</string>\n" +
		"    <key>ProgramArguments</key>\n" +
		"    <array>\n" +
		"      <string>" + xmlEscape(exe) + "</string>\n" +
		"    </array>\n" +
		"    <key>RunAtLoad</key>\n" +
		"    <true/>\n" +
		"    <key>KeepAlive</key>\n" +
		"    <true/>\n" +
		"    <key>StandardOutPath</key>\n" +
		"    <string>" + xmlEscape(logPath) + "</string>\n" +
		"    <key>StandardErrorPath</key>\n" +
		"    <string>" + xmlEscape(logPath) + "</string>\n" +
		"    <key>EnvironmentVariables</key>\n" +
		"    <dict>\n" +
		"      <key>HOME</key>\n" +
		"      <string>" + xmlEscape(home()) + "</string>\n" +
		"      <key>PATH</key>\n" +
		"      <string>" + xmlEscape(launchAgentPathEnv()) + "</string>\n" +
		"    </dict>\n" +
		"  </dict>\n" +
		"</plist>\n"
	if err := os.WriteFile(path, []byte(plist), 0o644); err != nil {
		return fmt.Errorf("write plist: %w", err)
	}
	return nil
}

func xmlEscape(value string) string {
	var buf strings.Builder
	_ = xml.EscapeText(&buf, []byte(value))
	return buf.String()
}

func launchAgentPathEnv() string {
	dirs := []string{
		"/opt/homebrew/bin",
		"/opt/homebrew/sbin",
		"/usr/local/bin",
		"/usr/local/sbin",
		filepath.Join(home(), ".local", "bin"),
		"/usr/bin",
		"/bin",
		"/usr/sbin",
		"/sbin",
	}
	return strings.Join(dirs, ":")
}
