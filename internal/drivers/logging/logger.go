package logging

import (
	"fmt"
	"log"
	"log/slog"
	"strings"
)

// Level controls the level of the default logger. It exists to allow modifications to the level
// after the default logger is configured in NewAsDefault.
var Level slog.LevelVar

// Config contains the configuration for the logger.
type Config struct {
	// Level sets the level for the logger.
	Level slog.Level `koanf:"level" default:"info"`

	// LegacyLevel sets the level for logs from the "log" package. This is only applied if using
	// NewAsDefaultFromConfig.
	LegacyLevel slog.Level `koanf:"level" default:"debug"`
}

// NewAsDefaultFromConfig calls NewAsDefault with the given configuration.
func NewAsDefaultFromConfig(config Config, options ...Option) *slog.Logger {
	return NewAsDefault(config.Level, config.LegacyLevel, options...)
}

// NewAsDefault updates the global level with the given value and calls New with the new value.
//
// The level used by the "log" package is set with the given legacy level. If the global log level
// is less than or equal to than the legacy level, the logs emitted by log.Println, etc. will be
// output. Otherwise they will be dropped.
func NewAsDefault(level, legacyLevel slog.Level, options ...Option) *slog.Logger {
	Level.Set(level)
	logger := New(&Level, options...)

	// make it easier to see where legacy logs came from
	// needs to be called before slog.SetDefault to take effect
	log.SetFlags(log.Llongfile)

	slog.SetDefault(logger)
	slog.SetLogLoggerLevel(legacyLevel)

	return logger
}

// NewFromConfig calls New with the given configuration.
func NewFromConfig(config Config, options ...Option) *slog.Logger {
	return New(config.Level, options...)
}

// New creates a new logger with the given level using a JSON handler.
func New(level slog.Leveler, options ...Option) *slog.Logger {
	opts := defaultOptions()
	for _, option := range options {
		option.apply(opts)
	}

	inner := slog.NewJSONHandler(opts.writer, &slog.HandlerOptions{
		AddSource:   true,
		Level:       level,
		ReplaceAttr: replaceAttr(opts),
	})

	handler := Wrap(inner, opts.extractors)

	return slog.New(handler)
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
