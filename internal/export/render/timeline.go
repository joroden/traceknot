package render

import (
	"fmt"
	"strings"

	"traceknot/internal/export"
)

func timelineDoc(ordered []*export.Node, filenames map[string]string) string {
	var b strings.Builder
	b.WriteString("# Timeline\n\n")
	b.WriteString("| # | kind | name | cost | duration | file | preview |\n")
	b.WriteString("|---|---|---|---|---|---|---|\n")
	for i, node := range ordered {
		seq := i + 1
		var name string
		switch node.Kind {
		case "tool_call":
			name = deref(node.ToolName)
		case "agent":
			name = deref(node.AgentName)
		default:
			name = deref(node.Name)
		}
		file := filenames[node.NodeID]
		fmt.Fprintf(&b, "| %03d | %s | %s | %s | %s | [%s](nodes/%s) | %s |\n",
			seq, node.Kind, escapeCell(name), formatCost(node.Cost), formatDuration(node.DurationMs),
			file, file, escapeCell(preview(node)))
	}
	return b.String()
}

func preview(node *export.Node) string {
	if node.Detail == nil {
		return ""
	}
	switch {
	case node.Detail.Chat != nil:
		if out := deref(node.Detail.Chat.OutputText); strings.TrimSpace(out) != "" {
			return firstLine(out, 100)
		}
		return firstLine(deref(node.Detail.Chat.PromptText), 100)
	case node.Detail.ToolCall != nil:
		return firstLine(deref(node.Detail.ToolCall.ResultText), 100)
	case node.Detail.Agent != nil:
		return firstLine(deref(node.Detail.Agent.ResultSummary), 100)
	}
	return ""
}

func escapeCell(s string) string {
	return strings.ReplaceAll(s, "|", "\\|")
}
