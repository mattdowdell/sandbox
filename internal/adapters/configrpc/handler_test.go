package configrpc_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mattdowdell/sandbox/internal/adapters/configrpc"
	"github.com/mattdowdell/sandbox/mocks/adapters/mockconfigrpc"
)

type TestConfig struct {
	Foo string
}

func Test_New(t *testing.T) {
	// arrange
	loader := mockconfigrpc.NewLoader(t)

	// act
	handler := configrpc.New[TestConfig](loader)

	// assert
	assert.NotNil(t, handler)
}

func Test_Handler_Register(t *testing.T) {
	// arrange
	loader := mockconfigrpc.NewLoader(t)
	handler := configrpc.New[TestConfig](loader)
	mux := http.NewServeMux()

	// act
	handler.Register(mux, nil /*opts*/)

	// assert
	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodPost,
		"/config.v1.ConfigService/",
		http.NoBody,
	)
	require.NoError(t, err)

	_, pattern := mux.Handler(req)
	assert.Equal(t, "/config.v1.ConfigService/", pattern, "pattern not handled")
}
