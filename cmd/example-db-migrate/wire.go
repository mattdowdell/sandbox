//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/google/wire"

	"github.com/mattdowdell/sandbox/internal/drivers/config"
	"github.com/mattdowdell/sandbox/internal/drivers/config/flagoptions"
	"github.com/mattdowdell/sandbox/internal/drivers/logging"
	"github.com/mattdowdell/sandbox/internal/drivers/otelx"
	"github.com/mattdowdell/sandbox/internal/drivers/pgsql"
)

func ProvideApp(ctx context.Context) (*App, error) {
	wire.Build(
		// config
		flagoptions.New,
		config.New,
		LoadConfig,
		wire.FieldsOf(
			new(Config),
			"Database",
			"LoggerProvider",
			"Logging",
			"MeterProvider",
			"TracerProvider",
		),
		// observability
		baggageFilter,
		otelx.SetupTracerProviderFromConfig,
		otelx.SetupMeterProviderFromConfig,
		otelx.SetupLoggerProviderFromConfig,
		loggerOptions,
		logging.NewAsDefaultFromConfig,
		// providers
		pgsql.NewFromConfig,
		// app
		NewApp,
	)

	return &App{}, nil
}
