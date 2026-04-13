package herdconv_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"

	"github.com/mattdowdell/sandbox/pkg/herd/internal/herdconv"
)

func Test_Attributes(t *testing.T) {
	tests := map[string]struct {
		got  attribute.KeyValue
		want attribute.KeyValue
	}{
		"herd.version.before": {
			got:  herdconv.HerdVersionBefore(1),
			want: attribute.Int("herd.version.before", 1),
		},
		"herd.version.after": {
			got:  herdconv.HerdVersionAfter(1),
			want: attribute.Int("herd.version.after", 1),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.got)
		})
	}
}
