package otelx

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// ...
type TracerProviderConfig struct {
	// ...
	Endpoint string `koanf:"endpoint"`
}

// TracerProviderShutdown provides a dedicated type for the tracer provider shutdown function.
type TracerProviderShutdown func(context.Context) error

// NewMeterProviderFromConfig calls NewMeterProvider with the given configuration.
func NewTracerProviderFromConfig(
	ctx context.Context,
	conf TracerProviderConfig,
) (TracerProviderShutdown, error) {
	return NewTracerProvider(ctx, conf.Endpoint)
}

// NewTracerProvider creates a new [trace.TracerProvider] and sets it as the default using
// [otel.SetTracerProvider]. The returned function should be called when the process exits to
// publish any lingering spans.
//
// [trace.TracerProvider]: https://pkg.go.dev/go.opentelemetry.io/otel/trace#TracerProvider
// [otel.SetTracerProvider]: https://pkg.go.dev/go.opentelemetry.io/otel#SetTracerProvider
func NewTracerProvider(
	ctx context.Context,
	endpoint string,
) (TracerProviderShutdown, error) {
	exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(endpoint))
	if err != nil {
		return nil, err
	}

	res, err := newResource()
	if err != nil {
		return nil, err
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(provider)

	return provider.Shutdown, nil
}

// Tracer wraps [otel.Tracer] to provide a [trace.Tracer] with the package name and version
// automatically set based on the direct caller. It is advised to cache the result when possible to
// avoid computing the caller's package details unnecessarily.
//
// [otel.Tracer]: https://pkg.go.dev/go.opentelemetry.io/otel#Tracer
// [trace.Tracer]: https://pkg.go.dev/go.opentelemetry.io/otel/trace#Tracer
func Tracer() trace.Tracer {
	pkg := packageName(1 /*skip*/)
	ver := packageVersion(pkg)

	return otel.Tracer(pkg, trace.WithInstrumentationVersion(ver))
}
