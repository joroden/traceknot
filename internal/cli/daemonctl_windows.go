//go:build windows

package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func startDaemonBackground(args []string) error {
	ctx := context.Background()
	var lastErr error
	for range 3 {
		freePort(ctx, portFromURL(defaultServerURL))
		time.Sleep(300 * time.Millisecond)
		if err := spawnDaemonProcess(args); err != nil {
			lastErr = err
			continue
		}
		if waitHealthy(ctx, defaultServerURL) {
			return nil
		}
		lastErr = fmt.Errorf("daemon did not become healthy")
	}
	return lastErr
}

func spawnDaemonProcess(args []string) error {
	if err := os.MkdirAll(daemonLogDir(), 0o755); err != nil {
		return fmt.Errorf("create log dir: %w", err)
	}
	logFile, err := os.OpenFile(daemonLogPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open daemon log: %w", err)
	}
	command := exec.Command(args[0], args[1:]...)
	command.Stdout = logFile
	command.Stderr = logFile
	command.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x00000200}
	if err := command.Start(); err != nil {
		return fmt.Errorf("start daemon: %w", err)
	}
	return nil
}

func WriteDaemonPID() func() {
	if err := os.MkdirAll(daemonLogDir(), 0o755); err != nil {
		return func() {}
	}
	if err := os.WriteFile(daemonPidPath(), fmt.Append(nil, os.Getpid()), 0o644); err != nil {
		return func() {}
	}
	return func() { _ = os.Remove(daemonPidPath()) }
}
