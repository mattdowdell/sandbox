package otelx

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

// ...
type TracerProviderShutdown func(context.Context) error

// ...
//
// TODO: look at options
func NewTracerProvider() TracerProviderShutdown {
	provider := trace.NewTracerProvider()
	otel.SetTracerProvider(provider)

	return provider.Shutdown
}
