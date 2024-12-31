package examplerpc

import (
	"context"

	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/pkg/example/v1"
)

// ...
func (h *Handler) GetResource(
	ctx context.Context,
	req *connect.Request[examplev1.GetResourceRequest],
) (*connect.Response[examplev1.GetResourceResponse], error) {
	return nil, ErrUnimplemented
}
