package claude

import (
	"encoding/hex"
	"time"

	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
)

type SpanEvent struct {
	Name       string
	Attributes map[string]any
}

type Span struct {
	SessionID    string
	SpanID       string
	ParentSpanID string
	Name         string
	TimestampMs  int64
	Attributes   map[string]any
	Events       []SpanEvent
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
				attributes := otlpAttributes(span.Attributes)
				sessionID, _ := attributeString(attributes, "session.id")
				spans = append(spans, Span{
					SessionID:    sessionID,
					SpanID:       hex.EncodeToString(span.SpanId),
					ParentSpanID: hex.EncodeToString(span.ParentSpanId),
					Name:         span.Name,
					TimestampMs:  int64(span.StartTimeUnixNano / uint64(time.Millisecond)),
					Attributes:   attributes,
					Events:       extractSpanEvents(span.Events),
				})
			}
		}
	}
	return spans
}

func extractSpanEvents(events []*tracepb.Span_Event) []SpanEvent {
	output := make([]SpanEvent, 0, len(events))
	for _, event := range events {
		if event == nil {
			continue
		}
		output = append(output, SpanEvent{
			Name:       event.Name,
			Attributes: otlpAttributes(event.Attributes),
		})
	}
	return output
}
