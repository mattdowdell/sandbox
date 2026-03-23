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
	"github.com/knadh/koanf/v2"

	"github.com/mattdowdell/sandbox/internal/drivers/config/providers/k8smount"
)

// The delimiter to use for joining configuration keys.
const (
	delimiter = "."
)

// provider is a pairing of a provider and parser to save work when reloading configuration.
type provider struct {
	provider koanf.Provider
	parser   koanf.Parser
}

func envProvider(prefix string) *provider {
	return newProvider(
		env.Provider(delimiter, env.Opt{
			Prefix: prefix,
			TransformFunc: func(k, v string) (string, any) {
				return transformKey(k, prefix), v
			},
		}),
		nil, /*parser*/
	)
}

func fileProviders(paths []string) ([]*provider, error) {
	providers := make([]*provider, 0, len(paths))

	for _, path := range paths {
		p, err := fileProvider(path)
		if err != nil {
			return nil, err
		}

		providers = append(providers, p)
	}

	return providers, nil
}

func fileProvider(path string) (*provider, error) {
	parser, err := fileParser(path)
	if err != nil {
		return nil, err
	}

	return newProvider(file.Provider(path), parser), nil
}

func mountProviders(paths []string) []*provider {
	providers := make([]*provider, 0, len(paths))

	for _, path := range paths {
		providers = append(providers, mountProvider(path))
	}

	return providers
}

func mountProvider(path string) *provider {
	return newProvider(k8smount.Provider(path, "_" /*delimiter*/), nil)
}

func newProvider(prov koanf.Provider, parser koanf.Parser) *provider {
	return &provider{
		provider: prov,
		parser:   parser,
	}
}

func (p *provider) load(k *koanf.Koanf) error {
	return k.Load(p.provider, p.parser)
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
