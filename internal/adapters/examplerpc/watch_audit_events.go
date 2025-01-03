package examplerpc

import (
	"context"

	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/pkg/example/v1"
)

// ...
func (h *Handler) WatchAuditEvents(
	_ context.Context,
	_ *connect.Request[examplev1.WatchAuditEventsRequest],
	_ *connect.ServerStream[examplev1.WatchAuditEventsResponse],
) error {
	return ErrUnimplemented
}
