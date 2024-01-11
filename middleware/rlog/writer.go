package rlog

import (
	"net/http"
)

type customWriter struct {
	http.ResponseWriter
	statusCode int
}

func nw(w http.ResponseWriter) *customWriter {
	return &customWriter{
		ResponseWriter: w,
	}
}

func (cw *customWriter) WriteHeader(code int) {
	cw.statusCode = code
	cw.ResponseWriter.WriteHeader(code)
}
