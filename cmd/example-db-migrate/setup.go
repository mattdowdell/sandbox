package main

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/contrib/processors/baggagecopy"

	"github.com/mattdowdell/sandbox/internal/drivers/config"
	"github.com/mattdowdell/sandbox/internal/drivers/config/flagoptions"
	"github.com/mattdowdell/sandbox/internal/drivers/logging"
	"github.com/mattdowdell/sandbox/internal/drivers/otelx"
	"github.com/mattdowdell/sandbox/internal/drivers/pgsql"
)

func SetupApp(ctx context.Context) (*App, error) {
	options := flagoptions.New()

	conf, err := config.Load[Config](options)
	if err != nil {
		return nil, err
	}

	tpShutdown, mpShutdown, lpShutdown, logger, err := setupObservability(ctx, conf)
	if err != nil {
		return nil, err
	}

	db, err := pgsql.NewFromConfig(ctx, conf.Database)
	if err != nil {
		return nil, err
	}

	return NewApp(conf, logger, db, tpShutdown, mpShutdown, lpShutdown), nil
}

func setupObservability(
	ctx context.Context,
	conf *Config,
) (
	otelx.TracerProviderShutdown,
	otelx.MeterProviderShutdown,
	otelx.LoggerProviderShutdown,
	*slog.Logger,
	error,
) {
	filter := baggagecopy.AllowAllMembers

	tpShutdown, err := otelx.SetupTracerProviderFromConfig(ctx, conf.TracerProvider, filter)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	mpShutdown, err := otelx.SetupMeterProviderFromConfig(ctx, conf.MeterProvider)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	lpShutdown, err := otelx.SetupLoggerProviderFromConfig(ctx, conf.LoggerProvider, filter)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	extractor := otelx.NewExtractor(otelx.WithSpanID(true), otelx.WithSampled(true))
	logger := logging.NewAsDefaultFromConfig(conf.Logging, logging.WithExtractors(extractor))

	return tpShutdown, mpShutdown, lpShutdown, logger, nil
}
