package otelx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/drivers/otelx"
	"github.com/mattdowdell/sandbox/internal/drivers/otelx/otelt"
)

func Test_Meter(t *testing.T) {
	provider, _ := otelt.MeterProvider()

	tests := map[string]struct {
		options []otelx.MeterOption
	}{
		"no options": {},
		"with provider": {
			options: []otelx.MeterOption{
				otelx.WithMeterProvider(provider),
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange

			// act
			meter := otelx.Meter(tt.options...)

			// assert
			assert.NotNil(t, meter)
		})
	}
}
