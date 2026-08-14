package codex

import (
	"strconv"

	"traceknot/internal/normalize/shared"
	"traceknot/internal/ptr"
	"traceknot/internal/tokenize"
)

func attributeBool(attributes map[string]any, key string) (bool, bool) {
	return shared.AttributeBool(attributes, key)
}

func itoa(value int64) string {
	return strconv.FormatInt(value, 10)
}

func (builder *Builder) model(event Event) *string {
	if value, ok := attributeString(event.Attributes, "model"); ok {
		return ptr.String(value)
	}
	return nil
}

func estimateToolTokens(args string, result string, estimator *tokenize.Estimator) (int64, int64, *string) {
	return shared.EstimateToolTokens(args, result, estimator)
}
