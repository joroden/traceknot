//go:build linux

package linux

import (
	"context"
	"os"
	"os/exec"
	"strings"
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

func FreePortOnWindowsHost(ctx context.Context, port string) {
	if !IsWSL() {
		return
	}
	output, err := exec.CommandContext(ctx, "cmd.exe", "/c", "netstat", "-ano").Output()
	if err != nil {
		return
	}
	killed := make(map[string]bool)
	for line := range strings.SplitSeq(string(output), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 5 || fields[0] != "TCP" || fields[3] != "LISTENING" {
			continue
		}
		if !strings.HasSuffix(fields[1], ":"+port) {
			continue
		}
		pid := fields[4]
		if killed[pid] {
			continue
		}
		killed[pid] = true
		runQuiet(ctx, "taskkill.exe", "/F", "/PID", pid)
	}
}
