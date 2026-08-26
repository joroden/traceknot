package render

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	return "```\n" + collapsed + "\n```"
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
		return "```json\n" + buf.String() + "\n```"
	}
	return "```\n" + raw + "\n```"
}
