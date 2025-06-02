package otelx

import (
	"context"

	"go.opentelemetry.io/contrib/processors/baggagecopy"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// TracerProviderConfig contains configuration for configuring an OpenTelemetry tracer provider.
type TracerProviderConfig struct {
	// Use HTTP instead of HTTPS for the trace exporter's HTTP connection.
	Insecure bool `koanf:"insecure"`

	// The host and optional port for the trace exporter to connect to.
	Endpoint string `koanf:"endpoint"`

	// The URL path the trace exporter should send requests to.
	Path string `koanf:"path" default:"/v1/traces"`
}

// TracerProviderShutdown provides a dedicated type for the tracer provider shutdown function.
type TracerProviderShutdown func(context.Context) error

// SetupTracerProviderFromConfig calls SetupTracerProvider with the given configuration.
func SetupTracerProviderFromConfig(
	ctx context.Context,
	conf TracerProviderConfig,
	filter baggagecopy.Filter,
) (TracerProviderShutdown, error) {
	return SetupTracerProvider(ctx, conf.Insecure, conf.Endpoint, conf.Path, filter)
}

// SetupTracerProvider creates a new [trace.TracerProvider] and sets it as the default using
// [otel.SetTracerProvider]. The returned function should be called when the process exits to
// publish any lingering spans.
//
// Propagation of W3C trace contexts and W3C baggage is also configured for both incoming requests
// and outgoing requests when supported.
//
// [trace.TracerProvider]: https://pkg.go.dev/go.opentelemetry.io/otel/trace#TracerProvider
// [otel.SetTracerProvider]: https://pkg.go.dev/go.opentelemetry.io/otel#SetTracerProvider
func SetupTracerProvider(
	ctx context.Context,
	insecure bool,
	endpoint string,
	path string,
	filter baggagecopy.Filter,
) (TracerProviderShutdown, error) {
	opts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithURLPath(path),
	}
	if insecure {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	exporter, err := otlptracehttp.New(ctx, opts...)
	if err != nil {
		return nil, err
	}

	res, err := newResource()
	if err != nil {
		return nil, err
	}

	setupPropagator()

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(baggagecopy.NewSpanProcessor(filter)),
	)
	otel.SetTracerProvider(provider)

	return provider.Shutdown, nil
}
