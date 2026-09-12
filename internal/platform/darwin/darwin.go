//go:build darwin

package darwin

import (
	"os"
	"path/filepath"
)

func home() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return dir
}

func daemonLogPath() string {
	return filepath.Join(home(), ".traceknot", "daemon.log")
}
