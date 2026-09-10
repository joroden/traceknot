package settings

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"traceknot/internal/config"
	"traceknot/internal/install/agenthooks"
	"traceknot/internal/install/agentskills"
	"traceknot/internal/install/autostart"
)

func Apply(ctx context.Context, patch Patch) error {
	exe, err := resolveExe()
	if err != nil {
		return err
	}

	if patch.AutostartOn != nil {
		if err := applyAutostart(ctx, *patch.AutostartOn); err != nil {
			return fmt.Errorf("autostart: %w", err)
		}
	}
	if patch.RequireWorkItem != nil {
		if err := config.Save(config.Config{RequireWorkItem: *patch.RequireWorkItem}); err != nil {
			return fmt.Errorf("work item requirement: %w", err)
		}
	}
	if patch.Hooks != nil {
		if err := applyHooks(exe, patch.Hooks); err != nil {
			return err
		}
	}
	if patch.Skills != nil {
		if err := applySkills(patch.Skills); err != nil {
			return err
		}
	}
	return nil
}

func applyAutostart(ctx context.Context, on bool) error {
	if on {
		return autostart.Enable(ctx)
	}
	return autostart.Disable(ctx)
}

func applyHooks(exe string, wants map[string]bool) error {
	for _, provider := range agenthooks.Providers {
		want, ok := wants[provider.Binary]
		if !ok {
			continue
		}
		if want {
			if err := provider.Install(exe); err != nil {
				return fmt.Errorf("hook %s: %w", provider.Binary, err)
			}
			continue
		}
		if err := provider.Remove(exe); err != nil {
			return fmt.Errorf("hook %s: %w", provider.Binary, err)
		}
	}
	return nil
}

func applySkills(wants map[string]bool) error {
	for _, provider := range agentskills.Providers {
		want, ok := wants[provider.Binary]
		if !ok {
			continue
		}
		if want {
			if err := provider.Install(); err != nil {
				return fmt.Errorf("skill %s: %w", provider.Binary, err)
			}
			continue
		}
		if err := provider.Remove(); err != nil {
			return fmt.Errorf("skill %s: %w", provider.Binary, err)
		}
	}
	return nil
}

func resolveExe() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("resolve binary: %w", err)
	}
	resolved, err := filepath.Abs(exe)
	if err != nil {
		return exe, nil
	}
	return resolved, nil
}
