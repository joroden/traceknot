//go:build linux

package linux

import (
	"os"
	"strings"
)

func home() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return dir
}

func IsWSL() bool {
	if os.Getenv("WSL_DISTRO_NAME") != "" || os.Getenv("WSL_INTEROP") != "" {
		return true
	}
	version, err := os.ReadFile("/proc/version")
	if err != nil {
		return false
	}
	lower := strings.ToLower(string(version))
	return strings.Contains(lower, "microsoft") || strings.Contains(lower, "wsl")
}
