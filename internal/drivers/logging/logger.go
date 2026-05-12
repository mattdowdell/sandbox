package logging

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"strings"

	"go.opentelemetry.io/contrib/bridges/otelslog"
)

// Level controls the level of the default logger. It exists to allow modifications to the level
// after the default logger is configured in NewAsDefault.
var Level slog.LevelVar

// Config contains the configuration for the logger.
type Config struct {
	// Level sets the level for the logger. Logs emitted below the configured level will be dropped.
	Level slog.Level `koanf:"level" default:"info"`

	// LegacyLevel sets the level for logs from the "log" package. This is only applied if using
	// NewAsDefaultFromConfig.
	//
	// Unlike Level, this does not act as a filtering threshold. Instead, it creates the concept of
	// a level in the "log" package and causes any logs output by uses of that package to use the
	// configured level. These logs may subsequently be filtered out depending on the value of
	// Level.
	LegacyLevel slog.Level `koanf:"legacylevel" default:"debug"`

	// Enable the JSON logging handler. This causes logs to be output to stdout under the assumption
	// that they will be pulled from the Node the Pod is scheduled upon. At least one of
	// EnableJSONHandler or EnableOtelHandler must be set to true.
	EnableJSONHandler bool `koanf:"enablejsonhandler" default:"true"`

	// Enable the OpenTelemetry logging handler. This causes logs to be periodically pushed to the
	// configured endpoint. At least one of EnableJSONHandler or EnableOtelHandler must be set to
	// true.
	EnableOtelHandler bool `koanf:"enableotelhandler"`
}

// NewAsDefaultFromConfig calls NewAsDefault with the given configuration.
func NewAsDefaultFromConfig(config Config, options ...Option) (*slog.Logger, error) {
	options = append(
		options,
		EnableJSONHandler(config.EnableJSONHandler),
		EnableOtelHandler(config.EnableOtelHandler),
	)

	return NewAsDefault(config.Level, config.LegacyLevel, options...)
}

// NewAsDefault updates the global level with the given value and calls New with the new value.
//
// The level used by the "log" package is set with the given legacy level. If the global log level
// is less than or equal to than the legacy level, the logs emitted by log.Println, etc. will be
// output. Otherwise they will be dropped.
func NewAsDefault(level, legacyLevel slog.Level, options ...Option) (*slog.Logger, error) {
	Level.Set(level)

	logger, err := New(&Level, options...)
	if err != nil {
		return nil, err
	}

	// make it easier to see where legacy logs came from
	// needs to be called before slog.SetDefault to take effect
	log.SetFlags(log.Llongfile)

	slog.SetDefault(logger)
	slog.SetLogLoggerLevel(legacyLevel)

	return logger, nil
}

// NewFromConfig calls New with the given configuration.
func NewFromConfig(config Config, options ...Option) (*slog.Logger, error) {
	return New(config.Level, options...)
}

// New creates a new logger with the given level using a JSON and OpenTelemetry handler.
func New(level slog.Leveler, options ...Option) (*slog.Logger, error) {
	opts := defaultOptions()
	for _, option := range options {
		option.apply(opts)
	}

	var handlers []slog.Handler

	if opts.enableJSONHandler {
		handler := slog.NewJSONHandler(opts.writer, &slog.HandlerOptions{
			AddSource:   true,
			Level:       level,
			ReplaceAttr: replaceAttr(opts),
		})

		handlers = append(handlers, Wrap(handler, opts.extractors))
	}

	if opts.enableOtelHandler {
		handler := otelslog.NewHandler("", otelslog.WithSource(true))
		handlers = append(handlers, Wrap(handler, opts.extractors))
	}

	if len(handlers) == 0 {
		return nil, errors.New("one of json or otel log handlers must be enabled")
	}

	return slog.New(slog.NewMultiHandler(handlers...)), nil
}

func replaceAttr(opts *loggerOpts) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, attr slog.Attr) slog.Attr {
		if len(groups) > 0 {
			return attr
		}

		switch attr.Key {
		case slog.LevelKey:
			attr.Value = slog.StringValue(strings.ToLower(attr.Value.String()))

		case slog.SourceKey:
			if opts.suppressSource {
				return slog.Attr{}
			}

			if source, ok := attr.Value.Any().(*slog.Source); ok {
				attr.Value = slog.StringValue(fmt.Sprintf("%s:%d", source.File, source.Line))
			}

		case slog.TimeKey:
			if opts.suppressTime {
				return slog.Attr{}
			}
		}

		return attr
	}
}
