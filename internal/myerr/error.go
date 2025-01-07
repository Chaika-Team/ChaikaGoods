package myerr

import (
	"errors"
	"fmt"
)

// ErrorType defines the type of error for categorization
type ErrorType string

const (
	ErrorTypeNotFound     ErrorType = "NOT_FOUND"
	ErrorTypeValidation   ErrorType = "VALIDATION_ERROR"
	ErrorTypeInternal     ErrorType = "INTERNAL_ERROR"
	ErrorTypeDuplicate    ErrorType = "DUPLICATE_ERROR"
	ErrorTypeUnauthorized ErrorType = "UNAUTHORIZED"
	ErrorTypeForbidden    ErrorType = "FORBIDDEN"
	ErrorTypeConflict     ErrorType = "CONFLICT"
	ErrorTypeUnknown      ErrorType = "UNKNOWN"
)

// AppError represents a structured error with type and optional context
type AppError struct {
	Type    ErrorType              `json:"type"`              // High-level error type
	Message string                 `json:"message"`           // Human-readable message
	Cause   error                  `json:"-"`                 // Underlying cause of the error (optional)
	Context map[string]interface{} `json:"context,omitempty"` // Additional context (optional)
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap allows errors.Unwrap to extract the cause
func (e *AppError) Unwrap() error {
	return e.Cause
}

// New creates a new instance of AppError
func New(errType ErrorType, message string, cause error, context map[string]interface{}) *AppError {
	return &AppError{
		Type:    errType,
		Message: message,
		Cause:   cause,
		Context: context,
	}
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

// Wrap creates a new AppError wrapping an existing error
func Wrap(err error, errType ErrorType, message string, context map[string]interface{}) *AppError {
	if err == nil {
		return nil
	}
	return New(errType, message, err, context)
}

// PredefinedErrorTypes provides reusable error constructors for common cases
func NotFound(message string, cause error) *AppError {
	return New(ErrorTypeNotFound, message, cause, nil)
}

func Validation(message string, cause error) *AppError {
	return New(ErrorTypeValidation, message, cause, nil)
}

func Internal(message string, cause error) *AppError {
	return New(ErrorTypeInternal, message, cause, nil)
}

func Duplicate(message string, cause error) *AppError {
	return New(ErrorTypeDuplicate, message, cause, nil)
}

func Unauthorized(message string, cause error) *AppError {
	return New(ErrorTypeUnauthorized, message, cause, nil)
}

func Forbidden(message string, cause error) *AppError {
	return New(ErrorTypeForbidden, message, cause, nil)
}

func Conflict(message string, cause error) *AppError {
	return New(ErrorTypeConflict, message, cause, nil)
}

func Unknown(message string, cause error) *AppError {
	return New(ErrorTypeUnknown, message, cause, nil)
}

// WithContext adds context to an existing AppError
func WithContext(err error, context map[string]interface{}) *AppError {
	if appErr, ok := IsAppError(err); ok {
		for k, v := range context {
			if appErr.Context == nil {
				appErr.Context = make(map[string]interface{})
			}
			appErr.Context[k] = v
		}
		return appErr
	}
	return New(ErrorTypeUnknown, "unknown error with context", err, context)
}

// IsType checks if an error is of a specific type
func IsType(err error, errType ErrorType) bool {
	if appErr, ok := IsAppError(err); ok {
		return appErr.Type == errType
	}
	return false
}

func IsNotFound(err error) bool {
	return IsType(err, ErrorTypeNotFound)
}

func IsValidation(err error) bool {
	return IsType(err, ErrorTypeValidation)
}

func IsInternal(err error) bool {
	return IsType(err, ErrorTypeInternal)
}

func IsDuplicate(err error) bool {
	return IsType(err, ErrorTypeDuplicate)
}

func IsUnauthorized(err error) bool {
	return IsType(err, ErrorTypeUnauthorized)
}

func IsForbidden(err error) bool {
	return IsType(err, ErrorTypeForbidden)
}

func IsConflict(err error) bool {
	return IsType(err, ErrorTypeConflict)
}

func IsUnknown(err error) bool {
	return IsType(err, ErrorTypeUnknown)
}
