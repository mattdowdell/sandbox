package datastore

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/mattdowdell/sandbox/internal/domain/entities"
)

// ...
//
// TODO: decide how to implement this, use poller or notify?
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
