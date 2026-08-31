// Package coverage exposes a HTTP handler to execute [runtime/coverage] functions on-demand without
// having to terminate the process.
package coverage

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// Collect provides a HTTP handler that using [runtime/coverage] to write coverage to the directory
// in the GOCOVERDIR environment variable.
func Collect(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	if err := writeCoverage(); err != nil {
		slog.ErrorContext(ctx, "failed to write coverage", slogx.Err(err))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	files, err := listCoverage()
	if err != nil {
		slog.ErrorContext(ctx, "failed to list coverage files", slogx.Err(err))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(files)
	if err != nil {
		slog.ErrorContext(ctx, "failed to marshal response body", slogx.Err(err))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Add("Content-Type", "application/json")

	if _, err := writer.Write(body); err != nil {
		slog.WarnContext(ctx, "failed to write response body", slogx.Err(err))
	}
}
