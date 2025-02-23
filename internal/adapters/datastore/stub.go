package datastore

import (
	"github.com/go-jet/jet/v2/qrm"
)

type Stub struct {
	db qrm.DB
}

func NewStub() *Stub {
	return &Stub{}
}
