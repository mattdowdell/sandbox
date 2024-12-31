package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/creasty/defaults"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"

	"github.com/mattdowdell/sandbox/internal/drivers/config/providers/k8smount"
)

// The delimiter to use for splitting and joining configuration keys.
const (
	delimiter = "."
)

// Options provides the values to bootstrap configuration collection.
type Options struct {
	// The prefix of environment variables to read configuration from. Matching environment
	// variables have the prefix removed, are converted to lowercase and any underscores are
	// replaced with ".". For example, "APP_LOG_LEVEL" with a prefix of "APP_" would become
	// "log.level".
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

	// The directories of Kubernetes pod mounts to read configuration from. The filenames of the
	// mounted values become the configuration keys. For example, a configmap field of "log.level"
	// can be accessed at the same name.
	Mounts []string
}

// Config provides loading of configuration values for a service.
type Config struct {
	inner *koanf.Koanf
	opts  *Options
}

// New creates a new Config instance.
func New(opts *Options) *Config {
	return &Config{
		inner: koanf.New(delimiter),
		opts:  opts,
	}
}

// Load reads configuration, using it to populate a struct. See Options for information on the
// supported configuration sources.
//
// The struct should contains fields with "koanf" tags identifying the configuration value to
// assign. For example:
//
//	type LoggingConfig struct {
//		Level string `koanf:"logging.level"`
//	}
//
// The "koanf" struct tag value must use "." delimited keys for nested values. The field type can be
// anything supported by [mapstructure], or an implementation of [encoding.TextUnmarshaler].
//
// A default value can be set using a "default" struct tag with the desired value. For example:
//
//	type LoggingConfig struct {
//		Level string `koanf:"logging.level" default:"info"`
//	}
//
// A field type can be anything supported by [defaults], or an implementation of [defaults.Setter],
// [encoding.TextUnmarshaler], or [encoding/json.Unmarshaler].
//
// [mapstructure]: https://pkg.go.dev/github.com/go-viper/mapstructure/v2
// [encoding.TextUnmarshaler]: https://pkg.go.dev/encoding#TextUnmarshaler
// [defaults]: https://pkg.go.dev/github.com/creasty/defaults
// [defaults.Setter]: https://pkg.go.dev/github.com/creasty/defaults#Setter
// [encoding/json.Unmarshaler]: https://pkg.go.dev/encoding/json#Unmarshaler
func Load[T any](conf *Config) (T, error) {
	var val T

	if err := conf.loadEnv(); err != nil {
		return val, err
	}

	if err := conf.loadFiles(); err != nil {
		return val, err
	}

	if err := conf.loadMounts(); err != nil {
		return val, err
	}

	if err := defaults.Set(&val); err != nil {
		return val, err
	}

	if err := conf.inner.Unmarshal("", &val); err != nil {
		var zeroVal T
		return zeroVal, err
	}

	return val, nil
}

func (c *Config) loadEnv() error {
	return c.inner.Load(envProvider(c.opts.EnvPrefix), nil)
}

func (c *Config) loadFiles() error {
	for _, path := range c.opts.Files {
		parser, err := fileParser(path)
		if err != nil {
			return err
		}

		if err := c.inner.Load(file.Provider(path), parser); err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) loadMounts() error {
	for _, path := range c.opts.Mounts {
		if err := c.inner.Load(k8smount.Provider(path, delimiter), nil); err != nil {
			return err
		}
	}

	return nil
}

func envProvider(prefix string) *env.Env {
	return env.Provider(prefix, delimiter, func(s string) string {
		return strings.ReplaceAll(
			strings.ToLower(strings.TrimPrefix(s, prefix)),
			"_",
			delimiter,
		)
	})
}

func fileParser(path string) (koanf.Parser, error) {
	switch filepath.Ext(path) {
	case ".json":
		return json.Parser(), nil

	case ".yaml", ".yml":
		return yaml.Parser(), nil

	case ".toml":
		return toml.Parser(), nil

	default:
		return nil, fmt.Errorf("supported file extension for path: %q", path)
	}
}
