//go:build windows

package windows

import (
	"os"
	"os/exec"
	"syscall"
)

func SpawnDaemon(args []string, log *os.File) error {
	command := exec.Command(args[0], args[1:]...)
	command.Stdout = log
	command.Stderr = log
	command.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: createNoWindow | createNewProcessGroup,
	}
	return command.Start()
}
