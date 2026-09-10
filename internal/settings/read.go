package settings

import (
	"context"

	"traceknot/internal/config"
	"traceknot/internal/install/agenthooks"
	"traceknot/internal/install/agentskills"
	"traceknot/internal/install/autostart"
)

func Current(ctx context.Context) (State, error) {
	exe, err := resolveExe()
	if err != nil {
		return State{}, err
	}

	hooks := make([]ProviderState, 0, len(agenthooks.Providers))
	for _, provider := range agenthooks.Providers {
		hooks = append(hooks, ProviderState{Binary: provider.Binary, Enabled: provider.Installed(exe)})
	}

	skills := make([]ProviderState, 0, len(agentskills.Providers))
	for _, provider := range agentskills.Providers {
		skills = append(skills, ProviderState{Binary: provider.Binary, Enabled: provider.Installed()})
	}

	return State{
		AutostartOn:     autostart.Enabled(ctx),
		RequireWorkItem: config.Load().RequireWorkItem,
		Hooks:           hooks,
		Skills:          skills,
	}, nil
}
