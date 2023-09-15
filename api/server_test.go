package api

import (
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
	assert.Nil(t, srv.Handler)
	assert.Equal(t, ":8080", srv.config.Addr)
	assert.NotNil(t, srv.config.RuntimeErrorHandler)

	go func() {
		assert.NoError(t, srv.ServeHttp())
	}()

	// stop server
	srv.quitChn <- syscall.SIGTERM
}

func TestServerRuntime(t *testing.T) {
	config := ServerConfig{
		Addr: "invalid-config-addr",
		Timeout: ServerTimeoutConfig{
			Shutdown: time.Second * 10,
			Read:     time.Second * 20,
			Write:    time.Second * 20,
			Idle:     time.Second * 20,
		},
		RuntimeErrorHandler: nil,
	}

	srv := InitHttpServer(config, nil, WithHandleRuntimeErr(func(err error) {
		log.Panicf("server runtime err: %v", err)
	}))

	assert.IsType(t, new(Server), srv)
	assert.Nil(t, srv.Handler)
	assert.NotNil(t, srv.config.RuntimeErrorHandler)

	assert.Panics(t, srv.serve)
}
