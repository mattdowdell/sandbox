package examplev1client

import (
	"context"
	"errors"
)

type ctxKey int

const (
	clientCtxKey ctxKey = iota + 1
	cleanupsCtxKey
	scenarioCtxKey
)

// IntoContext returns a child context with the Client within.
func IntoContext(ctx context.Context, client *Client) context.Context {
	return context.WithValue(ctx, clientCtxKey, client)
}

// FromContext retrieves a Client from the given context.
func FromContext(ctx context.Context) (*Client, error) {
	client, ok := ctx.Value(clientCtxKey).(*Client)
	if !ok {
		return nil, errors.New("client not found in context")
	}

	return client, nil
}

// ...
func AppendCleanup(ctx context.Context, fn ...Cleanup) context.Context {
	if cleanups, ok := ctx.Value(cleanupsCtxKey).([]Cleanup); ok {
		return context.WithValue(ctx, cleanupsCtxKey, append(cleanups, fn...))
	}

	return context.WithValue(ctx, cleanupsCtxKey, fn)
}

// ...
func RunCleanups(ctx context.Context) error {
	cleanups, ok := ctx.Value(cleanupsCtxKey).([]Cleanup)
	if !ok {
		return nil
	}

	var errs []error

	for _, fn := range cleanups {
		if err := fn(ctx); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
