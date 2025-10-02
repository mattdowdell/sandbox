//nolint:dupl // similar to meter_test.go/tracer_test.go
package otelx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/drivers/otelx"
	"github.com/mattdowdell/sandbox/pkg/otelt"
)

func Test_Tracer(t *testing.T) {
	provider, _ := otelt.NewTracerProvider()

	tests := map[string]struct {
		options []otelx.TracerOption
	}{
		"no options": {},
		"with provider": {
			options: []otelx.TracerOption{
				otelx.WithTracerProvider(provider),
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange

			// act
			tracer := otelx.Tracer(tt.options...)

			// assert
			assert.NotNil(t, tracer)
		})
	}
}
