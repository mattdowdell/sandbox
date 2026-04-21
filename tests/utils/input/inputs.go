package input

import (
	"context"
	"errors"

	"github.com/cucumber/godog"
)

type ctxKey int

const (
	scenarioCtxKey ctxKey = iota + 1
	nameCtxKey
	idCtxKey
	newNameCtxKey
	limitCtxKey
	authnCtxKey
	configKeyCtxKey
)

// ...
func ScenarioIntoContext(ctx context.Context, scen *godog.Scenario) context.Context {
	return context.WithValue(ctx, scenarioCtxKey, scen.Name)
}

// ...
func ScenarioFromContext(ctx context.Context) (string, bool) {
	name, ok := ctx.Value(scenarioCtxKey).(string)
	return name, ok
}

// ...
func NameIntoContext(ctx context.Context, name string) context.Context {
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
func NewNameIntoContext(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, newNameCtxKey, name)
}

// ...
func NewNameFromContext(ctx context.Context) (string, error) {
	if name, ok := ctx.Value(newNameCtxKey).(string); ok {
		return name, nil
	}

	return "", errors.New("new name not found in context")
}

// IDIntoContext adds an ID to the given context and returns the result.
func IDIntoContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, idCtxKey, id)
}

// IDFromContext extracts an ID from the context, if one exists.
func IDFromContext(ctx context.Context) (string, error) {
	if id, ok := ctx.Value(idCtxKey).(string); ok {
		return id, nil
	}

	return "", errors.New("id not found in context")
}

// ...
func LimitIntoContext(ctx context.Context, limit int32) context.Context {
	return context.WithValue(ctx, limitCtxKey, limit)
}

// ...
func LimitFromContext(ctx context.Context) (int32, error) {
	if limit, ok := ctx.Value(limitCtxKey).(int32); ok {
		return limit, nil
	}

	return 0, errors.New("limit not found in context")
}

// ...
func AuthnIntoContext(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, authnCtxKey, value)
}

// ...
func AuthnFromContext(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(authnCtxKey).(string)
	return val, ok
}

// ...
func ConfigKeyIntoContext(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, configKeyCtxKey, value)
}

// ...
func ConfigKeyFromContext(ctx context.Context) (string, error) {
	if key, ok := ctx.Value(configKeyCtxKey).(string); ok {
		return key, nil
	}

	return "", errors.New("config key not found in context")
}
