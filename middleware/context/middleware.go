package context

import (
	"context"
	"github.com/with-module/go-http/middleware"
	"net/http"
)

type Config struct {
	Skipper middleware.Skipper
	Binder  func(r *http.Request) context.Context
}

func Middleware(config Config) middleware.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if fn := config.Skipper; fn != nil && fn(r) {
				next.ServeHTTP(w, r)
				return
			}
			if config.Binder == nil {
				next.ServeHTTP(w, r)
				return
			}
			ctx := config.Binder(r)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
