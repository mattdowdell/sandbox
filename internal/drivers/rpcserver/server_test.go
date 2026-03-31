package rpcserver_test

import (
	"net/http"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mattdowdell/sandbox/internal/drivers/rpcserver"
	"github.com/mattdowdell/sandbox/mocks/drivers/mockrpcserver"
)

var mockServeMux = mock.AnythingOfType("*http.ServeMux")

func Test_NewFromConfig(t *testing.T) {
	// arrange
	conf := rpcserver.Config{
		Host:        "127.0.0.1",
		Port:        0,
		EnablePprof: true,
	}

	var opts []connect.HandlerOption

	handler := mockrpcserver.NewHandler(t)
	handler.EXPECT().Register(mockServeMux, opts).Once()

	handlers := []rpcserver.Handler{
		handler,
	}

	// act
	server := rpcserver.NewFromConfig(conf, handlers, opts)

	// assert
	assert.NotNil(t, server)
}

func Test_New(t *testing.T) {
	// arrange
	var opts []connect.HandlerOption

	handler := mockrpcserver.NewHandler(t)
	handler.EXPECT().Register(mockServeMux, opts).Once()

	handlers := []rpcserver.Handler{
		handler,
	}

	// act
	server := rpcserver.New(
		"127.0.0.1", /*host*/
		0,           /*port*/
		true,        /*enablePprof*/
		true,        /*enableCoverage*/
		handlers,
		opts,
	)

	// assert
	assert.NotNil(t, server)
}

func Test_Server_Start(t *testing.T) {
	// arrange
	server := rpcserver.New(
		"127.0.0.1", /*host*/
		0,           /*port*/
		true,        /*enablePprof*/
		true,        /*enableCoverage*/
		nil,         /*handlers*/
		nil,         /*opts*/
	)

	// act
	go func() {
		assert.NoError(t, server.Start(t.Context()))
	}()

	defer func() {
		assert.NoError(t, server.Shutdown(t.Context()))
	}()

	// assert
	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodGet,
		server.URL("/debug/pprof/"),
		http.NoBody,
	)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
