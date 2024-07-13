package errors

import "net/http"

type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return e.Message
}

func NewNotFoundError(message string) *APIError {
	return &APIError{
		StatusCode: http.StatusNotFound,
		Message:    message,
	}
}

func NewBadRequestError(message string) *APIError {
	return &APIError{
		StatusCode: http.StatusBadRequest,
		Message:    message,
	}
}

func NewInternalServerError(message string) *APIError {
	return &APIError{
		StatusCode: http.StatusInternalServerError,
		Message:    message,
	}
}