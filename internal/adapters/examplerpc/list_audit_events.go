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
func (h *Handler) ListAuditEvents(
	ctx context.Context,
	_ *connect.Request[examplev1.ListAuditEventsRequest],
) (*connect.Response[examplev1.ListAuditEventsResponse], error) {
	output, err := h.auditEventLister.Execute(ctx)
	if err != nil {
		slog.DebugContext(ctx, "usecase error", slogx.Err(err))
		return nil, ErrInternal
	}

	return connect.NewResponse(&examplev1.ListAuditEventsResponse{
		Items: models.AuditEventsFromDomain(output),
	}), nil
}
