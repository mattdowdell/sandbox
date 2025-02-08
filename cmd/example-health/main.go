package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"buf.build/gen/go/grpc/grpc/connectrpc/go/grpc/health/v1/healthv1connect"
	"buf.build/gen/go/grpc/grpc/protocolbuffers/go/grpc/health/v1"
	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/internal/drivers/exit"
	"github.com/mattdowdell/sandbox/internal/drivers/logging"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

func main() {
	os.Exit(run(context.Background()))
}

func run(ctx context.Context) int {
	logger := logging.New(slog.LevelInfo) // TODO: make level configurable

	client := healthv1connect.NewHealthClient(
		http.DefaultClient,
		"http://localhost:5000", // TODO: make port configurable
	)

	resp, err := client.Check(ctx, connect.NewRequest(&healthv1.HealthCheckRequest{}))
	if err != nil {
		logger.ErrorContext(ctx, "health check failed", slogx.Err(err))
		return exit.Failure
	}

	status := resp.Msg.GetStatus()

	if status != healthv1.HealthCheckResponse_SERVING {
		logger.InfoContext(ctx, "service not serving", slogx.HealthStatus(status))
		return exit.Failure
	}

	logger.DebugContext(ctx, "service is serving", slogx.HealthStatus(status))
	return exit.Success
}
