package authnrpc_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mattdowdell/sandbox/internal/adapters/authnrpc"
)

func Test_New(t *testing.T) {
	// arrange
	issuer := authnrpc.NewMockIssuer(t)
	parser := authnrpc.NewMockParser(t)

	// act
	handler := authnrpc.New(issuer, parser)

	// assert
	assert.NotNil(t, handler)
}

func Test_Handler_Register(t *testing.T) {
	// arrange
	handler := authnrpc.New(
		authnrpc.NewMockIssuer(t),
		authnrpc.NewMockParser(t),
	)

	mux := http.NewServeMux()

	// act
	handler.Register(mux, nil /*opts*/)

	// assert
	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodPost,
		"/authn.v1.AuthnService/",
		http.NoBody,
	)
	require.NoError(t, err)

	_, pattern := mux.Handler(req)
	assert.Equal(t, "/authn.v1.AuthnService/", pattern, "pattern not handled")
}
