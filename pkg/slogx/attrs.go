package slogx

import (
	"fmt"
	"log/slog"
	"runtime/debug"

	"go.opentelemetry.io/otel/trace"
)

const (
	ConfigKey       = "config"
	ErrorKey        = "error"
	HealthStatusKey = "health_status"
	PanicKey        = "panic"
	StacktraceKey   = "stacktrace"
	TraceIDKey      = "trace_id"
	SpanIDKey       = "span_id"
	SampledKey      = "sampled"
	SubjectKey      = "subject"
	RPCSystemKey    = "rpc_system"
	RPCServiceKey   = "rpc_service"
	RPCMethodKey    = "rpc_method"
)

// ...
//
// TODO: document how to redact secrets.
func Config(conf any) slog.Attr {
	return slog.Any(ConfigKey, conf)
}

// ...
func Err(err error) slog.Attr {
	return slog.String(ErrorKey, err.Error())
}

// ...
func HealthStatus(status fmt.Stringer) slog.Attr {
	return slog.String(HealthStatusKey, status.String())
}

// ...
func Panic(val any) slog.Attr {
	return slog.Any(PanicKey, val)
}

// ...
func Stacktrace() slog.Attr {
	return slog.String(StacktraceKey, string(debug.Stack()))
}

// ...
func TraceID(span trace.Span) slog.Attr {
	return slog.String(TraceIDKey, span.SpanContext().TraceID().String())
}

// ...
func SpanID(span trace.Span) slog.Attr {
	return slog.String(SpanIDKey, span.SpanContext().SpanID().String())
}

// ...
func Sampled(span trace.Span) slog.Attr {
	return slog.Bool(SampledKey, span.SpanContext().IsSampled())
}

// ...
func Subject(subject string) slog.Attr {
	return slog.String(SubjectKey, subject)
}

// ...
func RPCSystem(system string) slog.Attr {
	return slog.String(RPCSystemKey, system)
}

// ...
func RPCService(service string) slog.Attr {
	return slog.String(RPCServiceKey, service)
}

// ...
func RPCMethod(method string) slog.Attr {
	return slog.String(RPCMethodKey, method)
}
