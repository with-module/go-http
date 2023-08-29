package api

import "net/http"

type (
	response struct {
		httpCode    int
		bodyContent any
	}

	ResponseState interface {
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
		Data T `json:"data" validate:"required"`
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

func Response(opts ...ResponseOption) ResponseState {
	resp := &response{
		httpCode:    http.StatusOK,
		bodyContent: nil,
	}
	for _, fn := range opts {
		fn.use(resp)
	}
	return resp
}

func ResponseErr(err error) ResponseState {
	resp := new(response)
	WithErr(err).use(resp)
	return resp
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
