package main

import (
	"context"
	"log/slog"
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"github.com/gofrs/uuid/v5"
	"go.opentelemetry.io/contrib/processors/baggagecopy"

	"github.com/mattdowdell/sandbox/gen/authn/v1/authnv1connect"
	"github.com/mattdowdell/sandbox/gen/config/v1/configv1connect"
	"github.com/mattdowdell/sandbox/gen/example/v1/examplev1connect"
	"github.com/mattdowdell/sandbox/internal/adapters/authnrpc"
	"github.com/mattdowdell/sandbox/internal/adapters/configrpc"
	"github.com/mattdowdell/sandbox/internal/adapters/datastore"
	"github.com/mattdowdell/sandbox/internal/adapters/examplerpc"
	"github.com/mattdowdell/sandbox/internal/adapters/healthrpc"
	"github.com/mattdowdell/sandbox/internal/adapters/reflectrpc"
	"github.com/mattdowdell/sandbox/internal/adapters/usecasefacades"
	"github.com/mattdowdell/sandbox/internal/domain/repositories"
	"github.com/mattdowdell/sandbox/internal/drivers/config"
	"github.com/mattdowdell/sandbox/internal/drivers/config/flagoptions"
	"github.com/mattdowdell/sandbox/internal/drivers/jwtx"
	"github.com/mattdowdell/sandbox/internal/drivers/logging"
	"github.com/mattdowdell/sandbox/internal/drivers/otelx"
	"github.com/mattdowdell/sandbox/internal/drivers/otelx/logx"
	"github.com/mattdowdell/sandbox/internal/drivers/otelx/metricx"
	"github.com/mattdowdell/sandbox/internal/drivers/otelx/tracex"
	"github.com/mattdowdell/sandbox/internal/drivers/pgsql"
	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver"
	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver/interceptors/authn"
	logginginterceptor "github.com/mattdowdell/sandbox/internal/drivers/rpcserver/interceptors/logging"
	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver/interceptors/otelconnectx"
	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver/interceptors/validatex"
	"github.com/mattdowdell/sandbox/internal/usecases"
	"github.com/mattdowdell/sandbox/pkg/timex"
)

func SetupApp(ctx context.Context) (*App, error) {
	options := flagoptions.New()
	loader := config.New[Config](options)

	conf, err := loader.Load()
	if err != nil {
		return nil, err
	}

	tpShutdown, mpShutdown, lpShutdown, logger, err := initObservability(ctx, conf)
	if err != nil {
		return nil, err
	}

	clock := timex.NewClock()
	uuidgen := uuid.NewGen()

	resource, auditEvent, err := initFacades(ctx, conf, clock, uuidgen)
	if err != nil {
		return nil, err
	}

	handlerOpts, err := initHandlerOptions(conf)
	if err != nil {
		return nil, err
	}

	handlers, err := initHandlers(conf, loader, clock, uuidgen, resource, auditEvent)
	if err != nil {
		return nil, err
	}

	server := rpcserver.NewFromConfig(conf.RPCServer, handlers, handlerOpts)

	return NewApp(conf, logger, server, tpShutdown, mpShutdown, lpShutdown), nil
}

//nolint:gocritic // result types are differentiated by package name.
func initObservability(
	ctx context.Context,
	conf *Config,
) (
	tracex.ProviderShutdown,
	metricx.ProviderShutdown,
	logx.ProviderShutdown,
	*slog.Logger,
	error,
) {
	filter := baggagecopy.AllowAllMembers

	tpShutdown, err := tracex.SetupTracerProviderFromConfig(ctx, conf.TracerProvider, filter)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	mpShutdown, err := metricx.SetupMeterProviderFromConfig(ctx, conf.MeterProvider)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	lpShutdown, err := logx.SetupLoggerProviderFromConfig(ctx, conf.LoggerProvider, filter)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	extractor := otelx.NewExtractor(otelx.WithSpanID(true), otelx.WithSampled(true))
	logger := logging.NewAsDefaultFromConfig(conf.Logging, logging.WithExtractors(extractor))

	return tpShutdown, mpShutdown, lpShutdown, logger, nil
}

func initFacades(
	ctx context.Context,
	conf *Config,
	clock repositories.Clock,
	uuidgen repositories.UUIDGenerator,
) (*usecasefacades.Resource, *usecasefacades.AuditEvent, error) {
	db, _, err := pgsql.NewFromConfig(ctx, conf.Database)
	if err != nil {
		return nil, nil, err
	}

	provider := datastore.NewProvider(db)

	createResource := usecases.NewCreateResource(clock, uuidgen)
	getResource := usecases.NewGetResource()
	listResources := usecases.NewListResources()
	updateResource := usecases.NewUpdateResource(clock)
	deleteResource := usecases.NewDeleteResource()
	listAuditevents := usecases.NewListAuditEvents()
	watchAuditEvents := usecases.NewWatchAuditEvents()

	resource := usecasefacades.NewResource(
		provider,
		createResource,
		getResource,
		listResources,
		updateResource,
		deleteResource,
	)

	auditEvent := usecasefacades.NewAuditEvent(provider, listAuditevents, watchAuditEvents)

	return resource, auditEvent, nil
}

func initHandlers(
	conf *Config,
	loader *config.Loader[Config],
	clock repositories.Clock,
	uuidgen repositories.UUIDGenerator,
	resource examplerpc.ResourceFacade,
	auditEvent examplerpc.AuditEventFacade,
) ([]rpcserver.Handler, error) {
	issuer, err := jwtx.NewIssuerFromConfig(clock, uuidgen, conf.Issuer)
	if err != nil {
		return nil, err
	}

	parser, err := jwtx.NewParserFromConfig(clock, conf.Parser)
	if err != nil {
		return nil, err
	}

	return []rpcserver.Handler{
		authnrpc.New(issuer, parser),
		examplerpc.New(resource, auditEvent),
		reflectrpc.New([]string{
			authnv1connect.AuthnServiceName,
			configv1connect.ConfigServiceName,
			examplev1connect.ExampleServiceName,
			grpchealth.HealthV1ServiceName,
		}),
		healthrpc.New(),
		configrpc.New(loader),
	}, nil
}

func initHandlerOptions(conf *Config) ([]connect.HandlerOption, error) {
	otelconnectInterceptor, err := otelconnectx.NewFromConfig(conf.OtelConnect)
	if err != nil {
		return nil, err
	}

	validateInterceptor := validatex.New()
	loggingInterceptor := logginginterceptor.New()

	client := authn.NewClientFromConfig(
		http.DefaultClient,
		conf.Authn,
		otelconnectInterceptor,
		validateInterceptor,
	)

	authnInterceptor := authn.New(client, authn.WithIgnoreService(
		grpchealth.HealthV1ServiceName,
		grpcreflect.ReflectV1ServiceName,
		grpcreflect.ReflectV1AlphaServiceName,
		authnv1connect.AuthnServiceName,
		configv1connect.ConfigServiceName, // TODO: remove
	))

	recoverer, err := rpcserver.NewRecoverer()
	if err != nil {
		return nil, err
	}

	return []connect.HandlerOption{
		connect.WithInterceptors(
			otelconnectInterceptor,
			validateInterceptor,
			loggingInterceptor,
			authnInterceptor,
		),
		connect.WithRecover(recoverer.Handle),
	}, nil
}
