package shared

import (
	"fmt"
	"strconv"
	"time"

	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	logspb "go.opentelemetry.io/proto/otlp/logs/v1"

	"traceknot/internal/tokenize"
)

func AttributeString(attributes map[string]any, key string) (string, bool) {
	value, ok := attributes[key]
	if !ok {
		return "", false
	}
	switch typed := value.(type) {
	case string:
		return typed, typed != ""
	case int64:
		return strconv.FormatInt(typed, 10), true
	default:
		return fmt.Sprintf("%v", typed), true
	}
}

func AttributeInt(attributes map[string]any, key string) (int64, bool) {
	value, ok := attributes[key]
	if !ok {
		return 0, false
	}
	switch typed := value.(type) {
	case int64:
		return typed, true
	case float64:
		return int64(typed), true
	case string:
		parsed, err := strconv.ParseInt(typed, 10, 64)
		return parsed, err == nil
	default:
		return 0, false
	}
}

func AttributeBool(attributes map[string]any, key string) (bool, bool) {
	value, ok := attributes[key]
	if !ok {
		return false, false
	}
	switch typed := value.(type) {
	case bool:
		return typed, true
	case string:
		return typed == "true", true
	default:
		return false, false
	}
}

func TimestampMillis(record *logspb.LogRecord) int64 {
	if record.TimeUnixNano != 0 {
		return int64(record.TimeUnixNano / uint64(time.Millisecond))
	}
	if record.ObservedTimeUnixNano != 0 {
		return int64(record.ObservedTimeUnixNano / uint64(time.Millisecond))
	}
	return 0
}

func OTLPAttributes(attributes []*commonpb.KeyValue) map[string]any {
	output := make(map[string]any, len(attributes))
	for _, attribute := range attributes {
		if attribute == nil {
			continue
		}
		output[attribute.Key] = OTLPValue(attribute.Value)
	}
	return output
}

func OTLPValue(value *commonpb.AnyValue) any {
	if value == nil {
		return nil
	}
	switch typed := value.Value.(type) {
	case *commonpb.AnyValue_StringValue:
		return typed.StringValue
	case *commonpb.AnyValue_IntValue:
		return typed.IntValue
	case *commonpb.AnyValue_BoolValue:
		return typed.BoolValue
	case *commonpb.AnyValue_DoubleValue:
		return typed.DoubleValue
	case *commonpb.AnyValue_ArrayValue:
		items := make([]any, 0, len(typed.ArrayValue.Values))
		for _, item := range typed.ArrayValue.Values {
			items = append(items, OTLPValue(item))
		}
		return items
	case *commonpb.AnyValue_KvlistValue:
		return OTLPAttributes(typed.KvlistValue.Values)
	default:
		return nil
	}
}

func EstimateToolTokens(args string, result string, estimator *tokenize.Estimator) (int64, int64, *string) {
	var estimatedInput int64
	var estimatedOutput int64
	if result != "" {
		estimatedInput = estimator.Count(result)
	}
	if args != "" {
		estimatedOutput = estimator.Count(args)
	}
	if estimatedInput == 0 && estimatedOutput == 0 {
		return 0, 0, nil
	}
	method := estimator.MethodName()
	return estimatedInput, estimatedOutput, &method
}
