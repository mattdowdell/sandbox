package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/mattdowdell/sandbox/pkg/exit"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

func main() {
	os.Exit(run(context.Background()))
}

func run(ctx context.Context) int {
	app, err := SetupApp(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to setup app", slogx.Err(err))
		return exit.Failure
	}

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app.Start(ctx, stop)

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), app.shutdownTimeout)
	defer cancel()

	app.Shutdown(ctx)

	return exit.Success
}
