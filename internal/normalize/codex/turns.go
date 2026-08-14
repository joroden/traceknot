package codex

import (
	"traceknot/internal/model"
	"traceknot/internal/normalize/shared"
)

type turnSpan struct {
	userNode *model.ChatSeed
	lastStep *model.ChatSeed
}

func hasUserPrompt(event Event) bool {
	prompt, _ := attributeString(event.Attributes, "prompt")
	return prompt != ""
}

func resolveRolloutMessages(turns []*turnSpan, sorted []Event, containerByCallID map[string]*model.ChatSeed) {
	if len(turns) == 0 {
		return
	}
	turnIndex := -1
	var pending []string
	flush := func() {
		if turnIndex < 0 || len(pending) == 0 {
			return
		}
		container := turns[turnIndex].lastStep
		if container == nil {
			container = turns[turnIndex].userNode
		}
		for _, text := range pending {
			appendOutputText(container, text)
		}
		pending = nil
	}
	for _, event := range sorted {
		switch event.Name {
		case eventUserPrompt:
			if hasUserPrompt(event) {
				flush()
				turnIndex++
			}
		case eventRolloutMessage:
			text, _ := attributeString(event.Attributes, "text")
			if text != "" {
				pending = append(pending, text)
			}
		case eventRolloutCall:
			if len(pending) == 0 {
				continue
			}
			callID, _ := attributeString(event.Attributes, "call_id")
			if container, ok := containerByCallID[callID]; ok {
				for _, text := range pending {
					appendOutputText(container, text)
				}
			}
			pending = nil
		}
	}
	flush()
}

func appendOutputText(container *model.ChatSeed, text string) {
	if container.OutputText != "" {
		container.OutputText += "\n" + text
	} else {
		container.OutputText = text
	}
	container.PreviewText = shared.Preview("assistant", container.OutputText)
}
