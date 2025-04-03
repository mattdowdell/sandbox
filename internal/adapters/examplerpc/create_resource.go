package examplerpc

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/gen/example/v1"
	"github.com/mattdowdell/sandbox/internal/adapters/examplerpc/models"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// ...
func (h *Handler) CreateResource(
	ctx context.Context,
	req *connect.Request[examplev1.CreateResourceRequest],
) (*connect.Response[examplev1.CreateResourceResponse], error) {
	input := models.ResourceCreateToDomain(req.Msg.GetResource())

	output, err := h.resource.Create(ctx, input)
	if err != nil {
		slog.DebugContext(ctx, "failed to create resource", slogx.Err(err))
		return nil, ErrInternal // TODO: use more granular errors
	}

	return connect.NewResponse(&examplev1.CreateResourceResponse{
		Resource: models.ResourceFromDomain(output),
	}), nil
}
