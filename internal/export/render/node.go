package render

import (
	"fmt"
	"strings"

	"traceknot/internal/export"
	"traceknot/internal/export/content"
	"traceknot/internal/store"
)

func nodeFilename(seq int, node *export.Node) string {
	label := node.Kind
	switch node.Kind {
	case "tool_call":
		if node.ToolName != nil {
			label = node.Kind + "-" + slug(*node.ToolName)
		}
	case "agent":
		if node.AgentName != nil {
			label = node.Kind + "-" + slug(*node.AgentName)
		}
	default:
		if node.Name != nil {
			label = node.Kind + "-" + slug(*node.Name)
		}
	}
	return fmt.Sprintf("%03d-%s.md", seq, label)
}

func renderNode(seq int, node *export.Node, dedup *content.Deduper, relPath string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %03d · %s", seq, node.Kind)
	if node.Name != nil {
		fmt.Fprintf(&b, " · %s", *node.Name)
	}
	b.WriteString("\n\n")
	fmt.Fprintf(&b, "- node_id: %s\n", node.NodeID)
	if node.Status != nil {
		fmt.Fprintf(&b, "- status: %s\n", *node.Status)
	}
	fmt.Fprintf(&b, "- started: %s\n", formatTime(node.StartedAtUnixMs))
	fmt.Fprintf(&b, "- duration: %s\n", formatDuration(node.DurationMs))
	fmt.Fprintf(&b, "- cost: %s\n", formatCost(node.Cost))
	if node.Kind == "agent" {
		fmt.Fprintf(&b, "- subtree cost: %s (across %d descendant nodes)\n", formatCost(node.AggCost), node.DescendantCount)
	}
	fmt.Fprintf(&b, "- tokens: in=%d (cached=%d) out=%d reasoning=%d\n",
		node.InputTokens, node.CachedInputTokens, node.OutputTokens, node.ReasoningTokens)
	b.WriteString("\n")

	if node.Detail == nil {
		b.WriteString("_(no content recorded for this node)_\n")
		return b.String()
	}

	switch {
	case node.Detail.Chat != nil:
		renderChatDetail(&b, node.Detail, dedup, relPath)
	case node.Detail.ToolCall != nil:
		renderToolCallDetail(&b, node.Detail, dedup, relPath)
	case node.Detail.Agent != nil:
		renderAgentDetail(&b, node.Detail, dedup, relPath)
	}
	return b.String()
}

func renderChatDetail(b *strings.Builder, detail *store.NodeDetail, dedup *content.Deduper, relPath string) {
	chat := detail.Chat
	if system := strings.TrimSpace(deref(chat.SystemText)); system != "" {
		b.WriteString("## System\n\n")
		b.WriteString(renderText(dedup, relPath+"#system", system))
		b.WriteString("\n\n")
	}
	b.WriteString("## Prompt\n\n")
	b.WriteString(renderText(dedup, relPath+"#prompt", deref(chat.PromptText)))
	b.WriteString("\n\n## Output\n\n")
	b.WriteString(renderText(dedup, relPath+"#output", deref(chat.OutputText)))
	if reasoning := strings.TrimSpace(deref(chat.ReasoningText)); reasoning != "" {
		b.WriteString("\n\n## Reasoning\n\n")
		b.WriteString(renderText(dedup, relPath+"#reasoning", reasoning))
	}
	b.WriteString("\n")
}

func renderToolCallDetail(b *strings.Builder, detail *store.NodeDetail, dedup *content.Deduper, relPath string) {
	tc := detail.ToolCall
	if name := deref(tc.ToolName); name != "" {
		fmt.Fprintf(b, "- tool: %s\n", name)
	}
	if decision := deref(tc.ApprovalDecision); decision != "" {
		fmt.Fprintf(b, "- approval: %s\n", decision)
	}
	b.WriteString("\n## Arguments\n\n")
	b.WriteString(renderJSON(deref(tc.ArgumentsJSON)))
	b.WriteString("\n\n## Result\n\n")
	b.WriteString(renderText(dedup, relPath+"#result", deref(tc.ResultText)))
	if errDetails := strings.TrimSpace(deref(tc.ErrorDetailsJSON)); errDetails != "" {
		b.WriteString("\n\n## Error details\n\n")
		b.WriteString(renderJSON(errDetails))
	}
	b.WriteString("\n")
}

func renderAgentDetail(b *strings.Builder, detail *store.NodeDetail, dedup *content.Deduper, relPath string) {
	agent := detail.Agent
	if agentType := deref(agent.AgentType); agentType != "" {
		fmt.Fprintf(b, "- agent_type: %s\n", agentType)
	}
	b.WriteString("\n## Spawn prompt\n\n")
	b.WriteString(renderText(dedup, relPath+"#spawn", deref(agent.SpawnPrompt)))
	b.WriteString("\n\n## Result summary\n\n")
	b.WriteString(renderText(dedup, relPath+"#summary", deref(agent.ResultSummary)))
	b.WriteString("\n")
}
