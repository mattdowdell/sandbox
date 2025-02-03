package otelx

import (
	"os"
	"path/filepath"
	"runtime/debug"

	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/semconv/v1.26.0"
)

func newResource() (*resource.Resource, error) {
	exec, err := os.Executable()
	if err != nil {
		return nil, err
	}

	version := "(devel)"

	info, ok := debug.ReadBuildInfo()
	if ok {
		// TODO: will only start working in Go 1.24 with `go build`
		version = info.Main.Version
	}

	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(filepath.Base(exec)),
		semconv.ServiceVersionKey.String(version),
	), nil
}
