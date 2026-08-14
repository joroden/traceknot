package claude

import (
	"encoding/json"
	"strings"
)

func systemPromptFromBody(bodyJSON string) string {
	var parsed struct {
		System []struct {
			Text string `json:"text"`
		} `json:"system"`
	}
	if err := json.Unmarshal([]byte(bodyJSON), &parsed); err != nil {
		return ""
	}
	best := ""
	for _, block := range parsed.System {
		if isBillingHeader(block.Text) {
			continue
		}
		if len(block.Text) > len(best) {
			best = block.Text
		}
	}
	return best
}

func isBillingHeader(text string) bool {
	return strings.HasPrefix(text, "x-anthropic-billing-header")
}
