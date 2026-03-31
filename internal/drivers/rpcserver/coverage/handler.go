// Package coverage exposes a HTTP handler to execute [runtime/coverage] functions on-demand without
// having to terminate the process.
package coverage

import (
	"log/slog"
	"net/http"
	"runtime/coverage"
	"os"

	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// Collect provides a HTTP handler that using [runtime/coverage] to write coverage to the directory
// in the GOCOVERDIR environment variable.
func Collect(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	dir := os.Getenv("GOCOVERDIR")
	if dir == "" {
		slog.WarnContext(ctx, "GOCOVERDIR not set")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := coverage.WriteCountersDir(dir); err != nil {
		slog.ErrorContext(ctx, "failed to write convcounters", slogx.Err(err))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := coverage.WriteMetaDir(dir); err != nil {
		slog.ErrorContext(ctx, "failed to write covmeta", slogx.Err(err))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := coverage.ClearCounters(); err != nil {
		slog.ErrorContext(ctx, "failed to clear coverage counters", slogx.Err(err))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
