package datastore

import (
	"context"
	"errors"
	"math/rand/v2"
	"time"

	"github.com/google/uuid"

	"github.com/mattdowdell/sandbox/internal/domain/entities"
)

type Stub struct{}

func NewStub() *Stub {
	return &Stub{}
}

func (s *Stub) CreateResource(context.Context, *entities.Resource) error {
	// temporarily introduce a small sleep to generate more interesting metrics for duration
	i := rand.IntN(1000)
	time.Sleep(time.Millisecond * time.Duration(i))

	return nil
}

func (s *Stub) GetResource(context.Context, uuid.UUID) (*entities.Resource, error) {
	return nil, errors.New("not implemented")
}

func (s *Stub) ListResources(context.Context) ([]*entities.Resource, error) {
	return nil, nil
}

func (s *Stub) UpdateResource(context.Context, *entities.Resource) error {
	return nil
}

func (s *Stub) DeleteResource(context.Context, uuid.UUID) error {
	return nil
}

func (s *Stub) CreateAuditEvent(context.Context, *entities.AuditEvent) error {
	return nil
}

func (s *Stub) ListAuditEvents(context.Context) ([]*entities.AuditEvent, error) {
	return nil, nil
}

func (s *Stub) WatchAuditEvents(ctx context.Context, ch chan<- *entities.AuditEvent) error {
	ticker := time.NewTicker(time.Second * 2) //nolint:mnd // ignore in stub

	for {
		select {
		case <-ctx.Done():
			close(ch)
			return nil

		case t := <-ticker.C:
			ch <- &entities.AuditEvent{
				ID:           uuid.Must(uuid.NewV7()),
				Operation:    entities.OperationCreated,
				CreatedAt:    t,
				Summary:      "summary",
				ResourceID:   uuid.Must(uuid.NewV7()),
				ResourceType: entities.ResourceTypeResource,
			}
		}
	}
}
