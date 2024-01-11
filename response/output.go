package response

import (
	"encoding/json"
	"github.com/with-module/go-http/use"
	"net/http"
)

type (
	Status string
	Output interface {
		Write(w http.ResponseWriter)
	}
	Success[T any] struct {
		Status Status `json:"status" validate:"required" example:"success"`
		Data   T      `json:"data" validate:"required"`
	}

	response struct {
		status int
		data   any
	}
)

var encode = json.Marshal

func With(httpCode int, data any) Output {
	return &response{
		status: httpCode,
		data:   data,
	}
}

const StatusSuccess Status = "success"

func New[T any](data T) Output {
	return With(http.StatusOK, Success[T]{
		Status: StatusSuccess,
		Data:   data,
	})
}
func (r *response) Write(w http.ResponseWriter) {
	if use.IsZero(r.data) {
		w.WriteHeader(r.status)
		return
	}
	payload, err := encode(r.data)
	if err != nil {
		internalServerError.WriteDetail("failed to encode json response", err.Error()).Write(w)
		return
	}
	w.WriteHeader(r.status)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if _, err = w.Write(payload); err != nil {
		internalServerError.WriteDetail("failed to write response", err.Error()).Write(w)
		return
	}
}
