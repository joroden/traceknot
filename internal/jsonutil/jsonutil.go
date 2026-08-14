package jsonutil

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

func sortJSONValue(value any) any {
	switch typed := value.(type) {
	case []any:
		output := make([]any, 0, len(typed))
		for _, entry := range typed {
			output = append(output, sortJSONValue(entry))
		}
		return output
	case map[string]any:
		keys := make([]string, 0, len(typed))
		for key := range typed {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		output := make(map[string]any, len(typed))
		for _, key := range keys {
			output[key] = sortJSONValue(typed[key])
		}
		return output
	default:
		return value
	}
}

func ToCanonicalJSON(value any) string {
	encoded, err := json.Marshal(sortJSONValue(value))
	if err != nil {
		return "null"
	}
	return string(encoded)
}

func ToSearchableText(value any) string {
	switch typed := value.(type) {
	case nil:
		return ""
	case string:
		return typed
	case json.Number:
		return typed.String()
	case float64:
		if typed == float64(int64(typed)) {
			return formatInt(int64(typed))
		}
		return formatFloat(typed)
	case bool:
		if typed {
			return "true"
		}
		return "false"
	case []any:
		var parts []string
		for _, entry := range typed {
			if text := ToSearchableText(entry); strings.TrimSpace(text) != "" {
				parts = append(parts, text)
			}
		}
		return strings.Join(parts, "\n")
	case map[string]any:
		if text, ok := typed["text"].(string); ok {
			return text
		}
		if content, ok := typed["content"].(string); ok {
			return content
		}
		if _, ok := typed["content"]; ok {
			if nested := ToSearchableText(typed["content"]); strings.TrimSpace(nested) != "" {
				return nested
			}
		}
		if _, ok := typed["message"]; ok {
			if nested := ToSearchableText(typed["message"]); strings.TrimSpace(nested) != "" {
				return nested
			}
		}
		encoded, err := json.Marshal(typed)
		if err == nil {
			return string(encoded)
		}
		return ""
	default:
		encoded, err := json.Marshal(typed)
		if err == nil {
			return string(encoded)
		}
		return ""
	}
}

func TruncateText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}
	cut := maxLength - 3
	for cut > 0 && !utf8.RuneStart(text[cut]) {
		cut--
	}
	return text[:cut] + "..."
}

func formatInt(value int64) string {
	return strconv.FormatInt(value, 10)
}

func formatFloat(value float64) string {
	return strconv.FormatFloat(value, 'g', -1, 64)
}
