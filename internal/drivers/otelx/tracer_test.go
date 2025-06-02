//nolint:dupl // similar to meter_test.go/tracer_test.go
package otelx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/drivers/otelx"
)

func Test_Tracer(t *testing.T) {
	provider, _ := otelx.TestTracerProvider()

	testCases := []struct {
		name    string
		options []otelx.TracerOption
	}{
		{
			name: "no options",
		},
		{
			name: "with provider",
			options: []otelx.TracerOption{
				otelx.WithTracerProvider(provider),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange

			// act
			tracer := otelx.Tracer(tc.options...)

			// assert
			assert.NotNil(t, tracer)
		})
	}
}
