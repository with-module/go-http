package health

import (
	"github.com/with-module/go-http/middleware"
	"net/http"
	"slices"
	"strings"
)

func Middleware(path string) middleware.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if slices.Contains([]string{"GET", "HEAD"}, r.Method) && strings.EqualFold(r.URL.Path, path) {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
