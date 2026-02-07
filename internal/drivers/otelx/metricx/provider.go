package metricx

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"

	"github.com/mattdowdell/sandbox/internal/drivers/otelx"
)

// ProviderConfig contains configuration for configuring an OpenTelemetry meter provider.
type ProviderConfig struct {
	// Use HTTP instead of HTTPS for the metric exporter's HTTP connection.
	Insecure bool `koanf:"insecure"`

	// The host and optional port for the metric exporter to connect to.
	Endpoint string `koanf:"endpoint"`

	// The URL path the metric exporter should send requests to.
	Path string `koanf:"path" default:"/v1/metrics"`
}

// MeterProviderShutdown provides a dedicated type for the meter provider shutdown function.
type ProviderShutdown func(context.Context) error

// SetupMeterProviderFromConfig calls SetupMeterProvider with the given configuration.
func SetupMeterProviderFromConfig(
	ctx context.Context,
	conf ProviderConfig,
) (ProviderShutdown, error) {
	return SetupMeterProvider(ctx, conf.Insecure, conf.Endpoint, conf.Path)
}

// SetupMeterProvider creates a new [metric.MeterProvider] and sets it as the default using
// [otel.SetMeterProvider]. The returned function should be called when the process exits to publish
// any lingering metrics.
//
// [metric.MeterProvider]: https://pkg.go.dev/go.opentelemetry.io/otel/metric#MeterProvider
// [otel.SetMeterProvider]: https://pkg.go.dev/go.opentelemetry.io/otel#SetMeterProvider
func SetupMeterProvider(
	ctx context.Context,
	insecure bool,
	endpoint string,
	path string,
) (ProviderShutdown, error) {
	opts := []otlpmetrichttp.Option{
		otlpmetrichttp.WithEndpoint(endpoint),
		otlpmetrichttp.WithURLPath(path),
	}
	if insecure {
		opts = append(opts, otlpmetrichttp.WithInsecure())
	}

	exporter, err := otlpmetrichttp.New(ctx, opts...)
	if err != nil {
		return nil, err
	}

	res, err := otelx.NewResource()
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
