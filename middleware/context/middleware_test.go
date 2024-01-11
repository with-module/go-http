package context

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/with-module/go-http/middleware"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

const ctxKey = "_/text/custom"
const ctxValue = "sample-text-value"

func TestBindContext(t *testing.T) {
	var request = httptest.NewRequest(http.MethodGet, "http://localhost:8080/profile", nil)
	var response = httptest.NewRecorder()
	config := Config{
		Skipper: middleware.SkipPaths("/health"),
		Binder: func(r *http.Request) context.Context {
			return context.WithValue(context.Background(), ctxKey, ctxValue)
		},
	}
	var handle http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		assert.Equal(t, ctx.Value(ctxKey), ctxValue)
		assert.Equal(t, http.MethodGet, request.Method)
		log.Printf("\nsample handle: %s", ctxValue)
	}
	Middleware(config)(handle).ServeHTTP(response, request)
}

func TestBindContextSkipped(t *testing.T) {
	var request = httptest.NewRequest(http.MethodGet, "http://localhost:8080/profile", nil)
	var response = httptest.NewRecorder()
	config := Config{
		Skipper: middleware.SkipPaths("/profile"),
		Binder: func(r *http.Request) context.Context {
			return context.WithValue(context.Background(), ctxKey, ctxValue)
		},
	}
	var handle http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		assert.Equal(t, ctx.Value(ctxKey), nil)
		assert.Equal(t, http.MethodGet, request.Method)
		log.Printf("\nsample handle: %s", ctxValue)
	}
	Middleware(config)(handle).ServeHTTP(response, request)

	nilBinderConfig := Config{
		Skipper: middleware.SkipPaths("/health"),
		Binder:  nil,
	}
	Middleware(nilBinderConfig)(handle).ServeHTTP(response, request)
}
