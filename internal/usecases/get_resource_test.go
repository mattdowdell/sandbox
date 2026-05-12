package usecases_test

import (
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/domain"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
	"github.com/mattdowdell/sandbox/internal/domain/repositories/mockrepositories"
	"github.com/mattdowdell/sandbox/internal/usecases"
	"github.com/mattdowdell/sandbox/pkg/slogt"
)

func Test_NewGetResource(t *testing.T) {
	// no arrange necessary

	// act
	usecase := usecases.NewGetResource()

	// assert
	assert.NotNil(t, usecase)
}

func Test_GetResource_Success(t *testing.T) {
	// arrange
	id := uuid.Must(uuid.NewV7())
	now := time.Now()

	usecase := usecases.NewGetResource()
	logger := slogt.New(t)

	expected := &entities.Resource{
		ID:        id,
		Name:      testResourceName,
		CreatedAt: now,
		UpdatedAt: now,
	}

	store := mockrepositories.NewMockResource(t)
	store.EXPECT().GetResource(t.Context(), id).Return(expected, nil).Once()

	// act
	resource, err := usecase.Execute(t.Context(), logger, store, id)

	// assert
	assert.Equal(t, expected, resource)
	assert.NoError(t, err)
}

func Test_GetResource_Error(t *testing.T) {
	tests := map[string]struct {
		err error
	}{
		"not found": {
			err: domain.ErrNotFound,
		},
		"internal": {
			err: domain.ErrInternal,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			id := uuid.Must(uuid.NewV7())

			usecase := usecases.NewGetResource()
			logger := slogt.New(t)

			store := mockrepositories.NewMockResource(t)
			store.EXPECT().GetResource(t.Context(), id).Return(nil, tt.err).Once()

			// act
			resource, err := usecase.Execute(t.Context(), logger, store, id)

			// assert
			assert.Nil(t, resource)
			assert.ErrorIs(t, err, tt.err)
		})
	}
}
