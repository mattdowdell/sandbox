package examplerpc_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mattdowdell/sandbox/internal/adapters/examplerpc"
)

func Test_New(t *testing.T) {
	// arrange
	resource := examplerpc.NewMockResourceFacade(t)
	auditEvent := examplerpc.NewMockAuditEventFacade(t)

	// act
	handler := examplerpc.New(resource, auditEvent)

	// assert
	assert.NotNil(t, handler)
}

func Test_Handler_Register(t *testing.T) {
	// arrange
	handler := examplerpc.New(
		examplerpc.NewMockResourceFacade(t),
		examplerpc.NewMockAuditEventFacade(t),
	)

	mux := http.NewServeMux()

	// act
	handler.Register(mux, nil /*opts*/)

	// assert
	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodPost,
		"/example.v1.ExampleService/",
		http.NoBody,
	)
	require.NoError(t, err)

	_, pattern := mux.Handler(req)
	assert.Equal(t, "/example.v1.ExampleService/", pattern, "pattern not handled")
}
