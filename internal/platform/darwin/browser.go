//go:build darwin

package darwin

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func CanOpenBrowser() bool { return true }

func OpenBrowser(target string) error {
	if browser := findChromium(); browser != "" {
		command := exec.Command("open", "-na", chromiumAppBundle(browser), "--args", "--app="+target)
		if err := command.Start(); err == nil {
			return nil
		}
	}
	return OpenDefault(target)
}

func OpenDefault(target string) error {
	return exec.Command("open", target).Start()
}

func chromiumAppBundle(executablePath string) string {
	if idx := strings.Index(executablePath, ".app/"); idx != -1 {
		return executablePath[:idx+len(".app")]
	}
	return executablePath
}

func findChromium() string {
	for _, candidate := range chromiumCandidates() {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return ""
}

func chromiumCandidates() []string {
	apps := []string{
		"Google Chrome.app/Contents/MacOS/Google Chrome",
		"Microsoft Edge.app/Contents/MacOS/Microsoft Edge",
		"Brave Browser.app/Contents/MacOS/Brave Browser",
		"Opera.app/Contents/MacOS/Opera",
	}
	userHome, err := os.UserHomeDir()
	if err != nil {
		userHome = ""
	}
	candidates := make([]string, 0, len(apps)*2)
	for _, app := range apps {
		if userHome != "" {
			candidates = append(candidates, filepath.Join(userHome, "Applications", app))
		}
		candidates = append(candidates, filepath.Join("/Applications", app))
	}
	return candidates
}
