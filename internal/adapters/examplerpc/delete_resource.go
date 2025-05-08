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

// DeleteResource handles requests for deleting resources.
//
// This assumes that validation has already been performed on the request. As such, syntactic errors
// in the input are considered internal errors to indicate that the system is not configured as
// expected.
func (h *Handler) DeleteResource(
	ctx context.Context,
	req *connect.Request[examplev1.DeleteResourceRequest],
) (*connect.Response[examplev1.DeleteResourceResponse], error) {
	logger := slogx.FromContext(ctx)

	id, err := models.ParseID(req.Msg)
	if err != nil {
		logger.ErrorContext(ctx, "failed to parse id", slogx.Err(err))
		return nil, ErrInternal
	}

	if err := h.resource.Delete(ctx, logger, id); err != nil {
		logger.DebugContext(ctx, "failed to delete resource", slogx.Err(err))

		if errors.Is(err, domain.ErrNotFound) {
			return nil, ErrResourceNotFound
		}

		return nil, ErrInternal
	}

	return connect.NewResponse(&examplev1.DeleteResourceResponse{}), nil
}
