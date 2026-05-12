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
	buildInfo     *debug.BuildInfo
	nowFunc       func() time.Time
	targetVersion int64
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

type optionFn func(*herdOpts)

func (f optionFn) apply(o *herdOpts) {
	f(o)
}

// WithBuildInfo overrides the soure of the code version and revision used in migration records.
// Defaults to the output of [runtime/debug.ReadBuildInfo].
//
// This option is intended to support unit testing, or when build info is otherwise unavailable.
func WithBuildInfo(info *debug.BuildInfo) Option {
	return optionFn(func(o *herdOpts) {
		o.buildInfo = info
	})
}

// WithBuildInfoValues wraps WithBuildInfo, constructing a [runtime/debug.BuildInfo] instance from
// the given values.
//
//   - version should be the version of the application using Herd, e.g. a git tag.
//   - revision should be the VCS revision at build time, e.g. a git commit.
//
// Both values must be non-empty.
func WithBuildInfoValues(version, revision string) Option {
	return optionFn(func(o *herdOpts) {
		o.buildInfo = &debug.BuildInfo{
			Main: debug.Module{
				Version: version,
			},
			Settings: []debug.BuildSetting{
				{
					Key:   "vcs.revision",
					Value: revision,
				},
			},
		}
	})
}

// WithNowFunc overrides the use of [time.Now] for recording when a migration was applied.
//
// This option is intended to support unit testing only.
func WithNowFunc(fn func() time.Time) Option {
	return optionFn(func(o *herdOpts) {
		o.nowFunc = fn
	})
}

// WithTargetVersion causes pending migrations with a version greater than the given value to be
// skipped. By default, all pending migrations are applied, regardless of version.
//
// This option is intended to allow a number of migrations to be applied, before inserting data and
// applying the final migration(s). This is useful when simulating a migration on a production
// database using representative data.
func WithTargetVersion(version int64) Option {
	return optionFn(func(o *herdOpts) {
		o.targetVersion = version
	})
}
