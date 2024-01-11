package health

import (
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	var request = httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	var response = httptest.NewRecorder()
	var handle http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		assert.Equal(t, "/", request.URL.Path)
		log.Printf("sample handle")
	}

	Middleware("/health")(handle).ServeHTTP(response, request)

	var hit = false
	var skipHandle http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
		hit = true
	}
	Middleware("/")(skipHandle).ServeHTTP(response, request)
	assert.Equal(t, false, hit)
}
