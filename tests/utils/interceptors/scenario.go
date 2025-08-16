package interceptors

import (
	"context"
	"fmt"
	"strings"

	"connectrpc.com/connect"
	"github.com/cucumber/godog"
	"go.opentelemetry.io/otel/baggage"
)

type scenarioCtxKey struct{}

// ...
func ScenarioIntoContext(ctx context.Context, scen *godog.Scenario) context.Context {
	return context.WithValue(ctx, scenarioCtxKey{}, scen.Name)
}

// ...
func ScenarioFromContext(ctx context.Context) (string, bool) {
	name, ok := ctx.Value(scenarioCtxKey{}).(string)
	return name, ok
}

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
