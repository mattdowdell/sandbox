package datastore

import (
	"github.com/go-jet/jet/v2/qrm"

	"github.com/mattdowdell/sandbox/internal/domain/repositories"
)

type Stub struct {
	db qrm.DB
}

func NewStub() *Stub {
	return &Stub{}
}

// ...
type Datastore struct {
	db qrm.DB
}

// ...
func NewResource(db qrm.DB) repositories.Resource {
	return &Datastore{
		db: db,
	}
}

// ...
func NewAuditEvent(db qrm.DB) repositories.AuditEvent {
	return &Datastore{
		db: db,
	}
}
