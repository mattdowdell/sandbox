package flagoptions

import (
	"flag"
	"os"
	"slices"
	"strings"

	"github.com/mattdowdell/sandbox/internal/drivers/config"
)

// Command-line flags for bootstrapping the service configuration collection.
var (
	envPrefix = flag.String(
		"config.envprefix",
		envOrDefault("CONFIG_ENVPREFIX", "APP_"),
		"The environment variable prefix to filter by for configuration.",
	)
	files = flag.String(
		"config.files",
		os.Getenv("CONFIG_FILES"),
		"The file paths to use for configuration.",
	)
	mounts = flag.String(
		"config.mounts",
		os.Getenv("CONFIG_MOUNTS"),
		"The kubernetes pod mounts to use for configuration.",
	)
)

// New creates a new Options instance populated with value from CLI flags.
//
// - -config.envprefix is the prefix for environment variables to read configuration from.
// - -config.files is a comma delimited list of files to read configuration from.
// - -config.mounts is a comma delimited list of kubernetes pod mounts to read configuration from.
func New() *config.Options {
	flag.Parse()

	return &config.Options{
		EnvPrefix: *envPrefix,
		Files:     splitPaths(*files),
		Mounts:    splitPaths(*mounts),
	}
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
	parts := strings.Split(input, ",")

	return slices.DeleteFunc(parts, func(p string) bool {
		return p == ""
	})
}
