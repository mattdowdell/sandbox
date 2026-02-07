package interceptors

import (
	"context"

	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/tests/utils/input"
)

// ...
func AuthnUnaryInterceptor(next connect.UnaryFunc) connect.UnaryFunc {
	return connect.UnaryFunc(func(
		ctx context.Context,
		req connect.AnyRequest,
	) (connect.AnyResponse, error) {
		if authn, ok := input.AuthnFromContext(ctx); ok && authn != "" {
			req.Header().Set("Authorization", authn)
		}

		return next(ctx, req)
	})
}
