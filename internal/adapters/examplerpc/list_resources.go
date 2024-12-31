package examplerpc

import (
	"context"

	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/pkg/example/v1"
)

// ...
func (h *Handler) ListResources(
	ctx context.Context,
	req *connect.Request[examplev1.ListResourcesRequest],
) (*connect.Response[examplev1.ListResourcesResponse], error) {
	return nil, ErrUnimplemented
}
