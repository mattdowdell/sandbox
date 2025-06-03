package flagoptions_test

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/drivers/config"
	"github.com/mattdowdell/sandbox/internal/drivers/config/flagoptions"
)

func setArgs(t *testing.T, args ...string) {
	t.Helper()

	old := os.Args
	os.Args = args

	t.Cleanup(func() {
		os.Args = old
	})
}

func Test_New(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		env  map[string]string
		want *config.Options
	}{
		{
			name: "defaults",
			args: []string{t.Name()},
			want: &config.Options{
				EnvPrefix: "APP_",
				Files:     []string{},
				Mounts:    []string{},
			},
		},
		{
			name: "arguments",
			args: []string{
				t.Name(),
				"-config.envprefix=ARG_",
				"-config.files=/etc/arg1,/etc/arg2",
				"-config.mounts=/etc/arg3,/etc/arg4",
			},
			want: &config.Options{
				EnvPrefix: "ARG_",
				Files:     []string{"/etc/arg1", "/etc/arg2"},
				Mounts:    []string{"/etc/arg3", "/etc/arg4"},
			},
		},
		{
			name: "environment",
			env: map[string]string{
				"CONFIG_ENVPREFIX": "ENV_",
				"CONFIG_FILES":     "/etc/env1,/etc/env2",
				"CONFIG_MOUNTS":    "/etc/env3,/etc/env4",
			},
			want: &config.Options{
				EnvPrefix: "ENV_",
				Files:     []string{"/etc/env1", "/etc/env2"},
				Mounts:    []string{"/etc/env3", "/etc/env4"},
			},
		},
		{
			name: "both",
			args: []string{
				t.Name(),
				"-config.envprefix=ARG_",
				"-config.files=/etc/arg1,/etc/arg2",
				"-config.mounts=/etc/arg3,/etc/arg4",
			},
			env: map[string]string{
				"CONFIG_ENVPREFIX": "ENV_",
				"CONFIG_FILES":     "/etc/env1,/etc/env2",
				"CONFIG_MOUNTS":    "/etc/env3,/etc/env4",
			},
			want: &config.Options{
				EnvPrefix: "ARG_",
				Files:     []string{"/etc/arg1", "/etc/arg2"},
				Mounts:    []string{"/etc/arg3", "/etc/arg4"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			setArgs(t, tc.args...)

			for key, val := range tc.env {
				t.Setenv(key, val)
			}

			// act
			got := flagoptions.New()

			// assert
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_NewWithFlagSet_Success(t *testing.T) {
	// arrange
	flagset := flag.NewFlagSet(t.Name(), flag.ContinueOnError)

	setArgs(t, t.Name())

	// act
	got, err := flagoptions.NewWithFlagSet(flagset)

	// assert
	want := &config.Options{
		EnvPrefix: "APP_",
		Files:     []string{},
		Mounts:    []string{},
	}

	assert.Equal(t, want, got)
	assert.NoError(t, err)
}

func Test_NewWithFlagSet_Error(t *testing.T) {
	// arrange
	flagset := flag.NewFlagSet(t.Name(), flag.ContinueOnError)

	setArgs(t, t.Name(), "-invalid=invalid")

	// act
	got, err := flagoptions.NewWithFlagSet(flagset)

	// assert
	assert.Nil(t, got)
	assert.EqualError(t, err, "flag provided but not defined: -invalid")
}
