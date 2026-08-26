package content

import (
	"encoding/json"
	"math"
	"sort"
	"strconv"
	"strings"
)

const (
	minTableRows    = 3
	maxTableColumns = 15
)

func CompactTabularJSON(raw string) (string, bool) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return raw, false
	}
	var rows []map[string]any
	if err := json.Unmarshal([]byte(trimmed), &rows); err != nil || len(rows) < minTableRows {
		return raw, false
	}
	columns := commonColumns(rows)
	if len(columns) == 0 || len(columns) > maxTableColumns {
		return raw, false
	}
	return renderTable(columns, rows), true
}

func commonColumns(rows []map[string]any) []string {
	seen := make(map[string]bool)
	var columns []string
	for _, row := range rows {
		for key := range row {
			if !seen[key] {
				seen[key] = true
				columns = append(columns, key)
			}
		}
	}
	sort.Strings(columns)
	return columns
}

func renderTable(columns []string, rows []map[string]any) string {
	var b strings.Builder
	b.WriteString("| ")
	b.WriteString(strings.Join(columns, " | "))
	b.WriteString(" |\n|")
	b.WriteString(strings.Repeat(" --- |", len(columns)))
	b.WriteString("\n")
	for _, row := range rows {
		cells := make([]string, len(columns))
		for i, col := range columns {
			cells[i] = cellText(row[col])
		}
		b.WriteString("| ")
		b.WriteString(strings.Join(cells, " | "))
		b.WriteString(" |\n")
	}
	return b.String()
}

func cellText(value any) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return strings.ReplaceAll(v, "|", "\\|")
	case float64:
		return formatFloat(v)
	case bool:
		return strconv.FormatBool(v)
	default:
		encoded, err := json.Marshal(v)
		if err != nil {
			return ""
		}
		return strings.ReplaceAll(string(encoded), "|", "\\|")
	}
}

func formatFloat(v float64) string {
	if v == math.Trunc(v) && math.Abs(v) < 1e15 {
		return strconv.FormatInt(int64(v), 10)
	}
	return strconv.FormatFloat(v, 'g', -1, 64)
}
