package examplerpc

import (
	"context"

	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/gen/example/v1"
	"github.com/mattdowdell/sandbox/internal/adapters/examplerpc/models"
	"github.com/mattdowdell/sandbox/internal/domain/repositories"
	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver/rpcerrors"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// ...
func (h *Handler) ListResources(
	ctx context.Context,
	req *connect.Request[examplev1.ListResourcesRequest],
) (*connect.Response[examplev1.ListResourcesResponse], error) {
	logger := slogx.FromContext(ctx)
	pager := repositories.Pager{
		Limit: int(req.Msg.GetLimit()),
	}

	output, err := h.resource.List(ctx, logger, pager)
	if err != nil {
		logger.DebugContext(ctx, "failed to list resources", slogx.Err(err))
		return nil, rpcerrors.ErrInternal
	}

	return connect.NewResponse(&examplev1.ListResourcesResponse{
		Items: models.ResourcesFromDomain(output.Items),
	}), nil
}
