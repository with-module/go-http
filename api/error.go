package api

import (
	"errors"
	"fmt"
	"net/http"
)

type (
	HttpError struct {
		Http int `json:"-" yaml:"http" config:"http" validate:"required"`
		// unique code to identify the error
		Code string `json:"code" yaml:"code" validate:"required" config:"code" example:"INVALID_REQUEST"`
		// detail text to describe the error, may be varied by localization
		Message string `json:"message" yaml:"message" config:"message" validate:"required" example:"invalid request input"`
	}
)

var internalServerError = HttpError{
	Http:    http.StatusInternalServerError,
	Code:    "INTERNAL_SERVER_ERROR",
	Message: "internal server error",
}

func (e HttpError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func ISE() HttpError {
	return internalServerError
}
func mapErr(err error) HttpError {
	var he HttpError
	if errors.As(err, &he) {
		return he
	}
	return internalServerError
}
