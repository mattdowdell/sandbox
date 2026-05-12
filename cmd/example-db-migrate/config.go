package main

import (
	"github.com/mattdowdell/sandbox/internal/drivers/logging"
	"github.com/mattdowdell/sandbox/internal/drivers/otelx/logx"
	"github.com/mattdowdell/sandbox/internal/drivers/otelx/metricx"
	"github.com/mattdowdell/sandbox/internal/drivers/otelx/tracex"
	"github.com/mattdowdell/sandbox/internal/drivers/pgsql"
)

type Config struct {
	App            AppConfig              `koanf:",squash"`
	Logging        logging.Config         `koanf:"logging"`
	Database       pgsql.Config           `koanf:"database"`
	MeterProvider  metricx.ProviderConfig `koanf:"meterprovider"`
	TracerProvider tracex.ProviderConfig  `koanf:"tracerprovider"`
	LoggerProvider logx.ProviderConfig    `koanf:"loggerprovider"`
}
