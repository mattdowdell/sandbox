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

// ...
type TracerProviderShutdown func(context.Context) error

// ...
func NewTracerProvider(
	ctx context.Context,
	conf TracerProviderConfig,
) (TracerProviderShutdown, error) {
	exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(conf.Endpoint))
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

// ...
func Tracer() trace.Tracer {
	pkg := packageName(1 /*skip*/)
	ver := packageVersion(pkg)

	return otel.Tracer(pkg, trace.WithInstrumentationVersion(ver))
}
