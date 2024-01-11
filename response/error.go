package response

import (
	"errors"
	"net/http"
	"strings"
)

type (
	Error struct {
		HttpCode int    `json:"-" yaml:"httpCode"`
		Code     string `json:"code" yaml:"code" validate:"required" example:"INTERNAL_SERVER_ERROR"`
		Message  string `json:"message" yaml:"message" validate:"required" example:"internal server error"`
		Detail   string `json:"detail,omitempty" yaml:"detail,omitempty" example:"database connection error"`
	}

	Failure struct {
		Status Status `json:"status" validate:"required" example:"error"`
		Error  Error  `json:"error" validate:"required"`
	}
)

const StatusError Status = "error"

var errorDetailEnabled = false

var internalServerError = Error{
	HttpCode: http.StatusInternalServerError,
	Code:     "INTERNAL_SERVER_ERROR",
	Message:  "internal server error",
}

func EnableErrorDetail(enabled bool) {
	errorDetailEnabled = enabled
}

func Err(err error) Error {
	var e Error
	if errors.As(err, &e) {
		return e
	}
	return internalServerError.WriteDetail("undefined error", err.Error())
}

func (e Error) Error() string {
	var sb strings.Builder
	sb.WriteString(e.Code)
	sb.WriteString(": ")
	sb.WriteString(e.Message)
	if errorDetailEnabled {
		sb.WriteString("; ")
		sb.WriteString(e.Detail)
	}
	return sb.String()
}

func (e Error) Write(w http.ResponseWriter) {
	if !errorDetailEnabled {
		e.Detail = ""
	}
	payload, err := encode(Failure{Status: StatusError, Error: e})
	if err != nil {
		http.Error(w, internalServerError.Error(), internalServerError.HttpCode)
		return
	}
	w.WriteHeader(e.HttpCode)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if _, err := w.Write(payload); err != nil {
		http.Error(w, internalServerError.Error(), internalServerError.HttpCode)
		return
	}
}

func (e Error) WriteDetail(args ...string) Error {
	if e.Detail != "" {
		args = append([]string{e.Detail}, args...)
	}
	e.Detail = strings.Join(args, "; ")
	return e
}
