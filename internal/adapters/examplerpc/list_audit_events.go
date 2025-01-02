package examplerpc

import (
	"context"

	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/pkg/example/v1"
)

// ...
func (h *Handler) ListAuditEvents(
	_ context.Context,
	_ *connect.Request[examplev1.ListAuditEventsRequest],
) (*connect.Response[examplev1.ListAuditEventsResponse], error) {
	return nil, ErrUnimplemented
}
