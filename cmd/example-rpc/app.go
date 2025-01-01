package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/mattdowdell/sandbox/internal/drivers/logging"
	"github.com/mattdowdell/sandbox/internal/drivers/otelx"
	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver"
)

type AppConfig struct {
	ShutdownTimeout time.Duration `koanf:"app.shutdowntimeout" default:"30s"`
}

// ...
type App struct {
	shutdownTimeout time.Duration
	logger          *slog.Logger
	server          *rpcserver.Server
	tpShutdown      otelx.TracerProviderShutdown
	mpShutdown      otelx.MeterProviderShutdown
}

// ...
func NewApp(
	config AppConfig,
	logger *slog.Logger,
	server *rpcserver.Server,
	tpShutdown otelx.TracerProviderShutdown,
	mpShutdown otelx.MeterProviderShutdown,
) *App {
	return &App{
		shutdownTimeout: config.ShutdownTimeout,
		logger:          logger,
		server:          server,
		tpShutdown:      tpShutdown,
		mpShutdown:      mpShutdown,
	}
}

// ...
func (a *App) Start(ctx context.Context, stop context.CancelFunc) {
	a.logger.InfoContext(ctx, "starting")

	go func() {
		if err := a.server.Start(); err != nil {
			a.logger.ErrorContext(ctx, "failed to start server", logging.Error(err))
		}

		stop()
	}()
}

// ...
func (a *App) Shutdown(ctx context.Context) {
	a.logger.InfoContext(ctx, "stopping")

	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.WarnContext(ctx, "failed to shutdown server", logging.Error(err))
	}

	if err := a.tpShutdown(ctx); err != nil {
		a.logger.WarnContext(ctx, "failed to shutdown tracer provider", logging.Error(err))
	}

	if err := a.mpShutdown(ctx); err != nil {
		a.logger.WarnContext(ctx, "failed to shutdown meter provider", logging.Error(err))
	}
}
