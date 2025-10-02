package usecases_test

import (
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/domain"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
	"github.com/mattdowdell/sandbox/internal/usecases"
	"github.com/mattdowdell/sandbox/mocks/domain/mockrepositories"
	"github.com/mattdowdell/sandbox/pkg/slogt"
)

func Test_NewUpdateResource(t *testing.T) {
	// arrange
	timer := mockrepositories.NewTimer(t)

	// act
	usecase := usecases.NewUpdateResource(timer)

	// assert
	assert.NotNil(t, usecase)
}

func Test_UpdateResource_Success(t *testing.T) {
	// arrange
	id := uuid.Must(uuid.NewV7())
	now := time.Now().UTC().Round(time.Second)

	timer := mockrepositories.NewTimer(t)
	timer.EXPECT().UTCNow().Return(now).Once()

	usecase := usecases.NewUpdateResource(timer)
	logger := slogt.New(t)

	expected := &entities.Resource{
		ID:        id,
		Name:      testResourceName,
		UpdatedAt: now,
	}

	store := mockrepositories.NewResource(t)
	store.EXPECT().UpdateResource(t.Context(), expected).Return(expected, nil).Once()

	changes := &entities.Resource{
		ID:   id,
		Name: testResourceName,
	}

	// act
	resource, err := usecase.Execute(t.Context(), logger, store, changes)

	// assert
	assert.Equal(t, expected, resource)
	assert.NoError(t, err)
}

func Test_UpdateResource_Error(t *testing.T) {
	tests := map[string]struct {
		err error
	}{
		"not found": {
			err: domain.ErrNotFound,
		},
		"already exists": {
			err: domain.ErrAlreadyExists,
		},
		"internal": {
			err: domain.ErrInternal,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			id := uuid.Must(uuid.NewV7())
			now := time.Now().UTC().Round(time.Second)

			timer := mockrepositories.NewTimer(t)
			timer.EXPECT().UTCNow().Return(now).Once()

			usecase := usecases.NewUpdateResource(timer)
			logger := slogt.New(t)

			expected := &entities.Resource{
				ID:        id,
				Name:      testResourceName,
				UpdatedAt: now,
			}

			store := mockrepositories.NewResource(t)
			store.EXPECT().UpdateResource(t.Context(), expected).Return(nil, tt.err).Once()

			changes := &entities.Resource{
				ID:   id,
				Name: testResourceName,
			}

			// act
			resource, err := usecase.Execute(t.Context(), logger, store, changes)

			// assert
			assert.Nil(t, resource)
			assert.ErrorIs(t, err, tt.err)
		})
	}
}
