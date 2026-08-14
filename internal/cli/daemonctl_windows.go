//go:build windows

package cli

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func startDaemonBackground(args []string) error {
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
	if err := os.WriteFile(daemonPidPath(), []byte(fmt.Sprint(command.Process.Pid)), 0o644); err != nil {
		return fmt.Errorf("write pid file: %w", err)
	}
	return nil
}
