package reflectrpc_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mattdowdell/sandbox/internal/adapters/reflectrpc"
)

func Test_New(t *testing.T) {
	// arrange

	// act
	handler := reflectrpc.New()

	// assert
	assert.NotNil(t, handler)
}

func Test_Handler_Register(t *testing.T) {
	// arrange
	handler := reflectrpc.New()
	mux := http.NewServeMux()

	// act
	handler.Register(mux, nil /*opts*/)

	// assert
	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodPost,
		"/grpc.reflection.v1.ServerReflection/",
		http.NoBody,
	)
	require.NoError(t, err)

	_, pattern := mux.Handler(req)
	assert.Equal(t, "/grpc.reflection.v1.ServerReflection/", pattern, "pattern not handled")
}
