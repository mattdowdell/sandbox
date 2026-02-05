package reflectrpc

import (
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"
)

// Handler implements the ServerReflection RPC.
type Handler struct {
	services []string
}

// New creates a new Handler.
func New(services []string) *Handler {
	return &Handler{
		services: services,
	}
}

// Register adds the handler to the given multiplexer.
func (h *Handler) Register(mux *http.ServeMux, opts []connect.HandlerOption) {
	reflector := grpcreflect.NewStaticReflector(h.services...)

	mux.Handle(grpcreflect.NewHandlerV1(reflector, opts...))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector, opts...))
}
