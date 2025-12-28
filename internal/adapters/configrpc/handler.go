package configrpc

import (
	"net/http"

	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/gen/config/v1/configv1connect"
	"github.com/mattdowdell/sandbox/internal/drivers/config"
)

// Non-allocating compile-time check for interface implementation.
var _ configv1connect.ConfigServiceHandler = (*Handler[struct{}])(nil)

// Handler implements the ConfigService RPC.
type Handler[T any] struct {
	loader *config.Config[T]
}

// New creates a new Handler.
func New[T any](loader *config.Config[T]) *Handler[T] {
	return &Handler[T]{
		loader: loader,
	}
}

// Register adds the handler to the given multiplexer.
func (h *Handler[T]) Register(mux *http.ServeMux, opts []connect.HandlerOption) {
	mux.Handle(configv1connect.NewConfigServiceHandler(h, opts...))
}
