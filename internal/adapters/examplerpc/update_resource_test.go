package examplerpc_test

import (
	"context"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/mattdowdell/sandbox/gen/example/v1"
	"github.com/mattdowdell/sandbox/internal/adapters/examplerpc"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
	"github.com/mattdowdell/sandbox/mocks/adapters/mockexamplerpc"
)

func Test_Handler_UpdateResource_Success(t *testing.T) {
	// arrange
	id := uuid.New()
	now := time.Now()

	facade := mockexamplerpc.NewResourceFacade(t)
	facade.
		EXPECT().
		Update(t.Context(), mock.AnythingOfType("*entities.Resource")).
		RunAndReturn(func(
			_ context.Context,
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
	resp, err := handler.UpdateResource(t.Context(), req)

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
