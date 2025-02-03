package slogx

import (
	"log/slog"
	"runtime/debug"

	"go.opentelemetry.io/otel/trace"
)

// ...
func Err(err error) slog.Attr {
	return slog.String("error", err.Error())
}

// ...
func Panic(val any) slog.Attr {
	return slog.Any("panic", val)
}

// ...
func Stacktrace() slog.Attr {
	return slog.String("stacktrace", string(debug.Stack()))
}

// ...
func TraceID(span trace.Span) slog.Attr {
	return slog.String("trace_id", span.SpanContext().TraceID().String())
}

// ...
func RPCService(service string) slog.Attr {
	return slog.String("rpc_service", service)
}

// ...
func RPCMethod(method string) slog.Attr {
	return slog.String("rpc_method", method)
}
