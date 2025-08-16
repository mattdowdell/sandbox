package main

import (
	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/otelconnect"
	"connectrpc.com/validate"
	"go.opentelemetry.io/contrib/processors/baggagecopy"

	"github.com/mattdowdell/sandbox/gen/authn/v1/authnv1connect"
	"github.com/mattdowdell/sandbox/internal/adapters/authnrpc"
	"github.com/mattdowdell/sandbox/internal/adapters/examplerpc"
	"github.com/mattdowdell/sandbox/internal/adapters/healthrpc"
	"github.com/mattdowdell/sandbox/internal/adapters/reflectrpc"
	"github.com/mattdowdell/sandbox/internal/drivers/logging"
	"github.com/mattdowdell/sandbox/internal/drivers/otelx"
	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver"
	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver/interceptors/authn"
	logginginterceptor "github.com/mattdowdell/sandbox/internal/drivers/rpcserver/interceptors/logging"
)

// ...
func baggageFilter() baggagecopy.Filter {
	return baggagecopy.AllowAllMembers
}

// loggerOptions provides logger configuration options.
func loggerOptions() []logging.Option {
	extractor := otelx.NewExtractor(otelx.WithSpanID(true), otelx.WithSampled(true))

	return []logging.Option{
		logging.WithExtractors(extractor),
	}
}

// collectHandlers merges multiple rpcserver.Handler implementations into a slice.
//
// While wire can cast a struct to an interface, it gets confused if multiple instances of a type
// are present. For more details, see [google/wire#207].
//
// [google/wire#207]: https://github.com/google/wire/issues/207.
func collectHandlers(
	authnHandler *authnrpc.Handler,
	exampleHandler *examplerpc.Handler,
	reflectHandler *reflectrpc.Handler,
	healthHandler *healthrpc.Handler,
) []rpcserver.Handler {
	return []rpcserver.Handler{
		authnHandler,
		exampleHandler,
		reflectHandler,
		healthHandler,
	}
}

// collectHandlers merges multiple [connect.Interceptor] implementations into a slice.
//
// While wire can cast a struct to an interface, it gets confused if multiple instances of a type
// are present. For more details, see [google/wire#207].
//
// The ordering of interceptors in the returned slice aligns with the order they are called. To
// avoid losing observability, the otelconnect interceptor must be called first.
//
// [connect.Interceptor]: https://pkg.go.dev/connectrpc.com/connect#Interceptor
// [google/wire#207]: https://github.com/google/wire/issues/207
func collectInterceptors(
	otelconnec *otelconnect.Interceptor,
	validat *validate.Interceptor,
	loggin *logginginterceptor.Interceptor,
	auth *authn.Interceptor,
) []connect.Interceptor {
	return []connect.Interceptor{
		otelconnec,
		validat,
		loggin,
		auth,
	}
}

// collectHandlerOptions merges multiple [connect.HandlerOption] implementations into a slice.
//
// While wire can cast a struct to an interface, it gets confused if multiple instances of a type
// are present. For more details, see [google/wire#207].
//
// The panic recovery interceptor is applied last meaning it exclusively applies to the called
// handler and not any other interceptors. See [connectrpc/connect-go#816] for further discussion.
//
// [connect.HandlerOption]: https://pkg.go.dev/connectrpc.com/connect#HandlerOption
// [google/wire#207]: https://github.com/google/wire/issues/207
// [connectrpc/connect-go#816]: https://github.com/connectrpc/connect-go/issues/816
func collectHandlerOptions(
	interceptors []connect.Interceptor,
	recoverer *rpcserver.Recoverer,
) []connect.HandlerOption {
	return []connect.HandlerOption{
		connect.WithInterceptors(interceptors...),
		connect.WithRecover(recoverer.Handle),
	}
}

// ...
func recovererOptions() []rpcserver.RecovererOption {
	return nil
}

// ...
func authnOptions() []authn.Option {
	return []authn.Option{
		authn.WithIgnoreService(
			grpchealth.HealthV1ServiceName,
			grpcreflect.ReflectV1ServiceName,
			grpcreflect.ReflectV1AlphaServiceName,
			authnv1connect.AuthnServiceName,
		),
	}
}
