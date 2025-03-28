package rpcserver

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strconv"
	"time"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const (
	readHeaderTimeout = time.Second * 3
)

// Config contains the configuration for creating a Server instance.
type Config struct {
	Host string `koanf:"host" default:"localhost"`
	Port uint16 `koanf:"port" default:"5000"`
}

// Handler implementations can register themselves to be hosted by the server.
type Handler interface {
	Register(*http.ServeMux, []connect.HandlerOption)
}

// Server provides a HTTP/2 server for one or more HTTP handlers.
type Server struct {
	server *http.Server
}

// New creates a new Server instance from the given configuration.
func NewFromConfig(config Config, handlers []Handler, opts []connect.HandlerOption) *Server {
	return New(config.Host, config.Port, handlers, opts)
}

// New creates a new Server instance.
func New(host string, port uint16, handlers []Handler, opts []connect.HandlerOption) *Server {
	mux := http.NewServeMux()

	for _, h := range handlers {
		h.Register(mux, opts)
	}

	return &Server{
		server: &http.Server{
			Addr:              net.JoinHostPort(host, strconv.FormatUint(uint64(port), 10 /*base*/)),
			Handler:           h2c.NewHandler(mux, &http2.Server{}),
			ReadHeaderTimeout: readHeaderTimeout,
		},
	}
}

// Start starts the server. This blocks until the server stops.
func (s *Server) Start() error {
	err := s.server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
