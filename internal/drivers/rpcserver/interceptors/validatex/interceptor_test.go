package validatex_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver/interceptors/validatex"
)

func Test_NewFromConfig(t *testing.T) {
	// arrange

	// act
	interceptor := validatex.New()

	// assert
	assert.NotNil(t, interceptor)
}
