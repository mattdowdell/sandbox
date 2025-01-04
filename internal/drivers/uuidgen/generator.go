package uuidgen

import (
	"github.com/google/uuid"

	"github.com/mattdowdell/sandbox/internal/domain/repositories"
)

// ...
var _ repositories.UUIDGenerator = (*Generator)(nil)

// ...
type Generator struct{}

// ...
func New() *Generator {
	return &Generator{}
}

// ...
func (g *Generator) NewV7() (uuid.UUID, error) {
	return uuid.NewV7()
}
