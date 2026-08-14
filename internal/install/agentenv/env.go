package agentenv

import (
	"os"
	"runtime"
	"strings"
)

func ApplyEnv(binDir string) error {
	if runtime.GOOS == "windows" {
		return applyWindowsEnv(binDir)
	}
	return applyUnixEnv(binDir)
}

func RemoveEnv() error {
	if runtime.GOOS == "windows" {
		return removeWindowsEnv()
	}
	return removeUnixEnv()
}

func homeDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return dir
}

type envPair struct {
	name  string
	value string
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func unixEnvPairs(block []byte) []envPair {
	var pairs []envPair
	for _, line := range strings.Split(string(block), "\n") {
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, "export ")
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, found := strings.Cut(line, "=")
		if found {
			pairs = append(pairs, envPair{name: strings.TrimSpace(key), value: trimQuotes(strings.TrimSpace(value))})
		}
	}
	return pairs
}

func windowsEnvPairs(block []byte) []envPair {
	var pairs []envPair
	for _, line := range strings.Split(string(block), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, found := strings.Cut(line, "=")
		if found {
			pairs = append(pairs, envPair{name: strings.TrimSpace(key), value: trimQuotes(strings.TrimSpace(value))})
		}
	}
	return pairs
}

func trimQuotes(value string) string {
	if len(value) >= 2 {
		first, last := value[0], value[len(value)-1]
		if (first == '\'' && last == '\'') || (first == '"' && last == '"') {
			return value[1 : len(value)-1]
		}
	}
	return value
}
