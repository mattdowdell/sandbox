package slogx

import (
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

// ...
func Err(err error) slog.Attr {
	return slog.String("error", err.Error())
}

// ...
func TraceID(span trace.Span) slog.Attr {
	return slog.String("trace_id", span.SpanContext().TraceID().String())
}
