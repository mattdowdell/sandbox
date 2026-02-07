package logging

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_replaceAttr(t *testing.T) {
	tests := map[string]struct {
		groups []string
		attr   slog.Attr
		want   slog.Attr
	}{
		"with group": {
			groups: []string{"example"},
			attr:   slog.Any(slog.LevelKey, slog.LevelInfo),
			want:   slog.Any(slog.LevelKey, slog.LevelInfo),
		},
		"lowercase level": {
			attr: slog.Any(slog.LevelKey, slog.LevelInfo),
			want: slog.String(slog.LevelKey, "info"),
		},
		"format source": {
			attr: slog.Any(slog.SourceKey, &slog.Source{
				File: "example.go",
				Line: 1,
			}),
			want: slog.String(slog.SourceKey, "example.go:1"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			opts := defaultOptions()
			fn := replaceAttr(opts)

			// act
			output := fn(tt.groups, tt.attr)

			// assert
			assert.Equal(t, tt.want, output)
		})
	}
}
