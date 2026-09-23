package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"strings"

	"traceknot/internal/export/content"
)

func renderText(dedup *content.Deduper, label, text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return "_(empty)_"
	}
	if firstLabel, isDuplicate := dedup.Check(text, label); isDuplicate {
		return fmt.Sprintf("_(identical to %s)_", firstLabel)
	}
	stripped := content.StripBinaryBlobs(text)
	collapsed := content.CollapseRepeatedLines(stripped)
	return fencedText(collapsed, "")
}

func renderJSON(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "_(empty)_"
	}
	if table, ok := content.CompactTabularJSON(raw); ok {
		return table
	}
	var buf bytes.Buffer
	if err := json.Indent(&buf, []byte(raw), "", "  "); err == nil {
		return fencedText(buf.String(), "json")
	}
	return fencedText(raw, "")
}

func fencedText(value, language string) string {
	longest, run := 2, 0
	for i := 0; i < len(value); i++ {
		if value[i] == '`' {
			run++
			if run > longest {
				longest = run
			}
		} else {
			run = 0
		}
	}
	fence := strings.Repeat("`", longest+1)
	return fence + language + "\n" + value + "\n" + fence
}

func escapeMarkdownInline(value string) string {
	value = strings.NewReplacer("\r", " ", "\n", " ").Replace(value)
	value = html.EscapeString(value)
	return strings.NewReplacer(
		"\\", "\\\\", "`", "\\`", "*", "\\*", "_", "\\_",
		"[", "\\[", "]", "\\]", "!", "\\!", "|", "\\|", "~", "\\~",
	).Replace(value)
}
