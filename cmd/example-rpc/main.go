package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/mattdowdell/sandbox/internal/drivers/exit"
	"github.com/mattdowdell/sandbox/internal/drivers/logging"
)

func main() {
	os.Exit(run(context.Background()))
}

func run(ctx context.Context) int {
	app, err := ProvideApp()
	if err != nil {
		slog.ErrorContext(ctx, "failed to build app", logging.Error(err))
		return exit.Failure
	}

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer stop()

	app.logger.InfoContext(ctx, "starting")

	go func() {
		if err := app.server.Start(); err != nil {
			slog.ErrorContext(ctx, "failed to start server", logging.Error(err))
		}

		stop()
	}()

	<-ctx.Done()

	app.logger.InfoContext(ctx, "stopping")

	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), app.shutdownTimeout)
	defer cancel()

	if err := app.server.Shutdown(ctx); err != nil {
		app.logger.WarnContext(ctx, "failed to shutdown server", logging.Error(err))
	}

	if err := app.tpShutdown(ctx); err != nil {
		app.logger.WarnContext(ctx, "failed to shutdown tracer provider", logging.Error(err))
	}

	return exit.Success
}
