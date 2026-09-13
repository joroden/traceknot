package cli

import (
	"encoding/json"
	"strconv"
)

func hookSessionID(payload []byte) (string, bool) {
	var fields map[string]any
	if err := json.Unmarshal(payload, &fields); err != nil {
		return "", false
	}
	sessionID, _ := fields["session_id"].(string)
	if sessionID == "" {
		sessionID, _ = fields["sessionId"].(string)
	}
	return sessionID, sessionID != ""
}

func hookPrompt(payload []byte) string {
	var fields map[string]any
	if err := json.Unmarshal(payload, &fields); err != nil {
		return ""
	}
	prompt, _ := fields["prompt"].(string)
	return prompt
}

func hookTurnKey(payload []byte) string {
	var fields map[string]any
	if err := json.Unmarshal(payload, &fields); err != nil {
		return ""
	}
	switch timestamp := fields["timestamp"].(type) {
	case string:
		return timestamp
	case float64:
		return strconv.FormatFloat(timestamp, 'f', -1, 64)
	default:
		return ""
	}
}
