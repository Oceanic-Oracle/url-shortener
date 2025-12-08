package service

import (
	"errors"
	"net/http"
)

var (
	ErrURLCollision = errors.New("URL collision")
)

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
	cause   error  `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.cause
}

func WrapErrNotFound(err error) error {
	return &AppError{
		Code:    "NOT_FOUND",
		Message: "resource not found",
		Status:  http.StatusNotFound,
		cause:   err,
	}
}

func WrapErrInternalServer(err error) error {
	return &AppError{
		Code:    "INTERNAL",
		Message: "internal server error",
		Status:  http.StatusInternalServerError,
		cause:   err,
	}
}

func WrapErrBadRequest(err error) error {
	return &AppError{
		Code:    "BAD_REQUEST",
		Message: "bad request",
		Status:  http.StatusBadRequest,
		cause:   err,
	}
}
