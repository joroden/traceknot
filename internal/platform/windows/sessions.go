//go:build windows

package windows

import (
	"context"
	"os/exec"
	"strings"

	"github.com/shirou/gopsutil/v4/process"
)

var vscodeProcessNames = []string{"Code", "Code - Insiders", "Code - OSS", "Codium"}

func CloseInteractiveSessions(ctx context.Context) (int, error) {
	procs, err := process.ProcessesWithContext(ctx)
	if err != nil {
		return 0, err
	}
	closed := 0
	vscodeClosed := false
	for _, proc := range procs {
		name, err := proc.NameWithContext(ctx)
		if err == nil && matchesVSCodeName(name) {
			_ = proc.KillWithContext(ctx)
			closed++
			vscodeClosed = true
			continue
		}
		cmdline, err := proc.CmdlineWithContext(ctx)
		if err == nil && strings.Contains(strings.ToLower(cmdline), "copilot") {
			_ = proc.KillWithContext(ctx)
			closed++
		}
	}
	if vscodeClosed {
		relaunchVSCode(ctx)
	}
	return closed, nil
}

func matchesVSCodeName(name string) bool {
	for _, candidate := range vscodeProcessNames {
		if strings.EqualFold(name, candidate) {
			return true
		}
	}
	return false
}

func relaunchVSCode(ctx context.Context) {
	path, err := exec.LookPath("code")
	if err != nil {
		return
	}
	command := exec.CommandContext(ctx, path)
	hidden(command)
	_ = command.Start()
}
