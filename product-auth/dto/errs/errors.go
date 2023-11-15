package errs

import "net/http"

type AppError struct {
	Code int `json:",omitempty"`
	Message string `json:"message"`
	Error error `json:",omitempty`
}

func NewNotFoundError(message string, err error) *AppError {
	return &AppError{
		Message: message,
		Code: http.StatusNotFound,
		Error: err,
	}
}

func NewUnexpectedError(message string, err error) *AppError {
	return &AppError{
		Message: message,
		Code: http.StatusInternalServerError,
		Error: err,
	}
}

func NewAuthenticationError(message string, err error) *AppError {
     return &AppError{
		Message: message,
		Code: http.StatusUnauthorized,
		Error: err,
	 }
}

func NewAuthorizationError(message string, err error) *AppError {
	return &AppError{
		Message: message,
		Code: http.StatusForbidden,
		Error: err,
	}
}