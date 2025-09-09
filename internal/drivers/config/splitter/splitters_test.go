package splitter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/drivers/config"
	"github.com/mattdowdell/sandbox/internal/drivers/config/splitter"
)

type MyConfig struct {
	Foo splitter.Comma `koanf:"foo"`
	Bar splitter.Space `koanf:"bar"`
}

func Test_Comma_UnmarshalText(t *testing.T) {
	testCases := []struct {
		name string
		have string
		want []string
	}{
		{
			name: "empty",
			have: "",
			want: []string{},
		},
		{
			name: "non-empty",
			have: "foo,bar,baz",
			want: []string{"foo", "bar", "baz"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			c := splitter.Comma{}

			// act
			err := c.UnmarshalText([]byte(tc.have))

			// assert
			assert.Equal(t, tc.want, c.Unwrap())
			assert.NoError(t, err)
		})
	}
}

func Test_Space_UnmarshalText(t *testing.T) {
	testCases := []struct {
		name string
		have string
		want []string
	}{
		{
			name: "empty",
			have: "",
			want: []string{},
		},
		{
			name: "non-empty",
			have: "foo bar baz",
			want: []string{"foo", "bar", "baz"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			s := splitter.Space{}

			// act
			err := s.UnmarshalText([]byte(tc.have))

			// assert
			assert.Equal(t, tc.want, s.Unwrap())
			assert.NoError(t, err)
		})
	}
}

func Test_Load(t *testing.T) {
	// arrange
	t.Setenv("APP_FOO", "foo,bar,baz")
	t.Setenv("APP_BAR", "foo bar baz")

	options := &config.Options{
		EnvPrefix: "APP_",
	}

	// act
	got, err := config.Load[MyConfig](options)

	// assert
	want := &MyConfig{
		Foo: splitter.Comma{"foo", "bar", "baz"},
		Bar: splitter.Space{"foo", "bar", "baz"},
	}

	assert.Equal(t, want, got)
	assert.NoError(t, err)
}
