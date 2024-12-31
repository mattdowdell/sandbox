//go:build wireinject
// +build wireinject

package main

import (
	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
	"connectrpc.com/validate"
	"github.com/google/wire"

	"github.com/mattdowdell/sandbox/internal/adapters/examplerpc"
	"github.com/mattdowdell/sandbox/internal/adapters/healthrpc"
	"github.com/mattdowdell/sandbox/internal/adapters/reflectrpc"
	"github.com/mattdowdell/sandbox/internal/drivers/config"
	"github.com/mattdowdell/sandbox/internal/drivers/config/flagoptions"
	"github.com/mattdowdell/sandbox/internal/drivers/logging"
	"github.com/mattdowdell/sandbox/internal/drivers/otelx"
	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver"
	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver/interceptors/otelconnectx"
)

func ProvideApp() (*App, error) {
	wire.Build(
		// config
		flagoptions.New,
		config.New,
		LoadConfig,
		wire.FieldsOf(new(Config), "App", "Logging", "OtelConnect", "RPCServer"),
		// observability
		logging.NewAsDefaultFromConfig,
		otelx.NewTracerProvider,
		// middleware
		validateInterceptor,
		otelconnectx.NewFromConfig,
		collectInterceptors,
		collectHandlerOptions,
		// handlers
		examplerpc.New,
		reflectrpc.New,
		healthrpc.New,
		// servers
		collectHandlers,
		rpcserver.NewFromConfig,
		// app
		NewApp,
	)

	return &App{}, nil
}

// collectHandlers merges multiple Handler implementations into a slice.
//
// While wire can cast a struct to an interface, it gets confused if multiple instances of a type
// are present. For more details, see https://github.com/google/wire/issues/207.
func collectHandlers(
	example *examplerpc.Handler,
	reflect *reflectrpc.Handler,
	health *healthrpc.Handler,
) []rpcserver.Handler {
	return []rpcserver.Handler{
		example,
		reflect,
		health,
	}
}

// ...
func validateInterceptor() (*validate.Interceptor, error) {
	return validate.NewInterceptor()
}

// ...
func collectInterceptors(
	validat *validate.Interceptor,
	otelconnec *otelconnect.Interceptor,
) []connect.Interceptor {
	return []connect.Interceptor{
		validat,
		otelconnec,
	}
}

// ...
func collectHandlerOptions(
	interceptors []connect.Interceptor,
) []connect.HandlerOption {
	return []connect.HandlerOption{
		connect.WithInterceptors(interceptors...),
	}
}
