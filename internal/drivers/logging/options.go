package logging

import (
	"io"
	"os"
)

// ...
type Option interface {
	apply(*loggerOpts)
}

type loggerOpts struct {
	writer         io.Writer
	suppressTime   bool
	suppressSource bool
}

func defaultOptions() *loggerOpts {
	return &loggerOpts{
		writer: os.Stdout,
	}
}

// ...
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

// ...
func WithSuppressTime(suppress bool) Option {
	return suppressTimeOpt(suppress)
}

type suppressTimeOpt bool

func (o suppressTimeOpt) apply(options *loggerOpts) {
	options.suppressTime = bool(o)
}

// ...
func WithSuppressSource(suppress bool) Option {
	return suppressSourceOpt(suppress)
}

type suppressSourceOpt bool

func (o suppressSourceOpt) apply(options *loggerOpts) {
	options.suppressSource = bool(o)
}
