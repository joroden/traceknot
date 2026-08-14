package cli

import (
	"context"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func mustHome() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return home
}

func resolveExe() string {
	exe, err := os.Executable()
	if err != nil {
		return "traceknot"
	}
	if resolved, resolveErr := filepath.Abs(exe); resolveErr == nil {
		exe = resolved
	}
	return exe
}

func isTTY(file *os.File) bool {
	info, err := file.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

func IsInteractive() bool {
	return isTTY(os.Stdin) && isTTY(os.Stdout)
}

func daemonHealthy(ctx context.Context, server string) bool {
	client := &http.Client{Timeout: 2 * time.Second}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(server, "/")+"/healthz", nil)
	if err != nil {
		return false
	}
	response, err := client.Do(request)
	if err != nil {
		return false
	}
	defer response.Body.Close()
	return response.StatusCode == http.StatusOK
}

func runQuiet(ctx context.Context, name string, args ...string) {
	_ = exec.CommandContext(ctx, name, args...).Run()
}
