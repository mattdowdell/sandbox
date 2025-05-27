package slogx

import (
	"context"
	"log/slog"
)

type ctxKey struct{}

// FromContext returns the logger from the context, or the default logger if none exists.
func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok && logger != nil {
		return logger
	}

	return slog.Default()
}

// IntoContext adds the given logger to the context, returning the new context.
func IntoContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}
