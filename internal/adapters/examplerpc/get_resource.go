package examplerpc

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/gen/example/v1"
	"github.com/mattdowdell/sandbox/internal/adapters/examplerpc/models"
	"github.com/mattdowdell/sandbox/internal/domain/apperrors"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// ...
func (h *Handler) GetResource(
	ctx context.Context,
	req *connect.Request[examplev1.GetResourceRequest],
) (*connect.Response[examplev1.GetResourceResponse], error) {
	logger := slogx.FromContext(ctx)

	id, err := models.ParseID(req.Msg)
	if err != nil {
		logger.ErrorContext(ctx, "failed to parse id", slogx.Err(err))
		return nil, ErrInternal
	}

	output, err := h.resource.Get(ctx, logger, id)
	if err != nil {
		logger.DebugContext(ctx, "failed to get resource", slogx.Err(err))

		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, ErrResourceNotFound
		}

		return nil, ErrInternal
	}

	return connect.NewResponse(&examplev1.GetResourceResponse{
		Resource: models.ResourceFromDomain(output),
	}), nil
}
