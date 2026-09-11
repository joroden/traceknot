//go:build !darwin

package autostart

func launchAgentExists() bool { return false }

func launchAgentPath() string { return "" }

func writeLaunchAgent(_ string) error { return nil }
