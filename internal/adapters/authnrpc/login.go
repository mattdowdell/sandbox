package authnrpc

import (
	"context"

	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/gen/authn/v1"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

const (
	expiresInSeconds = 7200
	tokenType        = "Bearer"
)

// Login returns a JWT for the ID in the request.
//
// No validation of the caller's identity is performed, nor whether the given secret is in any way
// valid. While this is normally indefensible, this method exists only as a means to experiment with
// authorization implementations, and is thus acceptable.
//
// In a production environment, this would likely involve checking the identity exists in a
// database. Additionally, validation that the secret is acceptable for that identity would be
// performed.
func (h *Handler) Login(
	ctx context.Context,
	req *connect.Request[authnv1.LoginRequest],
) (*connect.Response[authnv1.LoginResponse], error) {
	logger := slogx.FromContext(ctx)

	token, err := h.issuer.Issue(req.Msg.GetId(), expiresInSeconds)
	if err != nil {
		logger.InfoContext(ctx, "failed to issue token", slogx.Err(err))
		return nil, ErrInternal
	}

	return connect.NewResponse(&authnv1.LoginResponse{
		AccessToken: token,
		TokenType:   tokenType,
		ExpiresIn:   expiresInSeconds,
	}), nil
}
