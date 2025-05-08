package examplerpc

import (
	"context"

	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/gen/example/v1"
	"github.com/mattdowdell/sandbox/internal/adapters/examplerpc/models"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// ...
func (h *Handler) WatchAuditEvents(
	ctx context.Context,
	_ *connect.Request[examplev1.WatchAuditEventsRequest],
	stream *connect.ServerStream[examplev1.WatchAuditEventsResponse],
) error {
	logger := slogx.FromContext(ctx)

	ch := h.auditEvent.Watch(ctx, logger)

	for {
		select {
		case <-ctx.Done():
			logger.DebugContext(ctx, "connection closed by client")
			return nil

		case event := <-ch:
			if err := stream.Send(&examplev1.WatchAuditEventsResponse{
				AuditEvent: models.AuditEventFromDomain(event),
			}); err != nil {
				logger.ErrorContext(ctx, "failed send", slogx.Err(err))
			}
		}
	}
}
