package agentenv

import (
	"os"
	"path/filepath"
	"strings"
)

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

func vscodeServerEnvSetupPath() string {
	return filepath.Join(homeDir(), ".vscode-server", "server-env-setup")
}
