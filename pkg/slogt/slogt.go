package slogt

import (
	"log/slog"
	"testing"
)

// New is an alias for Text.
func New(tb testing.TB) *slog.Logger {
	tb.Helper()
	return Text(tb)
}

// Text creates a new logger using a text handler and default options.
func Text(tb testing.TB) *slog.Logger {
	tb.Helper()
	return TextWithOptions(tb, nil)
}

// TextWithOptions creates a new logger using a text handler and the given options.
func TextWithOptions(tb testing.TB, options *slog.HandlerOptions) *slog.Logger {
	tb.Helper()
	return slog.New(slog.NewTextHandler(tb.Output(), options))
}

// JSON creates a new logger using a JSON handler and default options.
func JSON(tb testing.TB) *slog.Logger {
	tb.Helper()
	return JSONWithOptions(tb, nil)
}

// JSONWithOptions creates a new logger using a JSON handler and default options.
func JSONWithOptions(tb testing.TB, options *slog.HandlerOptions) *slog.Logger {
	tb.Helper()
	return slog.New(slog.NewTextHandler(tb.Output(), options))
}
