package apperror

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrNotFound = NewAppError(
		http.StatusNotFound,
		"requested resource is not found",
		"help")
	ErrEmptyString        = errors.New("empty string")
	ErrInvalidRequestBody = errors.New("invalid request body")
)

type AppError struct {
	Err              error  `json:"-"`
	Message          string `json:"message,omitempty"`
	DeveloperMessage string `json:"developer_message,omitempty"`
	Code             int    `json:"code,omitempty"`
}

func NewAppError(code int, developerMessage, message string) *AppError {
	return &AppError{
		Err:              fmt.Errorf(message),
		Message:          message,
		DeveloperMessage: developerMessage,
		Code:             code,
	}
}

func (e *AppError) Error() string {
	return e.Err.Error()
}
