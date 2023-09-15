package api

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestResponseWithOptions(t *testing.T) {
	responseWithData := ResponseWith(WithHttp(http.StatusOK), WithData("string-value"))
	assert.IsType(t, new(response), responseWithData)
	assert.Equal(t, http.StatusOK, responseWithData.StatusCode())
	assert.IsType(t, ResponseData[any]{}, responseWithData.Body())

	err := mapErr(errors.New("undefined error"))
	t.Logf("error: %s", err.Error())
	responseWithErr := ResponseWith(WithErr(err))
	assert.Equal(t, http.StatusInternalServerError, responseWithErr.StatusCode())

	anotherResponseErr := ResponseErr(ISE())
	assert.Equal(t, http.StatusInternalServerError, anotherResponseErr.StatusCode())
}
