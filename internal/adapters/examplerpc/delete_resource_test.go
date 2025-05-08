package examplerpc_test

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mattdowdell/sandbox/gen/example/v1"
	"github.com/mattdowdell/sandbox/internal/adapters/examplerpc"
	"github.com/mattdowdell/sandbox/internal/domain"
	"github.com/mattdowdell/sandbox/mocks/adapters/mockexamplerpc"
)

func Test_Handler_DeleteResource_Success(t *testing.T) {
	// arrange
	id := uuid.New()

	facade := mockexamplerpc.NewResourceFacade(t)
	facade.
		EXPECT().
		Delete(t.Context(), mock.AnythingOfType("*slog.Logger"), id).
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

func Test_Handler_DeleteResource_InvalidID(t *testing.T) {
	// arrange
	handler := examplerpc.New(
		mockexamplerpc.NewResourceFacade(t),
		mockexamplerpc.NewAuditEventFacade(t),
	)

	req := connect.NewRequest(&examplev1.DeleteResourceRequest{
		Id: "invalid",
	})

	// act
	resp, err := handler.DeleteResource(t.Context(), req)

	// assert
	assert.Nil(t, resp)
	assert.EqualError(t, err, "internal: internal error")
}

func Test_Handler_DeleteResource_NotFound(t *testing.T) {
	// arrange
	id := uuid.New()

	facade := mockexamplerpc.NewResourceFacade(t)
	facade.
		EXPECT().
		Delete(t.Context(), mock.AnythingOfType("*slog.Logger"), id).
		Return(domain.ErrNotFound).
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
	assert.Nil(t, resp)
	assert.EqualError(t, err, "not_found: resource does not exist")
}

func Test_Handler_DeleteResource_Internal(t *testing.T) {
	// arrange
	id := uuid.New()

	facade := mockexamplerpc.NewResourceFacade(t)
	facade.
		EXPECT().
		Delete(t.Context(), mock.AnythingOfType("*slog.Logger"), id).
		Return(domain.ErrInternal).
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
	assert.Nil(t, resp)
	assert.EqualError(t, err, "internal: internal error")
}
