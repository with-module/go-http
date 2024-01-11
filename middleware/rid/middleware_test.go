package rid

import (
	"github.com/stretchr/testify/assert"
	"github.com/with-module/go-http/middleware"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestID(t *testing.T) {
	var request = httptest.NewRequest(http.MethodGet, "http://localhost:8080/profile", nil)
	var response = httptest.NewRecorder()
	var handle http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		rid := Get(ctx)
		assert.NotEmpty(t, rid)
		log.Printf("\nsample handle: %s", rid)
	}
	Middleware(handle).ServeHTTP(response, request)
}

func TestRequestIDWithConfig(t *testing.T) {
	var request = httptest.NewRequest(http.MethodGet, "http://localhost:8080/profile", nil)
	var response = httptest.NewRecorder()
	mdw := MiddlewareUseConfig(Config{
		Skipper: middleware.SkipPaths("/health", "/ping"),
		Generator: func() string {
			return "0000"
		},
		HeaderKey: "request-id",
	})

	var handle http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		rid := Get(ctx)
		assert.Equal(t, "0000", rid)
	}

	mdw(handle).ServeHTTP(response, request)

	var skipHandle http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		rid := Get(ctx)
		assert.Empty(t, rid)
	}
	var healthCheckRequest = httptest.NewRequest(http.MethodGet, "http://localhost:8080/health", nil)
	mdw(skipHandle).ServeHTTP(response, healthCheckRequest)
}
