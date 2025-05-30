package otelx

import (
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The output of "go version -m ./dist/example-health" with the "./dist/example-health: " prefix
// replaced with "go\t" and all whitespace stripped from the start of each line.
//
// The use of tabs as delimiters is required by the parser. Furthermore, each line cannot be
// prefixed with whitespace or it will be ignored by the parser.
const buildInfoData = `go	go1.24.0
path	github.com/mattdowdell/sandbox/cmd/example-health
mod	github.com/mattdowdell/sandbox	v0.0.120
dep	buf.build/gen/go/grpc/grpc/connectrpc/go	v1.18.1-20250429200738-0ee95b84c2c7.1
dep	buf.build/gen/go/grpc/grpc/protocolbuffers/go	v1.36.6-20250429200738-0ee95b84c2c7.1
dep	connectrpc.com/connect	v1.18.1
dep	github.com/creasty/defaults	v1.8.0
dep	github.com/fsnotify/fsnotify	v1.9.0
dep	github.com/go-logr/logr	v1.4.2
dep	github.com/go-logr/stdr	v1.2.2
dep	github.com/go-viper/mapstructure/v2	v2.2.1
dep	github.com/knadh/koanf/maps	v0.1.2
dep	github.com/knadh/koanf/parsers/json	v1.0.0
dep	github.com/knadh/koanf/parsers/toml	v0.1.0
dep	github.com/knadh/koanf/parsers/yaml	v1.0.0
dep	github.com/knadh/koanf/providers/env	v1.1.0
dep	github.com/knadh/koanf/providers/file	v1.2.0
dep	github.com/knadh/koanf/v2	v2.2.0
dep	github.com/mitchellh/copystructure	v1.2.0
dep	github.com/mitchellh/reflectwalk	v1.0.2
dep	github.com/pelletier/go-toml	v1.9.5
dep	github.com/samber/lo	v1.49.1
dep	github.com/samber/slog-multi	v1.4.0
dep	go.opentelemetry.io/auto/sdk	v1.1.0
dep	go.opentelemetry.io/contrib/bridges/otelslog	v0.11.0
dep	go.opentelemetry.io/otel	v1.36.0
dep	go.opentelemetry.io/otel/log	v0.12.2
dep	go.opentelemetry.io/otel/metric	v1.36.0
dep	go.opentelemetry.io/otel/trace	v1.36.0
dep	golang.org/x/sys	v0.33.0
dep	golang.org/x/text	v0.25.0
dep	google.golang.org/protobuf	v1.36.6
dep	gopkg.in/yaml.v3	v3.0.1
build	-buildmode=exe
build	-compiler=gc
build	-trimpath=true
build	CGO_ENABLED=0
build	GOARCH=arm64
build	GOOS=darwin
build	GOARM64=v8.0
build	vcs=git
build	vcs.revision=3e787e3884c31c32e4e3d83fb8cd87c0d4e28a88
build	vcs.time=2025-05-27T09:23:42Z
build	vcs.modified=true`

func Test_extractVersion(t *testing.T) {
	testCases := []struct {
		name string
		have string
		want string
	}{
		{
			name: "main",
			have: "github.com/mattdowdell/sandbox/cmd/example-health",
			want: "v0.0.120",
		},
		{
			name: "dependency",
			have: "go.opentelemetry.io/otel",
			want: "v1.36.0",
		},
		{
			name: "dependency subpackage",
			have: "go.opentelemetry.io/otel/subpackage",
			want: "v1.36.0",
		},
		{
			name: "submodule",
			have: "go.opentelemetry.io/otel/log",
			want: "v0.12.2",
		},
		{
			name: "submodule subpackage",
			have: "go.opentelemetry.io/otel/log/subpackage",
			want: "v0.12.2",
		},
		{
			name: "not found",
			have: "github.com/does/not/exist",
			want: "(unknown)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			info, err := debug.ParseBuildInfo(buildInfoData)
			require.NoError(t, err)

			// act
			got := extractVersion(tc.have, info)

			// assert
			assert.Equal(t, tc.want, got)
		})
	}
}
