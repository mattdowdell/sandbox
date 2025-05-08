package slogx

import (
	"fmt"
	"log/slog"
	"runtime/debug"

	"go.opentelemetry.io/otel/trace"
)

// ...
//
// TODO: document how to redact secrets.
func Config(conf any) slog.Attr {
	return slog.Any("config", conf)
}

// ...
func Err(err error) slog.Attr {
	return slog.String("error", err.Error())
}

// ...
func HealthStatus(status fmt.Stringer) slog.Attr {
	return slog.String("health_status", status.String())
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
func RPCSystem(system string) slog.Attr {
	return slog.String("rpc_system", system)
}

// ...
func RPCService(service string) slog.Attr {
	return slog.String("rpc_service", service)
}

// ...
func RPCMethod(method string) slog.Attr {
	return slog.String("rpc_method", method)
}
