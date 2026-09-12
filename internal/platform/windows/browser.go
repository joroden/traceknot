//go:build windows

package windows

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func OpenBrowser(target string) error {
	if browser := findChromium(); browser != "" {
		command := exec.Command(browser, "--app="+target)
		if err := command.Start(); err == nil {
			return nil
		}
	}
	return OpenDefault(target)
}

func OpenDefault(target string) error {
	command := exec.Command("cmd", "/c", "start", "", target)
	hidden(command)
	return command.Start()
}

func findChromium() string {
	for _, candidate := range preferDefault(chromiumCandidates(), defaultBrowserProgID()) {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return ""
}

func chromiumCandidates() []string {
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

func defaultBrowserProgID() string {
	command := exec.Command("reg", "query",
		`HKCU\Software\Microsoft\Windows\Shell\Associations\UrlAssociations\http\UserChoice`, "/v", "ProgId")
	hidden(command)
	output, err := command.Output()
	if err != nil {
		return ""
	}
	fields := strings.Fields(string(output))
	if len(fields) == 0 {
		return ""
	}
	return fields[len(fields)-1]
}
