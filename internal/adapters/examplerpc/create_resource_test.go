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

const (
	testResourceName = "example"
)

func Test_Handler_CreateResource_Success(t *testing.T) {
	// arrange
	id := uuid.New()
	now := time.Now()

	facade := mockexamplerpc.NewResourceFacade(t)
	facade.
		EXPECT().
		Create(t.Context(), mock.AnythingOfType("*entities.Resource")).
		RunAndReturn(func(
			_ context.Context,
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
	resp, err := handler.CreateResource(t.Context(), req)

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
