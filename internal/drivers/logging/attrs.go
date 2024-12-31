package logging

import (
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

// ...
//
// TODO: move this to a more clean arch friendly package, maybe under internal/domain/?
func Error(err error) slog.Attr {
	return slog.String("error", err.Error())
}

// ...
func TraceID(span trace.Span) slog.Attr {
	return slog.String("trace_id", span.SpanContext().TraceID().String())
}
