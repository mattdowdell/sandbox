// Package config provides support for loading configuration from various sources into a struct. It
// is a thin wrapper around [Koanf].
//
// This is designed for use with microservices running in Kubernetes, but may work well elsewhere
// too. Configuration may come from the following sources:
//
//   - Environment variables, e.g. those added to a Pod's container directly or via a ConfigMap.
//   - JSON, YAML or TOML files, e.g. volume mounts on a Pod via a ConfigMap or Secret.
//   - Single value files, e.g. simple key-value pairs from a volume mount.
//
// [Koanf]: https://pkg.go.dev/github.com/knadh/koanf/v2
package config

import (
	"github.com/creasty/defaults"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/v2"

	"github.com/mattdowdell/sandbox/internal/drivers/config/providers/k8smount"
)

// Options provides the values to bootstrap configuration loading.
//
// This struct is intended to be populated on startup based on environment variables or CLI options.
// Support for this is provided by the [flagoptions] package which is based on the [flag] package.
// However, this can trivially be replaced if an alternative library is desired instead.
//
// [flagoptions]: https://pkg.go.dev/github.com/mattdowdell/sandbox/internal/drivers/config/flagoptions
type Options struct {
	// The prefix of environment variables to read configuration from. Matching environment
	// variables have the prefix removed, are converted to lowercase and any underscores ("_") are
	// replaced with periods ("."). For example, "APP_LOG_LEVEL" with a prefix of "APP_" would
	// become "log.level".
	EnvPrefix string

	// The file paths to read configuration from. Supported file formats and extensions are:
	//
	// - JSON: .json
	// - YAML: .yml, .yaml
	// - TOML: .toml
	//
	// The file content is read into a map where nested map keys are joined using ".". For example,
	// a JSON file containing {"log":{"level":"info"}} would become "log.level".
	Files []string

	// The directories of Kubernetes volume mounts to read configuration from. The filenames of the
	// mounted values become the configuration keys. Any underscores ("_") in the key are replaced
	// with periods ("."). For example, a configmap field of "log_level" would become "log.level".
	Mounts []string
}

// Loader provides loading of configuration values for a service into a struct.
type Loader struct {
	inner  *koanf.Koanf
	opts   *Options
	env    *env.Env
	files  []*fileWithParser
	mounts []*k8smount.K8SMount
}

// New creates a new Loader instance.
func New(opts *Options) (*Loader, error) {
	files, err := fileProviders(opts.Files)
	if err != nil {
		return nil, err
	}

	return &Loader{
		inner:  koanf.New(delimiter),
		opts:   opts,
		env:    envProvider(opts.EnvPrefix),
		files:  files,
		mounts: mountProviders(opts.Mounts),
	}, nil
}

// Load creates a new Loader instance and immediately calls its Load method.
func Load[T any](opts *Options) (*T, error) {
	loader, err := New(opts)
	if err != nil {
		return nil, err
	}

	conf := new(T)

	if err := loader.Load(conf); err != nil {
		return nil, err
	}

	return conf, nil
}

// Load reads configuration, using it to populate the given struct pointer.
//
// Configuration is first loaded from environment variables, then files, and finally Kubernetes
// mounts. The last loaded configuration value wins if any conflicts occur.
//
// The struct should contain fields with "koanf" tags identifying the configuration key to
// assign. For example:
//
//	type LoggingConfig struct {
//		// reads "level" into the field
//		Level string `koanf:"level"`
//	}
//
// The field type can be anything supported by [mapstructure], or an implementation of
// [encoding.TextUnmarshaler]. If configuration keys contain ".", then nested structs must be used
// to access the value. For example:
//
//	type Config struct {
//		Log struct{
//			// reads "log.level" into the field
//			Level string `koanf:"level"`
//		} `koanf:"log"`
//	}
//
// A default value can be set using a "default" struct tag with the desired value. For example:
//
//	type LoggingConfig struct {
//		Level string `koanf:"level" default:"info"`
//	}
//
// A default-able field type can be anything supported by [defaults], or an implementation of
// [defaults.Setter] or [encoding.TextUnmarshaler]. An invalid default value will be skipped and
// will not cause an error to be returned.
//
// [mapstructure]: https://pkg.go.dev/github.com/go-viper/mapstructure/v2
// [encoding.TextUnmarshaler]: https://pkg.go.dev/encoding#TextUnmarshaler
// [defaults]: https://pkg.go.dev/github.com/creasty/defaults
// [defaults.Setter]: https://pkg.go.dev/github.com/creasty/defaults#Setter
func (l *Loader) Load(target any) error {
	if err := l.inner.Load(l.env, nil /*parser*/); err != nil {
		return err
	}

	for _, file := range l.files {
		if err := l.inner.Load(file.provider, file.parser); err != nil {
			return err
		}
	}

	for _, mount := range l.mounts {
		if err := l.inner.Load(mount, nil /*parser*/); err != nil {
			return err
		}
	}

	if err := defaults.Set(target); err != nil {
		return err
	}

	if err := l.inner.Unmarshal("", target); err != nil {
		return err
	}

	return nil
}

func (l *Loader) Watch(fn func(any, error)) error {
	for _, file := range l.files {
		if err := file.provider.Watch(fn); err != nil {
			return err
		}
	}

	for _, mount := range l.mounts {
		if err := mount.Watch(fn); err != nil {
			return err
		}
	}

	return nil
}
