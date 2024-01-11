package recover

import (
	"github.com/with-module/go-http/middleware"
	"net/http"
)

type Config struct {
	WriteResponse func(w http.ResponseWriter, r *http.Request)
	GetLogger     middleware.UseLogger
}

func MiddlewareUseConfig(config Config) middleware.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rc := recover(); rc != nil {
					if rc == http.ErrAbortHandler {
						panic(rc)
					}
					if config.GetLogger != nil {
						if log := config.GetLogger(r.Context()); log != nil {
							log.Error("request panic", "method", r.Method, "uri", r.RequestURI, "panic_recover", rc)
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
		})
	}
}

var defaultGetLogger = middleware.DefaultLogger

func Middleware(next http.Handler) http.Handler {
	return MiddlewareUseConfig(Config{
		WriteResponse: nil,
		GetLogger:     defaultGetLogger,
	})(next)
}
