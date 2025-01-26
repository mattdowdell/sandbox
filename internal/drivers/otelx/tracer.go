package otelx

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/trace"
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

	provider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(provider)

	return provider.Shutdown, nil
}
