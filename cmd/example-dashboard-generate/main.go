package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/mattdowdell/sandbox/internal/drivers/exit"
	"github.com/mattdowdell/sandbox/pkg/dashboards"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

func main() {
	os.Exit(run(context.Background()))
}

func run(ctx context.Context) int {
	connectRPC, err := dashboards.ConnectRPC()
	if err != nil {
		slog.ErrorContext(ctx, "failed to generate connect rpc dashboard", slogx.Err(err))
	}

	fmt.Println(connectRPC)

	return exit.Success
}
