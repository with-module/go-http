package recover

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/with-module/go-http/middleware"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddleware(t *testing.T) {
	var request = httptest.NewRequest(http.MethodGet, "http://localhost:8080/profile", nil)
	var response = httptest.NewRecorder()
	var log = new(middleware.TestLogger)
	defaultGetLogger = func(_ context.Context) middleware.Logger {
		return log
	}
	err := errors.New("undefined error")
	var handle http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, "http://localhost:8080/profile", request.RequestURI)
		log.Info("into test handler")
		panic(err)
	}
	assert.NotPanics(t, func() {
		Middleware(handle).ServeHTTP(response, request)
	})
	output := log.String()
	assert.Contains(t, output, "into test handler")
	assert.Contains(t, output, "request panic")
	assert.Contains(t, output, err.Error())
}

func TestMiddlewareAbort(t *testing.T) {
	var request = httptest.NewRequest(http.MethodGet, "http://localhost:8080/profile", nil)
	var response = httptest.NewRecorder()
	var log = new(middleware.TestLogger)
	defaultGetLogger = func(_ context.Context) middleware.Logger {
		return log
	}
	var handle http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, "http://localhost:8080/profile", request.RequestURI)
		log.Info("into test handler")
		panic(http.ErrAbortHandler)
	}
	assert.PanicsWithError(t, http.ErrAbortHandler.Error(), func() {
		Middleware(handle).ServeHTTP(response, request)
	})
	assert.Contains(t, log.String(), "INFO into test handler")
}

func TestMiddlewareUseConfig(t *testing.T) {
	var request = httptest.NewRequest(http.MethodGet, "http://localhost:8080/profile", nil)
	var response = httptest.NewRecorder()
	var log = new(middleware.TestLogger)
	err := errors.New("undefined error")
	var handle http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, "http://localhost:8080/profile", request.RequestURI)
		log.Info("into test handler")
		panic(err)
	}
	assert.NotPanics(t, func() {
		MiddlewareUseConfig(Config{
			WriteResponse: func(w http.ResponseWriter, r *http.Request) {
				log.Error("got panic and need to write response")
			},
			GetLogger: func(_ context.Context) middleware.Logger {
				return log
			},
		})(handle).ServeHTTP(response, request)
	})
	output := log.String()
	assert.Contains(t, output, "into test handler")
	assert.Contains(t, output, "request panic")
	assert.Contains(t, output, "got panic and need to write response")
	assert.Contains(t, output, err.Error())
}
