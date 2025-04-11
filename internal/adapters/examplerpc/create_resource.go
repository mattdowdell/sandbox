package examplerpc

import (
	"context"
	"errors"
	"log/slog"

	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/gen/example/v1"
	"github.com/mattdowdell/sandbox/internal/adapters/examplerpc/models"
	"github.com/mattdowdell/sandbox/internal/domain/apperrors"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// CreateResource handles requests for creating resources.
func (h *Handler) CreateResource(
	ctx context.Context,
	req *connect.Request[examplev1.CreateResourceRequest],
) (*connect.Response[examplev1.CreateResourceResponse], error) {
	input := models.ResourceCreateToDomain(req.Msg.GetResource())

	output, err := h.resource.Create(ctx, input)
	if err != nil {
		slog.DebugContext(ctx, "failed to create resource", slogx.Err(err))

		if errors.Is(err, apperrors.ErrAlreadyExists) {
			return nil, ErrResourceAlreadyExists(input.Name)
		}

		return nil, ErrInternal
	}

	return connect.NewResponse(&examplev1.CreateResourceResponse{
		Resource: models.ResourceFromDomain(output),
	}), nil
}
