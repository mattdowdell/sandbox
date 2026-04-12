package examplerpc_test

import (
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
	"github.com/mattdowdell/sandbox/internal/domain/repositories"
	"github.com/mattdowdell/sandbox/pkg/slogt"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

const (
	testLimit = 50
)

func Test_Handler_ListResources_Success(t *testing.T) {
	// arrange
	ctx := slogx.IntoContext(t.Context(), slogt.New(t))

	id := uuid.Must(uuid.NewV7())
	now := time.Now()

	pager := repositories.Pager{
		Limit: testLimit,
	}

	facade := examplerpc.NewMockResourceFacade(t)
	facade.
		EXPECT().
		List(ctx, mockLogger, pager).
		Return(
			&repositories.Paged[*entities.Resource]{
				Items: []*entities.Resource{
					{
						ID:        id,
						Name:      testResourceName,
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
			},
			nil,
		).
		Once()

	handler := examplerpc.New(
		facade,
		examplerpc.NewMockAuditEventFacade(t),
	)

	req := connect.NewRequest(&examplev1.ListResourcesRequest{
		Limit: testLimit,
	})

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

	pager := repositories.Pager{
		Limit: testLimit,
	}

	facade := examplerpc.NewMockResourceFacade(t)
	facade.
		EXPECT().
		List(ctx, mockLogger, pager).
		Return(nil, domain.ErrInternal).
		Once()

	handler := examplerpc.New(
		facade,
		examplerpc.NewMockAuditEventFacade(t),
	)

	req := connect.NewRequest(&examplev1.ListResourcesRequest{
		Limit: testLimit,
	})

	// act
	resp, err := handler.ListResources(ctx, req)

	// assert
	assert.Nil(t, resp)
	assert.EqualError(t, err, "internal: internal error")
}
