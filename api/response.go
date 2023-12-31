package api

import "net/http"

type (
	response struct {
		httpCode    int
		bodyContent any
	}

	Response interface {
		StatusCode() int
		Body() any
	}

	ResponseStatus string

	ResponseOption interface {
		use(resp *response)
	}

	useOption func(resp *response)

	ResponseData[T any] struct {
		// status response, can be `success` or `error`
		Status ResponseStatus `json:"status" validate:"required" enums:"success" example:"success"`
		// data response after processing request, available if `status == "success"`
		Data T `json:"data,omitempty"`
	}

	ResponseError struct {
		// status response, can be `success` or `error`
		Status ResponseStatus `json:"status" validate:"required" enums:"error" example:"error"`
		// error encountered while processing request, available if `status == "error"`
		Error HttpError `json:"error" validate:"required"`
	}
)

const (
	StatusSuccess ResponseStatus = "success"
	StatusError   ResponseStatus = "error"
)

func ResponseWith(opts ...ResponseOption) Response {
	resp := &response{
		httpCode:    http.StatusOK,
		bodyContent: nil,
	}
	for _, fn := range opts {
		fn.use(resp)
	}
	return resp
}

func ResponseErr(err error) Response {
	return ResponseWith(WithErr(err))
}

func (resp *response) StatusCode() int {
	return resp.httpCode
}

func (resp *response) Body() any {
	return resp.bodyContent
}

func (fn useOption) use(r *response) {
	fn(r)
}

func WithErr(err error) ResponseOption {
	return useOption(func(r *response) {
		he := mapErr(err)
		r.httpCode = he.Http
		r.bodyContent = ResponseError{
			Status: StatusError,
			Error:  he,
		}
	})
}

func WithHttp(code int) ResponseOption {
	return useOption(func(r *response) {
		r.httpCode = code
	})
}

func WithData(data any) ResponseOption {
	return useOption(func(r *response) {
		r.bodyContent = ResponseData[any]{
			Status: StatusSuccess,
			Data:   data,
		}
	})
}
