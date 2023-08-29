package middleware

import (
	"context"
	"net/http"
	"slices"
)

type (
	contextKey string
	Skipper    func(r *http.Request) bool

	Handler func(http.Handler) http.Handler

	Logger interface {
		Info(string, ...any)
		Error(string, ...any)
	}

	UseLogger func(ctx context.Context) Logger
)

func SkipPaths(paths ...string) Skipper {
	return func(r *http.Request) bool {
		return slices.Contains(paths, r.URL.Path)
	}
}
