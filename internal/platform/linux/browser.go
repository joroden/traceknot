//go:build linux

package linux

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func CanOpenBrowser() bool {
	if IsWSL() {
		return true
	}
	return os.Getenv("DISPLAY") != "" || os.Getenv("WAYLAND_DISPLAY") != ""
}

func OpenBrowser(target string) error {
	if browser := findChromium(); browser != "" {
		command := exec.Command(browser, "--app="+target)
		if workDir := chromiumWorkDir(); workDir != "" {
			command.Dir = workDir
		}
		if err := command.Start(); err == nil {
			return nil
		}
	}
	return OpenDefault(target)
}

func OpenDefault(target string) error {
	if IsWSL() {
		if err := exec.Command("wslview", target).Start(); err == nil {
			return nil
		}
		return exec.Command("explorer.exe", target).Start()
	}
	return exec.Command("xdg-open", target).Start()
}

func findChromium() string {
	candidates := chromiumCandidates()
	if IsWSL() {
		candidates = preferDefault(candidates, wslDefaultBrowserProgID())
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return ""
}

func chromiumCandidates() []string {
	if IsWSL() {
		return wslChromiumCandidates()
	}
	names := []string{
		"google-chrome",
		"google-chrome-stable",
		"chromium",
		"chromium-browser",
		"microsoft-edge",
		"brave-browser",
		"opera",
	}
	candidates := make([]string, 0, len(names))
	for _, name := range names {
		if path, err := exec.LookPath(name); err == nil {
			candidates = append(candidates, path)
		}
	}
	return candidates
}

func wslChromiumCandidates() []string {
	return []string{
		"/mnt/c/Program Files/Google/Chrome/Application/chrome.exe",
		"/mnt/c/Program Files (x86)/Google/Chrome/Application/chrome.exe",
		"/mnt/c/Program Files (x86)/Microsoft/Edge/Application/msedge.exe",
		"/mnt/c/Program Files/Microsoft/Edge/Application/msedge.exe",
		"/mnt/c/Program Files/BraveSoftware/Brave-Browser/Application/brave.exe",
		"/mnt/c/Program Files (x86)/BraveSoftware/Brave-Browser/Application/brave.exe",
		"/mnt/c/Program Files (x86)/Opera Software/Opera/opera.exe",
	}
}

func preferDefault(candidates []string, progID string) []string {
	needle := ""
	switch progID {
	case "ChromeHTML":
		needle = "chrome"
	case "MSEdgeHTM":
		needle = "edge"
	case "BraveHTML":
		needle = "brave"
	case "OperaStable":
		needle = "opera"
	}
	if needle == "" {
		return candidates
	}
	preferred := make([]string, 0, len(candidates))
	rest := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		if strings.Contains(strings.ToLower(filepath.Base(candidate)), needle) {
			preferred = append(preferred, candidate)
		} else {
			rest = append(rest, candidate)
		}
	}
	return append(preferred, rest...)
}

func wslDefaultBrowserProgID() string {
	output, err := exec.Command("cmd.exe", "/c", "reg", "query",
		`HKCU\Software\Microsoft\Windows\Shell\Associations\UrlAssociations\http\UserChoice`,
		"/v", "ProgId").Output()
	if err != nil {
		return ""
	}
	fields := strings.Fields(string(output))
	if len(fields) == 0 {
		return ""
	}
	return fields[len(fields)-1]
}

func chromiumWorkDir() string {
	if IsWSL() {
		return "/mnt/c"
	}
	return ""
}
