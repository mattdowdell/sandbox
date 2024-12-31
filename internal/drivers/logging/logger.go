package logging

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

// ...
type Config struct {
	// ...
	Level slog.Level `koanf:"logging.level" default:"info"`
}

// ...
func NewAsDefaultFromConfig(config Config) *slog.Logger {
	return NewAsDefault(config.Level)
}

// ...
func NewAsDefault(level slog.Level) *slog.Logger {
	logger := New(level)
	slog.SetDefault(logger)

	return logger
}

// ...
func NewFromConfig(config Config) *slog.Logger {
	return New(config.Level)
}

// ...
func New(level slog.Level) *slog.Logger {
	// TODO: add trace/span id to log attrs (requires custom handler?)
	handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if len(groups) > 0 {
				return a
			}

			switch a.Key {
			case slog.SourceKey:
				source := a.Value.Any().(*slog.Source)
				a.Value = slog.StringValue(fmt.Sprintf("%s:%d", source.File, source.Line))

			case slog.LevelKey:
				a.Value = slog.StringValue(strings.ToLower(a.Value.String()))
			}

			return a
		},
	})

	return slog.New(handler)
}
