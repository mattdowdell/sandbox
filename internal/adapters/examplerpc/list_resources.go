package examplerpc

import (
	"context"

	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/gen/example/v1"
	"github.com/mattdowdell/sandbox/internal/adapters/examplerpc/models"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// ...
func (h *Handler) ListResources(
	ctx context.Context,
	_ *connect.Request[examplev1.ListResourcesRequest],
) (*connect.Response[examplev1.ListResourcesResponse], error) {
	logger := slogx.FromContext(ctx)

	output, err := h.resource.List(ctx, logger)
	if err != nil {
		logger.DebugContext(ctx, "failed to list resources", slogx.Err(err))
		return nil, ErrInternal
	}

	return connect.NewResponse(&examplev1.ListResourcesResponse{
		Items: models.ResourcesFromDomain(output),
	}), nil
}
