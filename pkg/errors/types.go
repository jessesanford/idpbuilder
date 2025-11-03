package errors

import (
	"fmt"
)

// BaseError is the base type for all custom errors.
// It provides common fields for error message, suggestion, and cause tracking.
type BaseError struct {
	Message    string
	Suggestion string
	Cause      error
}

// Error returns the error message formatted with optional suggestion.
func (e *BaseError) Error() string {
	if e.Suggestion != "" {
		return fmt.Sprintf("Error: %s\nSuggestion: %s", e.Message, e.Suggestion)
	}
	return fmt.Sprintf("Error: %s", e.Message)
}

// Unwrap returns the underlying cause of this error for error chain traversal.
func (e *BaseError) Unwrap() error {
	return e.Cause
}

// ValidationError represents input validation failures (exit code 1).
// This error type is used for invalid image names, registry URLs, or credentials.
type ValidationError struct {
	BaseError
	Field    string
	ExitCode int
}

// NewValidationError creates a new validation error with field context.
//
// Example:
//
//	err := NewValidationError("imageName", "invalid image format", "use format: name:tag")
func NewValidationError(field, message, suggestion string) *ValidationError {
	return &ValidationError{
		BaseError: BaseError{
			Message:    message,
			Suggestion: suggestion,
		},
		Field:    field,
		ExitCode: 1,
	}
}

// AuthenticationError represents authentication failures (exit code 2).
// This error type is used when registry authentication fails.
type AuthenticationError struct {
	BaseError
	Registry string
	ExitCode int
}

// NewAuthenticationError creates a new authentication error with registry context.
//
// Example:
//
//	err := NewAuthenticationError("docker.io", "authentication failed", "check credentials")
func NewAuthenticationError(registry, message, suggestion string) *AuthenticationError {
	return &AuthenticationError{
		BaseError: BaseError{
			Message:    message,
			Suggestion: suggestion,
		},
		Registry: registry,
		ExitCode: 2,
	}
}

// NetworkError represents network/connectivity failures (exit code 3).
// This error type is used for connection timeouts, TLS errors, or network failures.
type NetworkError struct {
	BaseError
	Target   string
	ExitCode int
}

// NewNetworkError creates a new network error with target context.
//
// Example:
//
//	err := NewNetworkError("registry.example.com", "connection timeout", "check network")
func NewNetworkError(target, message, suggestion string) *NetworkError {
	return &NetworkError{
		BaseError: BaseError{
			Message:    message,
			Suggestion: suggestion,
		},
		Target:   target,
		ExitCode: 3,
	}
}

// ImageNotFoundError represents missing image errors (exit code 4).
// This error type is used when an image cannot be found in Docker daemon or registry.
type ImageNotFoundError struct {
	BaseError
	ImageName string
	ExitCode  int
}

// NewImageNotFoundError creates a new image not found error with image name context.
//
// Example:
//
//	err := NewImageNotFoundError("alpine:latest", "image not found", "pull image first")
func NewImageNotFoundError(imageName, message, suggestion string) *ImageNotFoundError {
	return &ImageNotFoundError{
		BaseError: BaseError{
			Message:    message,
			Suggestion: suggestion,
		},
		ImageName: imageName,
		ExitCode:  4,
	}
}

// SSRFWarning represents a potential SSRF risk (warning, not error).
// This warning is issued for private IP ranges or localhost targets but does not block execution.
type SSRFWarning struct {
	Target     string
	Message    string
	Suggestion string
}

// Error returns the warning message formatted.
func (w *SSRFWarning) Error() string {
	return fmt.Sprintf("Warning: %s\nSuggestion: %s", w.Message, w.Suggestion)
}

// SecurityWarning represents security concerns (warning, not error).
// This warning is issued for weak credentials or insecure configurations but does not block execution.
type SecurityWarning struct {
	Message    string
	Suggestion string
}

// Error returns the warning message formatted.
func (w *SecurityWarning) Error() string {
	return fmt.Sprintf("Warning: %s\nSuggestion: %s", w.Message, w.Suggestion)
}

// IsWarning returns true if the error is a warning (should not stop execution).
// Warnings include SSRF and security warnings that inform users of potential risks.
func IsWarning(err error) bool {
	_, isSSRF := err.(*SSRFWarning)
	_, isSecurity := err.(*SecurityWarning)
	return isSSRF || isSecurity
}
