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

func Test_NewGetResource(t *testing.T) {
	// no arrange necessary

	// act
	usecase := usecases.NewGetResource()

	// assert
	assert.NotNil(t, usecase)
}

func Test_GetResource_Success(t *testing.T) {
	// arrange
	id := uuid.New()
	now := time.Now()

	usecase := usecases.NewGetResource()
	logger := slog.New(slog.DiscardHandler)

	expected := &entities.Resource{
		ID:        id,
		Name:      testResourceName,
		CreatedAt: now,
		UpdatedAt: now,
	}

	store := mockrepositories.NewResource(t)
	store.EXPECT().GetResource(t.Context(), id).Return(expected, nil).Once()

	// act
	resource, err := usecase.Execute(t.Context(), logger, store, id)

	// assert
	assert.Equal(t, expected, resource)
	assert.NoError(t, err)
}

func Test_GetResource_Error(t *testing.T) {
	testCases := []struct {
		name string
		err  error
	}{
		{
			name: "not found",
			err:  apperrors.ErrNotFound,
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

			usecase := usecases.NewGetResource()
			logger := slog.New(slog.DiscardHandler)

			store := mockrepositories.NewResource(t)
			store.EXPECT().GetResource(t.Context(), id).Return(nil, tc.err).Once()

			// act
			resource, err := usecase.Execute(t.Context(), logger, store, id)

			// assert
			assert.Nil(t, resource)
			assert.ErrorIs(t, err, tc.err)
		})
	}
}
