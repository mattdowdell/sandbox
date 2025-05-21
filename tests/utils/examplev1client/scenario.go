package examplev1client

import (
	"context"
	"fmt"
	"strings"

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
			s, err := baggage.NewMember("scenario.name", strings.ReplaceAll(scenario, " ", "_"))
			if err != nil {
				return nil, connect.NewError(
					connect.CodeCanceled,
					fmt.Errorf("failed to create baggage member: %w", err),
				)
			}

			b, err := baggage.New(s)
			if err != nil {
				return nil, connect.NewError(
					connect.CodeCanceled,
					fmt.Errorf("failed to create baggage: %w", err),
				)
			}

			ctx = baggage.ContextWithBaggage(ctx, b)
		}

		return next(ctx, req)
	})
}
