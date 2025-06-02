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
	// Use HTTP instead of HTTPS for the log exporter's HTTP connection.
	Insecure bool `koanf:"insecure"`

	// The host and optional port for the log exporter to connect to.
	Endpoint string `koanf:"endpoint"`

	// The URL path the log exporter should send requests to.
	Path string `koanf:"path" default:"/v1/logs"`
}

// LoggerProviderShutdown provides a dedicated type for the logger provider shutdown function.
type LoggerProviderShutdown func(context.Context) error

// SetupLoggerProviderFromConfig calls NewLoggerProvider with the given configuration.
func SetupLoggerProviderFromConfig(
	ctx context.Context,
	conf LoggerProviderConfig,
	filter baggagecopy.Filter,
) (LoggerProviderShutdown, error) {
	return SetupLoggerProvider(ctx, conf.Insecure, conf.Endpoint, conf.Path, filter)
}

// SetupLoggerProvider creates a new [log.LoggerProvider] and sets it as the default using
// [global.SetLoggerProvider]. The returned function should be called when the process exits to
// export any lingering logs.
//
// [log.LoggerProvider]: https://pkg.go.dev/go.opentelemetry.io/otel/log#LoggerProvider
// [global.SetLoggerProvider]: https://pkg.go.dev/go.opentelemetry.io/otel/log/global#SetLoggerProvider
func SetupLoggerProvider(
	ctx context.Context,
	insecure bool,
	endpoint string,
	path string,
	filter baggagecopy.Filter,
) (LoggerProviderShutdown, error) {
	opts := []otlploghttp.Option{
		otlploghttp.WithEndpoint(endpoint),
		otlploghttp.WithURLPath(path),
	}
	if insecure {
		opts = append(opts, otlploghttp.WithInsecure())
	}

	exporter, err := otlploghttp.New(ctx, opts...)
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
