package render

import (
	"fmt"
	"strings"

	"traceknot/internal/export/flags"
)

func flagsDoc(computed []flags.Flag, filenames map[string]string) string {
	var b strings.Builder
	b.WriteString("# Flags\n\n")
	if len(computed) == 0 {
		b.WriteString("No statistical outliers, errors, or repeated-call loops detected.\n")
		return b.String()
	}
	b.WriteString("Nodes worth checking first, computed from cost/duration statistics, tool-call status, ")
	b.WriteString("and repeated-call detection — start here before reading the full timeline.\n\n")
	for _, f := range computed {
		links := make([]string, 0, len(f.NodeIDs))
		for _, id := range f.NodeIDs {
			if file, ok := filenames[id]; ok {
				links = append(links, fmt.Sprintf("[%s](nodes/%s)", id, file))
			}
		}
		fmt.Fprintf(&b, "- **%s** — %s (%s)\n", f.Kind, f.Reason, strings.Join(links, ", "))
	}
	return b.String()
}
