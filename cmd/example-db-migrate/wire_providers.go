package main

import (
	"go.opentelemetry.io/contrib/processors/baggagecopy"

	"github.com/mattdowdell/sandbox/internal/drivers/logging"
	"github.com/mattdowdell/sandbox/internal/drivers/otelx"
)

// ...
func baggageFilter() baggagecopy.Filter {
	return baggagecopy.AllowAllMembers
}

// loggerOptions provides logger configuration options.
func loggerOptions() []logging.Option {
	extractor := otelx.NewExtractor(otelx.WithSpanID(true), otelx.WithSampled(true))

	return []logging.Option{
		logging.WithExtractors(extractor),
	}
}
