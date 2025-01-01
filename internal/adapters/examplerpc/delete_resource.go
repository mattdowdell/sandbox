package examplerpc

import (
	"context"

	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/pkg/example/v1"
)

// ...
func (h *Handler) DeleteResource(
	_ context.Context,
	_ *connect.Request[examplev1.DeleteResourceRequest],
) (*connect.Response[examplev1.DeleteResourceResponse], error) {
	return nil, ErrUnimplemented
}
