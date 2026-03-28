package main

import (
	"context"
	"log/slog"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/runtime"

	"github.com/mattdowdell/sandbox/internal/drivers/config"
	"github.com/mattdowdell/sandbox/internal/drivers/otelx/logx"
	"github.com/mattdowdell/sandbox/internal/drivers/otelx/metricx"
	"github.com/mattdowdell/sandbox/internal/drivers/otelx/tracex"
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
	conf            *Config
	shutdownTimeout time.Duration
	logger          *slog.Logger
	loader *config.Loader
	server          *rpcserver.Server
	tpShutdown      tracex.ProviderShutdown
	mpShutdown      metricx.ProviderShutdown
	lpShutdown      logx.ProviderShutdown
}

// ...
//
//nolint:gocritic // config is large(ish), but this is called once
func NewApp(
	conf *Config,
	logger *slog.Logger,
	loader *config.Loader,
	server *rpcserver.Server,
	tpShutdown tracex.ProviderShutdown,
	mpShutdown metricx.ProviderShutdown,
	lpShutdown logx.ProviderShutdown,
) *App {
	return &App{
		conf:            conf,
		shutdownTimeout: conf.App.ShutdownTimeout,
		logger:          logger,
		loader: loader,
		server:          server,
		tpShutdown:      tpShutdown,
		mpShutdown:      mpShutdown,
		lpShutdown:      lpShutdown,
	}
}

// ...
func (a *App) Start(ctx context.Context, stop context.CancelFunc) {
	encoded, err := config.Encode(a.conf, "_")
	if err != nil {
		a.logger.ErrorContext(ctx, "failed to encode config")
		stop()

		return
	}

	a.logger.InfoContext(ctx, "starting", slogx.Config(encoded))

	if err := runtime.Start(); err != nil {
		a.logger.ErrorContext(ctx, "failed to start runtime metrics", slogx.Err(err))
		stop()

		return
	}

	go func() {
		if err := a.server.Start(ctx); err != nil {
			a.logger.ErrorContext(ctx, "failed to start server", slogx.Err(err))
		}

		stop()
	}()

	if err := a.loader.Watch(func(err error) {
		if err != nil {
			a.logger.ErrorContext(ctx, "error during config watch", slogx.Err(err))
			stop()
		}

		reloaded := new(Config)
		if err := a.loader.Load(reloaded); err != nil {
			a.logger.ErrorContext(ctx, "config reload failed", slogx.Err(err))
			return
		}

		a.logger.InfoContext(ctx, "config reloaded", slogx.Config(reloaded))
	}); err != nil {
		a.logger.ErrorContext(ctx, "failed to start config watcher", slogx.Err(err))
		stop()
	}
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
