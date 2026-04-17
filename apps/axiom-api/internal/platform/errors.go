package platform

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"error"`
	Detail  string `json:"detail,omitempty"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func ErrNotFound(msg string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: msg}
}

func ErrUnauthorized(msg string) *AppError {
	return &AppError{Code: http.StatusUnauthorized, Message: msg}
}

func ErrForbidden(msg string) *AppError {
	return &AppError{Code: http.StatusForbidden, Message: msg}
}

func ErrBadRequest(msg string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: msg}
}

func ErrConflict(msg string) *AppError {
	return &AppError{Code: http.StatusConflict, Message: msg}
}

func ErrValidation(msg string, detail string) *AppError {
	return &AppError{Code: http.StatusUnprocessableEntity, Message: msg, Detail: detail}
}

func ErrInternal(msg string) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: msg}
}
