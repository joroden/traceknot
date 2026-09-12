//go:build windows

package platform

import (
	"context"
	"os"

	"traceknot/internal/platform/windows"
)

type windowsPlatform struct{ unsupported }

var Current Platform = windowsPlatform{}

func (windowsPlatform) SpawnDaemon(args []string, log *os.File) (int, error) {
	return 0, windows.SpawnDaemon(args, log)
}

func (windowsPlatform) KillProcess(ctx context.Context, pid string) {
	killProcess(ctx, pid)
}

func (windowsPlatform) FreePort(ctx context.Context, port string) {
	freePort(ctx, port)
}

func (windowsPlatform) AutostartEnabled(ctx context.Context) bool {
	return windows.AutostartEnabled(ctx)
}

func (windowsPlatform) AutostartEnable(ctx context.Context) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	return windows.AutostartEnable(ctx, exe)
}

func (windowsPlatform) AutostartDisable(ctx context.Context) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	return windows.AutostartDisable(ctx, exe)
}

func (windowsPlatform) ApplyEnv(vars map[string]string, binDir string) error {
	return windows.ApplyEnv(vars, binDir)
}

func (windowsPlatform) RemoveEnv(names []string, binDir string) error {
	return windows.RemoveEnv(names, binDir)
}

func (windowsPlatform) CloseInteractiveSessions(ctx context.Context) (int, error) {
	return windows.CloseInteractiveSessions(ctx)
}

func (windowsPlatform) CanOpenBrowser() bool { return true }

func (windowsPlatform) OpenBrowser(target string) error {
	return windows.OpenBrowser(target)
}

func (windowsPlatform) OpenDefaultHandler(target string) error {
	return windows.OpenDefault(target)
}

func (windowsPlatform) RemoveBinary(exe string) error {
	return windows.RemoveBinary(exe)
}
