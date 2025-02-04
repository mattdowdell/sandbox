package otelx

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// ...
type MeterProviderConfig struct {
	// ...
	Endpoint string `koanf:"endpoint"`
}

// ...
type MeterProviderShutdown func(context.Context) error

// ...
func NewMeterProvider(
	ctx context.Context,
	conf MeterProviderConfig,
) (MeterProviderShutdown, error) {
	exporter, err := otlpmetrichttp.New(ctx, otlpmetrichttp.WithEndpointURL(conf.Endpoint))
	if err != nil {
		return nil, err
	}

	res, err := newResource()
	if err != nil {
		return nil, err
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(provider)

	return provider.Shutdown, nil
}

// ...
func Meter() metric.Meter {
	pkg := packageName(1 /*skip*/)
	ver := packageVersion(pkg)

	return otel.Meter(pkg, metric.WithInstrumentationVersion(ver))
}
