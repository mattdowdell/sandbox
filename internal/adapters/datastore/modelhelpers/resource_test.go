package modelhelpers_test

import (
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/adapters/datastore/modelhelpers"
	"github.com/mattdowdell/sandbox/internal/adapters/datastore/models/postgres/public/model"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
)

const (
	testResourceName = "name"
)

func Test_ResourcesToDomain(t *testing.T) {
	// arrange
	id := uuid.Must(uuid.NewV7())
	now := time.Now()

	have := []*model.Resources{
		{
			ID:        modelhelpers.ToGoogleUUID(id),
			Name:      testResourceName,
			CreatedAt: now,
			UpdatedAt: now.Add(time.Hour),
		},
	}

	// act
	got := modelhelpers.ResourcesToDomain(have)

	// assert
	want := []*entities.Resource{
		{
			ID:        id,
			Name:      testResourceName,
			CreatedAt: now,
			UpdatedAt: now.Add(time.Hour),
		},
	}

	assert.Equal(t, want, got)
}

func Test_ResourceToDomain(t *testing.T) {
	// arrange
	id := uuid.Must(uuid.NewV7())
	now := time.Now()

	have := &model.Resources{
		ID:        modelhelpers.ToGoogleUUID(id),
		Name:      testResourceName,
		CreatedAt: now,
		UpdatedAt: now.Add(time.Hour),
	}

	// act
	got := modelhelpers.ResourceToDomain(have)

	// assert
	want := &entities.Resource{
		ID:        id,
		Name:      testResourceName,
		CreatedAt: now,
		UpdatedAt: now.Add(time.Hour),
	}

	assert.Equal(t, want, got)
}

func Test_ResourceFromDomain(t *testing.T) {
	// arrange
	id := uuid.Must(uuid.NewV7())
	now := time.Now()

	have := &entities.Resource{
		ID:        id,
		Name:      testResourceName,
		CreatedAt: now,
		UpdatedAt: now.Add(time.Hour),
	}

	// act
	got := modelhelpers.ResourceFromDomain(have)

	// assert
	want := &model.Resources{
		ID:        modelhelpers.ToGoogleUUID(id),
		Name:      testResourceName,
		CreatedAt: now,
		UpdatedAt: now.Add(time.Hour),
	}

	assert.Equal(t, want, got)
}
