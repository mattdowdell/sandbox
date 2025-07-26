package otelx

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"

	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// ExtractOption implementations modify the behaviour of Extractor.Extract.
type ExtractOption interface {
	apply(*Extractor)
}

// Extractor uses span metadata to create log attributes.
type Extractor struct {
	withSpanID  bool
	withSampled bool
}

// NewExtractor creates a new Extractor.
func NewExtractor(options ...ExtractOption) *Extractor {
	e := &Extractor{}

	for _, option := range options {
		option.apply(e)
	}

	return e
}

// Extract extracts span metadata from the given context to create log attributes.
func (e *Extractor) Extract(ctx context.Context) []slog.Attr {
	span := trace.SpanFromContext(ctx)

	if !span.SpanContext().IsValid() {
		return nil
	}

	attrs := []slog.Attr{
		slogx.TraceID(span),
	}

	if e.withSpanID {
		attrs = append(attrs, slogx.SpanID(span))
	}

	if e.withSampled {
		attrs = append(attrs, slogx.Sampled(span))
	}

	return attrs
}

// WithSpanID causes Extractor.Extract to include the current span ID in the return log attributes.
func WithSpanID(include bool) ExtractOption {
	return spanIDOpt(include)
}

type spanIDOpt bool

func (o spanIDOpt) apply(e *Extractor) {
	e.withSpanID = bool(o)
}

// WithSpanID causes Extractor.Extract to include the sampling status in the return log attributes.
func WithSampled(include bool) ExtractOption {
	return sampledOpt(include)
}

type sampledOpt bool

func (o sampledOpt) apply(e *Extractor) {
	e.withSampled = bool(o)
}
