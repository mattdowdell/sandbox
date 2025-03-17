//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/google/wire"

	"github.com/mattdowdell/sandbox/internal/adapters/common"
	"github.com/mattdowdell/sandbox/internal/adapters/datastore"
	"github.com/mattdowdell/sandbox/internal/adapters/examplerpc"
	"github.com/mattdowdell/sandbox/internal/adapters/healthrpc"
	"github.com/mattdowdell/sandbox/internal/adapters/reflectrpc"
	"github.com/mattdowdell/sandbox/internal/adapters/usecasefacades"
	"github.com/mattdowdell/sandbox/internal/domain/repositories"
	"github.com/mattdowdell/sandbox/internal/drivers/clock"
	"github.com/mattdowdell/sandbox/internal/drivers/config"
	"github.com/mattdowdell/sandbox/internal/drivers/config/flagoptions"
	"github.com/mattdowdell/sandbox/internal/drivers/logging"
	"github.com/mattdowdell/sandbox/internal/drivers/otelx"
	"github.com/mattdowdell/sandbox/internal/drivers/pgsql"
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
		wire.FieldsOf(new(Config), "App", "Database", "Logging", "Meter", "OtelConnect", "RPCServer", "Tracer"),
		// observability
		logging.NewAsDefaultFromConfig,
		otelx.NewTracerProviderFromConfig,
		otelx.NewMeterProviderFromConfig,
		// providers
		pgsql.NewFromConfig,
		datastore.NewProvider,
		wire.Bind(new(common.Provider), new(*datastore.Provider)),
		// repositories
		clock.New,
		wire.Bind(new(repositories.Clock), new(*clock.Clock)),
		uuidgen.New,
		wire.Bind(new(repositories.UUIDGenerator), new(*uuidgen.Generator)),
		// usecases
		usecases.NewCreateResource,
		wire.Bind(new(usecasefacades.ResourceCreator), new(*usecases.CreateResource)),
		usecases.NewGetResource,
		wire.Bind(new(usecasefacades.ResourceGetter), new(*usecases.GetResource)),
		usecases.NewListResources,
		wire.Bind(new(usecasefacades.ResourceLister), new(*usecases.ListResources)),
		usecases.NewUpdateResource,
		wire.Bind(new(usecasefacades.ResourceUpdater), new(*usecases.UpdateResource)),
		usecases.NewDeleteResource,
		wire.Bind(new(usecasefacades.ResourceDeleter), new(*usecases.DeleteResource)),
		usecases.NewListAuditEvents,
		wire.Bind(new(usecasefacades.AuditEventLister), new(*usecases.ListAuditEvents)),
		usecases.NewWatchAuditEvents,
		wire.Bind(new(usecasefacades.AuditEventWatcher), new(*usecases.WatchAuditEvents)),
		// facades
		usecasefacades.NewResource,
		wire.Bind(new(examplerpc.ResourceFacade), new(*usecasefacades.Resource)),
		usecasefacades.NewAuditEvent,
		wire.Bind(new(examplerpc.AuditEventFacade), new(*usecasefacades.AuditEvent)),
		// middleware
		validatex.New,
		otelconnectx.NewFromConfig,
		collectInterceptors,
		rpcserver.NewRecoverer,
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
