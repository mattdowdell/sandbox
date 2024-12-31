package main

import (
	"github.com/mattdowdell/sandbox/internal/drivers/config"
	"github.com/mattdowdell/sandbox/internal/drivers/logging"
	"github.com/mattdowdell/sandbox/internal/drivers/otelx"
	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver"
	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver/interceptors/otelconnectx"
)

// Config contains the service configuration.
type Config struct {
	App         AppConfig
	Logging     logging.Config
	Meter       otelx.MeterProviderConfig
	Tracer      otelx.TracerProviderConfig
	OtelConnect otelconnectx.Config
	RPCServer   rpcserver.Config
}

// LoadConfig loads the service configuration.
//
// This is effectively a workaround for wire not supporting generics. For more details, see
// https://github.com/google/wire/issues/354.
func LoadConfig(conf *config.Config) (Config, error) {
	return config.Load[Config](conf)
}
