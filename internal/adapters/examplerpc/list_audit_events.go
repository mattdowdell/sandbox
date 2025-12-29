package examplerpc

import (
	"context"

	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/gen/example/v1"
	"github.com/mattdowdell/sandbox/internal/adapters/examplerpc/models"
	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver/rpcerrors"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// ...
func (h *Handler) ListAuditEvents(
	ctx context.Context,
	_ *connect.Request[examplev1.ListAuditEventsRequest],
) (*connect.Response[examplev1.ListAuditEventsResponse], error) {
	logger := slogx.FromContext(ctx)

	output, err := h.auditEvent.List(ctx, logger)
	if err != nil {
		logger.DebugContext(ctx, "usecase error", slogx.Err(err))
		return nil, rpcerrors.ErrInternal
	}

	return connect.NewResponse(&examplev1.ListAuditEventsResponse{
		Items: models.AuditEventsFromDomain(output),
	}), nil
}
