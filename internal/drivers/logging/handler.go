package logging

import (
	"context"
	"log/slog"
)

// Extractor implementations are used to add attributes to a log record based on context metadata.
// One or more extractors can be added to the logger using the WithExtractors option.
type Extractor interface {
	Extract(context.Context) []slog.Attr
}

type handler struct {
	inner      slog.Handler
	extractors []Extractor
}

// Wrap wraps a slog handler with support for augmenting a log record with metadata from the
// context.
func Wrap(inner slog.Handler, extractors []Extractor) slog.Handler {
	return &handler{
		inner:      inner,
		extractors: extractors,
	}
}

func (h *handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

//nolint:gocritic // can't change parameter types of third-party interface
func (h *handler) Handle(ctx context.Context, record slog.Record) error {
	for _, e := range h.extractors {
		if attrs := e.Extract(ctx); len(attrs) > 0 {
			record.AddAttrs(attrs...)
		}
	}

	return h.inner.Handle(ctx, record)
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &handler{
		inner: h.inner.WithAttrs(attrs),
	}
}

func (h *handler) WithGroup(name string) slog.Handler {
	return &handler{
		inner: h.inner.WithGroup(name),
	}
}
