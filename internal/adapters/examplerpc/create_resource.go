package examplerpc

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"
	"go.opentelemetry.io/otel/trace"

	"github.com/mattdowdell/sandbox/internal/drivers/logging"
	"github.com/mattdowdell/sandbox/pkg/example/v1"
)

// ...
func (h *Handler) CreateResource(
	ctx context.Context,
	_ *connect.Request[examplev1.CreateResourceRequest],
) (*connect.Response[examplev1.CreateResourceResponse], error) {
	span := trace.SpanFromContext(ctx)
	slog.InfoContext(ctx, "create resource called", logging.TraceID(span))

	return nil, ErrUnimplemented
}
