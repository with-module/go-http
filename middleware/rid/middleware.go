package rid

import (
	"context"
	"github.com/google/uuid"
	"github.com/with-module/go-http/middleware"
	"github.com/with-module/go-http/use"
	"net/http"
)

type Config struct {
	Skipper   middleware.Skipper
	Generator func() string
	HeaderKey string
}

type contextKey string

const defaultRequestIDHeaderKey = "X-Request-Id"
const requestIDContextKey contextKey = "_http/ctx-key/request-id"

func MiddlewareUseConfig(config Config) middleware.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if fn := config.Skipper; fn != nil && fn(r) {
				next.ServeHTTP(w, r)
				return
			}
			headerKey := use.GetOrDefault(config.HeaderKey, defaultRequestIDHeaderKey)
			requestID := r.Header.Get(headerKey)
			if use.IsZero(requestID) {
				requestID = use.If(config.Generator != nil, config.Generator, uuid.NewString)()
			}

			w.Header().Set(headerKey, requestID)
			ctx := context.WithValue(r.Context(), requestIDContextKey, requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Middleware(next http.Handler) http.Handler {
	return MiddlewareUseConfig(Config{
		Skipper:   nil,
		Generator: uuid.NewString,
		HeaderKey: defaultRequestIDHeaderKey,
	})(next)
}

func Get(ctx context.Context) string {
	if rid, exist := ctx.Value(requestIDContextKey).(string); exist {
		return rid
	}
	return ""
}
