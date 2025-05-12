package main

import (
	"context"
	"log/slog"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/runtime"

	"github.com/mattdowdell/sandbox/internal/drivers/otelx"
	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// ...
type AppConfig struct {
	// ...
	ShutdownTimeout time.Duration `koanf:"shutdowntimeout" default:"30s"`
}

// ...
type App struct {
	conf            Config
	shutdownTimeout time.Duration
	logger          *slog.Logger
	server          *rpcserver.Server
	tpShutdown      otelx.TracerProviderShutdown
	mpShutdown      otelx.MeterProviderShutdown
	lpShutdown      otelx.LoggerProviderShutdown
}

// ...
//
//nolint:gocritic // config is large(ish), but this is called once
func NewApp(
	conf Config,
	logger *slog.Logger,
	server *rpcserver.Server,
	tpShutdown otelx.TracerProviderShutdown,
	mpShutdown otelx.MeterProviderShutdown,
	lpShutdown otelx.LoggerProviderShutdown,
) *App {
	return &App{
		conf:            conf,
		shutdownTimeout: conf.App.ShutdownTimeout,
		logger:          logger,
		server:          server,
		tpShutdown:      tpShutdown,
		mpShutdown:      mpShutdown,
		lpShutdown:      lpShutdown,
	}
}

// ...
func (a *App) Start(ctx context.Context, stop context.CancelFunc) {
	a.logger.InfoContext(ctx, "starting", slogx.Config(a.conf))

	if err := runtime.Start(); err != nil {
		a.logger.ErrorContext(ctx, "failed to start runtime metrics", slogx.Err(err))
		stop()

		return
	}

	go func() {
		if err := a.server.Start(); err != nil {
			a.logger.ErrorContext(ctx, "failed to start server", slogx.Err(err))
		}

		stop()
	}()
}

// ...
func (a *App) Shutdown(ctx context.Context) {
	a.logger.InfoContext(ctx, "stopping")

	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.WarnContext(ctx, "failed to shutdown server", slogx.Err(err))
	}

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
