package output

import (
	"context"
	"errors"

	"github.com/mattdowdell/sandbox/gen/example/v1"
)

type ctxKey int

const (
	emptyCtxKey ctxKey = iota + 1
	resourceCtxKey
	resourcesCtxKey
	nextCtxKey
	errCtxKey
	configValueCtxKey
)

// ...
func EmptyIntoContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, emptyCtxKey, 0)
}

// ...
func EmptyFromContext(ctx context.Context) error {
	if val, ok := ctx.Value(emptyCtxKey).(int); ok && val == 0 {
		return nil
	}

	return errors.New("empty not found in context")
}

// ...
func ResourceIntoContext(ctx context.Context, r *examplev1.Resource) context.Context {
	return context.WithValue(ctx, resourceCtxKey, r)
}

// ...
func ResourceFromContext(ctx context.Context) (*examplev1.Resource, error) {
	if r, ok := ctx.Value(resourceCtxKey).(*examplev1.Resource); ok && r != nil {
		return r, nil
	}

	return nil, errors.New("resource not found in context")
}

// ...
func ResourcesIntoContext(ctx context.Context, r []*examplev1.Resource) context.Context {
	return context.WithValue(ctx, resourcesCtxKey, r)
}

// ...
func ResourcesFromContext(ctx context.Context) ([]*examplev1.Resource, error) {
	if r, ok := ctx.Value(resourcesCtxKey).([]*examplev1.Resource); ok {
		return r, nil
	}

	return nil, errors.New("resources not found in context")
}

// ...
func NextIntoContext(ctx context.Context, n string) context.Context {
	return context.WithValue(ctx, nextCtxKey, n)
}

// ...
func NextFromContext(ctx context.Context) (string, error) {
	if n, ok := ctx.Value(nextCtxKey).(string); ok {
		return n, nil
	}

	return "", errors.New("next not found in context")
}

// ...
func ErrIntoContext(ctx context.Context, err error) context.Context {
	return context.WithValue(ctx, errCtxKey, err)
}

// ...
func ErrFromContext(ctx context.Context) (have, err error) {
	if have, ok := ctx.Value(errCtxKey).(error); ok && have != nil {
		return have, nil
	}

	return nil, errors.New("err not found in context")
}

// ...
func ConfigValueIntoContext(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, configValueCtxKey, value)
}

// ...
func ConfigValueFromContext(ctx context.Context) (string, error) {
	if val, ok := ctx.Value(configValueCtxKey).(string); ok && val != "" {
		return val, nil
	}

	return "", errors.New("config value not found in context")
}
