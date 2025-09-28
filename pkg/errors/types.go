package errors

import (
	"fmt"
	"time"
)

// ErrorCode represents a categorized error type
type ErrorCode string

const (
	// Build-related errors
	BuildFailed  ErrorCode = "BUILD_FAILED"
	BuildTimeout ErrorCode = "BUILD_TIMEOUT"

	// Registry errors
	RegistryUnreachable ErrorCode = "REGISTRY_UNREACHABLE"
	RegistryAuthFailed  ErrorCode = "REGISTRY_AUTH_FAILED"
	RegistryPushFailed  ErrorCode = "REGISTRY_PUSH_FAILED"

	// Certificate errors
	CertificateInvalid ErrorCode = "CERTIFICATE_INVALID"
	CertificateExpired ErrorCode = "CERTIFICATE_EXPIRED"

	// General errors
	ValidationFailed   ErrorCode = "VALIDATION_FAILED"
	ConfigurationError ErrorCode = "CONFIGURATION_ERROR"
	InternalError      ErrorCode = "INTERNAL_ERROR"
)

// StructuredError provides rich error context
type StructuredError struct {
	Code      ErrorCode
	Op        string    // Operation that failed
	Message   string    // Human-readable message
	Cause     error     // Underlying error
	Timestamp time.Time // When error occurred
}

// Error implements the error interface
func (e *StructuredError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %s (caused by: %v)", e.Code, e.Op, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s: %s", e.Code, e.Op, e.Message)
}

// Unwrap implements error unwrapping for errors.Is and errors.As
func (e *StructuredError) Unwrap() error {
	return e.Cause
}

// NewStructuredError creates a new structured error
func NewStructuredError(code ErrorCode, op, message string, cause error) *StructuredError {
	return &StructuredError{
		Code:      code,
		Op:        op,
		Message:   message,
		Cause:     cause,
		Timestamp: time.Now(),
	}
}

// IsRetryable determines if an error code represents a retryable condition
func (code ErrorCode) IsRetryable() bool {
	switch code {
	case RegistryUnreachable, BuildTimeout:
		return true
	default:
		return false
	}
}

// IsRecoverable determines if an error code represents a recoverable condition
func (code ErrorCode) IsRecoverable() bool {
	switch code {
	case ConfigurationError, ValidationFailed:
		return true
	default:
		return false
	}
}