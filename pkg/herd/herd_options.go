package herd

import (
	"runtime/debug"
	"time"
)

// ...
type Option interface {
	apply(*herdOpts)
}

type herdOpts struct {
	buildInfo *debug.BuildInfo
	nowFunc   func() time.Time
}

func defaultHerdOpts() *herdOpts {
	info, _ := debug.ReadBuildInfo()

	return &herdOpts{
		buildInfo: info,
		nowFunc:   time.Now,
	}
}

type buildInfoOpt struct {
	info *debug.BuildInfo
}

// WithBuildInfo overrides the soure of the code version and revision used in migration records.
// Defaults to the output of [runtime/debug.ReadBuildInfo].
func WithBuildInfo(info *debug.BuildInfo) Option {
	return &buildInfoOpt{info}
}

func (o *buildInfoOpt) apply(opts *herdOpts) {
	opts.buildInfo = o.info
}

type nowFuncOpt struct {
	nowFunc func() time.Time
}

// WithNowFunc
func WithNowFunc(fn func() time.Time) Option {
	return &nowFuncOpt{fn}
}

func (o *nowFuncOpt) apply(opts *herdOpts) {
	opts.nowFunc = o.nowFunc
}
