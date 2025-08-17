package main

import (
	"github.com/mattdowdell/sandbox/internal/drivers/config"
	"github.com/mattdowdell/sandbox/internal/drivers/jwtx"
	"github.com/mattdowdell/sandbox/internal/drivers/logging"
	"github.com/mattdowdell/sandbox/internal/drivers/otelx"
	"github.com/mattdowdell/sandbox/internal/drivers/pgsql"
	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver"
	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver/interceptors/authn"
	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver/interceptors/otelconnectx"
)

// Config contains the service configuration.
type Config struct {
	App            AppConfig                  `koanf:",squash"`
	Authn          authn.Config               `koanf:"authn"`
	Database       pgsql.Config               `koanf:"database"`
	Issuer         jwtx.IssuerConfig          `koanf:"issuer"`
	Logging        logging.Config             `koanf:"logging"`
	MeterProvider  otelx.MeterProviderConfig  `koanf:"meterprovider"`
	TracerProvider otelx.TracerProviderConfig `koanf:"tracerprovider"`
	LoggerProvider otelx.LoggerProviderConfig `koanf:"loggerprovider"`
	Parser         jwtx.ParserConfig          `koanf:"parser"`
	OtelConnect    otelconnectx.Config        `koanf:"otelconnect"`
	RPCServer      rpcserver.Config           `koanf:"rpcserver"`
}

// LoadConfig loads the service configuration.
//
// This is effectively a workaround for wire not supporting generics. For more details, see
// https://github.com/google/wire/issues/354.
func LoadConfig(conf *config.Config) (Config, error) {
	return config.Load[Config](conf)
}
