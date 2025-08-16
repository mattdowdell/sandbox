package authnrpc

import (
	"context"

	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/gen/authn/v1"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// Authenticate validates the given token and returns a subset of the claims for use in other logic.
func (h *Handler) Authenticate(
	ctx context.Context,
	req *connect.Request[authnv1.AuthenticateRequest],
) (*connect.Response[authnv1.AuthenticateResponse], error) {
	logger := slogx.FromContext(ctx)

	claims, err := h.parser.Parse(req.Msg.GetToken())
	if err != nil {
		logger.InfoContext(ctx, "failed to parse token", slogx.Err(err))
		return nil, ErrUnauthenticated
	}

	subject, err := claims.GetSubject()
	if err != nil {
		logger.ErrorContext(ctx, "failed to extract subject claim", slogx.Err(err))
		return nil, ErrInternal
	}

	return connect.NewResponse(&authnv1.AuthenticateResponse{
		Subject: subject,
	}), nil
}
