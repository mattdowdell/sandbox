package examplerpc_test

import (
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/google/uuid"
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

func Test_Handler_ListResources_Success(t *testing.T) {
	// arrange
	ctx := slogx.IntoContext(t.Context(), slogt.New(t))

	id := uuid.New()
	now := time.Now()

	facade := mockexamplerpc.NewResourceFacade(t)
	facade.
		EXPECT().
		List(ctx, mock.AnythingOfType("*slog.Logger")).
		Return(
			[]*entities.Resource{
				{
					ID:        id,
					Name:      testResourceName,
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			nil,
		).
		Once()

	handler := examplerpc.New(
		facade,
		mockexamplerpc.NewAuditEventFacade(t),
	)

	req := connect.NewRequest(&examplev1.ListResourcesRequest{})

	// act
	resp, err := handler.ListResources(ctx, req)

	// assert
	expected := connect.NewResponse(&examplev1.ListResourcesResponse{
		Items: []*examplev1.Resource{
			{
				Id:        id.String(),
				Name:      testResourceName,
				CreatedAt: timestamppb.New(now),
				UpdatedAt: timestamppb.New(now),
			},
		},
	})

	assert.Equal(t, expected, resp)
	assert.NoError(t, err)
}

func Test_Handler_ListResources_Error(t *testing.T) {
	// arrange
	ctx := slogx.IntoContext(t.Context(), slogt.New(t))

	facade := mockexamplerpc.NewResourceFacade(t)
	facade.
		EXPECT().
		List(ctx, mock.AnythingOfType("*slog.Logger")).
		Return(nil, domain.ErrInternal).
		Once()

	handler := examplerpc.New(
		facade,
		mockexamplerpc.NewAuditEventFacade(t),
	)

	req := connect.NewRequest(&examplev1.ListResourcesRequest{})

	// act
	resp, err := handler.ListResources(ctx, req)

	// assert
	assert.Nil(t, resp)
	assert.EqualError(t, err, "internal: internal error")
}
