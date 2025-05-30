//nolint:dupl // similar to meter_test.go/tracer_test.go
package otelx_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/contrib/processors/baggagecopy"

	"github.com/mattdowdell/sandbox/internal/drivers/otelx"
)

func urlToTracerProviderConfig(t *testing.T, input string) otelx.TracerProviderConfig {
	t.Helper()

	u, err := url.Parse(input)
	require.NoError(t, err)

	return otelx.TracerProviderConfig{
		Insecure: isHTTP(u.Scheme),
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

	conf := urlToTracerProviderConfig(t, server.URL)

	// act
	shutdown, err := otelx.SetupTracerProviderFromConfig(t.Context(), conf, baggagecopy.AllowAllMembers)

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

	conf := otelx.TracerProviderConfig{}

	// act
	shutdown, err := otelx.SetupTracerProviderFromConfig(ctx, conf, baggagecopy.AllowAllMembers)

	// assert
	assert.Nil(t, shutdown)
	assert.EqualError(t, err, "context canceled")
}

func Test_Tracer(t *testing.T) {
	// arrange

	// act
	meter := otelx.Tracer()

	// assert
	assert.NotNil(t, meter)
}
