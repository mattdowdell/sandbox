package examplerpc_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mattdowdell/sandbox/internal/adapters/examplerpc"
	"github.com/mattdowdell/sandbox/mocks/adapters/mockexamplerpc"
)

func Test_New(t *testing.T) {
	// arrange
	resource := mockexamplerpc.NewResourceFacade(t)
	auditEvent := mockexamplerpc.NewAuditEventFacade(t)

	// act
	handler := examplerpc.New(resource, auditEvent)

	// assert
	assert.NotNil(t, handler)
}

func Test_Handler_Register(t *testing.T) {
	// arrange
	handler := examplerpc.New(
		mockexamplerpc.NewResourceFacade(t),
		mockexamplerpc.NewAuditEventFacade(t),
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
