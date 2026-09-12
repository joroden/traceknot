//go:build windows

package windows

import (
	"os/exec"
	"syscall"
)

const (
	createNoWindow        = 0x08000000
	createNewProcessGroup = 0x00000200
)

func hidden(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: createNoWindow}
}
