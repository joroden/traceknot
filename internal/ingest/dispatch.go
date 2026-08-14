package ingest

import (
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	logspb "go.opentelemetry.io/proto/otlp/logs/v1"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
)

func logsServiceName(data *logspb.LogsData) string {
	for _, resourceLogs := range data.ResourceLogs {
		if resourceLogs == nil || resourceLogs.Resource == nil {
			continue
		}
		if name, ok := serviceName(resourceLogs.Resource.Attributes); ok {
			return name
		}
	}
	return ""
}

func tracesServiceName(data *tracepb.TracesData) string {
	for _, resourceSpans := range data.ResourceSpans {
		if resourceSpans == nil || resourceSpans.Resource == nil {
			continue
		}
		if name, ok := serviceName(resourceSpans.Resource.Attributes); ok {
			return name
		}
	}
	return ""
}

func serviceName(attributes []*commonpb.KeyValue) (string, bool) {
	for _, attribute := range attributes {
		if attribute == nil || attribute.Key != "service.name" {
			continue
		}
		if stringValue, ok := attribute.Value.GetValue().(*commonpb.AnyValue_StringValue); ok {
			return stringValue.StringValue, stringValue.StringValue != ""
		}
	}
	return "", false
}
