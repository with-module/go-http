package middleware

import (
	"net/http"
	"slices"
)

type (
	Skipper func(r *http.Request) bool

	Handler func(http.Handler) http.Handler
)

func SkipPaths(paths ...string) Skipper {
	return func(r *http.Request) bool {
		return slices.Contains(paths, r.URL.Path)
	}
}
