package usecases_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/neilotoole/slogt"
	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/domain"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
	"github.com/mattdowdell/sandbox/internal/usecases"
	"github.com/mattdowdell/sandbox/mocks/domain/mockrepositories"
)

func Test_NewListResources(t *testing.T) {
	// no arrange necessary

	// act
	usecase := usecases.NewListResources()

	// assert
	assert.NotNil(t, usecase)
}

func Test_ListResources_Success(t *testing.T) {
	// arrange
	id := uuid.New()
	now := time.Now()

	usecase := usecases.NewListResources()
	logger := slogt.New(t)

	expected := []*entities.Resource{
		{
			ID:        id,
			Name:      testResourceName,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	store := mockrepositories.NewResource(t)
	store.EXPECT().ListResources(t.Context()).Return(expected, nil).Once()

	// act
	resources, err := usecase.Execute(t.Context(), logger, store)

	// assert
	assert.Equal(t, expected, resources)
	assert.NoError(t, err)
}

func Test_ListResources_Error(t *testing.T) {
	// arrange
	usecase := usecases.NewListResources()
	logger := slogt.New(t)

	store := mockrepositories.NewResource(t)
	store.EXPECT().ListResources(t.Context()).Return(nil, domain.ErrInternal).Once()

	// act
	resources, err := usecase.Execute(t.Context(), logger, store)

	// assert
	assert.Empty(t, resources)
	assert.ErrorIs(t, err, domain.ErrInternal)
}
