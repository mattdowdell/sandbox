//nolint:dupl // similar to meter_test.go/tracer_test.go
package otelx_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/contrib/processors/baggagecopy"

	"github.com/mattdowdell/sandbox/internal/drivers/otelx"
)

func isHTTP(scheme string) bool {
	return scheme == "http"
}

func urlToLoggerProviderConfig(t *testing.T, input string) otelx.LoggerProviderConfig {
	t.Helper()

	u, err := url.Parse(input)
	require.NoError(t, err)

	return otelx.LoggerProviderConfig{
		Insecure: isHTTP(u.Scheme),
		Endpoint: u.Host,
		Path:     u.Path,
	}
}

func Test_SetupLoggerProviderFromConfig_Success(t *testing.T) {
	// arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	conf := urlToLoggerProviderConfig(t, server.URL)

	// act
	shutdown, err := otelx.SetupLoggerProviderFromConfig(t.Context(), conf, baggagecopy.AllowAllMembers)

	// assert
	if assert.NotNil(t, shutdown) {
		assert.NoError(t, shutdown(t.Context()))
	}

	assert.NoError(t, err)
}

func Test_SetupLoggerProviderFromConfig_Error(t *testing.T) {
	// arrange
	conf := otelx.LoggerProviderConfig{
		Endpoint: "\t",
	}

	// act
	shutdown, err := otelx.SetupLoggerProviderFromConfig(t.Context(), conf, baggagecopy.AllowAllMembers)

	// assert
	assert.Nil(t, shutdown)
	assert.EqualError(t, err, `parse "https://%09": invalid URL escape "%09"`)
}
