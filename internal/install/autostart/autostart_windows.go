//go:build windows

package autostart

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

const scheduledTaskName = "traceknot"

func Enabled(ctx context.Context) bool {
	return scheduledTaskExists(ctx)
}

func Enable(ctx context.Context) error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("resolve binary: %w", err)
	}
	return setScheduledTask(ctx, exe)
}

func Disable(ctx context.Context) error {
	if !scheduledTaskExists(ctx) {
		return nil
	}
	deleteScheduledTask(ctx)
	return nil
}

func SystemdUnitExists() bool { return false }

func LaunchAgentExists() bool { return false }

func LaunchAgentPath() string { return "" }

func RepairLaunchAgent(_ string) error { return nil }

func scheduledTaskExists(ctx context.Context) bool {
	command := exec.CommandContext(ctx, "schtasks", "/query", "/tn", scheduledTaskName)
	return command.Run() == nil
}

func setScheduledTask(ctx context.Context, exe string) error {
	command := exec.CommandContext(ctx, "schtasks", "/create", "/tn", scheduledTaskName,
		"/tr", "\""+exe+"\"", "/sc", "onlogon", "/rl", "limited", "/f")
	if err := command.Run(); err != nil {
		return fmt.Errorf("create scheduled task: %w", err)
	}
	return nil
}

func deleteScheduledTask(ctx context.Context) {
	_ = exec.CommandContext(ctx, "schtasks", "/delete", "/tn", scheduledTaskName, "/f").Run()
}
