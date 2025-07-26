package rpcserver

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/pprof"
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
	Host        string `koanf:"host" default:"localhost"`
	Port        uint16 `koanf:"port" default:"5000"`
	EnablePprof bool   `koanf:"enablepprof"`
}

// Handler implementations can register themselves to be hosted by the server.
type Handler interface {
	Register(*http.ServeMux, []connect.HandlerOption)
}

// Server provides a HTTP/2 server for one or more HTTP handlers.
type Server struct {
	server *http.Server
	url    string
	ch     chan struct{}
}

// New creates a new Server instance from the given configuration.
func NewFromConfig(config Config, handlers []Handler, opts []connect.HandlerOption) *Server {
	return New(config.Host, config.Port, config.EnablePprof, handlers, opts)
}

// New creates a new Server instance.
func New(
	host string,
	port uint16,
	enablePprof bool,
	handlers []Handler,
	opts []connect.HandlerOption,
) *Server {
	mux := http.NewServeMux()

	for _, h := range handlers {
		h.Register(mux, opts)
	}

	if enablePprof {
		mux.HandleFunc("GET /debug/pprof/", pprof.Index)
		mux.HandleFunc("GET /debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("GET /debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("GET /debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("GET /debug/pprof/trace", pprof.Trace)
	}

	return &Server{
		server: &http.Server{
			Addr:              net.JoinHostPort(host, strconv.FormatUint(uint64(port), 10 /*base*/)),
			Handler:           h2c.NewHandler(mux, &http2.Server{}),
			ReadHeaderTimeout: readHeaderTimeout,
		},
		ch: make(chan struct{}),
	}
}

// Start starts the server and blocks until the server stops.
func (s *Server) Start(ctx context.Context) error {
	var lc net.ListenConfig
	listener, err := lc.Listen(ctx, "tcp", s.server.Addr)
	if err != nil {
		return err
	}

	s.url = listener.Addr().String()
	close(s.ch)

	if err := s.server.Serve(listener); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
