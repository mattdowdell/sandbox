package examplerpc_test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/mattdowdell/sandbox/gen/example/v1"
	"github.com/mattdowdell/sandbox/internal/adapters/examplerpc"
	"github.com/mattdowdell/sandbox/internal/domain"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
	"github.com/mattdowdell/sandbox/mocks/adapters/mockexamplerpc"
	"github.com/mattdowdell/sandbox/pkg/slogt"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

func Test_Handler_UpdateResource_Success(t *testing.T) {
	// arrange
	ctx := slogx.IntoContext(t.Context(), slogt.New(t))

	id := uuid.Must(uuid.NewV7())
	now := time.Now()

	facade := mockexamplerpc.NewResourceFacade(t)
	facade.
		EXPECT().
		Update(ctx, mockLogger, mockResource).
		RunAndReturn(func(
			_ context.Context,
			_ *slog.Logger,
			r *entities.Resource,
		) (*entities.Resource, error) {
			return &entities.Resource{
				ID:        r.ID,
				Name:      r.Name,
				CreatedAt: now.Add(time.Hour * -1),
				UpdatedAt: now,
			}, nil
		}).
		Once()

	handler := examplerpc.New(
		facade,
		mockexamplerpc.NewAuditEventFacade(t),
	)

	req := connect.NewRequest(&examplev1.UpdateResourceRequest{
		Resource: &examplev1.ResourceUpdate{
			Id:   id.String(),
			Name: testResourceName + "2",
		},
	})

	// act
	resp, err := handler.UpdateResource(ctx, req)

	// assert
	expected := connect.NewResponse(&examplev1.UpdateResourceResponse{
		Resource: &examplev1.Resource{
			Id:        id.String(),
			Name:      testResourceName + "2",
			CreatedAt: timestamppb.New(now.Add(time.Hour * -1)),
			UpdatedAt: timestamppb.New(now),
		},
	})

	assert.Equal(t, expected, resp)
	assert.NoError(t, err)
}

func Test_Handler_UpdateResource_InvalidID(t *testing.T) {
	// arrange
	ctx := slogx.IntoContext(t.Context(), slogt.New(t))

	handler := examplerpc.New(
		mockexamplerpc.NewResourceFacade(t),
		mockexamplerpc.NewAuditEventFacade(t),
	)

	req := connect.NewRequest(&examplev1.UpdateResourceRequest{
		Resource: &examplev1.ResourceUpdate{
			Id:   "invalid",
			Name: testResourceName + "2",
		},
	})

	// act
	resp, err := handler.UpdateResource(ctx, req)

	// assert
	assert.Nil(t, resp)
	assert.EqualError(t, err, "internal: internal error")
}

func Test_Handler_UpdateResource_UsecaseError(t *testing.T) {
	tests := map[string]struct {
		have error
		want string
	}{
		"not found": {
			have: domain.ErrNotFound,
			want: "not_found: resource does not exist",
		},
		"already exists": {
			have: domain.ErrAlreadyExists,
			want: "already_exists: resource name already in use",
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
				Update(ctx, mockLogger, mockResource).
				Return(nil, tt.have).
				Once()

			handler := examplerpc.New(
				facade,
				mockexamplerpc.NewAuditEventFacade(t),
			)

			req := connect.NewRequest(&examplev1.UpdateResourceRequest{
				Resource: &examplev1.ResourceUpdate{
					Id:   id.String(),
					Name: testResourceName + "2",
				},
			})

			// act
			resp, err := handler.UpdateResource(ctx, req)

			// assert
			assert.Nil(t, resp)
			assert.EqualError(t, err, tt.want)
		})
	}
}
