package authnv1client

import (
	"context"
	"errors"
)

type clientCtxKey struct{}

// IntoContext returns a child context with the Client within.
func IntoContext(ctx context.Context, client *Client) context.Context {
	return context.WithValue(ctx, clientCtxKey{}, client)
}

// FromContext retrieves a Client from the given context.
func FromContext(ctx context.Context) (*Client, error) {
	client, ok := ctx.Value(clientCtxKey{}).(*Client)
	if !ok {
		return nil, errors.New("client not found in context")
	}

	return client, nil
}
