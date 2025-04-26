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
	errCtxKey
)

// ...
func AddEmptyToContext(ctx context.Context) context.Context {
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
func AddResourceToContext(ctx context.Context, r *examplev1.Resource) context.Context {
	return context.WithValue(ctx, resourceCtxKey, r)
}

// ...
func ResourceFromContext(ctx context.Context) (*examplev1.Resource, error) {
	if r, ok := ctx.Value(resourceCtxKey).(*examplev1.Resource); ok && r != nil {
		return r, nil
	}

	return nil, errors.New("empty not found in context")
}

// ...
func AddErrToContext(ctx context.Context, err error) context.Context {
	return context.WithValue(ctx, errCtxKey, err)
}

// ...
func ErrFromContext(ctx context.Context) (have, err error) {
	if have, ok := ctx.Value(errCtxKey).(error); ok && have != nil {
		return have, nil
	}

	return nil, errors.New("empty not found in context")
}
