package otelx

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// Span creates a span using the give trace ID, span ID and sampled value, panicking on an invalid
// trace ID or span ID. It can be used to produce a span with a deterministic value and is intended
// for test code only.
func MustSpan(ctx context.Context, traceID, spanID string, sampled bool) context.Context {
	ctx, err := Span(ctx, traceID, spanID, sampled)
	if err != nil {
		panic(err)
	}

	return ctx
}

// Span creates a span using the give trace ID, span ID and sampled value. It can be used to produce
// a span with a deterministic value and is intended for test code only.
func Span(ctx context.Context, traceID, spanID string, sampled bool) (context.Context, error) {
	tID, err := trace.TraceIDFromHex(traceID)
	if err != nil {
		return nil, err
	}

	sID, err := trace.SpanIDFromHex(spanID)
	if err != nil {
		return nil, err
	}

	var flags trace.TraceFlags
	if sampled {
		flags = trace.FlagsSampled
	}

	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    tID,
		SpanID:     sID,
		TraceFlags: flags,
	})

	return trace.ContextWithRemoteSpanContext(ctx, spanCtx), nil
}
