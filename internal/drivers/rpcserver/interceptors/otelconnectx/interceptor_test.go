package otelconnectx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver/interceptors/otelconnectx"
)

func Test_NewFromConfig(t *testing.T) {
	// arrange
	conf := otelconnectx.Config{
		TrustRemote: true,
	}

	// act
	interceptor, err := otelconnectx.NewFromConfig(conf)

	// assert
	assert.NotNil(t, interceptor)
	assert.NoError(t, err)
}
