package main

import (
	"context"
	"database/sql"
	"log/slog"

	"go.opentelemetry.io/contrib/instrumentation/runtime"

	"github.com/mattdowdell/sandbox/internal/drivers/otelx"
	"github.com/mattdowdell/sandbox/internal/drivers/otelx/logx"
	"github.com/mattdowdell/sandbox/internal/drivers/otelx/metricx"
	"github.com/mattdowdell/sandbox/internal/drivers/otelx/tracex"
	"github.com/mattdowdell/sandbox/pkg/herd"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// ...
type App struct {
	conf       *Config
	logger     *slog.Logger
	db         *sql.DB
	herder     *herd.Herd
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
	herder *herd.Herd,
	tpShutdown tracex.ProviderShutdown,
	mpShutdown metricx.ProviderShutdown,
	lpShutdown logx.ProviderShutdown,
) *App {
	return &App{
		conf:       conf,
		logger:     logger,
		db:         db,
		herder:     herder,
		tpShutdown: tpShutdown,
		mpShutdown: mpShutdown,
		lpShutdown: lpShutdown,
	}
}

// ...
//
//nolint:sloglint // little benefit defining reusable keys for a one-off application.
func (a *App) Run(ctx context.Context) error {
	ctx, span := otelx.Tracer().Start(ctx, "DB Migrate")
	defer span.End()

	a.logger.InfoContext(ctx, "starting", slogx.Config(a.conf))

	if err := runtime.Start(); err != nil {
		span.RecordError(err)
		a.logger.ErrorContext(ctx, "failed to start runtime metrics", slogx.Err(err))
		return err
	}

	result, err := a.herder.Migrate(ctx, a.db)
	if err != nil {
		span.RecordError(err)
		a.logger.ErrorContext(ctx, "failed to migrate db", slogx.Err(err))
		return err
	}

	a.logger.InfoContext(
		ctx,
		"migrated db",
		slog.Int64("before", result.Before),
		slog.Int64("after", result.After),
	)
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
