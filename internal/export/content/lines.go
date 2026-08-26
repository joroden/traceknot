package content

import (
	"fmt"
	"strings"
)

const minRepeatRun = 3

func CollapseRepeatedLines(text string) string {
	lines := strings.Split(text, "\n")
	var out []string
	i := 0
	for i < len(lines) {
		j := i + 1
		for j < len(lines) && lines[j] == lines[i] {
			j++
		}
		runLen := j - i
		out = append(out, lines[i])
		if runLen >= minRepeatRun {
			out = append(out, fmt.Sprintf("... (previous line repeated %d more times) ...", runLen-1))
		} else if runLen == 2 {
			out = append(out, lines[i])
		}
		i = j
	}
	return strings.Join(out, "\n")
}
