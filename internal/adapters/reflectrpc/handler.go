package reflectrpc

import (
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"

	"github.com/mattdowdell/sandbox/pkg/example/v1/examplev1connect"
)

// ...
type Handler struct{}

// ...
func New() *Handler {
	return &Handler{}
}

// Register adds the handler to the given multiplexer.
func (*Handler) Register(mux *http.ServeMux, opts []connect.HandlerOption) {
	reflector := grpcreflect.NewStaticReflector(
		examplev1connect.ExampleServiceName,
		grpchealth.HealthV1ServiceName,
	)

	mux.Handle(grpcreflect.NewHandlerV1(reflector, opts...))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector, opts...))
}
