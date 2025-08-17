package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/mattdowdell/sandbox/internal/drivers/exit"
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

	defer app.Shutdown(ctx)

	if err := app.Run(ctx); err != nil {
		return exit.Failure
	}

	return exit.Success
}
