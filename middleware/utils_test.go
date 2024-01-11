package middleware

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSkipPaths(t *testing.T) {
	var fn = SkipPaths("/profile")
	assert.NotNil(t, fn)
	var request = httptest.NewRequest(http.MethodGet, "http://localhost:8080/profile", nil)
	assert.Equal(t, true, fn(request))
}
