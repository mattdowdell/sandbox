package usecases_test

import (
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/domain"
	"github.com/mattdowdell/sandbox/internal/usecases"
	"github.com/mattdowdell/sandbox/mocks/domain/mockrepositories"
	"github.com/mattdowdell/sandbox/pkg/slogt"
)

func Test_NewDeleteResource(t *testing.T) {
	// no arrange necessary

	// act
	usecase := usecases.NewDeleteResource()

	// assert
	assert.NotNil(t, usecase)
}

func Test_DeleteResource_Success(t *testing.T) {
	// arrange
	id := uuid.Must(uuid.NewV7())

	usecase := usecases.NewDeleteResource()
	logger := slogt.New(t)

	store := mockrepositories.NewResource(t)
	store.EXPECT().DeleteResource(t.Context(), id).Return(nil).Once()

	// act
	err := usecase.Execute(t.Context(), logger, store, id)

	// assert
	assert.NoError(t, err)
}

func Test_DeleteResource_Error(t *testing.T) {
	testCases := []struct {
		name string
		err  error
	}{
		{
			name: "not found",
			err:  domain.ErrNotFound,
		},
		{
			name: "internal",
			err:  domain.ErrInternal,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			id := uuid.Must(uuid.NewV7())

			usecase := usecases.NewDeleteResource()
			logger := slogt.New(t)

			store := mockrepositories.NewResource(t)
			store.EXPECT().DeleteResource(t.Context(), id).Return(tc.err).Once()

			// act
			err := usecase.Execute(t.Context(), logger, store, id)

			// assert
			assert.ErrorIs(t, err, tc.err)
		})
	}
}
