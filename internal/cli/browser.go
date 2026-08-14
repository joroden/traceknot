package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func openBrowser(target string) error {
	if browser := findChromium(); browser != "" {
		command := exec.Command(browser, "--app="+target)
		if workDir := chromiumWorkDir(); workDir != "" {
			command.Dir = workDir
		}
		if err := command.Start(); err == nil {
			return nil
		}
	}
	return openWithDefaultHandler(target)
}

func findChromium() string {
	candidates := chromiumCandidates()
	if runtime.GOOS == "windows" || isWSL() {
		candidates = preferDefault(candidates, windowsDefaultProgID())
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return ""
}

func chromiumCandidates() []string {
	switch runtime.GOOS {
	case "darwin":
		return macChromiumCandidates()
	case "windows":
		return windowsChromiumCandidates()
	default:
		if isWSL() {
			return windowsChromiumCandidates()
		}
		return linuxChromiumCandidates()
	}
}

func macChromiumCandidates() []string {
	apps := []string{
		"Google Chrome.app/Contents/MacOS/Google Chrome",
		"Microsoft Edge.app/Contents/MacOS/Microsoft Edge",
		"Brave Browser.app/Contents/MacOS/Brave Browser",
		"Opera.app/Contents/MacOS/Opera",
	}
	home, err := os.UserHomeDir()
	if err != nil {
		home = ""
	}
	candidates := make([]string, 0, len(apps)*2)
	for _, app := range apps {
		if home != "" {
			candidates = append(candidates, filepath.Join(home, "Applications", app))
		}
		candidates = append(candidates, filepath.Join("/Applications", app))
	}
	return candidates
}

func windowsChromiumCandidates() []string {
	if runtime.GOOS == "windows" {
		return []string{
			`C:\Program Files\Google\Chrome\Application\chrome.exe`,
			`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
			`C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`,
			`C:\Program Files\Microsoft\Edge\Application\msedge.exe`,
			`C:\Program Files\BraveSoftware\Brave-Browser\Application\brave.exe`,
			`C:\Program Files (x86)\BraveSoftware\Brave-Browser\Application\brave.exe`,
			`C:\Program Files (x86)\Opera Software\Opera\opera.exe`,
		}
	}
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

func linuxChromiumCandidates() []string {
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

func windowsDefaultProgID() string {
	args := []string{
		"query",
		`HKCU\Software\Microsoft\Windows\Shell\Associations\UrlAssociations\http\UserChoice`,
		"/v",
		"ProgId",
	}
	var output []byte
	var err error
	if runtime.GOOS == "windows" {
		output, err = exec.Command("reg", args...).Output()
	} else {
		output, err = exec.Command("cmd.exe", append([]string{"/c", "reg"}, args...)...).Output()
	}
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
	if runtime.GOOS != "windows" && isWSL() {
		return "/mnt/c"
	}
	return ""
}

func openWithDefaultHandler(target string) error {
	switch runtime.GOOS {
	case "darwin":
		return launchBrowser("open", target)
	case "windows":
		return launchBrowser("cmd", "/c", "start", "", `"`+target+`"`)
	default:
		if isWSL() {
			if err := launchBrowser("wslview", target); err == nil {
				return nil
			}
			return launchBrowser("explorer.exe", target)
		}
		return launchBrowser("xdg-open", target)
	}
}

func launchBrowser(name string, args ...string) error {
	command := exec.Command(name, args...)
	return command.Start()
}

func isWSL() bool {
	osrelease, err := os.ReadFile("/proc/sys/kernel/osrelease")
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(osrelease)), "microsoft")
}
