package response

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type writer struct {
	status   int
	writeErr error
	header   http.Header
}

func (w *writer) WriteHeader(statusCode int) {
	w.status = statusCode
}

func (w *writer) Write([]byte) (int, error) {
	if w.writeErr != nil {
		return 0, w.writeErr
	}
	return 0, nil
}

func (w *writer) Header() http.Header {
	if w.header == nil {
		w.header = make(http.Header)
	}
	return w.header
}

func TestNew(t *testing.T) {
	data := []string{"a", "b", "c"}
	output := New(data)
	w := new(writer)
	output.Write(w)
	assert.Equal(t, http.StatusOK, w.status)
}
