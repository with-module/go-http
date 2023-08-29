package api

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type (
	ServerConfig struct {
		Addr    string `json:"addr" yaml:"addr" config:"addr"`
		Timeout struct {
			Shutdown time.Duration `json:"shutdown" yaml:"shutdown" config:"shutdown"`
			Read     time.Duration `json:"read" yaml:"read" config:"read"`
			Write    time.Duration `json:"write" yaml:"write" config:"write"`
			Idle     time.Duration `json:"idle" yaml:"idle" config:"idle"`
		} `json:"timeout" yaml:"timeout" config:"timeout"`
	}

	RouteConfig struct {
		Enabled bool   `config:"enabled" json:"enabled" yaml:"enabled"`
		Method  string `config:"method" json:"method" yaml:"method"`
		Path    string `config:"path" json:"path" yaml:"path"`
	}

	Server struct {
		*http.Server
		config ServerConfig
	}
)

var (
	httpSvr      *Server
	httpSvrSetup sync.Once
)

func InitHttpServer(config ServerConfig, handler http.Handler) {
	httpSvrSetup.Do(func() {
		httpSvr = &Server{
			Server: &http.Server{
				Handler:      handler,
				Addr:         config.Addr,
				ReadTimeout:  config.Timeout.Read,
				WriteTimeout: config.Timeout.Write,
				IdleTimeout:  config.Timeout.Idle,
			},
			config: config,
		}
	})
}

func ServeHttp() error {
	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)

	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	var shutdownErr error
	var svrShutdown sync.WaitGroup
	// start the server connection
	go func() {
		<-quit
		svrShutdown.Add(1)
		defer svrShutdown.Done()

		// cancelable context to help cancel the halted shutdown process
		ctx, cancel := context.WithTimeout(context.Background(), httpSvr.config.Timeout.Shutdown)
		defer cancel()

		// perform shutdown then wait until the server finished shutdown process or the timeout had been reached
		shutdownErr = httpSvr.Shutdown(ctx)
	}()

	if err := httpSvr.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return shutdownErr
}
