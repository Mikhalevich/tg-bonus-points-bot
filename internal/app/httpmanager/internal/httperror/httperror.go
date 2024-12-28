package httperror

import (
	"fmt"
	"net/http"
)

type ErrorHTTPResponse struct {
	Code    int
	Message string
	Err     error
}

func New(code int, msg string) *ErrorHTTPResponse {
	return &ErrorHTTPResponse{
		Code:    code,
		Message: msg,
	}
}

func (e *ErrorHTTPResponse) WithError(err error) *ErrorHTTPResponse {
	e.Err = err
	return e
}

func InternalServerError(wrapContext string, err error) *ErrorHTTPResponse {
	return &ErrorHTTPResponse{
		Code:    http.StatusInternalServerError,
		Message: "internal error",
		Err:     fmt.Errorf("%s: %w", wrapContext, err),
	}
}

func BadRequest(msg string) *ErrorHTTPResponse {
	return New(http.StatusBadRequest, msg)
}

func NotFound(msg string) *ErrorHTTPResponse {
	return New(http.StatusNotFound, msg)
}
