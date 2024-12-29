package myerr

import (
	"errors"
	"fmt"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// ErrorType defines the type of error for categorization
type ErrorType int

const (
	ErrorTypeNotFound ErrorType = iota + 1
	ErrorTypeValidation
	ErrorTypeInternal
	ErrorTypeDuplicate
)

// AppError represents a custom error structure with an error type
type AppError struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewAppError(t ErrorType, msg string, err error) error {
	return &AppError{
		Type:    t,
		Message: msg,
		Err:     err,
	}
}

// ToAppError converts an any error to an AppError or returns nil if the error is nil
func ToAppError(log log.Logger, err error, description string) error {
	if err == nil {
		return nil
	}
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	_ = level.Error(log).Log("msg", description, "err", err)
	// Log the error and return an internal error, without exposing the original error to security risks
	return NewAppError(ErrorTypeInternal, "internal error", nil)
}
