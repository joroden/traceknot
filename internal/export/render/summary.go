package render

import (
	"fmt"
	"strings"

	"traceknot/internal/store"
)

func summaryDoc(meta *store.SessionMeta) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Session %s\n\n", meta.SessionID)
	fmt.Fprintf(&b, "- provider: %s\n", meta.Provider)
	fmt.Fprintf(&b, "- title: %s\n", meta.Title)
	if meta.ServiceName != nil {
		fmt.Fprintf(&b, "- service: %s\n", *meta.ServiceName)
	}
	if meta.Status != nil {
		fmt.Fprintf(&b, "- status: %s\n", *meta.Status)
	}
	fmt.Fprintf(&b, "- started: %s\n", formatTime(meta.StartedAtUnixMs))
	fmt.Fprintf(&b, "- ended: %s\n", formatTime(meta.EndedAtUnixMs))
	fmt.Fprintf(&b, "- duration: %s\n", formatDuration(meta.DurationMs))
	fmt.Fprintf(&b, "- nodes: %d (tool calls: %d, subagent runs: %d)\n", meta.NodeCount, meta.ToolCallCount, meta.AgentRunCount)
	fmt.Fprintf(&b, "- cost: %s\n", formatCost(meta.Cost))
	fmt.Fprintf(&b, "- tokens: in=%d (cached=%d, cache_write=%d) out=%d reasoning=%d\n",
		meta.InputTokens, meta.CachedInputTokens, meta.CacheWriteTokens, meta.OutputTokens, meta.ReasoningTokens)
	if models := strings.TrimSpace(meta.ModelsJSON); models != "" && models != "[]" && models != "null" {
		fmt.Fprintf(&b, "- models: %s\n", models)
	}
	b.WriteString("\nSee timeline.md for the full ordered list of turns/tool calls, and flags.md for the nodes most worth reading first.\n")
	return b.String()
}
