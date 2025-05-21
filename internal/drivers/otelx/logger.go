package otelx

import (
	"context"

	"go.opentelemetry.io/contrib/processors/baggagecopy"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

// LoggerProviderConfig contains configuration for configuring an OpenTelemetry logger provider.
type LoggerProviderConfig struct {
	// The URL of the OTLP HTTP endpoint to export logs to.
	Endpoint string `koanf:"endpoint"`
}

// LoggerProviderShutdown provides a dedicated type for the logger provider shutdown function.
type LoggerProviderShutdown func(context.Context) error

// SetupLoggerProviderFromConfig calls NewLoggerProvider with the given configuration.
func SetupLoggerProviderFromConfig(
	ctx context.Context,
	conf LoggerProviderConfig,
	filter baggagecopy.Filter,
) (LoggerProviderShutdown, error) {
	return SetupLoggerProvider(ctx, conf.Endpoint, filter)
}

// SetupLoggerProvider creates a new [log.LoggerProvider] and sets it as the default using
// [global.SetLoggerProvider]. The returned function should be called when the process exits to
// export any lingering logs.
//
// [log.LoggerProvider]: https://pkg.go.dev/go.opentelemetry.io/otel/log#LoggerProvider
// [global.SetLoggerProvider]: https://pkg.go.dev/go.opentelemetry.io/otel/log/global#SetLoggerProvider
func SetupLoggerProvider(
	ctx context.Context,
	endpoint string,
	filter baggagecopy.Filter,
) (LoggerProviderShutdown, error) {
	exporter, err := otlploghttp.New(ctx, otlploghttp.WithEndpointURL(endpoint))
	if err != nil {
		return nil, err
	}

	res, err := newResource()
	if err != nil {
		return nil, err
	}

	provider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(baggagecopy.NewLogProcessor(filter)),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
		sdklog.WithResource(res),
	)
	global.SetLoggerProvider(provider)

	return provider.Shutdown, nil
}
