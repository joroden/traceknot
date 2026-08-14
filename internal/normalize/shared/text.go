package shared

import "traceknot/internal/jsonutil"

func Preview(role string, text string) string {
	if text == "" {
		return role
	}
	cleaned := jsonutil.ToSearchableText(text)
	return jsonutil.TruncateText(cleaned, 120)
}

func Title(text string) string {
	cleaned := jsonutil.ToSearchableText(text)
	return jsonutil.TruncateText(cleaned, 80)
}

func EventMetadata(eventName string) string {
	return jsonutil.ToCanonicalJSON(map[string]any{"eventName": eventName})
}
