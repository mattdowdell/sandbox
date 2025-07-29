package authnrpc

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/gen/authn/v1"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// ...
func (h *Handler) Authenticate(
	ctx context.Context,
	req *connect.Request[authnv1.AuthenticateRequest],
) (*connect.Response[authnv1.AuthenticateResponse], error) {
	if err := h.parser.Parse(req.Msg.GetToken()); err != nil {
		slog.InfoContext(ctx, "failed to parse token", slogx.Err(err))
		return nil, ErrUnauthenticated
	}

	return connect.NewResponse(&authnv1.AuthenticateResponse{
		// TODO: subject
		// +other claims?
	}), nil
}
