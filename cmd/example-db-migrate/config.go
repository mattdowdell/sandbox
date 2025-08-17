package main

import (
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
