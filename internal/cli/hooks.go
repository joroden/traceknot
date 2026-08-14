package cli

import (
	"fmt"
	"os"

	"traceknot/internal/install/agenthooks"
)

type providerChoice struct {
	provider agenthooks.Provider
	want     bool
}

func providerChoices(exe string) []providerChoice {
	choices := make([]providerChoice, 0, len(agenthooks.Providers))
	for _, provider := range agenthooks.Providers {
		choices = append(choices, providerChoice{
			provider: provider,
			want:     provider.Installed(exe),
		})
	}
	return choices
}

func applyChoices(choices []providerChoice, exe string) {
	for _, choice := range choices {
		if choice.want {
			if err := choice.provider.Install(exe); err != nil {
				fmt.Fprintln(os.Stderr, "hooks: "+choice.provider.Binary+":", err)
				continue
			}
			fmt.Println("hooks: enabled for " + choice.provider.Binary)
			continue
		}
		if err := choice.provider.Remove(exe); err != nil {
			fmt.Fprintln(os.Stderr, "hooks: "+choice.provider.Binary+":", err)
			continue
		}
		fmt.Println("hooks: disabled for " + choice.provider.Binary)
	}
	for _, choice := range choices {
		if !choice.want {
			continue
		}
		switch choice.provider.Binary {
		case "codex":
			fmt.Println("Next: run /hooks in the codex CLI and trust the SessionStart hook")
		case "copilot":
			fmt.Println("Next: restart the copilot CLI so hooks load")
		}
	}
}
