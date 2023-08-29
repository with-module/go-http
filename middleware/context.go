package middleware

import (
	"context"
	"net/http"
)

type (
	BindContextConfig struct {
		Skipper Skipper
		Binder  func(r *http.Request) context.Context
	}
)

func BindContext(config BindContextConfig) Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
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
		}
		return http.HandlerFunc(fn)
	}
}
