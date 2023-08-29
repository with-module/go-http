package middleware

import (
	"net/http"
	"slices"
	"strings"
)

func HealthCheck(path string) Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if slices.Contains([]string{"GET", "HEAD"}, r.Method) && strings.EqualFold(r.URL.Path, path) {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
