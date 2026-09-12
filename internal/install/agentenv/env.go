package agentenv

import (
	"os"
	"strings"
)

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

func pairMap(pairs []envPair) map[string]string {
	vars := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		vars[pair.name] = pair.value
	}
	return vars
}

func pairNames(pairs []envPair) []string {
	names := make([]string, len(pairs))
	for i, pair := range pairs {
		names[i] = pair.name
	}
	return names
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
