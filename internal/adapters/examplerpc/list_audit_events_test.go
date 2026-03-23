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
	"github.com/mattdowdell/sandbox/mocks/adapters/mockexamplerpc"
	"github.com/mattdowdell/sandbox/pkg/slogt"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

func Test_Handler_ListAuditEvents_Success(t *testing.T) {
	// arrange
	ctx := slogx.IntoContext(t.Context(), slogt.New(t))

	id := uuid.Must(uuid.NewV7())
	now := time.Now()
	resourceID := uuid.Must(uuid.NewV7())

	facade := mockexamplerpc.NewAuditEventFacade(t)
	facade.
		EXPECT().
		List(ctx, mockLogger).
		Return(
			[]*entities.AuditEvent{
				{
					ID:           id,
					Operation:    entities.OperationCreated,
					CreatedAt:    now,
					Summary:      "summary",
					ResourceID:   resourceID,
					ResourceType: entities.ResourceTypeResource,
				},
			},
			nil,
		).
		Once()

	handler := examplerpc.New(
		mockexamplerpc.NewResourceFacade(t),
		facade,
	)

	req := connect.NewRequest(&examplev1.ListAuditEventsRequest{})

	// act
	resp, err := handler.ListAuditEvents(ctx, req)

	// assert
	expected := connect.NewResponse(&examplev1.ListAuditEventsResponse{
		Items: []*examplev1.AuditEvent{
			{
				Id:           id.String(),
				Operation:    examplev1.Operation_OPERATION_CREATED,
				CreatedAt:    timestamppb.New(now),
				Summary:      "summary",
				ResourceId:   resourceID.String(),
				ResourceType: "example.v1.Resource",
			},
		},
	})

	assert.Equal(t, expected, resp)
	assert.NoError(t, err)
}

func Test_Handler_ListAuditEvents_Error(t *testing.T) {
	// arrange
	ctx := slogx.IntoContext(t.Context(), slogt.New(t))

	facade := mockexamplerpc.NewAuditEventFacade(t)
	facade.
		EXPECT().
		List(ctx, mockLogger).
		Return(nil, domain.ErrInternal).
		Once()

	handler := examplerpc.New(
		mockexamplerpc.NewResourceFacade(t),
		facade,
	)

	req := connect.NewRequest(&examplev1.ListAuditEventsRequest{})

	// act
	resp, err := handler.ListAuditEvents(ctx, req)

	// assert
	assert.Nil(t, resp)
	assert.EqualError(t, err, "internal: internal error")
}
