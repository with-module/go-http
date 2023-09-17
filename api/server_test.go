package api

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"log"
	"syscall"
	"testing"
	"time"
)

func TestInitiateSimpleHttpServer(t *testing.T) {
	config := ServerConfig{
		Addr: ":8080",
		Timeout: ServerTimeoutConfig{
			Shutdown: time.Second * 10,
			Read:     time.Second * 20,
			Write:    time.Second * 20,
			Idle:     time.Second * 20,
		},
		RuntimeErrorHandler: nil,
	}

	srv := InitHttpServer(config, nil, WithHandleRuntimeErr(func(err error) {
		log.Fatalf("server runtime err: %v", err)
	}))

	assert.IsType(t, new(Server), srv)
	assert.Equal(t, ":8080", srv.config.Addr)
	go func() {
		time.Sleep(time.Second * 2)
		// stop server
		srv.quitChn <- syscall.SIGTERM
	}()
	assert.NoError(t, srv.ServeHttp())
}

type mockSrv struct{}

func (mock mockSrv) ListenAndServe() error {
	return errors.New("server runtime error")
}

func TestServerRuntime(t *testing.T) {
	var runtimeErr error
	OnServe(new(mockSrv), func(err error) {
		runtimeErr = err
	})

	assert.ErrorContains(t, runtimeErr, "server runtime error")
}
