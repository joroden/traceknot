package platform

import (
	"context"
	"strings"

	"github.com/shirou/gopsutil/v4/process"
)

func closeCmdlineSessions(ctx context.Context) (closed int, vscodeClosed bool) {
	procs, err := process.ProcessesWithContext(ctx)
	if err != nil {
		return 0, false
	}
	for _, proc := range procs {
		cmdline, err := proc.CmdlineWithContext(ctx)
		if err != nil || cmdline == "" {
			continue
		}
		lowered := strings.ToLower(cmdline)
		switch {
		case strings.Contains(lowered, "github-copilot-cli") ||
			strings.Contains(lowered, " copilot ") ||
			strings.Contains(lowered, "/copilot"):
			_ = proc.KillWithContext(ctx)
			closed++
		case strings.Contains(lowered, "visual studio code") ||
			strings.Contains(lowered, "/code") ||
			strings.Contains(lowered, "code-insiders") ||
			strings.Contains(lowered, "/codium"):
			_ = proc.KillWithContext(ctx)
			closed++
			vscodeClosed = true
		}
	}
	return closed, vscodeClosed
}
