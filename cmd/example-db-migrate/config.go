package main

import (
	"github.com/mattdowdell/sandbox/internal/drivers/config"
	"github.com/mattdowdell/sandbox/internal/drivers/logging"
	"github.com/mattdowdell/sandbox/internal/drivers/otelx"
	"github.com/mattdowdell/sandbox/internal/drivers/pgsql"
)

type Config struct {
	Logging        logging.Config             `koanf:"logging"`
	Database       pgsql.Config               `koanf:"database"`
	MeterProvider  otelx.MeterProviderConfig  `koanf:"meterprovider"`
	TracerProvider otelx.TracerProviderConfig `koanf:"tracerprovider"`
	LoggerProvider otelx.LoggerProviderConfig `koanf:"loggerprovider"`
}

// LoadConfig loads the service configuration.
//
// This is effectively a workaround for wire not supporting generics. For more details, see
// https://github.com/google/wire/issues/354.
func LoadConfig(conf *config.Config) (Config, error) {
	return config.Load[Config](conf)
}
