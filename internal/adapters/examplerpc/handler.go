package examplerpc

import (
	"net/http"

	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/gen/example/v1/examplev1connect"
)

// Non-allocating compile-time check for interface implementation.
var _ examplev1connect.ExampleServiceHandler = (*Handler)(nil)

// Handler implements the ExampleService RPC.
type Handler struct {
	resource   ResourceFacade
	auditEvent AuditEventFacade
}

// New creates a new Handler.
func New(
	resource ResourceFacade,
	auditEvent AuditEventFacade,
) *Handler {
	return &Handler{
		resource:   resource,
		auditEvent: auditEvent,
	}
}

// Register adds the handler to the given multiplexer.
func (h *Handler) Register(mux *http.ServeMux, opts []connect.HandlerOption) {
	mux.Handle(examplev1connect.NewExampleServiceHandler(h, opts...))
}
