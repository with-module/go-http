package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
)

type RecoverConfig struct {
	WriteResponse func(w http.ResponseWriter, r *http.Request)
	GetLogger     UseLogger
}

func RecoverWithConfig(config RecoverConfig) Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rc := recover(); rc != nil {
					if rc == http.ErrAbortHandler {
						panic(rc)
					}
					if config.GetLogger != nil {
						if log := config.GetLogger(r.Context()); log != nil {
							log.Error(fmt.Sprintf("handler of %s request %s crashed and panic", r.Method, r.RequestURI), "panic_recover", rc)
						}
					}

					if r.Header.Get("Connection") != "Upgrade" {
						if fn := config.WriteResponse; fn != nil {
							fn(w, r)
						} else {
							w.WriteHeader(http.StatusInternalServerError)
						}
					}
				}
			}()
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func Recover(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rc := recover(); rc != nil {
				if rc == http.ErrAbortHandler {
					panic(rc)
				}
				slog.Error("request handler crashed and panic", slog.Any("panic_recover", rc))
				if r.Header.Get("Connection") != "Upgrade" {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
