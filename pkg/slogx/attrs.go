package slogx

import (
	"fmt"
	"log/slog"
	"runtime/debug"

	"go.opentelemetry.io/otel/trace"
)

// Keys to use in log attributes.
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

	ResourceIDKey   = "resource_id"
	ResourceNameKey = "resource_name"
)

// Config returns the given value with the "config" key.
func Config(conf any) slog.Attr {
	return slog.Any(ConfigKey, conf)
}

// Err returns a string representation of the error with the "error" key.
func Err(err error) slog.Attr {
	return slog.String(ErrorKey, err.Error())
}

// HealthStatus returns the given value with the "health_status" key.
func HealthStatus(status fmt.Stringer) slog.Attr {
	return slog.String(HealthStatusKey, status.String())
}

// Panic returns the given value with the "panic" key.
func Panic(val any) slog.Attr {
	return slog.Any(PanicKey, val)
}

// Stacktrace returns the current stacktrace with the "stacktrace" key. It should primarily be used
// during panic recovery.
func Stacktrace() slog.Attr {
	return slog.String(StacktraceKey, string(debug.Stack()))
}

// TraceID returns the trace ID from the given span with the "trace_id" key.
func TraceID(span trace.Span) slog.Attr {
	return slog.String(TraceIDKey, span.SpanContext().TraceID().String())
}

// TraceID returns the span ID from the given span with the "span_id" key.
func SpanID(span trace.Span) slog.Attr {
	return slog.String(SpanIDKey, span.SpanContext().SpanID().String())
}

// Sampled returns the sampled status from the given span with the "sampled" key.
func Sampled(span trace.Span) slog.Attr {
	return slog.Bool(SampledKey, span.SpanContext().IsSampled())
}

// Subject returns the given value with the "subject" key. It should be used to capture the "sub"
// claim of a JWT.
func Subject(subject string) slog.Attr {
	return slog.String(SubjectKey, subject)
}

// RPCSystem returns the given value with the "rpc_system" key. It should be used to capture the
// protocol a RPC request uses, e.g. gRPC, ConnectRPC, etc.
func RPCSystem(system string) slog.Attr {
	return slog.String(RPCSystemKey, system)
}

// RPCService returns the given value with the "rpc_service" key.
func RPCService(service string) slog.Attr {
	return slog.String(RPCServiceKey, service)
}

// RPCMethod returns the given value with the "rpc_method" key.
func RPCMethod(method string) slog.Attr {
	return slog.String(RPCMethodKey, method)
}

// ...
func ResourceID(id fmt.Stringer) slog.Attr {
	return slog.String(ResourceIDKey, id.String())
}

// ...
func ResourceName(name string) slog.Attr {
	return slog.String(ResourceNameKey, name)
}
