package middleware

import (
	"fmt"
	"net/http"
	"time"
)

type (
	RequestLoggerConfig struct {
		Skipper   Skipper
		GetLogger UseLogger
		GetStatus func(w http.ResponseWriter, r *http.Request) int
	}
)

func RequestLogger(config RequestLoggerConfig) Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if fn := config.Skipper; fn != nil && fn(r) {
				next.ServeHTTP(w, r)
				return
			}
			startTime := time.Now()
			next.ServeHTTP(w, r)
			log := config.GetLogger(r.Context())
			if log == nil {
				return
			}
			status := config.GetStatus(w, r)
			fn := log.Info
			if status > 399 {
				fn = log.Error
			}

			fn(fmt.Sprintf("%s request %s has been completed in %s with status %d", r.Method, r.RequestURI, time.Since(startTime), status))
		}
		return http.HandlerFunc(fn)
	}
}
