package flagoptions

import (
	"flag"
	"os"
	"path/filepath"
	"slices"

	"github.com/mattdowdell/sandbox/internal/drivers/config"
)

// New creates a new Options instance populated with values from CLI flags.
//
//   - -config.envprefix is the prefix for environment variables to read configuration from.
//   - -config.files is a list of file paths to read configuration from.
//   - -config.mounts is a list of kubernetes pod mount paths to read configuration from.
//
// Internally, a new [flag.FlagSet] instance is created with just the above flags allowed.
// Unrecognised flags will cause the process to exit. This is usually acceptable for a binary, but
// may cause problems if relying on flags defined in other modules/packages. In these cases, use
// NewWithFlagSet to either set the global FlagSet instance or another instance.
func New() *config.Options {
	// copied from flag, apparently execl can cause os.Args to be empty
	var name string
	if len(os.Args) > 0 {
		name = os.Args[0]
	}

	// use a flagset instead of flag.String, etc. so that multiple calls don't cause an error.
	flagset := flag.NewFlagSet(name, flag.ExitOnError)

	options, _ := NewWithFlagSet(flagset)
	return options
}

// NewWithFlagSet creates a new Options instance populated with values from CLI flags.
//
//   - -config.envprefix is the prefix for environment variables to read configuration from.
//   - -config.files is a list of file paths to read configuration from.
//   - -config.mounts is a list of kubernetes pod mount paths to read configuration from.
//
// New creates a separate [flag.FlagSet] instance which is appropriate for normal binaries. However,
// this separation means that flags defined with [flag.String], etc. in other packages will result
// in an error because they're unknown. This can happen during tests. To enable use of other flags,
// the default flagset can be passed here instead:
//
//	// ignore errors, configured with flag.ExitOnError
//	options, _ := flagoptions.NewWithFlagSet(flag.Commandline)
func NewWithFlagSet(flagset *flag.FlagSet) (*config.Options, error) {
	// define flags within the function so environment variables can var across unit tests
	envPrefix := flagset.String(
		"config.envprefix",
		envOrDefault("CONFIG_ENVPREFIX", "APP_"),
		"The environment variable prefix to filter by for configuration.",
	)
	files := flagset.String(
		"config.files",
		os.Getenv("CONFIG_FILES"),
		"The file paths to use for configuration.",
	)
	mounts := flagset.String(
		"config.mounts",
		os.Getenv("CONFIG_MOUNTS"),
		"The kubernetes pod mounts to use for configuration.",
	)

	// apparently execl can cause os.Args to be empty
	var args []string
	if len(os.Args) > 1 {
		args = os.Args[1:]
	}

	if err := flagset.Parse(args); err != nil {
		return nil, err
	}

	return &config.Options{
		EnvPrefix: *envPrefix,
		Files:     splitPaths(*files),
		Mounts:    splitPaths(*mounts),
	}, nil
}

// envOrDefault reads the value of the given environment variable, using the given fallback of the
// value is empty.
func envOrDefault(name, fallback string) string {
	if val := os.Getenv(name); val != "" {
		return val
	}

	return fallback
}

// splitPaths converts a comma delimited list of paths into a slice, discarding any empty elements.
func splitPaths(input string) []string {
	parts := filepath.SplitList(input)

	return slices.DeleteFunc(parts, func(p string) bool {
		return p == ""
	})
}
