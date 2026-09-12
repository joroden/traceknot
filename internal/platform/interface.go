package platform

import (
	"context"
	"os"
)

type Platform interface {
	SpawnDaemon(args []string, log *os.File) (int, error)
	KillProcess(ctx context.Context, pid string)
	FreePort(ctx context.Context, port string)

	AutostartEnabled(ctx context.Context) bool
	AutostartEnable(ctx context.Context) error
	AutostartDisable(ctx context.Context) error
	StartManagedDaemon(ctx context.Context) string
	StopManagedDaemon(ctx context.Context) bool

	ApplyEnv(vars map[string]string, binDir string) error
	RemoveEnv(names []string, binDir string) error
	SetLaunchctlEnv(vars map[string]string)
	UnsetLaunchctlEnv(names []string)

	CloseInteractiveSessions(ctx context.Context) (int, error)

	CanOpenBrowser() bool
	OpenBrowser(target string) error
	OpenDefaultHandler(target string) error

	RemoveBinary(exe string) error

	IsWSL() bool
}
