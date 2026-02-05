package authnrpc

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/golang-jwt/jwt/v5"

	"github.com/mattdowdell/sandbox/gen/authn/v1/authnv1connect"
)

// Non-allocating compile-time check for interface implementation.
var _ authnv1connect.AuthnServiceHandler = (*Handler)(nil)

// Parser is used to parse a JWT.
type Parser interface {
	Parse(string) (jwt.Claims, error)
}

// Issuer is used to issue JWTs.
type Issuer interface {
	Issue(string, uint32) (string, error)
}

// Handler implements the AuthnService RPC.
type Handler struct {
	issuer Issuer
	parser Parser
}

// New creates a new Handler.
func New(issuer Issuer, parser Parser) *Handler {
	return &Handler{
		issuer: issuer,
		parser: parser,
	}
}

// Register adds the handler to the given multiplexer.
func (h *Handler) Register(mux *http.ServeMux, opts []connect.HandlerOption) {
	mux.Handle(authnv1connect.NewAuthnServiceHandler(h, opts...))
}
