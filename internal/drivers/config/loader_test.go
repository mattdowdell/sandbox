package config_test

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mattdowdell/sandbox/internal/drivers/config"
)

type TestConfig struct {
	Foo int           `koanf:"foo" default:"1"`
	Bar time.Duration `koanf:"bar" default:"1s"`
	Baz string        `koanf:"baz" default:"baz"`
	Sub TestSubConfig `koanf:"sub"`
}

type TestSubConfig struct {
	Foo  int           `koanf:"foo"`
	Bar  time.Duration `koanf:"bar"`
	Baz  string        `koanf:"baz"`
	Quux slog.Level    `koanf:"quux"`
}

type BadDefaultConfig struct {
	Foo int `koanf:"foo" default:"invalid"`
}

func Test_New(t *testing.T) {
	// arrange
	options := &config.Options{}

	// act
	loader, err := config.New[TestConfig](options)

	// assert
	assert.NotNil(t, loader)
	assert.NoError(t, err)
}

func Test_Load_FromEnv(t *testing.T) {
	tests := map[string]struct {
		envPrefix string
		env       map[string]string
		files     map[string]string
		want      *TestConfig
	}{
		"defaults": {
			want: &TestConfig{
				Foo: 1,
				Bar: time.Second,
				Baz: "baz",
			},
		},
		"with env (uppercase)": {
			env: map[string]string{
				"FOO":      "2",
				"BAR":      "1h",
				"BAZ":      "baz2",
				"SUB_FOO":  "4",
				"SUB_BAR":  "1m",
				"SUB_BAZ":  "baz3",
				"SUB_QUUX": "DEBUG",
			},
			want: &TestConfig{
				Foo: 2,
				Bar: time.Hour,
				Baz: "baz2",
				Sub: TestSubConfig{
					Foo:  4,
					Bar:  time.Minute,
					Baz:  "baz3",
					Quux: slog.LevelDebug,
				},
			},
		},
		"with env (lowercase)": {
			env: map[string]string{
				"foo":      "2",
				"bar":      "1h",
				"baz":      "baz2",
				"sub_foo":  "4",
				"sub_bar":  "1m",
				"sub_baz":  "baz3",
				"sub_quux": "DEBUG",
			},
			want: &TestConfig{
				Foo: 2,
				Bar: time.Hour,
				Baz: "baz2",
				Sub: TestSubConfig{
					Foo:  4,
					Bar:  time.Minute,
					Baz:  "baz3",
					Quux: slog.LevelDebug,
				},
			},
		},
		"with env and env prefix": {
			envPrefix: "APP_",
			env: map[string]string{
				"APP_FOO":      "2",
				"APP_BAR":      "1h",
				"APP_BAZ":      "baz2",
				"APP_SUB_FOO":  "4",
				"APP_SUB_BAR":  "1m",
				"APP_SUB_BAZ":  "baz3",
				"APP_SUB_QUUX": "DEBUG",
			},
			want: &TestConfig{
				Foo: 2,
				Bar: time.Hour,
				Baz: "baz2",
				Sub: TestSubConfig{
					Foo:  4,
					Bar:  time.Minute,
					Baz:  "baz3",
					Quux: slog.LevelDebug,
				},
			},
		},
		"with json file (typed value)": {
			files: map[string]string{
				"example.json": `{"foo":2, "sub":{"foo": 3}}`,
			},
			want: &TestConfig{
				Foo: 2,
				Bar: time.Second,
				Baz: "baz",
				Sub: TestSubConfig{
					Foo: 3,
				},
			},
		},
		"with json file (string value)": {
			files: map[string]string{
				"example.json": `{"foo":"2", "sub":{"foo": "3"}}`,
			},
			want: &TestConfig{
				Foo: 2,
				Bar: time.Second,
				Baz: "baz",
				Sub: TestSubConfig{
					Foo: 3,
				},
			},
		},
		"with yaml file": {
			files: map[string]string{
				"example.yaml": `{"foo":2, "sub":{"foo": 3}}`,
			},
			want: &TestConfig{
				Foo: 2,
				Bar: time.Second,
				Baz: "baz",
				Sub: TestSubConfig{
					Foo: 3,
				},
			},
		},
		"with yml file": {
			files: map[string]string{
				"example.yml": `{"foo":"2", "sub":{"foo": "3"}}`,
			},
			want: &TestConfig{
				Foo: 2,
				Bar: time.Second,
				Baz: "baz",
				Sub: TestSubConfig{
					Foo: 3,
				},
			},
		},
		"with toml file": {
			files: map[string]string{
				"example.toml": "foo = 2",
			},
			want: &TestConfig{
				Foo: 2,
				Bar: time.Second,
				Baz: "baz",
			},
		},
		"file overrides env": {
			env: map[string]string{
				"FOO": "5",
			},
			files: map[string]string{
				"example.json": `{"foo":10}`,
			},
			want: &TestConfig{
				Foo: 10,
				Bar: time.Second,
				Baz: "baz",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			for key, value := range tt.env {
				t.Setenv(key, value)
			}

			fileDir := t.TempDir()
			files := make([]string, 0, len(tt.files))

			for filename, content := range tt.files {
				path := filepath.Join(fileDir, filename)
				files = append(files, path)

				require.NoError(t, os.WriteFile(path, []byte(content), 0o400))
			}

			options := &config.Options{
				EnvPrefix: tt.envPrefix,
				Files:     files,
			}

			// act
			got, err := config.Load[TestConfig](options)

			// assert
			assert.Equal(t, tt.want, got)
			assert.NoError(t, err)
		})
	}
}

func Test_Load_FileError(t *testing.T) {
	tests := map[string]struct {
		have []string
		want string
	}{
		"file does not exist": {
			have: []string{"example.json"},
			want: "open example.json: no such file or directory",
		},
		"unhandled file extension": {
			have: []string{"example.txt"},
			want: `unsupported file extension for path: "example.txt"`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			options := &config.Options{
				Files: tt.have,
			}

			// act
			got, err := config.Load[TestConfig](options)

			// assert
			assert.Nil(t, got)
			assert.EqualError(t, err, tt.want)
		})
	}
}

func Test_Load_UnmarshalError(t *testing.T) {
	tests := map[string]struct {
		have map[string]string
		want string
	}{
		"single parse error": {
			have: map[string]string{
				"FOO": "invalid",
			},
			want: "decoding failed due to the following error(s):\n\n" +
				"'foo' cannot parse value as 'int': strconv.ParseInt: invalid syntax",
		},
		"multiple parse error": {
			have: map[string]string{
				"BAR":     "invalid",
				"SUB_FOO": "invalid",
			},
			want: "decoding failed due to the following error(s):\n\n" +
				"'bar' time: invalid duration\n" +
				"'sub.foo' cannot parse value as 'int': strconv.ParseInt: invalid syntax",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			options := &config.Options{}

			for key, value := range tt.have {
				t.Setenv(key, value)
			}

			// act
			got, err := config.Load[TestConfig](options)

			// assert
			assert.Nil(t, got)
			assert.EqualError(t, err, tt.want)
		})
	}
}
