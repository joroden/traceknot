package cli

import (
	"context"
	"runtime"

	"traceknot/internal/install/autostart"
)

func autostartStatus(ctx context.Context) (bool, string) {
	switch runtime.GOOS {
	case "linux":
		if autostart.SystemdUnitExists() {
			return true, "Autostart: enabled (systemd user unit)"
		}
		return false, "Autostart: disabled"
	case "darwin":
		if autostart.LaunchAgentExists() {
			return true, "Autostart: enabled (LaunchAgent)"
		}
		return false, "Autostart: disabled"
	case "windows":
		if autostart.Enabled(ctx) {
			return true, "Autostart: enabled (HKCU Run key)"
		}
		return false, "Autostart: disabled"
	default:
		return false, "Autostart: unsupported"
	}
}
