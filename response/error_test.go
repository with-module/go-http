package response

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestEnableErrorDetail(t *testing.T) {
	EnableErrorDetail(false)
	assert.Equal(t, false, errorDetailEnabled)
	EnableErrorDetail(true)
	assert.Equal(t, true, errorDetailEnabled)
}

func TestErr(t *testing.T) {
	inputErr := errors.New("database connection error")
	err := Err(inputErr).WriteDetail("connection timeout")
	assert.NotNil(t, err)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.HttpCode)
	assert.Equal(t, "INTERNAL_SERVER_ERROR", err.Code)
	EnableErrorDetail(true)
	assert.Equal(t, true, errorDetailEnabled)
	assert.Contains(t, err.Error(), inputErr.Error())
	httpErr := Error{
		HttpCode: http.StatusBadRequest,
		Code:     "INVALID_REQUEST",
		Message:  "invalid request body",
	}
	httpErr = Err(httpErr)
	assert.Equal(t, http.StatusBadRequest, httpErr.HttpCode)
	assert.Equal(t, "INVALID_REQUEST", httpErr.Code)
}

func TestError_Write(t *testing.T) {
	EnableErrorDetail(false)
	assert.Equal(t, false, errorDetailEnabled)
	err := Error{
		HttpCode: http.StatusUnauthorized,
		Code:     "INVALID_CREDENTIAL",
		Message:  "user not found",
	}
	assert.Error(t, err)
	w := new(writer)
	err.Write(w)
	assert.Equal(t, http.StatusUnauthorized, w.status)

	w.writeErr = errors.New("write response error")
	err.Write(w)
	assert.Equal(t, http.StatusInternalServerError, w.status)

	encode = func(v any) ([]byte, error) {
		return nil, errors.New("encode data error")
	}
	defer func() {
		encode = json.Marshal
	}()
	err.Write(w)
	assert.Equal(t, http.StatusInternalServerError, w.status)
}
