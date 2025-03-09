package entities

import (
	"time"

	"github.com/google/uuid"
)

// ...
type Resource struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ...
func (r *Resource) Init(id uuid.UUID, now time.Time) {
	r.ID = id
	r.CreatedAt = now.Round(time.Second)
	r.UpdatedAt = now.Round(time.Second)
}

// ...
func (r *Resource) Update(now time.Time) {
	r.UpdatedAt = now.Round(time.Second)
}
