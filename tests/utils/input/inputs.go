package input

import (
	"context"
	"errors"
)

type ctxKey int

const (
	nameCtxKey ctxKey = iota + 1
	idCtxKey
	newNameCtxKey
)

// ...
func AddNameToContext(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, nameCtxKey, name)
}

// ...
func NameFromContext(ctx context.Context) (string, error) {
	if name, ok := ctx.Value(nameCtxKey).(string); ok {
		return name, nil
	}

	return "", errors.New("name not found in context")
}

// ...
func AddNewNameToContext(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, newNameCtxKey, name)
}

// ...
func NewNameFromContext(ctx context.Context) (string, error) {
	if name, ok := ctx.Value(newNameCtxKey).(string); ok {
		return name, nil
	}

	return "", errors.New("new name not found in context")
}

// ...
func AddIDToContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, idCtxKey, id)
}

// ...
func IDFromContext(ctx context.Context) (string, error) {
	if id, ok := ctx.Value(idCtxKey).(string); ok && id != "" {
		return id, nil
	}

	return "", errors.New("id not found in context")
}
