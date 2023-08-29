package middleware

import (
	"context"
	"github.com/google/uuid"
	"gitlab.com/with-junbach/go-modules/use"
	"net/http"
)

type RequestIDConfig struct {
	Skipper   Skipper
	Generator func() string
	HeaderKey string
}

const defaultRequestIDHeaderKey = "X-Request-Id"
const requestIDContextKey contextKey = "CTX_KEY_REQUEST_ID"

func RequestIDWithConfig(config RequestIDConfig) Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if fn := config.Skipper; fn != nil && fn(r) {
				next.ServeHTTP(w, r)
				return
			}
			headerKey := use.GetOrDefault(config.HeaderKey, defaultRequestIDHeaderKey)
			requestID := r.Header.Get(headerKey)
			if requestID != "" {
				fn := use.GetOrDefault(config.Generator, uuid.NewString)
				requestID = fn()
			}

			w.Header().Set(headerKey, requestID)
			ctx := context.WithValue(r.Context(), requestIDContextKey, requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func RequestID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		requestID := use.GetOrDefaultFunc(r.Header.Get(defaultRequestIDHeaderKey), uuid.NewString)
		w.Header().Set(defaultRequestIDHeaderKey, requestID)
		ctx := context.WithValue(r.Context(), requestIDContextKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func GetRequestID(ctx context.Context) string {
	return ctx.Value(requestIDContextKey).(string)
}
