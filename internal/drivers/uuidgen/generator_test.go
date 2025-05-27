package uuidgen_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/drivers/uuidgen"
)

func Test_Generator_NewV7(t *testing.T) {
	// arrange
	generator := uuidgen.New()

	// act
	id, err := generator.NewV7()

	// assert
	assert.Equal(t, 7, int(id.Version()))
	assert.NoError(t, err)
}
