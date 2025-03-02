package examplerpc

import (
	"net/http"

	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/gen/example/v1/examplev1connect"
	"github.com/mattdowdell/sandbox/internal/adapters/common"
)

// Non-allocating compile-time check for interface implementation.
var _ examplev1connect.ExampleServiceHandler = (*Handler)(nil)

// Handler implements the ExampleService RPC.
type Handler struct {
	provider          common.Provider
	resourceCreator   ResourceCreator
	resourceGetter    ResourceGetter
	resourceLister    ResourceLister
	resourceUpdater   ResourceUpdater
	resourceDeleter   ResourceDeleter
	auditEventLister  AuditEventLister
	auditEventWatcher AuditEventWatcher
}

// New creates a new Handler.
func New(
	provider common.Provider,
	resourceCreator ResourceCreator,
	resourceGetter ResourceGetter,
	resourceLister ResourceLister,
	resourceUpdater ResourceUpdater,
	resourceDeleter ResourceDeleter,
	auditEventLister AuditEventLister,
	auditEventWatcher AuditEventWatcher,
) *Handler {
	return &Handler{
		provider:          provider,
		resourceCreator:   resourceCreator,
		resourceGetter:    resourceGetter,
		resourceLister:    resourceLister,
		resourceUpdater:   resourceUpdater,
		resourceDeleter:   resourceDeleter,
		auditEventLister:  auditEventLister,
		auditEventWatcher: auditEventWatcher,
	}
}

// Register adds the handler to the given multiplexer.
func (h *Handler) Register(mux *http.ServeMux, opts []connect.HandlerOption) {
	mux.Handle(examplev1connect.NewExampleServiceHandler(h, opts...))
}
