package models_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	examplev1 "github.com/mattdowdell/sandbox/gen/example/v1"
	"github.com/mattdowdell/sandbox/internal/adapters/examplerpc/models"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
)

const (
	testResourceName = "example"
)

func Test_ResourceCreateToDomain(t *testing.T) {
	// arrange
	input := &examplev1.ResourceCreate{
		Name: testResourceName,
	}

	// act
	output := models.ResourceCreateToDomain(input)

	// assert
	expected := &entities.Resource{
		Name: testResourceName,
	}

	assert.Equal(t, expected, output)
}

func Test_ResourceUpdateToDomain_Success(t *testing.T) {
	// arrange
	id := uuid.New()

	input := &examplev1.ResourceUpdate{
		Id:   id.String(),
		Name: testResourceName,
	}

	// act
	output, err := models.ResourceUpdateToDomain(input)

	// assert
	expected := &entities.Resource{
		ID:   id,
		Name: testResourceName,
	}

	assert.Equal(t, expected, output)
	assert.NoError(t, err)
}

func Test_ResourceUpdateToDomain_Error(t *testing.T) {
	// arrange
	input := &examplev1.ResourceUpdate{
		Id:   "invalid",
		Name: testResourceName,
	}

	// act
	output, err := models.ResourceUpdateToDomain(input)

	// assert
	assert.Nil(t, output)
	assert.EqualError(t, err, `failed to parse id "invalid": invalid UUID length: 7`)
}

func Test_ResourcesFromDomain(t *testing.T) {
	// arrange
	id := uuid.New()
	now := time.Now()

	input := []*entities.Resource{
		{
			ID:        id,
			Name:      testResourceName,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	// act
	output := models.ResourcesFromDomain(input)

	// assert
	expected := []*examplev1.Resource{
		{
			Id:        id.String(),
			Name:      testResourceName,
			CreatedAt: timestamppb.New(now),
			UpdatedAt: timestamppb.New(now),
		},
	}

	assert.Equal(t, expected, output)
}

func Test_ResourceFromDomain(t *testing.T) {
	// arrange
	id := uuid.New()
	now := time.Now()

	input := &entities.Resource{
		ID:        id,
		Name:      testResourceName,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// act
	output := models.ResourceFromDomain(input)

	// assert
	expected := &examplev1.Resource{
		Id:        id.String(),
		Name:      testResourceName,
		CreatedAt: timestamppb.New(now),
		UpdatedAt: timestamppb.New(now),
	}

	assert.Equal(t, expected, output)
}
