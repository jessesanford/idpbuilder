package errors

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// ErrorCategory represents the category of an OCI error
type ErrorCategory int

const (
	// ErrorCategoryBuild represents build-related errors
	ErrorCategoryBuild ErrorCategory = iota + 1
	// ErrorCategoryRegistry represents registry-related errors
	ErrorCategoryRegistry
	// ErrorCategoryConfiguration represents configuration errors
	ErrorCategoryConfiguration
	// ErrorCategoryStack represents stack-related errors
	ErrorCategoryStack
	// ErrorCategoryAuthentication represents authentication errors
	ErrorCategoryAuthentication
	// ErrorCategorySystem represents system-level errors
	ErrorCategorySystem
)

// String returns the string representation of ErrorCategory
func (ec ErrorCategory) String() string {
	switch ec {
	case ErrorCategoryBuild:
		return "Build"
	case ErrorCategoryRegistry:
		return "Registry"
	case ErrorCategoryConfiguration:
		return "Configuration"
	case ErrorCategoryStack:
		return "Stack"
	case ErrorCategoryAuthentication:
		return "Authentication"
	case ErrorCategorySystem:
		return "System"
	default:
		return "Unknown"
	}
}

// OCIError represents a structured error for OCI operations
type OCIError struct {
	// Code is the specific error code
	Code string
	// Message is the human-readable error message
	Message string
	// Component is the component that generated the error
	Component string
	// Operation is the operation that failed
	Operation string
	// Category is the category of the error
	Category ErrorCategory
	// Details contains additional context information
	Details map[string]interface{}
	// Cause is the underlying error that caused this error
	Cause error
	// RequestID is the request ID for tracing
	RequestID string
	// Timestamp is when the error occurred
	Timestamp time.Time
	// RetryAfter indicates when the operation can be retried
	RetryAfter *time.Duration
}

// NewOCIError creates a new OCIError with the given parameters
func NewOCIError(code, component, operation, message string) *OCIError {
	return &OCIError{
		Code:      code,
		Component: component,
		Operation: operation,
		Message:   message,
		Category:  GetCategoryFromCode(code),
		Details:   make(map[string]interface{}),
		Timestamp: time.Now(),
	}
}

// Error implements the error interface
func (e *OCIError) Error() string {
	var parts []string
	parts = append(parts, fmt.Sprintf("[%s]", e.Code))
	
	if e.Component != "" && e.Operation != "" {
		parts = append(parts, fmt.Sprintf("%s.%s", e.Component, e.Operation))
	} else if e.Component != "" {
		parts = append(parts, e.Component)
	}
	
	parts = append(parts, e.Message)
	
	if e.Cause != nil {
		parts = append(parts, fmt.Sprintf("caused by: %v", e.Cause))
	}
	
	return strings.Join(parts, " - ")
}

// Unwrap returns the underlying error for error unwrapping
func (e *OCIError) Unwrap() error {
	return e.Cause
}

// Is implements error comparison for errors.Is()
func (e *OCIError) Is(target error) bool {
	if target == nil {
		return false
	}
	
	var ociErr *OCIError
	if errors.As(target, &ociErr) {
		return e.Code == ociErr.Code
	}
	
	return false
}

// WithCause adds a cause to the error
func (e *OCIError) WithCause(cause error) *OCIError {
	e.Cause = cause
	return e
}

// WithDetails adds details to the error
func (e *OCIError) WithDetails(details map[string]interface{}) *OCIError {
	for k, v := range details {
		e.Details[k] = v
	}
	return e
}

// WithDetail adds a single detail to the error
func (e *OCIError) WithDetail(key string, value interface{}) *OCIError {
	e.Details[key] = value
	return e
}

// WithRequestID sets the request ID for the error
func (e *OCIError) WithRequestID(requestID string) *OCIError {
	e.RequestID = requestID
	return e
}

// WithRetryAfter sets the retry after duration
func (e *OCIError) WithRetryAfter(duration time.Duration) *OCIError {
	e.RetryAfter = &duration
	return e
}

// IsRetryable returns whether the error is retryable based on its code
func (e *OCIError) IsRetryable() bool {
	return IsRetryable(e.Code)
}

// GetCategory returns the category of the error
func (e *OCIError) GetCategory() ErrorCategory {
	return e.Category
}

// GetCategoryFromCode returns the error category based on the error code
func GetCategoryFromCode(code string) ErrorCategory {
	if len(code) < 4 {
		return ErrorCategorySystem
	}
	
	prefix := code[:4]
	switch prefix {
	case "1000", "1001", "1002", "1003", "1004", "1005":
		return ErrorCategoryBuild
	case "2000", "2001", "2002", "2003", "2004", "2005":
		return ErrorCategoryRegistry
	case "3000", "3001", "3002", "3003", "3004":
		return ErrorCategoryConfiguration
	case "4000", "4001", "4002", "4003", "4004":
		return ErrorCategoryStack
	case "5000", "5001", "5002", "5003", "5004":
		return ErrorCategoryAuthentication
	case "6000", "6001", "6002", "6003", "6004":
		return ErrorCategorySystem
	default:
		return ErrorCategorySystem
	}
}

// IsRetryable returns whether an error code represents a retryable error
func IsRetryable(code string) bool {
	retryableCodes := map[string]bool{
		// Registry errors that are retryable
		"2001": true, // Connection timeout
		"2004": true, // Rate limited
		"2005": true, // Temporary failure
		
		// System errors that are retryable
		"6001": true, // Disk full (might be cleared)
		"6003": true, // Network issues
		"6004": true, // Resource exhaustion
		
		// Build errors that might be retryable
		"1005": true, // Temporary build failure
	}
	
	return retryableCodes[code]
}