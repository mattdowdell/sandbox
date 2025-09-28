package otelt

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// MustSpan wraps NewSpan, panicking if an error occurs. Errors can only occur if an invalid trace
// ID or span ID is provided.
func MustSpan(ctx context.Context, traceID, spanID string, sampled bool) context.Context {
	ctx, err := NewSpan(ctx, traceID, spanID, sampled)
	if err != nil {
		panic(err)
	}

	return ctx
}

// NewSpan creates a [trace.Span] using the give trace ID, span ID and sampled value and adds it to
// the given context. It can be used to produce a span with a deterministic value.
//
// [trace.Span]: https://pkg.go.dev/go.opentelemetry.io/otel/trace#Span
func NewSpan(ctx context.Context, traceID, spanID string, sampled bool) (context.Context, error) {
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
