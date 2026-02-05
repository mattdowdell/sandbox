package splitter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/drivers/config/splitter"
)

func Test_Comma_UnmarshalText(t *testing.T) {
	tests := map[string]struct {
		have string
		want []string
	}{
		"empty": {
			have: "",
			want: []string{},
		},
		"non-empty": {
			have: "foo,bar,baz",
			want: []string{"foo", "bar", "baz"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			c := splitter.Comma{}

			// act
			err := c.UnmarshalText([]byte(tt.have))

			// assert
			assert.Equal(t, tt.want, c.Unwrap())
			assert.NoError(t, err)
		})
	}
}

func Test_Comma_MarshalText(t *testing.T) {
	tests := map[string]struct {
		have splitter.Comma
		want string
	}{
		"nil": {
			have: splitter.Comma{},
			want: "",
		},
		"empty": {
			have: splitter.Comma{},
			want: "",
		},
		"non-empty": {
			have: splitter.Comma{"foo", "bar", "baz"},
			want: "foo,bar,baz",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange

			// act
			got, err := tt.have.MarshalText()

			// assert
			assert.Equal(t, tt.want, string(got))
			assert.NoError(t, err)
		})
	}
}

func Test_Space_UnmarshalText(t *testing.T) {
	tests := map[string]struct {
		have string
		want []string
	}{
		"empty": {
			have: "",
			want: []string{},
		},
		"non-empty": {
			have: "foo bar baz",
			want: []string{"foo", "bar", "baz"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			s := splitter.Space{}

			// act
			err := s.UnmarshalText([]byte(tt.have))

			// assert
			assert.Equal(t, tt.want, s.Unwrap())
			assert.NoError(t, err)
		})
	}
}

func Test_Space_MarshalText(t *testing.T) {
	tests := map[string]struct {
		have splitter.Space
		want string
	}{
		"nil": {
			have: splitter.Space{},
			want: "",
		},
		"empty": {
			have: splitter.Space{},
			want: "",
		},
		"non-empty": {
			have: splitter.Space{"foo", "bar", "baz"},
			want: "foo bar baz",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange

			// act
			got, err := tt.have.MarshalText()

			// assert
			assert.Equal(t, tt.want, string(got))
			assert.NoError(t, err)
		})
	}
}
