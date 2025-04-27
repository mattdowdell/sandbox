package examplev1client

import (
	"context"

	"connectrpc.com/connect"
	"go.opentelemetry.io/otel/baggage"
)

// ...
func ScenarioUnaryInterceptor(next connect.UnaryFunc) connect.UnaryFunc {
	return connect.UnaryFunc(func(
		ctx context.Context,
		req connect.AnyRequest,
	) (connect.AnyResponse, error) {
		if scenario, ok := ScenarioFromContext(ctx); ok && scenario != "" {
			s, err := baggage.NewMember("scenario.name", scenario)
			if err != nil {
				return nil, connect.NewError(connect.CodeCanceled, err)
			}

			b, err := baggage.New(s)
			if err != nil {
				return nil, connect.NewError(connect.CodeCanceled, err)
			}

			ctx = baggage.ContextWithBaggage(ctx, b)
		}

		return next(ctx, req)
	})
}
