package configrpc

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/gen/config/v1"
	"github.com/mattdowdell/sandbox/internal/drivers/config"
	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver/rpcerrors"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// GetConfig returns the current configuration values.
func (h *Handler[T]) GetConfig(
	ctx context.Context,
	_ *connect.Request[configv1.GetConfigRequest],
) (*connect.Response[configv1.GetConfigResponse], error) {
	conf, err := h.loader.Load()
	if err != nil {
		slog.ErrorContext(ctx, "failed to load configuration", slogx.Err(err))
		return nil, rpcerrors.ErrInternal
	}

	encoded, err := config.Encode(conf, "." /*delim*/)
	if err != nil {
		slog.ErrorContext(ctx, "failed to encode configuration", slogx.Err(err))
		return nil, rpcerrors.ErrInternal
	}

	return connect.NewResponse(&configv1.GetConfigResponse{
		Config: encoded,
	}), nil
}
