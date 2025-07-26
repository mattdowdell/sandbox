package examplerpc_test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/gofrs/uuid/v5"
	"github.com/neilotoole/slogt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/mattdowdell/sandbox/gen/example/v1"
	"github.com/mattdowdell/sandbox/internal/adapters/examplerpc"
	"github.com/mattdowdell/sandbox/internal/domain"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
	"github.com/mattdowdell/sandbox/mocks/adapters/mockexamplerpc"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

func Test_Handler_GetResource_Success(t *testing.T) {
	// arrange
	ctx := slogx.IntoContext(t.Context(), slogt.New(t))

	id := uuid.Must(uuid.NewV7())
	now := time.Now()

	facade := mockexamplerpc.NewResourceFacade(t)
	facade.
		EXPECT().
		Get(ctx, mock.AnythingOfType("*slog.Logger"), id).
		RunAndReturn(func(
			_ context.Context,
			_ *slog.Logger,
			id uuid.UUID,
		) (*entities.Resource, error) {
			return &entities.Resource{
				ID:        id,
				Name:      testResourceName,
				CreatedAt: now,
				UpdatedAt: now,
			}, nil
		}).
		Once()

	handler := examplerpc.New(
		facade,
		mockexamplerpc.NewAuditEventFacade(t),
	)

	req := connect.NewRequest(&examplev1.GetResourceRequest{
		Id: id.String(),
	})

	// act
	resp, err := handler.GetResource(ctx, req)

	// assert
	expected := connect.NewResponse(&examplev1.GetResourceResponse{
		Resource: &examplev1.Resource{
			Id:        id.String(),
			Name:      testResourceName,
			CreatedAt: timestamppb.New(now),
			UpdatedAt: timestamppb.New(now),
		},
	})

	assert.Equal(t, expected, resp)
	assert.NoError(t, err)
}

func Test_Handler_GetResource_InvalidID(t *testing.T) {
	// arrange
	ctx := slogx.IntoContext(t.Context(), slogt.New(t))

	handler := examplerpc.New(
		mockexamplerpc.NewResourceFacade(t),
		mockexamplerpc.NewAuditEventFacade(t),
	)

	req := connect.NewRequest(&examplev1.GetResourceRequest{
		Id: "invalid",
	})

	// act
	resp, err := handler.GetResource(ctx, req)

	// assert
	assert.Nil(t, resp)
	assert.EqualError(t, err, "internal: internal error")
}

func Test_Handler_GetResource_UsecaseError(t *testing.T) {
	testCases := []struct {
		name string
		have error
		want string
	}{
		{
			name: "not found",
			have: domain.ErrNotFound,
			want: "not_found: resource does not exist",
		},
		{
			name: "internal",
			have: domain.ErrInternal,
			want: "internal: internal error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			ctx := slogx.IntoContext(t.Context(), slogt.New(t))
			id := uuid.Must(uuid.NewV7())

			facade := mockexamplerpc.NewResourceFacade(t)
			facade.
				EXPECT().
				Get(ctx, mock.AnythingOfType("*slog.Logger"), id).
				Return(nil, tc.have).
				Once()

			handler := examplerpc.New(
				facade,
				mockexamplerpc.NewAuditEventFacade(t),
			)

			req := connect.NewRequest(&examplev1.GetResourceRequest{
				Id: id.String(),
			})

			// act
			resp, err := handler.GetResource(ctx, req)

			// assert
			assert.Nil(t, resp)
			assert.EqualError(t, err, tc.want)
		})
	}
}
