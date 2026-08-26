package agentskills

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type Provider struct {
	Binary string
	path   func() (string, error)
}

var Providers = []Provider{
	{Binary: "claude", path: claudeSkillPath},
	{Binary: "codex", path: codexSkillPath},
	{Binary: "copilot", path: copilotSkillPath},
}

func (p Provider) Install() error {
	dir, err := p.path()
	if err != nil {
		return err
	}
	files, err := SkillFiles()
	if err != nil {
		return fmt.Errorf("load skill files: %w", err)
	}
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("clear existing skill dir: %w", err)
	}
	return copyFS(files, dir)
}

func (p Provider) Remove() error {
	dir, err := p.path()
	if err != nil {
		return err
	}
	return os.RemoveAll(dir)
}

func (p Provider) Installed() bool {
	dir, err := p.path()
	if err != nil {
		return false
	}
	info, err := os.Stat(filepath.Join(dir, "SKILL.md"))
	return err == nil && !info.IsDir()
}

func copyFS(files fs.FS, dir string) error {
	return fs.WalkDir(files, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		target := filepath.Join(dir, path)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := fs.ReadFile(files, path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}
		return os.WriteFile(target, data, 0o644)
	})
}
