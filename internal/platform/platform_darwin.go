//go:build darwin

package platform

import (
	"context"
	"fmt"
	"os"

	"traceknot/internal/platform/darwin"
)

type darwinPlatform struct{ unsupported }

var Current Platform = darwinPlatform{}

func (darwinPlatform) SpawnDaemon(args []string, log *os.File) (int, error) {
	return darwin.SpawnDaemon(args, log)
}

func (darwinPlatform) KillProcess(ctx context.Context, pid string) {
	killProcess(ctx, pid)
}

func (darwinPlatform) FreePort(ctx context.Context, port string) {
	freePort(ctx, port)
}

func (darwinPlatform) AutostartEnabled(ctx context.Context) bool {
	return darwin.AutostartEnabled(ctx)
}

func (darwinPlatform) AutostartEnable(ctx context.Context) error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("resolve binary: %w", err)
	}
	return darwin.AutostartEnable(ctx, exe)
}

func (darwinPlatform) AutostartDisable(ctx context.Context) error {
	return darwin.AutostartDisable(ctx)
}

func (darwinPlatform) StartManagedDaemon(ctx context.Context) string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	return darwin.StartManaged(ctx, exe)
}

func (darwinPlatform) StopManagedDaemon(ctx context.Context) bool {
	return darwin.StopManaged(ctx)
}

func (darwinPlatform) SetLaunchctlEnv(vars map[string]string) {
	darwin.SetLaunchctlEnv(vars)
}

func (darwinPlatform) UnsetLaunchctlEnv(names []string) {
	darwin.UnsetLaunchctlEnv(names)
}

func (darwinPlatform) CloseInteractiveSessions(ctx context.Context) (int, error) {
	closed, vscodeClosed := closeCmdlineSessions(ctx)
	if vscodeClosed {
		darwin.RelaunchVSCode(ctx)
	}
	return closed, nil
}

func (darwinPlatform) CanOpenBrowser() bool { return darwin.CanOpenBrowser() }

func (darwinPlatform) OpenBrowser(target string) error {
	return darwin.OpenBrowser(target)
}

func (darwinPlatform) OpenDefaultHandler(target string) error {
	return darwin.OpenDefault(target)
}

func (darwinPlatform) RemoveBinary(exe string) error {
	return darwin.RemoveBinary(exe)
}
