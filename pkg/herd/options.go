package herd

import "runtime/debug"

// ...
type Option interface {
	apply(*migratorOptions)
}

type migratorOptions struct {
	dryRun      bool
	codeVersion string
}

type dryRunOpt bool

// ...
func WithDryRun(enable bool) Option {
	return dryRunOpt(enable)
}

func (o dryRunOpt) apply(options *migratorOptions) {
	options.dryRun = bool(o)
}

type codeVersionOpt string

// ...
func WithCodeVersion(value string) Option {
	return codeVersionOpt(value)
}

func (o codeVersionOpt) apply(options *migratorOptions) {
	options.codeVersion = string(o)
}

// ...
func WithVCSRevisionCodeVersion() Option {
	info, ok := debug.ReadBuildInfo()
	if ok {
		for _, s := range info.Settings {
			if s.Key == "vcs.revision" {
				return WithCodeVersion(s.Value)
			}
		}
	}

	return WithCodeVersion("(unknown)")
}

// ...
func WithMainVersionCodeVersion() Option {
	info, ok := debug.ReadBuildInfo()
	if ok {
		return WithCodeVersion(info.Main.Version)
	}

	return WithCodeVersion("(unknown)")
}
