package slogt

import (
	"log/slog"
	"testing"
)

// ...
func New(tb testing.TB) *slog.Logger {
	tb.Helper()
	return Text(tb)
}

// ...
func Text(tb testing.TB) *slog.Logger {
	tb.Helper()
	return TextWithOptions(tb, nil)
}

// ...
func TextWithOptions(tb testing.TB, options *slog.HandlerOptions) *slog.Logger {
	tb.Helper()
	return slog.New(slog.NewTextHandler(tb.Output(), options))
}

// ...
func JSON(tb testing.TB) *slog.Logger {
	tb.Helper()
	return JSONWithOptions(tb, nil)
}

// ...
func JSONWithOptions(tb testing.TB, options *slog.HandlerOptions) *slog.Logger {
	tb.Helper()
	return slog.New(slog.NewTextHandler(tb.Output(), options))
}
