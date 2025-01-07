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
func (h *Handler) DeleteResource(
	ctx context.Context,
	req *connect.Request[examplev1.DeleteResourceRequest],
) (*connect.Response[examplev1.DeleteResourceResponse], error) {
	id, err := models.ParseID(req.Msg)
	if err != nil {
		slog.ErrorContext(ctx, "failed to parse id", slogx.Err(err))
		return nil, ErrInternal
	}

	if err := h.resourceDeleter.Execute(ctx, id); err != nil {
		return nil, ErrInternal
	}

	return connect.NewResponse(&examplev1.DeleteResourceResponse{}), nil
}
