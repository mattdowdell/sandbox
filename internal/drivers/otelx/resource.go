package otelx

import (
	"os"
	"path/filepath"

	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/semconv/v1.26.0"
)

func newResource() (*resource.Resource, error) {
	exec, err := os.Executable()
	if err != nil {
		return nil, err
	}

	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(filepath.Base(exec)),
		// semconv.ServiceVersionKey.String("1.0.0"), // TODO: revisit in Go 1.24
	), nil
}
