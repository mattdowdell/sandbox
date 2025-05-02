package usecases_test

import (
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/domain/apperrors"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
	"github.com/mattdowdell/sandbox/internal/usecases"
	"github.com/mattdowdell/sandbox/mocks/domain/mockrepositories"
)

func Test_NewUpdateResource(t *testing.T) {
	// arrange
	clock := mockrepositories.NewClock(t)

	// act
	usecase := usecases.NewUpdateResource(clock)

	// assert
	assert.NotNil(t, usecase)
}

func Test_UpdateResource_Success(t *testing.T) {
	// arrange
	id := uuid.New()
	now := time.Now().UTC().Round(time.Second)

	clock := mockrepositories.NewClock(t)
	clock.EXPECT().Now().Return(now).Once()

	usecase := usecases.NewUpdateResource(clock)
	logger := slog.New(slog.DiscardHandler)

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
	testCases := []struct {
		name string
		err  error
	}{
		{
			name: "not found",
			err:  apperrors.ErrNotFound,
		},
		{
			name: "already exists",
			err:  apperrors.ErrAlreadyExists,
		},
		{
			name: "internal",
			err:  apperrors.ErrInternal,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			id := uuid.New()
			now := time.Now().UTC().Round(time.Second)

			clock := mockrepositories.NewClock(t)
			clock.EXPECT().Now().Return(now).Once()

			usecase := usecases.NewUpdateResource(clock)
			logger := slog.New(slog.DiscardHandler)

			expected := &entities.Resource{
				ID:        id,
				Name:      testResourceName,
				UpdatedAt: now,
			}

			store := mockrepositories.NewResource(t)
			store.EXPECT().UpdateResource(t.Context(), expected).Return(nil, tc.err).Once()

			changes := &entities.Resource{
				ID:   id,
				Name: testResourceName,
			}

			// act
			resource, err := usecase.Execute(t.Context(), logger, store, changes)

			// assert
			assert.Nil(t, resource)
			assert.ErrorIs(t, err, tc.err)
		})
	}
}
