package authn_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver/interceptors/authn"
)

const (
	testAddress = "http://localhost:5000"
)

func Test_NewClientFromConfig(t *testing.T) {
	// arrange
	conf := authn.Config{
		Address: testAddress,
	}

	// act
	client := authn.NewClientFromConfig(http.DefaultClient, conf)

	// assert
	assert.NotNil(t, client)
}

func Test_NewClient(t *testing.T) {
	// arrange

	// act
	client := authn.NewClient(http.DefaultClient, testAddress)

	// assert
	assert.NotNil(t, client)
}
