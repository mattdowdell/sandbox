package authnrpc

import (
	"net/http"

	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/gen/authn/v1/authnv1connect"
)

// Non-allocating compile-time check for interface implementation.
var _ authnv1connect.AuthnServiceHandler = (*Handler)(nil)

// ...
type Parser interface {
	Parse(string) error
}

// ...
type Issuer interface {
	Issue(string) (string, error)
}

// Handler implements the AuthnService RPC.
type Handler struct {
	parser Parser
	issuer Issuer
}

// New creates a new Handler.
func New() *Handler {
	return &Handler{}
}

// Register adds the handler to the given multiplexer.
func (h *Handler) Register(mux *http.ServeMux, opts []connect.HandlerOption) {
	mux.Handle(authnv1connect.NewAuthnServiceHandler(h, opts...))
}
