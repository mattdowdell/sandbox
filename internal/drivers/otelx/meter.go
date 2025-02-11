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

// MeterProviderShutdown provides a dedicated type for the meter provider shutdown function.
type MeterProviderShutdown func(context.Context) error

// NewMeterProviderFromConfig calls NewMeterProvider with the given configuration.
func NewMeterProviderFromConfig(
	ctx context.Context,
	conf MeterProviderConfig,
) (MeterProviderShutdown, error) {
	return NewMeterProvider(ctx, conf.Endpoint)
}

// NewMeterProvider creates a new [metric.MeterProvider] and sets it as the default using
// [otel.SetMeterProvider]. The returned function should be called when the process exits to publish
// any lingering metrics.
//
// [metric.MeterProvider]: https://pkg.go.dev/go.opentelemetry.io/otel/metric#MeterProvider
// [otel.SetMeterProvider]: https://pkg.go.dev/go.opentelemetry.io/otel#SetMeterProvider
func NewMeterProvider(
	ctx context.Context,
	endpoint string,
) (MeterProviderShutdown, error) {
	exporter, err := otlpmetrichttp.New(ctx, otlpmetrichttp.WithEndpointURL(endpoint))
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

// Meter wraps [otel.Meter] to provide a [metric.Meter] with the package name and version
// automatically set based on the direct caller. It is advised to cache the result when possible to
// avoid computing the caller's package details unnecessarily.
//
// [otel.Meter]: https://pkg.go.dev/go.opentelemetry.io/otel#Meter
// [metric.Meter]: https://pkg.go.dev/go.opentelemetry.io/otel/metric#Meter
func Meter() metric.Meter {
	pkg := packageName(1 /*skip*/)
	ver := packageVersion(pkg)

	return otel.Meter(pkg, metric.WithInstrumentationVersion(ver))
}
