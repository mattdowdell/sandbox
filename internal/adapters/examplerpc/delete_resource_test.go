package examplerpc_test

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/gen/example/v1"
	"github.com/mattdowdell/sandbox/internal/adapters/examplerpc"
	"github.com/mattdowdell/sandbox/mocks/adapters/mockexamplerpc"
)

func Test_Handler_DeleteResource_Success(t *testing.T) {
	// arrange
	id := uuid.New()

	facade := mockexamplerpc.NewResourceFacade(t)
	facade.
		EXPECT().
		Delete(t.Context(), id).
		Return(nil).
		Once()

	handler := examplerpc.New(
		facade,
		mockexamplerpc.NewAuditEventFacade(t),
	)

	req := connect.NewRequest(&examplev1.DeleteResourceRequest{
		Id: id.String(),
	})

	// act
	resp, err := handler.DeleteResource(t.Context(), req)

	// assert
	expected := connect.NewResponse(&examplev1.DeleteResourceResponse{})

	assert.Equal(t, expected, resp)
	assert.NoError(t, err)
}
