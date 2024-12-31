package examplerpc

import (
	"context"

	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/pkg/example/v1"
)

// ...
func (h *Handler) UpdateResource(
	ctx context.Context,
	req *connect.Request[examplev1.UpdateResourceRequest],
) (*connect.Response[examplev1.UpdateResourceResponse], error) {
	return nil, ErrUnimplemented
}
