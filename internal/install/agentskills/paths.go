package agentskills

import (
	"fmt"
	"os"
	"path/filepath"
)

const skillName = "traceknot-analyze"

func claudeSkillPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home: %w", err)
	}
	return filepath.Join(home, ".claude", "skills", skillName), nil
}

func codexSkillPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home: %w", err)
	}
	return filepath.Join(home, ".agents", "skills", skillName), nil
}

func copilotSkillPath() (string, error) {
	copilotHome := os.Getenv("COPILOT_HOME")
	if copilotHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("resolve home: %w", err)
		}
		copilotHome = filepath.Join(home, ".copilot")
	}
	return filepath.Join(copilotHome, "skills", skillName), nil
}
