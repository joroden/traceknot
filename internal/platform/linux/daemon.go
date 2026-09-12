//go:build linux

package linux

import (
	"os"
	"os/exec"
	"syscall"
)

func SpawnDaemon(args []string, log *os.File) (int, error) {
	command := exec.Command(args[0], args[1:]...)
	command.Stdout = log
	command.Stderr = log
	command.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := command.Start(); err != nil {
		return 0, err
	}
	return command.Process.Pid, nil
}
