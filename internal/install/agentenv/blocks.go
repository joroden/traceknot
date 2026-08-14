package agentenv

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func upsertBlock(path string, block string) error {
	content, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read %s: %w", path, err)
	}
	start, end := blockMarkers([]byte(block))
	lines := strings.Split(string(content), "\n")
	var kept []string
	inBlock := false
	for _, line := range lines {
		switch {
		case strings.TrimSpace(line) == start:
			inBlock = true
		case strings.TrimSpace(line) == end:
			inBlock = false
		case !inBlock:
			kept = append(kept, line)
		}
	}
	for len(kept) > 0 && strings.TrimSpace(kept[len(kept)-1]) == "" {
		kept = kept[:len(kept)-1]
	}
	var out strings.Builder
	out.WriteString(strings.Join(kept, "\n"))
	if len(kept) > 0 {
		out.WriteString("\n")
	}
	out.WriteString(strings.TrimRight(block, "\n"))
	out.WriteString("\n")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(out.String()), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

func removeBlock(path string, start string, end string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	lines := strings.Split(string(content), "\n")
	var kept []string
	inBlock := false
	changed := false
	for _, line := range lines {
		switch {
		case strings.TrimSpace(line) == start:
			inBlock = true
			changed = true
		case strings.TrimSpace(line) == end:
			inBlock = false
		case !inBlock:
			kept = append(kept, line)
		}
	}
	if !changed {
		return nil
	}
	for len(kept) > 0 && strings.TrimSpace(kept[len(kept)-1]) == "" {
		kept = kept[:len(kept)-1]
	}
	if err := os.WriteFile(path, []byte(strings.Join(kept, "\n")+"\n"), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

func blockMarkers(block []byte) (string, string) {
	lines := strings.Split(strings.TrimRight(string(block), "\n"), "\n")
	return strings.TrimSpace(lines[0]), strings.TrimSpace(lines[len(lines)-1])
}

func profileCandidates() []string {
	home := homeDir()
	shell := os.Getenv("SHELL")
	var candidates []string
	candidates = append(candidates, filepath.Join(home, ".profile"))
	for _, name := range []string{".bash_profile", ".bash_login", ".bashrc"} {
		path := filepath.Join(home, name)
		if fileExists(path) {
			candidates = append(candidates, path)
		}
	}
	if strings.HasSuffix(shell, "zsh") || fileExists(filepath.Join(home, ".zprofile")) {
		candidates = append(candidates, filepath.Join(home, ".zprofile"))
	}
	if fileExists(filepath.Join(home, ".zshrc")) {
		candidates = append(candidates, filepath.Join(home, ".zshrc"))
	}
	return candidates
}
