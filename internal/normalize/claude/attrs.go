package claude

import (
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	logspb "go.opentelemetry.io/proto/otlp/logs/v1"

	"traceknot/internal/normalize/shared"
)

func attributeString(attributes map[string]any, key string) (string, bool) {
	return shared.AttributeString(attributes, key)
}

func attributeInt(attributes map[string]any, key string) (int64, bool) {
	return shared.AttributeInt(attributes, key)
}

func attributeBool(attributes map[string]any, key string) (bool, bool) {
	return shared.AttributeBool(attributes, key)
}

func timestampMillis(record *logspb.LogRecord) int64 {
	return shared.TimestampMillis(record)
}

func otlpAttributes(attributes []*commonpb.KeyValue) map[string]any {
	return shared.OTLPAttributes(attributes)
}

func otlpValue(value *commonpb.AnyValue) any {
	return shared.OTLPValue(value)
}
