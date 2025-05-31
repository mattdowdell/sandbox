package otelx

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

// MeterOption implementations customise the behaviour of Meter.
type MeterOption interface {
	apply(*meterOpts)
}

type meterOpts struct {
	provider metric.MeterProvider
}

type meterProviderOpt struct {
	provider metric.MeterProvider
}

// WithMeterProvider overrides the global MeterProvider. This can be useful for testing the result
// of recording metrics.
func WithMeterProvider(provider metric.MeterProvider) MeterOption {
	return &meterProviderOpt{
		provider: provider,
	}
}

func (o *meterProviderOpt) apply(options *meterOpts) {
	options.provider = o.provider
}

// Meter provides a [metric.Meter] with the package name and version automatically set based on the
// direct caller. It is advised to cache the result when possible to avoid computing the caller's
// package details unnecessarily.
//
// [metric.Meter]: https://pkg.go.dev/go.opentelemetry.io/otel/metric#Meter
func Meter(options ...MeterOption) metric.Meter {
	opts := &meterOpts{}
	for _, option := range options {
		option.apply(opts)
	}

	provider := opts.provider
	if provider == nil {
		provider = otel.GetMeterProvider()
	}

	pkg := packageName(1)
	ver := packageVersion(pkg)

	return provider.Meter(pkg, metric.WithInstrumentationVersion(ver))
}
