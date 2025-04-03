package workerpool_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/drivers/workerpool"
)

func Test_OK(t *testing.T) {
	// arrange

	// act
	result := workerpool.OK(5)

	// assert
	assert.Equal(t, 5, result.OK)
	assert.NoError(t, result.Err)

	assert.True(t, result.IsOK())
	assert.False(t, result.IsErr())
}

func Test_Err(t *testing.T) {
	// arrange
	err := errors.New("example")

	// act
	result := workerpool.Err[int](err)

	// assert
	assert.Empty(t, result.OK)
	assert.EqualError(t, result.Err, "example")

	assert.False(t, result.IsOK())
	assert.True(t, result.IsErr())
}
