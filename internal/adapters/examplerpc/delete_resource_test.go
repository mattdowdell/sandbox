package examplerpc_test

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mattdowdell/sandbox/gen/example/v1"
	"github.com/mattdowdell/sandbox/internal/adapters/examplerpc"
	"github.com/mattdowdell/sandbox/internal/domain"
	"github.com/mattdowdell/sandbox/mocks/adapters/mockexamplerpc"
	"github.com/mattdowdell/sandbox/pkg/slogt"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

func Test_Handler_DeleteResource_Success(t *testing.T) {
	// arrange
	ctx := slogx.IntoContext(t.Context(), slogt.New(t))
	id := uuid.Must(uuid.NewV7())

	facade := mockexamplerpc.NewResourceFacade(t)
	facade.
		EXPECT().
		Delete(ctx, mock.AnythingOfType("*slog.Logger"), id).
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
	resp, err := handler.DeleteResource(ctx, req)

	// assert
	expected := connect.NewResponse(&examplev1.DeleteResourceResponse{})

	assert.Equal(t, expected, resp)
	assert.NoError(t, err)
}

func Test_Handler_DeleteResource_InvalidID(t *testing.T) {
	// arrange
	ctx := slogx.IntoContext(t.Context(), slogt.New(t))

	handler := examplerpc.New(
		mockexamplerpc.NewResourceFacade(t),
		mockexamplerpc.NewAuditEventFacade(t),
	)

	req := connect.NewRequest(&examplev1.DeleteResourceRequest{
		Id: "invalid",
	})

	// act
	resp, err := handler.DeleteResource(ctx, req)

	// assert
	assert.Nil(t, resp)
	assert.EqualError(t, err, "internal: internal error")
}

func Test_Handler_DeleteResource_UsecaseError(t *testing.T) {
	tests := map[string]struct {
		have error
		want string
	}{
		"not found": {
			have: domain.ErrNotFound,
			want: "not_found: resource does not exist",
		},
		"internal": {
			have: domain.ErrInternal,
			want: "internal: internal error",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			ctx := slogx.IntoContext(t.Context(), slogt.New(t))
			id := uuid.Must(uuid.NewV7())

			facade := mockexamplerpc.NewResourceFacade(t)
			facade.
				EXPECT().
				Delete(ctx, mock.AnythingOfType("*slog.Logger"), id).
				Return(tt.have).
				Once()

			handler := examplerpc.New(
				facade,
				mockexamplerpc.NewAuditEventFacade(t),
			)

			req := connect.NewRequest(&examplev1.DeleteResourceRequest{
				Id: id.String(),
			})

			// act
			resp, err := handler.DeleteResource(ctx, req)

			// assert
			assert.Nil(t, resp)
			assert.EqualError(t, err, tt.want)
		})
	}
}
