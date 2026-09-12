package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"traceknot/internal/platform"
)

const defaultServerURL = "http://127.0.0.1:4318"

func daemonRunning(ctx context.Context) bool {
	return daemonHealthy(ctx, defaultServerURL)
}

func startDaemonNow(ctx context.Context) error {
	platform.Current.FreePort(ctx, portFromURL(defaultServerURL))
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot resolve binary: %w", err)
	}
	if managed := platform.Current.StartManagedDaemon(ctx); managed != "" {
		if waitHealthy(ctx, defaultServerURL) {
			fmt.Println("daemon started (" + managed + ")")
			waitReady(ctx, defaultServerURL)
			return nil
		}
	}
	if err := startDaemonBackground([]string{exe}); err != nil {
		return err
	}
	if waitHealthy(ctx, defaultServerURL) {
		fmt.Println("daemon started (background, log at " + daemonLogPath() + ")")
		waitReady(ctx, defaultServerURL)
		return nil
	}
	return fmt.Errorf("daemon failed to start; check %s", daemonLogPath())
}

func startDaemonBackground(args []string) error {
	if err := os.MkdirAll(daemonLogDir(), 0o755); err != nil {
		return fmt.Errorf("create log dir: %w", err)
	}
	logFile, err := os.OpenFile(daemonLogPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open daemon log: %w", err)
	}

	if runtime.GOOS != "windows" {
		pid, err := platform.Current.SpawnDaemon(args, logFile)
		if err != nil {
			return fmt.Errorf("start daemon: %w", err)
		}
		return os.WriteFile(daemonPidPath(), []byte(fmt.Sprint(pid)), 0o644)
	}

	ctx := context.Background()
	var lastErr error
	for range 3 {
		platform.Current.FreePort(ctx, portFromURL(defaultServerURL))
		time.Sleep(300 * time.Millisecond)
		if _, err := platform.Current.SpawnDaemon(args, logFile); err != nil {
			lastErr = err
			continue
		}
		if waitHealthy(ctx, defaultServerURL) {
			return nil
		}
		lastErr = fmt.Errorf("daemon did not become healthy")
	}
	return lastErr
}

func WriteDaemonPID() func() {
	if runtime.GOOS != "windows" {
		return func() {}
	}
	if err := os.MkdirAll(daemonLogDir(), 0o755); err != nil {
		return func() {}
	}
	if err := os.WriteFile(daemonPidPath(), fmt.Append(nil, os.Getpid()), 0o644); err != nil {
		return func() {}
	}
	return func() { _ = os.Remove(daemonPidPath()) }
}

func stopDaemon(ctx context.Context, server string) string {
	if !daemonHealthy(ctx, server) {
		if pid := readPidFile(); pid != "" {
			platform.Current.KillProcess(ctx, pid)
			_ = os.Remove(daemonPidPath())
			return "stopped"
		}
		return "not_running"
	}
	if platform.Current.StopManagedDaemon(ctx) {
		return afterStopStatus(ctx, server)
	}
	if pid := readPidFile(); pid != "" {
		platform.Current.KillProcess(ctx, pid)
		_ = os.Remove(daemonPidPath())
		return "stopped"
	}
	return "manual"
}

func afterStopStatus(ctx context.Context, server string) string {
	if !daemonHealthy(ctx, server) {
		return "stopped"
	}
	if pid := readPidFile(); pid != "" {
		platform.Current.KillProcess(ctx, pid)
		_ = os.Remove(daemonPidPath())
		if !daemonHealthy(ctx, server) {
			return "stopped"
		}
	}
	return "manual"
}

func daemonLogDir() string {
	return filepath.Join(mustHome(), ".traceknot")
}

func daemonPidPath() string {
	return filepath.Join(daemonLogDir(), "daemon.pid")
}

func daemonLogPath() string {
	return filepath.Join(daemonLogDir(), "daemon.log")
}

func readPidFile() string {
	content, err := os.ReadFile(daemonPidPath())
	if err != nil {
		return ""
	}
	return string(content)
}

func waitHealthy(ctx context.Context, server string) bool {
	for range 10 {
		if daemonHealthy(ctx, server) {
			return true
		}
		time.Sleep(500 * time.Millisecond)
	}
	return false
}

type readyzResponse struct {
	Ready      bool     `json:"ready"`
	Rebuilding []string `json:"rebuilding"`
}

func daemonRebuilding(ctx context.Context, server string) ([]string, bool) {
	client := &http.Client{Timeout: 2 * time.Second}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(server, "/")+"/readyz", nil)
	if err != nil {
		return nil, false
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, false
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, false
	}
	var body readyzResponse
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		return nil, false
	}
	return body.Rebuilding, true
}

func waitReady(ctx context.Context, server string) {
	const maxWait = 60 * time.Second
	deadline := time.Now().Add(maxWait)
	announced := false
	for time.Now().Before(deadline) {
		rebuilding, ok := daemonRebuilding(ctx, server)
		if !ok {
			return
		}
		if len(rebuilding) == 0 {
			if announced {
				fmt.Println("session data update finished")
			}
			return
		}
		if !announced {
			fmt.Println("updating session data (" + strings.Join(rebuilding, ", ") + ")... this can take a few minutes")
			announced = true
		}
		time.Sleep(500 * time.Millisecond)
	}
	if announced {
		fmt.Println("still updating session data in the background; check " + daemonLogPath())
	}
}

func portFromURL(server string) string {
	parsed, err := url.Parse(server)
	if err != nil {
		return "4318"
	}
	if parsed.Port() != "" {
		return parsed.Port()
	}
	return "4318"
}
