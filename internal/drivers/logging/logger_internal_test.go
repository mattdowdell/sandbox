package logging

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_replaceAttr(t *testing.T) {
	testCases := []struct {
		name   string
		groups []string
		attr   slog.Attr
		want   slog.Attr
	}{
		{
			name:   "with group",
			groups: []string{"example"},
			attr:   slog.Any(slog.LevelKey, slog.LevelInfo),
			want:   slog.Any(slog.LevelKey, slog.LevelInfo),
		},
		{
			name: "lowercase level",
			attr: slog.Any(slog.LevelKey, slog.LevelInfo),
			want: slog.String(slog.LevelKey, "info"),
		},
		{
			name: "format source",
			attr: slog.Any(slog.SourceKey, &slog.Source{
				File: "example.go",
				Line: 1,
			}),
			want: slog.String(slog.SourceKey, "example.go:1"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			opts := defaultOptions()
			fn := replaceAttr(opts)

			// act
			output := fn(tc.groups, tc.attr)

			// assert
			assert.Equal(t, tc.want, output)
		})
	}
}
