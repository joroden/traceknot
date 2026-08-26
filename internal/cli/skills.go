package cli

import (
	"fmt"
	"os"
	"slices"

	"github.com/charmbracelet/huh"

	"traceknot/internal/install/agentskills"
)

type skillChoice struct {
	provider agentskills.Provider
	want     bool
}

func skillChoices() []skillChoice {
	choices := make([]skillChoice, 0, len(agentskills.Providers))
	for _, provider := range agentskills.Providers {
		choices = append(choices, skillChoice{
			provider: provider,
			want:     provider.Installed(),
		})
	}
	return choices
}

func skillOptions(choices []skillChoice) []huh.Option[string] {
	options := make([]huh.Option[string], 0, len(choices))
	for _, choice := range choices {
		options = append(options, huh.NewOption("session analysis skill for "+choice.provider.Binary, choice.provider.Binary))
	}
	return options
}

func selectedSkills(choices []skillChoice) []string {
	selected := make([]string, 0, len(choices))
	for _, choice := range choices {
		if choice.want {
			selected = append(selected, choice.provider.Binary)
		}
	}
	return selected
}

func applySkillSelection(choices []skillChoice, selected []string) {
	for i := range choices {
		choices[i].want = slices.Contains(selected, choices[i].provider.Binary)
	}
	for _, choice := range choices {
		if choice.want {
			if err := choice.provider.Install(); err != nil {
				fmt.Fprintln(os.Stderr, "skill: "+choice.provider.Binary+":", err)
				continue
			}
			fmt.Println("skill: enabled for " + choice.provider.Binary)
			continue
		}
		if err := choice.provider.Remove(); err != nil {
			fmt.Fprintln(os.Stderr, "skill: "+choice.provider.Binary+":", err)
			continue
		}
		fmt.Println("skill: disabled for " + choice.provider.Binary)
	}
}
