package herd

import (
	"errors"
	"runtime/debug"
	"time"
)

// Option adjust the behaviour of Herd.
type Option interface {
	apply(*herdOpts)
}

type herdOpts struct {
	buildInfo *debug.BuildInfo
	nowFunc   func() time.Time
}

func (o *herdOpts) codeInfo() (version, revision string, err error) {
	if o.buildInfo == nil {
		return "", "", errors.New("build info is unavailable")
	}

	version = o.buildInfo.Main.Version
	if version == "" {
		return "", "", errors.New("unable to extract code version from build info")
	}

	for _, setting := range o.buildInfo.Settings {
		if setting.Key == "vcs.revision" {
			revision = setting.Value
		}
	}

	if revision == "" {
		return "", "", errors.New("unable to extract code revision from build info")
	}

	return version, revision, nil
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

// WithNowFunc overrides the use of [time.Now] for recording when a migration was applied.
func WithNowFunc(fn func() time.Time) Option {
	return &nowFuncOpt{fn}
}

func (o *nowFuncOpt) apply(opts *herdOpts) {
	opts.nowFunc = o.nowFunc
}
