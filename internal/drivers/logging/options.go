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
	writer         io.Writer
	suppressTime   bool
	suppressSource bool
	extractors     []Extractor
}

func defaultOptions() *loggerOpts {
	return &loggerOpts{
		writer: os.Stdout,
	}
}

// WithWriter sets the output of a logger. Defaults to os.Stdout.
func WithWriter(w io.Writer) Option {
	return &writerOpt{
		w: w,
	}
}

type writerOpt struct {
	w io.Writer
}

func (o *writerOpt) apply(options *loggerOpts) {
	options.writer = o.w
}

// WithSuppressTime suppresses the time field of a log record. This is intended for testing where a
// deterministic log record is required.
func WithSuppressTime(suppress bool) Option {
	return suppressTimeOpt(suppress)
}

type suppressTimeOpt bool

func (o suppressTimeOpt) apply(options *loggerOpts) {
	options.suppressTime = bool(o)
}

// WithSuppressSource suppresses the source field of a log record. This is intended for testing
// where a deterministic log record is required.
func WithSuppressSource(suppress bool) Option {
	return suppressSourceOpt(suppress)
}

type suppressSourceOpt bool

func (o suppressSourceOpt) apply(options *loggerOpts) {
	options.suppressSource = bool(o)
}

// WithExtractors adds context extractors to the logger.
func WithExtractors(extractors ...Extractor) Option {
	return &extractorsOpt{
		e: extractors,
	}
}

type extractorsOpt struct {
	e []Extractor
}

func (o *extractorsOpt) apply(options *loggerOpts) {
	options.extractors = append(options.extractors, o.e...)
}
