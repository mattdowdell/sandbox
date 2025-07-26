package flagoptions_test

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/drivers/config"
	"github.com/mattdowdell/sandbox/internal/drivers/config/flagoptions"
)

var (
	testArgPath1 = filepath.Join("etc", "arg1")
	testArgPath2 = filepath.Join("etc", "arg2")
	testArgPath3 = filepath.Join("etc", "arg3")
	testArgPath4 = filepath.Join("etc", "arg4")

	testEnvPath1 = filepath.Join("etc", "env1")
	testEnvPath2 = filepath.Join("etc", "env2")
	testEnvPath3 = filepath.Join("etc", "env3")
	testEnvPath4 = filepath.Join("etc", "env4")
)

func setArgs(t *testing.T, args ...string) {
	t.Helper()

	old := os.Args
	os.Args = args

	t.Cleanup(func() {
		os.Args = old
	})
}

func joinPathList(paths []string) string {
	return strings.Join(paths, string(os.PathListSeparator))
}

func Test_New(t *testing.T) {
	argFiles := []string{testArgPath1, testArgPath2}
	argMounts := []string{testArgPath3, testArgPath4}

	envFiles := []string{testEnvPath1, testEnvPath2}
	envMounts := []string{testEnvPath3, testEnvPath4}

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
				"-config.files=" + joinPathList(argFiles),
				"-config.mounts=" + joinPathList(argMounts),
			},
			want: &config.Options{
				EnvPrefix: "ARG_",
				Files:     argFiles,
				Mounts:    argMounts,
			},
		},
		{
			name: "environment",
			env: map[string]string{
				"CONFIG_ENVPREFIX": "ENV_",
				"CONFIG_FILES":     joinPathList(envFiles),
				"CONFIG_MOUNTS":    joinPathList(envMounts),
			},
			want: &config.Options{
				EnvPrefix: "ENV_",
				Files:     envFiles,
				Mounts:    envMounts,
			},
		},
		{
			name: "both",
			args: []string{
				t.Name(),
				"-config.envprefix=ARG_",
				"-config.files=" + joinPathList(argFiles),
				"-config.mounts=" + joinPathList(argMounts),
			},
			env: map[string]string{
				"CONFIG_ENVPREFIX": "ENV_",
				"CONFIG_FILES":     joinPathList(envFiles),
				"CONFIG_MOUNTS":    joinPathList(envMounts),
			},
			want: &config.Options{
				EnvPrefix: "ARG_",
				Files:     argFiles,
				Mounts:    argMounts,
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
