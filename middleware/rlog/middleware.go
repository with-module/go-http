package rlog

import (
	"github.com/with-module/go-http/middleware"
	"github.com/with-module/go-http/use"
	"net/http"
	"time"
)

type Config struct {
	Skipper   middleware.Skipper
	GetLogger middleware.UseLogger
}

func MiddlewareUseConfig(config Config) middleware.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if fn := config.Skipper; fn != nil && fn(r) {
				next.ServeHTTP(w, r)
				return
			}
			if config.GetLogger == nil {
				next.ServeHTTP(w, r)
				return
			}
			var logger = config.GetLogger(r.Context())
			if logger == nil {
				next.ServeHTTP(w, r)
				return
			}
			startTime := time.Now()
			writer := nw(w)
			next.ServeHTTP(writer, r)
			printLog := use.If(writer.statusCode > 399, logger.Error, logger.Info)
			printLog(
				"request completed",
				"status", writer.statusCode,
				"latency", time.Now().Sub(startTime),
				"method", r.Method,
				"uri", r.RequestURI,
			)
		})
	}
}

var defaultGetLogger = middleware.DefaultLogger

func Middleware(next http.Handler) http.Handler {
	return MiddlewareUseConfig(Config{
		Skipper:   nil,
		GetLogger: defaultGetLogger,
	})(next)
}
