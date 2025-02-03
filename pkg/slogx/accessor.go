package slogx

import (
	"context"
	"log/slog"
)

type loggerCtxKey struct{}

// ...
func LoggerFromContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(loggerCtxKey{}).(*slog.Logger)
	if !ok {
		return slog.Default()
	}

	return logger
}
