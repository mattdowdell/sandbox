package otelx

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"

	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/semconv/v1.26.0"
)

func newResource() (*resource.Resource, error) {
	exec, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to read executable path: %w", err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to read hostname: %w", err)
	}

	version := "(devel)"

	info, ok := debug.ReadBuildInfo()
	if ok {
		// TODO: will only start working in Go 1.24 with `go build`
		version = info.Main.Version
	}

	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(filepath.Base(exec)),
		semconv.ServiceVersion(version),
		semconv.ServiceInstanceID(hostname),
	), nil
}
