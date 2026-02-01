package configrpc

import (
	"net/http"

	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/gen/config/v1/configv1connect"
)

// Non-allocating compile-time check for interface implementation.
var _ configv1connect.ConfigServiceHandler = (*Handler[struct{}])(nil)

type Loader interface {
	Load(any) error
}

// Handler implements the ConfigService RPC.
type Handler[T any] struct {
	loader Loader
}

// New creates a new Handler.
func New[T any](loader Loader) *Handler[T] {
	return &Handler[T]{
		loader: loader,
	}
}

// Register adds the handler to the given multiplexer.
func (h *Handler[T]) Register(mux *http.ServeMux, opts []connect.HandlerOption) {
	mux.Handle(configv1connect.NewConfigServiceHandler(h, opts...))
}
