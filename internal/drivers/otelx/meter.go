package otelx

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/metric"
)

// ...
type MeterProviderConfig struct {
	// ...
	Endpoint string `koanf:"meterprovider.endpoint"`
}

// ...
type MeterProviderShutdown func(context.Context) error

// ...
func NewMeterProvider(
	ctx context.Context,
	conf MeterProviderConfig,
) (MeterProviderShutdown, error) {
	exporter, err := otlpmetrichttp.New(ctx, otlpmetrichttp.WithEndpoint(conf.Endpoint))
	if err != nil {
		return nil, err
	}

	provider := metric.NewMeterProvider(metric.WithReader(metric.NewPeriodicReader(exporter)))
	otel.SetMeterProvider(provider)

	return provider.Shutdown, nil
}
