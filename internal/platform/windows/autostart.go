//go:build windows

package windows

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const scheduledTaskName = "traceknot"

func AutostartEnabled(ctx context.Context) bool {
	command := exec.CommandContext(ctx, "schtasks", "/query", "/tn", scheduledTaskName)
	hidden(command)
	return command.Run() == nil
}

func AutostartEnable(ctx context.Context, exe string) error {
	launcher := launcherPath(exe)
	if err := os.WriteFile(launcher, []byte(launcherScript(exe)), 0o644); err != nil {
		return fmt.Errorf("write launcher: %w", err)
	}
	command := exec.CommandContext(ctx, "schtasks", "/create", "/tn", scheduledTaskName,
		"/tr", "wscript.exe //B \""+launcher+"\"", "/sc", "onlogon", "/rl", "limited", "/f")
	hidden(command)
	if err := command.Run(); err != nil {
		return fmt.Errorf("create scheduled task: %w", err)
	}
	return nil
}

func AutostartDisable(ctx context.Context, exe string) error {
	if !AutostartEnabled(ctx) {
		return nil
	}
	command := exec.CommandContext(ctx, "schtasks", "/delete", "/tn", scheduledTaskName, "/f")
	hidden(command)
	_ = command.Run()
	_ = os.Remove(launcherPath(exe))
	return nil
}

func launcherPath(exe string) string {
	return filepath.Join(filepath.Dir(exe), "autostart-launch.vbs")
}

func launcherScript(exe string) string {
	return "CreateObject(\"WScript.Shell\").Run Chr(34) & \"" + exe + "\" & Chr(34) & \" daemon\", 0, False\n"
}
