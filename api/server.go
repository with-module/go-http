package api

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type (
	ServerConfig struct {
		Addr    string              `json:"addr" yaml:"addr" config:"addr"`
		Timeout ServerTimeoutConfig `json:"timeout" yaml:"timeout" config:"timeout"`

		RuntimeErrorHandler func(error)
	}

	ServerTimeoutConfig struct {
		Shutdown time.Duration `json:"shutdown" yaml:"shutdown" config:"shutdown"`
		Read     time.Duration `json:"read" yaml:"read" config:"read"`
		Write    time.Duration `json:"write" yaml:"write" config:"write"`
		Idle     time.Duration `json:"idle" yaml:"idle" config:"idle"`
	}

	Server struct {
		*http.Server
		config ServerConfig
		// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
		// Use a buffered channel to avoid missing signals as recommended for signal.Notify
		quitChn chan os.Signal
	}

	OverrideServerConfig func(config ServerConfig) ServerConfig
)

func InitHttpServer(config ServerConfig, handler http.Handler, fns ...OverrideServerConfig) *Server {
	for _, fn := range fns {
		config = fn(config)
	}
	return &Server{
		Server: &http.Server{
			Handler:      handler,
			Addr:         config.Addr,
			ReadTimeout:  config.Timeout.Read,
			WriteTimeout: config.Timeout.Write,
			IdleTimeout:  config.Timeout.Idle,
		},
		config:  config,
		quitChn: make(chan os.Signal, 1),
	}
}

func (s *Server) ServeHttp() error {
	// start the server connection
	go s.serve()

	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can't be caught, so don't need to add it
	signal.Notify(s.quitChn, syscall.SIGINT, syscall.SIGTERM)
	<-s.quitChn

	// cancelable context to help cancel the halted shutdown process
	ctx, cancel := context.WithTimeout(context.Background(), s.config.Timeout.Shutdown)
	defer cancel()

	// perform shutdown then wait until the server finished shutdown process or the timeout had been reached
	return s.Shutdown(ctx)
}

func (s *Server) serve() {
	err := s.ListenAndServe()
	if fn := s.config.RuntimeErrorHandler; fn != nil && err != nil && !errors.Is(err, http.ErrServerClosed) {
		fn(err)
	}
}

func WithHandleRuntimeErr(fn func(error)) OverrideServerConfig {
	return func(config ServerConfig) ServerConfig {
		config.RuntimeErrorHandler = fn
		return config
	}
}
