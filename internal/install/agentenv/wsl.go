package agentenv

import (
	"path/filepath"
)

func vscodeServerEnvSetupPath() string {
	return filepath.Join(homeDir(), ".vscode-server", "server-env-setup")
}
