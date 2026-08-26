package render

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"traceknot/internal/export"
	"traceknot/internal/export/content"
	"traceknot/internal/export/flags"
)

type sessionSummary struct {
	dir       string
	session   *export.Session
	flagCount int
}

func WriteWorkItem(wi *export.WorkItem, outDir string) error {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("create export directory: %w", err)
	}

	ordered := make([]*export.Session, len(wi.Sessions))
	copy(ordered, wi.Sessions)
	sort.SliceStable(ordered, func(i, j int) bool {
		return sessionStartOrZero(ordered[i]) < sessionStartOrZero(ordered[j])
	})

	dedup := content.NewDeduper()
	summaries := make([]sessionSummary, 0, len(ordered))
	for i, sess := range ordered {
		dir := fmt.Sprintf("%02d-%s", i+1, sess.Meta.SessionID)
		if err := writeSession(sess, filepath.Join(outDir, dir), dedup, dir+"/"); err != nil {
			return fmt.Errorf("write session %s: %w", sess.Meta.SessionID, err)
		}
		summaries = append(summaries, sessionSummary{
			dir:       dir,
			session:   sess,
			flagCount: len(flags.Compute(sess)),
		})
	}

	if err := os.WriteFile(filepath.Join(outDir, "overview.md"), []byte(workItemOverviewDoc(wi, summaries)), 0o644); err != nil {
		return fmt.Errorf("write overview.md: %w", err)
	}
	return nil
}

func sessionStartOrZero(sess *export.Session) int64 {
	if sess.Meta.StartedAtUnixMs == nil {
		return 0
	}
	return *sess.Meta.StartedAtUnixMs
}

func workItemOverviewDoc(wi *export.WorkItem, summaries []sessionSummary) string {
	var totalCost float64
	var totalNodes int64
	for _, s := range summaries {
		totalCost += s.session.Meta.Cost
		totalNodes += s.session.Meta.NodeCount
	}

	var b strings.Builder
	fmt.Fprintf(&b, "# Work item %s/%s\n\n", wi.Provider, wi.Key)
	fmt.Fprintf(&b, "- title: %s\n", wi.Title)
	fmt.Fprintf(&b, "- sessions: %d\n", len(summaries))
	fmt.Fprintf(&b, "- total cost: %s\n", formatCost(totalCost))
	fmt.Fprintf(&b, "- total nodes: %d\n\n", totalNodes)

	b.WriteString("Sessions most worth reading first (by number of flagged nodes):\n\n")
	byFlags := make([]sessionSummary, len(summaries))
	copy(byFlags, summaries)
	sort.SliceStable(byFlags, func(i, j int) bool { return byFlags[i].flagCount > byFlags[j].flagCount })
	for _, s := range byFlags {
		if s.flagCount == 0 {
			continue
		}
		fmt.Fprintf(&b, "- [%s](%s/summary.md) — %d flagged nodes, see `%s/flags.md`\n", s.session.Meta.Title, s.dir, s.flagCount, s.dir)
	}

	b.WriteString("\n## All sessions, in order\n\n")
	b.WriteString("| # | session | started | cost | duration | flags | dir |\n")
	b.WriteString("|---|---|---|---|---|---|---|\n")
	for i, s := range summaries {
		meta := s.session.Meta
		fmt.Fprintf(&b, "| %02d | %s | %s | %s | %s | %d | [%s](%s/summary.md) |\n",
			i+1, escapeCell(meta.Title), formatTime(meta.StartedAtUnixMs), formatCost(meta.Cost),
			formatDuration(meta.DurationMs), s.flagCount, s.dir, s.dir)
	}
	return b.String()
}
