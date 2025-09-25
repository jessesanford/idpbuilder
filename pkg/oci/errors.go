package oci

import (
	"fmt"
	"net/http"
)

// RegistryError represents an error that occurred during registry operations.
// It includes both HTTP status codes and descriptive error messages for
// better error handling and debugging.
type RegistryError struct {
	// Code is the HTTP status code or custom error code
	Code int `json:"code"`

	// Message is the human-readable error message
	Message string `json:"message"`

	// Operation is the operation that failed
	Operation string `json:"operation,omitempty"`
}

// Error implements the error interface, returning a formatted error message
// that includes both the error code and descriptive message.
func (e *RegistryError) Error() string {
	if e.Operation != "" {
		return fmt.Sprintf("registry error %d in %s: %s", e.Code, e.Operation, e.Message)
	}
	return fmt.Sprintf("registry error %d: %s", e.Code, e.Message)
}

// NewRegistryError creates a new RegistryError with the specified code and message.
// This constructor ensures consistent error creation throughout the OCI client.
func NewRegistryError(code int, msg string) *RegistryError {
	return &RegistryError{
		Code:    code,
		Message: msg,
	}
}

// NewRegistryErrorWithOperation creates a RegistryError with operation context.
func NewRegistryErrorWithOperation(code int, msg, operation string) *RegistryError {
	return &RegistryError{
		Code:      code,
		Message:   msg,
		Operation: operation,
	}
}

// IsUnauthorized returns true if the error represents an authentication failure.
func (e *RegistryError) IsUnauthorized() bool {
	return e.Code == http.StatusUnauthorized
}