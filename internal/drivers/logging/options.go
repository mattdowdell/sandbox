package logging

import (
	"io"
	"os"
)

// Option implementations are used to apply optional configuration to a logger.
type Option interface {
	apply(*loggerOpts)
}

type loggerOpts struct {
	writer            io.Writer
	suppressTime      bool
	suppressSource    bool
	extractors        []Extractor
	enableJSONHandler bool
	enableOtelHandler bool
}

func defaultOptions() *loggerOpts {
	return &loggerOpts{
		writer:            os.Stdout,
		enableJSONHandler: true,
	}
}

type optionFn func(*loggerOpts)

func (f optionFn) apply(o *loggerOpts) {
	f(o)
}

// WithWriter sets the output of a JSON logger. Defaults to os.Stdout.
func WithWriter(w io.Writer) Option {
	return optionFn(func(o *loggerOpts) {
		o.writer = w
	})
}

// WithSuppressTime suppresses the time field of a log record. This is intended for testing where a
// deterministic log record is required.
func WithSuppressTime(suppress bool) Option {
	return optionFn(func(o *loggerOpts) {
		o.suppressTime = suppress
	})
}

// WithSuppressSource suppresses the source field of a log record. This is intended for testing
// where a deterministic log record is required.
func WithSuppressSource(suppress bool) Option {
	return optionFn(func(o *loggerOpts) {
		o.suppressSource = suppress
	})
}

// WithExtractors adds context extractors to the logger.
func WithExtractors(extractors ...Extractor) Option {
	return optionFn(func(o *loggerOpts) {
		o.extractors = append(o.extractors, extractors...)
	})
}

// EnableJSONHandler enables the output of JSON logs to the configured writer.
func EnableJSONHandler(enabled bool) Option {
	return optionFn(func(o *loggerOpts) {
		o.enableJSONHandler = enabled
	})
}

// EnableOtelHandler enables the OpenTelemetry log handler which exports logs to the configured
// endpoint.
func EnableOtelHandler(enabled bool) Option {
	return optionFn(func(o *loggerOpts) {
		o.enableOtelHandler = enabled
	})
}
