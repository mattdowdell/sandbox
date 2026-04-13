package usecases_test

import (
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/domain"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
	"github.com/mattdowdell/sandbox/internal/domain/repositories"
	"github.com/mattdowdell/sandbox/internal/domain/repositories/mockrepositories"
	"github.com/mattdowdell/sandbox/internal/usecases"
	"github.com/mattdowdell/sandbox/pkg/slogt"
)

const (
	testLimit = 50
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
	id := uuid.Must(uuid.NewV7())
	now := time.Now()

	usecase := usecases.NewListResources()
	logger := slogt.New(t)
	pager := repositories.Pager{
		Limit: testLimit,
	}

	expected := &repositories.Paged[*entities.Resource]{
		Items: []*entities.Resource{
			{
				ID:        id,
				Name:      testResourceName,
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
	}

	store := mockrepositories.NewMockResource(t)
	store.EXPECT().ListResources(t.Context(), pager).Return(expected, nil).Once()

	// act
	resources, err := usecase.Execute(t.Context(), logger, store, pager)

	// assert
	assert.Equal(t, expected, resources)
	assert.NoError(t, err)
}

func Test_ListResources_Error(t *testing.T) {
	// arrange
	usecase := usecases.NewListResources()
	logger := slogt.New(t)
	pager := repositories.Pager{
		Limit: testLimit,
	}

	store := mockrepositories.NewMockResource(t)
	store.EXPECT().ListResources(t.Context(), pager).Return(nil, domain.ErrInternal).Once()

	// act
	resources, err := usecase.Execute(t.Context(), logger, store, pager)

	// assert
	assert.Empty(t, resources)
	assert.ErrorIs(t, err, domain.ErrInternal)
}
