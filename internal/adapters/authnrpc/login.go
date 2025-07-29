package authnrpc

import (
	"context"

	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/gen/authn/v1"
)

const (
	// 1 hour in seconds.
	expiresIn = 3600
)

// ...
func (h *Handler) Login(
	_ context.Context,
	req *connect.Request[authnv1.LoginRequest],
) (*connect.Response[authnv1.LoginResponse], error) {
	token, err := h.issuer.Issue(req.Msg.GetId())
	if err != nil {
		return nil, ErrInternal
	}

	return connect.NewResponse(&authnv1.LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn, // TODO: make dynamic
	}), nil
}
