// Package errorhandling provides custom error types and utilities for consistent
// error handling throughout the PokédexCLI application.
package errorhandling

import (
	"fmt"
	"net/http"
)

// ErrorType represents categories of errors that can occur in the application.
type ErrorType string

// Predefined error types for the application.
const (
	NotFound            ErrorType = "NOT_FOUND"
	InvalidInput        ErrorType = "INVALID_INPUT"
	NetworkError        ErrorType = "NETWORK_ERROR"
	ResourceUnavailable ErrorType = "RESOURCE_UNAVAILABLE"
	InternalError       ErrorType = "INTERNAL_ERROR"
)

// AppError is a custom error type that provides context about errors.
// It includes error type classification, status code for API errors,
// and a descriptive message with context about what went wrong.
type AppError struct {
	Type       ErrorType   // The category of error
	StatusCode int         // HTTP status code (for API errors)
	Message    string      // Human-readable error message
	Err        error       // The original error (optional)
	Context    interface{} // Additional context (optional)
}

// Error returns a string representation of the error.
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%s)", e.Type, e.Message, e.Err.Error())
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap returns the underlying error.
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewNotFoundError creates a new error for when a resource is not found.
func NewNotFoundError(resourceType, resourceName string, err error) *AppError {
	return &AppError{
		Type:       NotFound,
		StatusCode: http.StatusNotFound,
		Message:    fmt.Sprintf("The %s '%s' was not found", resourceType, resourceName),
		Err:        err,
		Context: map[string]string{
			"resourceType": resourceType,
			"resourceName": resourceName,
		},
	}
}

// NewInvalidInputError creates a new error for invalid user input.
func NewInvalidInputError(message string, err error) *AppError {
	return &AppError{
		Type:       InvalidInput,
		StatusCode: http.StatusBadRequest,
		Message:    message,
		Err:        err,
	}
}

// NewNetworkError creates a new error for network-related issues.
func NewNetworkError(message string, err error) *AppError {
	return &AppError{
		Type:       NetworkError,
		StatusCode: http.StatusServiceUnavailable,
		Message:    message,
		Err:        err,
	}
}

// NewAPIError creates a new error based on an HTTP response status code.
func NewAPIError(statusCode int, endpoint string, err error) *AppError {
	var errType ErrorType
	var message string

	switch statusCode {
	case http.StatusNotFound:
		errType = NotFound
		message = fmt.Sprintf("The requested resource at '%s' was not found", endpoint)
	case http.StatusBadRequest:
		errType = InvalidInput
		message = fmt.Sprintf("Bad request to '%s'", endpoint)
	case http.StatusTooManyRequests:
		errType = ResourceUnavailable
		message = "Rate limit exceeded. Please try again later"
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable:
		errType = ResourceUnavailable
		message = "The Pokémon API service is temporarily unavailable"
	default:
		errType = InternalError
		message = fmt.Sprintf("Unexpected API error (status code: %d)", statusCode)
	}

	return &AppError{
		Type:       errType,
		StatusCode: statusCode,
		Message:    message,
		Err:        err,
		Context: map[string]string{
			"endpoint": endpoint,
		},
	}
}

// NewInternalError creates a new error for internal application errors.
func NewInternalError(message string, err error) *AppError {
	return &AppError{
		Type:       InternalError,
		StatusCode: http.StatusInternalServerError,
		Message:    message,
		Err:        err,
	}
}

// IsNotFoundError checks if an error is a NotFound error.
func IsNotFoundError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Type == NotFound
	}
	return false
}

// IsInvalidInputError checks if an error is an InvalidInput error.
func IsInvalidInputError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Type == InvalidInput
	}
	return false
}

// FormatUserMessage formats an error for display to the user.
// Removes technical details and provides a user-friendly message.
func FormatUserMessage(err error) string {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Message
	}
	return err.Error()
}
