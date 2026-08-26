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

func WriteSession(sess *export.Session, outDir string, dedup *content.Deduper) error {
	return writeSession(sess, outDir, dedup, "")
}

func writeSession(sess *export.Session, outDir string, dedup *content.Deduper, labelPrefix string) error {
	if dedup == nil {
		dedup = content.NewDeduper()
	}
	if err := os.MkdirAll(filepath.Join(outDir, "nodes"), 0o755); err != nil {
		return fmt.Errorf("create export directory: %w", err)
	}

	ordered := orderedContentNodes(sess)
	filenames := make(map[string]string, len(ordered))
	for i, node := range ordered {
		filenames[node.NodeID] = nodeFilename(i+1, node)
	}

	for i, node := range ordered {
		filename := filenames[node.NodeID]
		body := renderNode(i+1, node, dedup, labelPrefix+"nodes/"+filename)
		if err := os.WriteFile(filepath.Join(outDir, "nodes", filename), []byte(body), 0o644); err != nil {
			return fmt.Errorf("write node file %s: %w", filename, err)
		}
	}

	computedFlags := flags.Compute(sess)

	if err := os.WriteFile(filepath.Join(outDir, "summary.md"), []byte(summaryDoc(sess.Meta)), 0o644); err != nil {
		return fmt.Errorf("write summary.md: %w", err)
	}
	if err := os.WriteFile(filepath.Join(outDir, "timeline.md"), []byte(timelineDoc(ordered, filenames)), 0o644); err != nil {
		return fmt.Errorf("write timeline.md: %w", err)
	}
	if err := os.WriteFile(filepath.Join(outDir, "flags.md"), []byte(flagsDoc(computedFlags, filenames)), 0o644); err != nil {
		return fmt.Errorf("write flags.md: %w", err)
	}
	return nil
}

func orderedContentNodes(sess *export.Session) []*export.Node {
	ordered := make([]*export.Node, 0, len(sess.Nodes))
	for _, node := range sess.Nodes {
		if isSyntheticRoot(node) {
			continue
		}
		ordered = append(ordered, node)
	}
	sort.SliceStable(ordered, func(i, j int) bool {
		return startOrZero(ordered[i]) < startOrZero(ordered[j])
	})
	return ordered
}

func isSyntheticRoot(node *export.Node) bool {
	return strings.HasPrefix(node.NodeID, "synthetic:root:")
}

func startOrZero(node *export.Node) int64 {
	if node.StartedAtUnixMs == nil {
		return 0
	}
	return *node.StartedAtUnixMs
}
