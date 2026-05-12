package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/k8smount"
	"github.com/knadh/koanf/v2"
)

// The delimiter to use for joining configuration keys.
const (
	delimiter = "."
)

// fileWithParser is a pairing of a file provider and parser to save work when reloading
// configuration.
type fileWithParser struct {
	provider *file.File
	parser   koanf.Parser
}

func envProvider(prefix string) *env.Env {
	return env.Provider(delimiter, env.Opt{
		Prefix: prefix,
		TransformFunc: func(k, v string) (string, any) {
			return transformKey(k, prefix), v
		},
	})
}

func fileProviders(paths []string) ([]*fileWithParser, error) {
	providers := make([]*fileWithParser, 0, len(paths))

	for _, path := range paths {
		p, err := fileProvider(path)
		if err != nil {
			return nil, err
		}

		providers = append(providers, p)
	}

	return providers, nil
}

func fileProvider(path string) (*fileWithParser, error) {
	parser, err := fileParser(path)
	if err != nil {
		return nil, err
	}

	return &fileWithParser{
		provider: file.Provider(path),
		parser:   parser,
	}, nil
}

func mountProviders(paths []string) []*k8smount.K8SMount {
	providers := make([]*k8smount.K8SMount, 0, len(paths))

	for _, path := range paths {
		providers = append(providers, mountProvider(path))
	}

	return providers
}

func mountProvider(path string) *k8smount.K8SMount {
	return k8smount.Provider(path, delimiter, k8smount.Opt{
		TransformFunc: func(k, v string) (string, any) {
			return transformKey(k, "" /*prefix*/), v
		},
	})
}

func transformKey(key, prefix string) string {
	return strings.ReplaceAll(
		strings.ToLower(strings.TrimPrefix(key, prefix)),
		"_",
		delimiter,
	)
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
		return nil, fmt.Errorf("unsupported file extension for path: %q", path)
	}
}
