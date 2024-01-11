package server

import (
	"context"
	"errors"
	"github.com/with-module/go-http/use"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type (
	S struct {
		*http.Server
		config Config
		// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
		// Use a buffered channel to avoid missing signals as recommended for signal.Notify
		quitChn chan os.Signal
	}
)

func New(config Config, handler http.Handler, opts ...Option) *S {
	for _, fn := range opts {
		config = fn.use(config)
	}
	return &S{
		Server: &http.Server{
			Handler:      handler,
			Addr:         use.GetOrDefault(config.Addr, ":8080"),
			ReadTimeout:  use.GetOrDefault(config.Timeout.Read, 10*time.Second),
			WriteTimeout: use.GetOrDefault(config.Timeout.Write, 10*time.Second),
			IdleTimeout:  use.GetOrDefault(config.Timeout.Idle, 10*time.Second),
		},
		config:  config,
		quitChn: make(chan os.Signal, 1),
	}
}

func (s *S) ServeHttp() error {
	// start the server connection
	go OnServe(s, s.config.moreOptions.RuntimeErrorHandle)

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

type HttpSrv interface {
	ListenAndServe() error
}

func OnServe(srv HttpSrv, onRuntimeErr func(error)) {
	if err := srv.ListenAndServe(); err != nil && onRuntimeErr != nil && !errors.Is(err, http.ErrServerClosed) {
		onRuntimeErr(err)
	}
}
