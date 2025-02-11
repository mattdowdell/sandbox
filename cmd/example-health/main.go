package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"buf.build/gen/go/grpc/grpc/connectrpc/go/grpc/health/v1/healthv1connect"
	"buf.build/gen/go/grpc/grpc/protocolbuffers/go/grpc/health/v1"
	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/internal/drivers/config"
	"github.com/mattdowdell/sandbox/internal/drivers/config/flagoptions"
	"github.com/mattdowdell/sandbox/internal/drivers/exit"
	"github.com/mattdowdell/sandbox/internal/drivers/logging"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// Config contains the service configuration.
type Config struct {
	Logging   logging.Config `koanf:"logging"`
	RPCServer struct {
		Port uint16 `koanf:"port" default:"5000"`
	} `koanf:"rpcserver"`
}

func main() {
	os.Exit(run(context.Background()))
}

func run(ctx context.Context) int {
	options := flagoptions.New()
	conf := config.New(options)

	c, err := config.Load[Config](conf)
	if err != nil {
		slog.ErrorContext(ctx, "failed to load configuration", slogx.Err(err))
		return exit.Failure
	}

	logger := logging.New(c.Logging.Level)

	client := healthv1connect.NewHealthClient(
		http.DefaultClient,
		fmt.Sprintf("http://localhost:%d", c.RPCServer.Port),
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
