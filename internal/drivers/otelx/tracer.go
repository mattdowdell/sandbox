package otelx

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// TracerOption implementations customise the behaviour of Tracer.
type TracerOption interface {
	apply(*tracerOpts)
}

type tracerOpts struct {
	provider trace.TracerProvider
}

type tracerProviderOpt struct {
	provider trace.TracerProvider
}

// WithTracerProvider overrides the global TracerProvider. This can be useful for testing the result
// of recording metrics.
func WithTracerProvider(provider trace.TracerProvider) TracerOption {
	return &tracerProviderOpt{
		provider: provider,
	}
}

func (o *tracerProviderOpt) apply(options *tracerOpts) {
	options.provider = o.provider
}

// Tracer provides a [trace.Tracer] with the package name and version automatically set based on the
// direct caller. It is advised to cache the result when possible to avoid computing the caller's
// package details unnecessarily.
//
// [trace.Tracer]: https://pkg.go.dev/go.opentelemetry.io/otel/trace#Tracer
func Tracer(options ...TracerOption) trace.Tracer {
	opts := &tracerOpts{}
	for _, option := range options {
		option.apply(opts)
	}

	provider := opts.provider
	if provider == nil {
		provider = otel.GetTracerProvider()
	}

	pkg := packageName(1 /*skip*/)
	ver := packageVersion(pkg)

	return provider.Tracer(pkg, trace.WithInstrumentationVersion(ver))
}
