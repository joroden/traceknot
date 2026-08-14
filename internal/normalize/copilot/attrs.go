package copilot

import (
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"

	"traceknot/internal/normalize/shared"
)

func attributeString(attributes map[string]any, key string) (string, bool) {
	return shared.AttributeString(attributes, key)
}

func attributeInt(attributes map[string]any, key string) (int64, bool) {
	return shared.AttributeInt(attributes, key)
}

func otlpAttributes(attributes []*commonpb.KeyValue) map[string]any {
	return shared.OTLPAttributes(attributes)
}

func otlpValue(value *commonpb.AnyValue) any {
	return shared.OTLPValue(value)
}
