package examplev1client

import (
	"context"

	"connectrpc.com/connect"
	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/protobuf/proto"
)

// ValidateUnaryInterceptor validates that responses are valid according to Protobuf annotations.
// Invalid responses result in an Internal error code with the validation failures as a message.
// Requests are not validated to prevent obscuring the validation implemented by the server.
//
// Server validation does not currently include response validation. See
// [connectrpc/validate-go#28] for details.
//
// [connectrpc/validate-go#28]: https://github.com/connectrpc/validate-go/issues/28
func ValidateUnaryInterceptor(next connect.UnaryFunc) connect.UnaryFunc {
	return connect.UnaryFunc(func(
		ctx context.Context,
		req connect.AnyRequest,
	) (connect.AnyResponse, error) {
		resp, err := next(ctx, req)
		if err != nil {
			return resp, err
		}

		msg, ok := resp.(proto.Message)
		if !ok {
			return resp, nil
		}

		if err := protovalidate.Validate(msg); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		return resp, nil
	})
}
