//go:build linux

package platform

import (
	"context"
	"fmt"
	"os"

	"traceknot/internal/platform/linux"
)

type linuxPlatform struct{ unsupported }

var Current Platform = linuxPlatform{}

func (linuxPlatform) SpawnDaemon(args []string, log *os.File) (int, error) {
	return linux.SpawnDaemon(args, log)
}

func (linuxPlatform) KillProcess(ctx context.Context, pid string) {
	killProcess(ctx, pid)
}

func (linuxPlatform) FreePort(ctx context.Context, port string) {
	freePort(ctx, port)
}

func (linuxPlatform) AutostartEnabled(ctx context.Context) bool {
	return linux.AutostartEnabled(ctx)
}

func (linuxPlatform) AutostartEnable(ctx context.Context) error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("resolve binary: %w", err)
	}
	return linux.AutostartEnable(ctx, exe)
}

func (linuxPlatform) AutostartDisable(ctx context.Context) error {
	return linux.AutostartDisable(ctx)
}

func (linuxPlatform) StartManagedDaemon(ctx context.Context) string {
	return linux.StartManaged(ctx)
}

func (linuxPlatform) StopManagedDaemon(ctx context.Context) bool {
	return linux.StopManaged(ctx)
}

func (linuxPlatform) CloseInteractiveSessions(ctx context.Context) (int, error) {
	closed, vscodeClosed := closeCmdlineSessions(ctx)
	if vscodeClosed && !linux.IsWSL() {
		linux.RelaunchVSCode(ctx)
	}
	return closed, nil
}

func (linuxPlatform) CanOpenBrowser() bool { return linux.CanOpenBrowser() }

func (linuxPlatform) OpenBrowser(target string) error {
	return linux.OpenBrowser(target)
}

func (linuxPlatform) OpenDefaultHandler(target string) error {
	return linux.OpenDefault(target)
}

func (linuxPlatform) RemoveBinary(exe string) error {
	return linux.RemoveBinary(exe)
}

func (linuxPlatform) IsWSL() bool { return linux.IsWSL() }
