package copilot

import (
	"encoding/json"
	"strings"
)

type messagePart struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type message struct {
	Role  string        `json:"role"`
	Parts []messagePart `json:"parts"`
}

func parseMessages(raw string) []message {
	if raw == "" {
		return nil
	}
	var parsed []message
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return nil
	}
	return parsed
}

func joinText(msg message) string {
	var parts []string
	for _, part := range msg.Parts {
		if part.Type == "text" && part.Content != "" {
			parts = append(parts, part.Content)
		}
	}
	return strings.Join(parts, "\n")
}

func lastRoleText(raw string, role string) string {
	messages := parseMessages(raw)
	for index := len(messages) - 1; index >= 0; index-- {
		if messages[index].Role != role {
			continue
		}
		if text := joinText(messages[index]); text != "" {
			return text
		}
	}
	return ""
}

func firstRoleText(raw string, role string) string {
	for _, msg := range parseMessages(raw) {
		if msg.Role != role {
			continue
		}
		if text := joinText(msg); text != "" {
			return text
		}
	}
	return ""
}

func stripInjectedContext(text string) string {
	if idx := strings.Index(text, "</current_datetime>"); idx != -1 {
		text = text[idx+len("</current_datetime>"):]
	}
	if idx := strings.Index(text, "<system_reminder>"); idx != -1 {
		text = text[:idx]
	}
	return strings.TrimSpace(text)
}
