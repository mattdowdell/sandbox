//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"connectrpc.com/connect"
	"github.com/google/wire"

	"github.com/mattdowdell/sandbox/internal/adapters/datastore"
	"github.com/mattdowdell/sandbox/internal/adapters/examplerpc"
	"github.com/mattdowdell/sandbox/internal/adapters/healthrpc"
	"github.com/mattdowdell/sandbox/internal/adapters/reflectrpc"
	"github.com/mattdowdell/sandbox/internal/domain/repositories"
	"github.com/mattdowdell/sandbox/internal/drivers/clock"
	"github.com/mattdowdell/sandbox/internal/drivers/config"
	"github.com/mattdowdell/sandbox/internal/drivers/config/flagoptions"
	"github.com/mattdowdell/sandbox/internal/drivers/logging"
	"github.com/mattdowdell/sandbox/internal/drivers/otelx"
	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver"
	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver/interceptors/otelconnectx"
	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver/interceptors/validatex"
	"github.com/mattdowdell/sandbox/internal/drivers/uuidgen"
	"github.com/mattdowdell/sandbox/internal/usecases"
)

func ProvideApp(ctx context.Context) (*App, error) {
	wire.Build(
		// config
		flagoptions.New,
		config.New,
		LoadConfig,
		wire.FieldsOf(new(Config), "App", "Logging", "Meter", "OtelConnect", "RPCServer", "Tracer"),
		// observability
		logging.NewAsDefaultFromConfig,
		otelx.NewTracerProvider,
		otelx.NewMeterProvider,
		// repositories
		clock.New,
		wire.Bind(new(repositories.Clock), new(*clock.Clock)),
		uuidgen.New,
		wire.Bind(new(repositories.UUIDGenerator), new(*uuidgen.Generator)),
		datastore.NewStub,
		wire.Bind(new(repositories.Resource), new(*datastore.Stub)),
		wire.Bind(new(repositories.AuditEvent), new(*datastore.Stub)),
		// usecases
		usecases.NewCreateResource,
		wire.Bind(new(examplerpc.ResourceCreator), new(*usecases.CreateResource)),
		usecases.NewGetResource,
		wire.Bind(new(examplerpc.ResourceGetter), new(*usecases.GetResource)),
		usecases.NewListResources,
		wire.Bind(new(examplerpc.ResourceLister), new(*usecases.ListResources)),
		usecases.NewUpdateResource,
		wire.Bind(new(examplerpc.ResourceUpdater), new(*usecases.UpdateResource)),
		usecases.NewDeleteResource,
		wire.Bind(new(examplerpc.ResourceDeleter), new(*usecases.DeleteResource)),
		usecases.NewListAuditEvents,
		wire.Bind(new(examplerpc.AuditEventLister), new(*usecases.ListAuditEvents)),
		usecases.NewWatchAuditEvents,
		wire.Bind(new(examplerpc.AuditEventWatcher), new(*usecases.WatchAuditEvents)),
		// middleware
		validatex.New,
		otelconnectx.NewFromConfig,
		collectInterceptors,
		collectHandlerOptions,
		// handlers
		examplerpc.New,
		reflectrpc.New,
		healthrpc.New,
		collectHandlers,
		// server
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
//
// TODO: document that the ordering here is important
func collectInterceptors(
	validate *validatex.Interceptor,
	otelconnect *otelconnectx.Interceptor,
) []connect.Interceptor {
	return []connect.Interceptor{
		otelconnect,
		validate,
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
