package agentenv

import (
	"context"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func CloseInteractiveSessions(ctx context.Context) (int, error) {
	if runtime.GOOS == "windows" {
		return closeSessionsWindows(ctx)
	}
	return closeSessionsPosix(ctx)
}

func closeSessionsPosix(ctx context.Context) (int, error) {
	output, err := exec.CommandContext(ctx, "ps", "-ax", "-o", "pid=", "-o", "command=").Output()
	if err != nil {
		return 0, err
	}

	closed := 0
	vscodeClosed := false
	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		pidField, commandLine, ok := strings.Cut(line, " ")
		if !ok {
			continue
		}
		pid, err := strconv.Atoi(pidField)
		if err != nil {
			continue
		}
		lowered := strings.ToLower(commandLine)

		switch {
		case strings.Contains(lowered, "github-copilot-cli") ||
			strings.Contains(lowered, " copilot ") ||
			strings.Contains(lowered, "/copilot"):
			killPID(pid)
			closed++
		case strings.Contains(lowered, "visual studio code") ||
			strings.Contains(lowered, "/code") ||
			strings.Contains(lowered, "code-insiders") ||
			strings.Contains(lowered, "/codium"):
			killPID(pid)
			closed++
			vscodeClosed = true
		}
	}

	if vscodeClosed && !IsWSL() {
		relaunchVSCodePosix(ctx)
	}
	return closed, nil
}

func killPID(pid int) {
	process, err := os.FindProcess(pid)
	if err != nil {
		return
	}
	_ = process.Kill()
}

func relaunchVSCodePosix(ctx context.Context) {
	for _, bin := range []string{"code", "code-insiders"} {
		path, err := exec.LookPath(bin)
		if err != nil {
			continue
		}
		_ = exec.CommandContext(ctx, path).Start()
		return
	}
	if runtime.GOOS == "darwin" {
		_ = exec.CommandContext(ctx, "open", "-a", "Visual Studio Code").Start()
	}
}

func closeSessionsWindows(ctx context.Context) (int, error) {
	script := `
$count = 0
foreach ($name in @('Code', 'Code - Insiders', 'Code - OSS', 'Codium')) {
  $procs = Get-Process -Name $name -ErrorAction SilentlyContinue
  $count += ($procs | Measure-Object).Count
  $procs | Stop-Process -Force -ErrorAction SilentlyContinue
}
try {
  Get-CimInstance Win32_Process | Where-Object {
    $_.CommandLine -and $_.CommandLine.ToLowerInvariant().Contains('copilot')
  } | ForEach-Object { $count++; Stop-Process -Id $_.ProcessId -Force -ErrorAction SilentlyContinue }
} catch {}
$codeCmd = Get-Command code -ErrorAction SilentlyContinue
if ($codeCmd) { Start-Process $codeCmd.Source | Out-Null }
Write-Output $count
`
	runCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	output, err := exec.CommandContext(runCtx, "powershell", "-NoProfile", "-NonInteractive", "-Command", script).Output()
	if err != nil {
		return 0, err
	}
	count, _ := strconv.Atoi(strings.TrimSpace(string(output)))
	return count, nil
}
