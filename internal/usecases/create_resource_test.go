package usecases_test

import (
	"errors"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/neilotoole/slogt"
	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/domain"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
	"github.com/mattdowdell/sandbox/internal/usecases"
	"github.com/mattdowdell/sandbox/mocks/domain/mockrepositories"
)

const (
	testResourceName = "example"
)

func Test_NewCreateResource(t *testing.T) {
	// arrange
	clock := mockrepositories.NewClock(t)
	uuidgen := mockrepositories.NewUUIDGenerator(t)

	// act
	usecase := usecases.NewCreateResource(clock, uuidgen)

	// assert
	assert.NotNil(t, usecase)
}

func Test_CreateResource_Success(t *testing.T) {
	// arrange
	now := time.Now().UTC().Round(time.Second)
	id := uuid.Must(uuid.NewV7())

	clock := mockrepositories.NewClock(t)
	clock.EXPECT().Now().Return(now).Once()

	uuidgen := mockrepositories.NewUUIDGenerator(t)
	uuidgen.EXPECT().NewV7().Return(id, nil).Once()

	usecase := usecases.NewCreateResource(clock, uuidgen)
	logger := slogt.New(t)

	expected := &entities.Resource{
		ID:        id,
		Name:      testResourceName,
		CreatedAt: now,
		UpdatedAt: now,
	}

	store := mockrepositories.NewResource(t)
	store.EXPECT().CreateResource(t.Context(), expected).Return(nil).Once()

	input := &entities.Resource{
		Name: testResourceName,
	}

	// act
	output, err := usecase.Execute(t.Context(), logger, store, input)

	// assert
	assert.Equal(t, expected, output)
	assert.NoError(t, err)
}

func Test_CreateResource_IDFailed(t *testing.T) {
	// arrange
	clock := mockrepositories.NewClock(t)

	uuidgen := mockrepositories.NewUUIDGenerator(t)
	uuidgen.EXPECT().NewV7().Return(uuid.Nil, errors.New("example")).Once()

	usecase := usecases.NewCreateResource(clock, uuidgen)
	logger := slogt.New(t)
	store := mockrepositories.NewResource(t)

	input := &entities.Resource{
		Name: testResourceName,
	}

	// act
	output, err := usecase.Execute(t.Context(), logger, store, input)

	// assert
	assert.Nil(t, output)
	assert.ErrorIs(t, err, domain.ErrInternal)
}

func Test_CreateResource_CreateFailed(t *testing.T) {
	testCases := []struct {
		name string
		err  error
	}{
		{
			name: "already exists",
			err:  domain.ErrAlreadyExists,
		},
		{
			name: "internal",
			err:  domain.ErrInternal,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			now := time.Now().UTC().Round(time.Second)
			id := uuid.Must(uuid.NewV7())

			clock := mockrepositories.NewClock(t)
			clock.EXPECT().Now().Return(now).Once()

			uuidgen := mockrepositories.NewUUIDGenerator(t)
			uuidgen.EXPECT().NewV7().Return(id, nil).Once()

			usecase := usecases.NewCreateResource(clock, uuidgen)
			logger := slogt.New(t)

			expected := &entities.Resource{
				ID:        id,
				Name:      testResourceName,
				CreatedAt: now,
				UpdatedAt: now,
			}

			store := mockrepositories.NewResource(t)
			store.
				EXPECT().
				CreateResource(t.Context(), expected).
				Return(tc.err).
				Once()

			input := &entities.Resource{
				Name: testResourceName,
			}

			// act
			output, err := usecase.Execute(t.Context(), logger, store, input)

			// assert
			assert.Nil(t, output)
			assert.ErrorIs(t, err, tc.err)
		})
	}
}
