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
func (h *Handler) GetResource(
	ctx context.Context,
	req *connect.Request[examplev1.GetResourceRequest],
) (*connect.Response[examplev1.GetResourceResponse], error) {
	id, err := models.ParseID(req.Msg)
	if err != nil {
		slog.ErrorContext(ctx, "failed to parse id", slogx.Err(err))
		return nil, ErrInternal
	}

	output, err := h.resourceGetter.Execute(ctx, h.provider.Datastore(), id)
	if err != nil {
		slog.DebugContext(ctx, "usecase error", slogx.Err(err))
		return nil, ErrInternal
	}

	return connect.NewResponse(&examplev1.GetResourceResponse{
		Resource: models.ResourceFromDomain(output),
	}), nil
}
