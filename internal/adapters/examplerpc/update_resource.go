package examplerpc

import (
	"context"
	"errors"
	"log/slog"

	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/gen/example/v1"
	"github.com/mattdowdell/sandbox/internal/adapters/examplerpc/models"
	"github.com/mattdowdell/sandbox/internal/domain"
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

	output, err := h.resource.Update(ctx, input)
	if err != nil {
		slog.DebugContext(ctx, "failed to update resource", slogx.Err(err))

		switch {
		case errors.Is(err, domain.ErrNotFound):
			return nil, ErrResourceNotFound

		case errors.Is(err, domain.ErrAlreadyExists):
			return nil, ErrResourceAlreadyExists

		default:
			return nil, ErrInternal
		}
	}

	return connect.NewResponse(&examplev1.UpdateResourceResponse{
		Resource: models.ResourceFromDomain(output),
	}), nil
}
