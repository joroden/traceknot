package copilot

import (
	"encoding/hex"
	"time"

	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
)

type Span struct {
	TraceID      string
	SpanID       string
	ParentSpanID string
	TimestampMs  int64
	Attributes   map[string]any
}

func ExtractSpans(request *tracepb.TracesData) []Span {
	var spans []Span
	for _, resourceSpans := range request.ResourceSpans {
		if resourceSpans == nil {
			continue
		}
		for _, scopeSpans := range resourceSpans.ScopeSpans {
			if scopeSpans == nil {
				continue
			}
			for _, span := range scopeSpans.Spans {
				if span == nil {
					continue
				}
				spans = append(spans, Span{
					TraceID:      hex.EncodeToString(span.TraceId),
					SpanID:       hex.EncodeToString(span.SpanId),
					ParentSpanID: hex.EncodeToString(span.ParentSpanId),
					TimestampMs:  int64(span.StartTimeUnixNano / uint64(time.Millisecond)),
					Attributes:   otlpAttributes(span.Attributes),
				})
			}
		}
	}
	return spans
}
