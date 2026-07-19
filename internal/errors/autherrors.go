package errors

import (
	"net/http"
)

type AuthError struct {
	responseError   string
	underlyingError error
}

func NewAuthError(responseError string, underlyingError error) *AuthError {
	return &AuthError{responseError, underlyingError}
}

func (e *AuthError) Error() string {
	return e.underlyingError.Error()
}

func (e *AuthError) ResponseError() string {
	return e.responseError
}

func (e *AuthError) Code() int {
	return http.StatusUnauthorized
}

type ServerError struct {
	underlyingError error
}

func NewServerError(underlyingError error) *ServerError {
	return &ServerError{underlyingError}
}

func (e *ServerError) Error() string {
	return e.underlyingError.Error()
}

func (e *ServerError) Code() int {
	return http.StatusInternalServerError
}

func (e *ServerError) ResponseError() string {
	return "An unexpected error occurred."
}

type BadRequestError struct {
	responseError   string
	underlyingError error
}

func NewBadRequestError(responseError string, underlyingError error) *BadRequestError {
	return &BadRequestError{responseError, underlyingError}
}

func (b *BadRequestError) Error() string {
	return b.underlyingError.Error()
}

func (b *BadRequestError) ResponseError() string {
	return b.responseError
}

func (b *BadRequestError) Code() int {
	return http.StatusBadRequest
}
