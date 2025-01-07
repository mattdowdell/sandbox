package healthrpc

import (
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"

	"github.com/mattdowdell/sandbox/gen/example/v1/examplev1connect"
)

// ...
type Handler struct{}

// ...
func New() *Handler {
	return &Handler{}
}

// Register adds the handler to the given multiplexer.
func (*Handler) Register(mux *http.ServeMux, opts []connect.HandlerOption) {
	checker := grpchealth.NewStaticChecker(
		examplev1connect.ExampleServiceName,
	)

	mux.Handle(grpchealth.NewHandler(checker, opts...))
}
