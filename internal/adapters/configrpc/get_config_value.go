package configrpc

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/gen/config/v1"
	"github.com/mattdowdell/sandbox/internal/adapters/configrpc/models"
	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver/rpcerrors"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// GetConfigValue returns a single current configuration value.
func (h *Handler[T]) GetConfigValue(
	ctx context.Context,
	req *connect.Request[configv1.GetConfigValueRequest],
) (*connect.Response[configv1.GetConfigValueResponse], error) {
	conf, err := h.loader.Load()
	if err != nil {
		slog.ErrorContext(ctx, "failed to load configuration", slogx.Err(err))
		return nil, rpcerrors.ErrInternal
	}

	encoded, err := models.Encode(conf, "." /*delim*/)
	if err != nil {
		slog.ErrorContext(ctx, "failed to encode configuration", slogx.Err(err))
		return nil, rpcerrors.ErrInternal
	}

	value, ok := encoded[req.Msg.GetKey()]
	if !ok {
		//nolint:sloglint // not much benefit from standardising this key
		slog.DebugContext(ctx, "requested config key was not found", slog.String("key", req.Msg.GetKey()))
		return nil, ErrValueNotFound
	}

	return connect.NewResponse(&configv1.GetConfigValueResponse{
		Value: value,
	}), nil
}
