//go:build linux

package linux

import (
	"context"
	"os/exec"
)

func RelaunchVSCode(ctx context.Context) {
	for _, bin := range []string{"code", "code-insiders"} {
		path, err := exec.LookPath(bin)
		if err != nil {
			continue
		}
		_ = exec.CommandContext(ctx, path).Start()
		return
	}
}
