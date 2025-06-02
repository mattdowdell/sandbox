package otelx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/drivers/otelx"
)

func Test_Meter(t *testing.T) {
	provider, _ := otelx.TestMeterProvider()

	testCases := []struct {
		name    string
		options []otelx.MeterOption
	}{
		{
			name: "no options",
		},
		{
			name: "with provider",
			options: []otelx.MeterOption{
				otelx.WithMeterProvider(provider),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange

			// act
			meter := otelx.Meter(tc.options...)

			// assert
			assert.NotNil(t, meter)
		})
	}
}
