package usecases_test

import (
	"errors"
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

const (
	testResourceName = "example"
)

func Test_NewCreateResource(t *testing.T) {
	// arrange
	clock := mockrepositories.NewMockClock(t)
	uuidgen := mockrepositories.NewMockUUIDGenerator(t)

	// act
	usecase := usecases.NewCreateResource(clock, uuidgen)

	// assert
	assert.NotNil(t, usecase)
}

func Test_CreateResource_Success(t *testing.T) {
	// arrange
	now := time.Now().UTC().Round(time.Second)
	id := uuid.Must(uuid.NewV7())

	clock := mockrepositories.NewMockClock(t)
	clock.EXPECT().UTCNow().Return(now).Once()

	uuidgen := mockrepositories.NewMockUUIDGenerator(t)
	uuidgen.EXPECT().NewV7().Return(id, nil).Once()

	usecase := usecases.NewCreateResource(clock, uuidgen)
	logger := slogt.New(t)

	expected := &entities.Resource{
		ID:        id,
		Name:      testResourceName,
		CreatedAt: now,
		UpdatedAt: now,
	}

	store := mockrepositories.NewMockResource(t)
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
	clock := mockrepositories.NewMockClock(t)

	uuidgen := mockrepositories.NewMockUUIDGenerator(t)
	uuidgen.EXPECT().NewV7().Return(uuid.Nil, errors.New("example")).Once()

	usecase := usecases.NewCreateResource(clock, uuidgen)
	logger := slogt.New(t)
	store := mockrepositories.NewMockResource(t)

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
	tests := map[string]struct {
		err error
	}{
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
			now := time.Now().UTC().Round(time.Second)
			id := uuid.Must(uuid.NewV7())

			clock := mockrepositories.NewMockClock(t)
			clock.EXPECT().UTCNow().Return(now).Once()

			uuidgen := mockrepositories.NewMockUUIDGenerator(t)
			uuidgen.EXPECT().NewV7().Return(id, nil).Once()

			usecase := usecases.NewCreateResource(clock, uuidgen)
			logger := slogt.New(t)

			expected := &entities.Resource{
				ID:        id,
				Name:      testResourceName,
				CreatedAt: now,
				UpdatedAt: now,
			}

			store := mockrepositories.NewMockResource(t)
			store.
				EXPECT().
				CreateResource(t.Context(), expected).
				Return(tt.err).
				Once()

			input := &entities.Resource{
				Name: testResourceName,
			}

			// act
			output, err := usecase.Execute(t.Context(), logger, store, input)

			// assert
			assert.Nil(t, output)
			assert.ErrorIs(t, err, tt.err)
		})
	}
}
