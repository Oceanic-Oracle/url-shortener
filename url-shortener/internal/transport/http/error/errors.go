package httperror

import (
	"errors"
	"net/http"
)

var (
	ErrBadRequest = errors.New("bad request")
)

type httpError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
	cause   error  `json:"-"`
}

func (e *httpError) Error() string {
	return e.Message
}

func (e *httpError) Unwrap() error {
	return e.cause
}

func wrapErrBadRequest(err error) *httpError {
	return &httpError{
		Code:    "BAD_REQUEST",
		Message: "bad request",
		Status:  http.StatusBadRequest,
		cause:   err,
	}
}

func wrapErrNotFound(err error) *httpError {
	return &httpError{
		Code:    "NOT_FOUND",
		Message: "resource not found",
		Status:  http.StatusNotFound,
		cause:   err,
	}
}

func wrapErrConflict(err error) *httpError {
	return &httpError{
		Code:    "CONFLICT",
		Message: "resource collision",
		Status:  http.StatusConflict,
		cause:   err,
	}
}

func wrapErrInternalServer(err error) *httpError {
	return &httpError{
		Code:    "INTERNAL",
		Message: "internal server error",
		Status:  http.StatusInternalServerError,
		cause:   err,
	}
}
