package platform

import "context"

type unsupported struct{}

func (unsupported) StartManagedDaemon(context.Context) string { return "" }

func (unsupported) StopManagedDaemon(context.Context) bool { return false }

func (unsupported) ApplyEnv(map[string]string, string) error { return nil }

func (unsupported) RemoveEnv([]string, string) error { return nil }

func (unsupported) SetLaunchctlEnv(map[string]string) {}

func (unsupported) UnsetLaunchctlEnv([]string) {}

func (unsupported) IsWSL() bool { return false }
