package examplerpc

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/gen/example/v1"
	"github.com/mattdowdell/sandbox/internal/adapters/common"
	"github.com/mattdowdell/sandbox/internal/adapters/examplerpc/models"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// ...
func (h *Handler) UpdateResource(
	ctx context.Context,
	req *connect.Request[examplev1.UpdateResourceRequest],
) (*connect.Response[examplev1.UpdateResourceResponse], error) {
	input, err := models.ResourceUpdateToDomain(req.Msg.GetResource())
	if err != nil {
		slog.ErrorContext(ctx, "failed to parse id", slogx.Err(err))
		return nil, ErrInternal
	}

	output, err := common.TxValue(ctx, h.provider, func(ds common.Datastore) (*entities.Resource, error) {
		return h.resourceCreator.Execute(ctx, ds, input)
	})
	if err != nil {
		slog.DebugContext(ctx, "failed to update resource", slogx.Err(err))
		return nil, ErrInternal // TODO: use more granular errors
	}

	return connect.NewResponse(&examplev1.UpdateResourceResponse{
		Resource: models.ResourceFromDomain(output),
	}), nil
}
