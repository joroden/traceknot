package render

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var slugDisallowed = regexp.MustCompile(`[^a-z0-9]+`)

func slug(s string) string {
	lowered := strings.ToLower(strings.TrimSpace(s))
	replaced := slugDisallowed.ReplaceAllString(lowered, "-")
	trimmed := strings.Trim(replaced, "-")
	if trimmed == "" {
		return "node"
	}
	return trimmed
}

func formatCost(value float64) string {
	return fmt.Sprintf("$%.4f", value)
}

func formatDuration(ms *float64) string {
	if ms == nil {
		return "n/a"
	}
	seconds := *ms / 1000
	if seconds < 1 {
		return fmt.Sprintf("%.0fms", *ms)
	}
	return fmt.Sprintf("%.1fs", seconds)
}

func formatTime(unixMs *int64) string {
	if unixMs == nil {
		return "n/a"
	}
	return time.UnixMilli(*unixMs).UTC().Format(time.RFC3339)
}

func firstLine(text string, maxLen int) string {
	text = strings.TrimSpace(text)
	if idx := strings.IndexByte(text, '\n'); idx >= 0 {
		text = text[:idx]
	}
	runes := []rune(text)
	if len(runes) > maxLen {
		return string(runes[:maxLen]) + "…"
	}
	return string(runes)
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
