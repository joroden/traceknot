//go:build darwin

package darwin

import (
	"context"
	"os/exec"
)

func SetLaunchctlEnv(vars map[string]string) {
	ctx := context.Background()
	for name, value := range vars {
		_ = exec.CommandContext(ctx, "launchctl", "setenv", name, value).Run()
	}
}

func UnsetLaunchctlEnv(names []string) {
	ctx := context.Background()
	for _, name := range names {
		_ = exec.CommandContext(ctx, "launchctl", "unsetenv", name).Run()
	}
}
