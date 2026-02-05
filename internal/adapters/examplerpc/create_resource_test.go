package examplerpc_test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/mattdowdell/sandbox/gen/example/v1"
	"github.com/mattdowdell/sandbox/internal/adapters/examplerpc"
	"github.com/mattdowdell/sandbox/internal/domain"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
	"github.com/mattdowdell/sandbox/mocks/adapters/mockexamplerpc"
	"github.com/mattdowdell/sandbox/pkg/slogt"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

const (
	testResourceName = "example"
)

func Test_Handler_CreateResource_Success(t *testing.T) {
	// arrange
	ctx := slogx.IntoContext(t.Context(), slogt.New(t))

	id := uuid.Must(uuid.NewV7())
	now := time.Now()

	facade := mockexamplerpc.NewResourceFacade(t)
	facade.
		EXPECT().
		Create(
			ctx,
			mock.AnythingOfType("*slog.Logger"),
			mock.AnythingOfType("*entities.Resource"),
		).
		RunAndReturn(func(
			_ context.Context,
			_ *slog.Logger,
			r *entities.Resource,
		) (*entities.Resource, error) {
			return &entities.Resource{
				ID:        id,
				Name:      r.Name,
				CreatedAt: now,
				UpdatedAt: now,
			}, nil
		}).
		Once()

	handler := examplerpc.New(
		facade,
		mockexamplerpc.NewAuditEventFacade(t),
	)

	req := connect.NewRequest(&examplev1.CreateResourceRequest{
		Resource: &examplev1.ResourceCreate{
			Name: testResourceName,
		},
	})

	// act
	resp, err := handler.CreateResource(ctx, req)

	// assert
	expected := connect.NewResponse(&examplev1.CreateResourceResponse{
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

func Test_Handler_CreateResource_AlreadyExists(t *testing.T) {
	tests := map[string]struct {
		have error
		want string
	}{
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

			facade := mockexamplerpc.NewResourceFacade(t)
			facade.
				EXPECT().
				Create(
					ctx,
					mock.AnythingOfType("*slog.Logger"),
					mock.AnythingOfType("*entities.Resource"),
				).
				Return(nil, tt.have).
				Once()

			handler := examplerpc.New(
				facade,
				mockexamplerpc.NewAuditEventFacade(t),
			)

			req := connect.NewRequest(&examplev1.CreateResourceRequest{
				Resource: &examplev1.ResourceCreate{
					Name: testResourceName,
				},
			})

			// act
			resp, err := handler.CreateResource(ctx, req)

			// assert
			assert.Nil(t, resp)
			assert.EqualError(t, err, tt.want)
		})
	}
}
