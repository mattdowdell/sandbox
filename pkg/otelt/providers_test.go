package otelt_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/pkg/otelt"
)

func Test_NewMeterProvider(t *testing.T) {
	// arrange

	// act
	provider, collect := otelt.NewMeterProvider()

	// assert
	assert.NotNil(t, provider)

	if assert.NotNil(t, collect) {
		metrics, err := collect(t.Context())

		assert.Empty(t, metrics.ScopeMetrics)
		assert.NoError(t, err)
	}
}

func Test_NewTracerProvider(t *testing.T) {
	// arrange

	// act
	provider, collect := otelt.NewTracerProvider()

	// assert
	assert.NotNil(t, provider)

	if assert.NotNil(t, collect) {
		spans := collect()
		assert.Empty(t, spans)
	}
}
