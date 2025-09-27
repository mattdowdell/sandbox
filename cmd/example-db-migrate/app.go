package main

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/pressly/goose/v3"
	"go.opentelemetry.io/contrib/instrumentation/runtime"

	"github.com/mattdowdell/sandbox/internal/adapters/datastore"
	"github.com/mattdowdell/sandbox/internal/drivers/otelx/logx"
	"github.com/mattdowdell/sandbox/internal/drivers/otelx/metricx"
	"github.com/mattdowdell/sandbox/internal/drivers/otelx/tracex"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// ...
type App struct {
	conf       *Config
	logger     *slog.Logger
	db         *sql.DB
	tpShutdown tracex.ProviderShutdown
	mpShutdown metricx.ProviderShutdown
	lpShutdown logx.ProviderShutdown
}

// ...
//
//nolint:gocritic // config is large(ish), but this is called once
func NewApp(
	conf *Config,
	logger *slog.Logger,
	db *sql.DB,
	tpShutdown tracex.ProviderShutdown,
	mpShutdown metricx.ProviderShutdown,
	lpShutdown logx.ProviderShutdown,
) *App {
	return &App{
		conf:       conf,
		logger:     logger,
		db:         db,
		tpShutdown: tpShutdown,
		mpShutdown: mpShutdown,
		lpShutdown: lpShutdown,
	}
}

// ...
func (a *App) Run(ctx context.Context) error {
	a.logger.InfoContext(ctx, "starting", slogx.Config(a.conf))

	if err := runtime.Start(); err != nil {
		a.logger.ErrorContext(ctx, "failed to start runtime metrics", slogx.Err(err))
		return err
	}

	goose.SetBaseFS(datastore.MigrationFS)

	if err := goose.SetDialect("postgres"); err != nil {
		a.logger.ErrorContext(ctx, "failed to configure dialect", slogx.Err(err))
		return err
	}

	if err := goose.UpContext(ctx, a.db, "migrations"); err != nil {
		a.logger.ErrorContext(ctx, "failed to migrate db", slogx.Err(err))
		return err
	}

	return nil
}

// ...
func (a *App) Shutdown(ctx context.Context) {
	a.logger.InfoContext(ctx, "stopping")

	if err := a.tpShutdown(ctx); err != nil {
		a.logger.WarnContext(ctx, "failed to shutdown tracer provider", slogx.Err(err))
	}

	if err := a.mpShutdown(ctx); err != nil {
		a.logger.WarnContext(ctx, "failed to shutdown meter provider", slogx.Err(err))
	}

	if err := a.lpShutdown(ctx); err != nil {
		a.logger.WarnContext(ctx, "failed to shutdown logger provider", slogx.Err(err))
	}
}
