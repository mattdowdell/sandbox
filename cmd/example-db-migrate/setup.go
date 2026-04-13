package main

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/contrib/processors/baggagecopy"

	"github.com/mattdowdell/sandbox/internal/adapters/datastore"
	"github.com/mattdowdell/sandbox/internal/drivers/config"
	"github.com/mattdowdell/sandbox/internal/drivers/config/flagoptions"
	"github.com/mattdowdell/sandbox/internal/drivers/logging"
	"github.com/mattdowdell/sandbox/internal/drivers/otelx"
	"github.com/mattdowdell/sandbox/internal/drivers/otelx/logx"
	"github.com/mattdowdell/sandbox/internal/drivers/otelx/metricx"
	"github.com/mattdowdell/sandbox/internal/drivers/otelx/tracex"
	"github.com/mattdowdell/sandbox/internal/drivers/pgsql"
	"github.com/mattdowdell/sandbox/pkg/herd"
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

	db, _, err := pgsql.NewFromConfig(ctx, conf.Database)
	if err != nil {
		return nil, err
	}

	migrations, err := herd.CollectFileMigrations(datastore.MigrationFS)
	if err != nil {
		return nil, err
	}

	herder, err := herd.New(migrations, herd.WithTargetVersion(conf.App.TargetVersion))
	if err != nil {
		return nil, err
	}

	return NewApp(conf, logger, db, herder, tpShutdown, mpShutdown, lpShutdown), nil
}

//nolint:gocritic // result types are differentiated by package name.
func setupObservability(
	ctx context.Context,
	conf *Config,
) (
	tracex.ProviderShutdown,
	metricx.ProviderShutdown,
	logx.ProviderShutdown,
	*slog.Logger,
	error,
) {
	filter := baggagecopy.AllowAllMembers

	tpShutdown, err := tracex.SetupTracerProviderFromConfig(ctx, conf.TracerProvider, filter)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	mpShutdown, err := metricx.SetupMeterProviderFromConfig(ctx, conf.MeterProvider)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	lpShutdown, err := logx.SetupLoggerProviderFromConfig(ctx, conf.LoggerProvider, filter)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	extractor := otelx.NewExtractor(otelx.WithSpanID(true), otelx.WithSampled(true))
	logger := logging.NewAsDefaultFromConfig(conf.Logging, logging.WithExtractors(extractor))

	return tpShutdown, mpShutdown, lpShutdown, logger, nil
}
