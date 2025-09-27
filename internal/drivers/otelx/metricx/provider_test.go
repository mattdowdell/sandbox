//nolint:dupl // similar to meter_test.go/tracer_test.go
package metricx_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mattdowdell/sandbox/internal/drivers/otelx/metricx"
)

func urlToProviderConfig(t *testing.T, input string) metricx.ProviderConfig {
	t.Helper()

	u, err := url.Parse(input)
	require.NoError(t, err)

	return metricx.ProviderConfig{
		Insecure: u.Scheme == "http",
		Endpoint: u.Host,
		Path:     u.Path,
	}
}

func Test_SetupMeterProviderFromConfig_Success(t *testing.T) {
	// arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	conf := urlToProviderConfig(t, server.URL)

	// act
	shutdown, err := metricx.SetupMeterProviderFromConfig(t.Context(), conf)

	// assert
	if assert.NotNil(t, shutdown) {
		assert.NoError(t, shutdown(t.Context()))
	}

	assert.NoError(t, err)
}

func Test_SetupMeterProviderFromConfig_Error(t *testing.T) {
	// arrange
	conf := metricx.ProviderConfig{
		Endpoint: "\t",
	}

	// act
	shutdown, err := metricx.SetupMeterProviderFromConfig(t.Context(), conf)

	// assert
	assert.Nil(t, shutdown)
	assert.EqualError(t, err, `parse "https://%09/v1/metrics": invalid URL escape "%09"`)
}
