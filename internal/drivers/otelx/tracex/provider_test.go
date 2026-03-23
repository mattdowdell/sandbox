//nolint:dupl // similar to meter_test.go/tracer_test.go
package tracex_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/contrib/processors/baggagecopy"

	"github.com/mattdowdell/sandbox/internal/drivers/otelx/tracex"
)

func urlToProviderConfig(t *testing.T, input string) tracex.ProviderConfig {
	t.Helper()

	u, err := url.Parse(input)
	require.NoError(t, err)

	return tracex.ProviderConfig{
		Insecure: u.Scheme == "http",
		Endpoint: u.Host,
		Path:     u.Path,
	}
}

func Test_SetupTracerProviderFromConfig_Success(t *testing.T) {
	// arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	conf := urlToProviderConfig(t, server.URL)

	// act
	shutdown, err := tracex.SetupTracerProviderFromConfig(t.Context(), conf, baggagecopy.AllowAllMembers)

	// assert
	if assert.NotNil(t, shutdown) {
		assert.NoError(t, shutdown(t.Context()))
	}

	assert.NoError(t, err)
}

func Test_SetupTracerProviderFromConfig_Error(t *testing.T) {
	// arrange
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	conf := tracex.ProviderConfig{}

	// act
	shutdown, err := tracex.SetupTracerProviderFromConfig(ctx, conf, baggagecopy.AllowAllMembers)

	// assert
	assert.Nil(t, shutdown)
	assert.EqualError(t, err, "context canceled")
}
