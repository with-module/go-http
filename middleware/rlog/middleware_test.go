package rlog

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/with-module/go-http/middleware"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddleware(t *testing.T) {
	log := new(middleware.TestLogger)
	defaultGetLogger = func(_ context.Context) middleware.Logger {
		return log
	}
	var request = httptest.NewRequest(http.MethodGet, "http://localhost:8080/profile", nil)
	var response = httptest.NewRecorder()
	var handle http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.WriteHeader(http.StatusNoContent)
	}
	Middleware(handle).ServeHTTP(response, request)
	output := log.String()
	assert.Contains(t, output, "status 204")
	assert.Contains(t, output, "INFO")
}

func TestMiddlewareUseConfig(t *testing.T) {
	var request = httptest.NewRequest(http.MethodGet, "http://localhost:8080/profile", nil)
	var response = httptest.NewRecorder()
	var handle http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.WriteHeader(http.StatusBadRequest)
	}
	log := new(middleware.TestLogger)
	MiddlewareUseConfig(Config{
		Skipper: nil,
		GetLogger: func(ctx context.Context) middleware.Logger {
			return log
		},
	})(handle).ServeHTTP(response, request)
	output := log.String()
	assert.Contains(t, output, "status 400")
	assert.Contains(t, output, "ERROR")
}

func TestMiddlewareNoLog(t *testing.T) {
	var request = httptest.NewRequest(http.MethodGet, "http://localhost:8080/profile", nil)
	var response = httptest.NewRecorder()
	var handle http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.WriteHeader(http.StatusNoContent)
	}
	log := new(middleware.TestLogger)
	MiddlewareUseConfig(Config{})(handle).ServeHTTP(response, request)
	MiddlewareUseConfig(Config{
		GetLogger: func(ctx context.Context) middleware.Logger {
			return nil
		},
	})(handle).ServeHTTP(response, request)
	assert.Equal(t, log.String(), "")

}

func TestMiddlewareWithSkipper(t *testing.T) {
	var request = httptest.NewRequest(http.MethodGet, "http://localhost:8080/profile", nil)
	var response = httptest.NewRecorder()
	var handle http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.WriteHeader(http.StatusNoContent)
	}
	log := new(middleware.TestLogger)
	MiddlewareUseConfig(Config{
		Skipper: middleware.SkipPaths("/profile"),
		GetLogger: func(ctx context.Context) middleware.Logger {
			return log
		},
	})(handle).ServeHTTP(response, request)
	assert.Equal(t, log.String(), "")
}
