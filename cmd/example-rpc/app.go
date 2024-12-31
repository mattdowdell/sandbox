package main

import (
	"log/slog"
	"time"

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
}

// ...
func NewApp(
	config AppConfig,
	logger *slog.Logger,
	server *rpcserver.Server,
	tpShutdown otelx.TracerProviderShutdown,
) *App {
	return &App{
		shutdownTimeout: config.ShutdownTimeout,
		logger:          logger,
		server:          server,
		tpShutdown:      tpShutdown,
	}
}
