// Package slogt integrates [testing] and [log/slog] together.
//
// Use of this package means that structured logs produced either by a test of the code under test
// can be emitted to the test's output. This improves debuggability, as the log output is now
// attached to the test that caused it to be emitted.
package slogt

import (
	"log/slog"
	"testing"
)

// New is an alias for [Text].
func New(tb testing.TB) *slog.Logger {
	tb.Helper()
	return Text(tb)
}

// Text wraps [TextWithOptions] with default (nil) options.
func Text(tb testing.TB) *slog.Logger {
	tb.Helper()
	return TextWithOptions(tb, nil)
}

// TextWithOptions creates a new logger using a text handler and the given options.
func TextWithOptions(tb testing.TB, options *slog.HandlerOptions) *slog.Logger {
	tb.Helper()
	return slog.New(slog.NewTextHandler(tb.Output(), options))
}

// JSON wraps [JSONWithOptions] with default (nil) options.
func JSON(tb testing.TB) *slog.Logger {
	tb.Helper()
	return JSONWithOptions(tb, nil)
}

// JSONWithOptions creates a new logger using a JSON handler and default options.
func JSONWithOptions(tb testing.TB, options *slog.HandlerOptions) *slog.Logger {
	tb.Helper()
	return slog.New(slog.NewTextHandler(tb.Output(), options))
}
