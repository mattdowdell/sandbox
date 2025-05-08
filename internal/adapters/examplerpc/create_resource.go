package examplerpc

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/gen/example/v1"
	"github.com/mattdowdell/sandbox/internal/adapters/examplerpc/models"
	"github.com/mattdowdell/sandbox/internal/domain"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// CreateResource handles requests for creating resources.
func (h *Handler) CreateResource(
	ctx context.Context,
	req *connect.Request[examplev1.CreateResourceRequest],
) (*connect.Response[examplev1.CreateResourceResponse], error) {
	logger := slogx.FromContext(ctx)

	input := models.ResourceCreateToDomain(req.Msg.GetResource())

	output, err := h.resource.Create(ctx, logger, input)
	if err != nil {
		logger.DebugContext(ctx, "failed to create resource", slogx.Err(err))

		if errors.Is(err, domain.ErrAlreadyExists) {
			return nil, ErrResourceAlreadyExists
		}

		return nil, ErrInternal
	}

	return connect.NewResponse(&examplev1.CreateResourceResponse{
		Resource: models.ResourceFromDomain(output),
	}), nil
}
