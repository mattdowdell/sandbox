package main

import (
	"github.com/mattdowdell/sandbox/internal/drivers/jwtx"
	"github.com/mattdowdell/sandbox/internal/drivers/logging"
	"github.com/mattdowdell/sandbox/internal/drivers/otelx"
	"github.com/mattdowdell/sandbox/internal/drivers/pgsql"
	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver"
	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver/interceptors/otelconnectx"
)

// Config contains the service configuration.
type Config struct {
	App            AppConfig                  `koanf:",squash"`
	Database       pgsql.Config               `koanf:"database"`
	Logging        logging.Config             `koanf:"logging"`
	MeterProvider  otelx.MeterProviderConfig  `koanf:"meterprovider"`
	TracerProvider otelx.TracerProviderConfig `koanf:"tracerprovider"`
	LoggerProvider otelx.LoggerProviderConfig `koanf:"loggerprovider"`
	OtelConnect    otelconnectx.Config        `koanf:"otelconnect"`
	RPCServer      rpcserver.Config           `koanf:"rpcserver"`
	Issuer         jwtx.IssuerConfig          `koanf:"issuer"`
	Parser         jwtx.ParserConfig          `koanf:"parser"`
}
