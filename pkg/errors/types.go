package errors

import (
	"fmt"
	"time"
)

// OCIError defines the standard interface for all OCI management errors.
// It extends the standard error interface with additional metadata for
// programmatic handling and enhanced error reporting.
type OCIError interface {
	// Error returns the error message (standard error interface)
	Error() string

	// Code returns the specific error code for programmatic handling
	Code() ErrorCode

	// Category returns the error category for handling strategy
	Category() ErrorCategory

	// Wrap wraps another error with additional context
	Wrap(error) OCIError

	// Unwrap returns the underlying error (Go 1.13+ compatible)
	Unwrap() error

	// Context returns additional error context data
	Context() map[string]interface{}

	// Timestamp returns when the error occurred
	Timestamp() time.Time

	// String returns a detailed string representation
	String() string
}

// BaseError provides a concrete implementation of the OCIError interface.
// It serves as the foundation for all OCI management errors with support
// for error wrapping, context, and structured error information.
type BaseError struct {
	// ErrorCode identifies the specific error type
	ErrorCode ErrorCode `json:"code"`

	// ErrorCategory indicates the error handling strategy
	ErrorCategory ErrorCategory `json:"category"`

	// Message provides a human-readable error description
	Message string `json:"message"`

	// Cause holds the underlying error that caused this error
	Cause error `json:"cause,omitempty"`

	// ErrorContext contains additional structured error information
	ErrorContext map[string]interface{} `json:"context,omitempty"`

	// ErrorTimestamp records when the error occurred
	ErrorTimestamp time.Time `json:"timestamp"`
}

// Error returns the error message, implementing the standard error interface.
func (e *BaseError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Code returns the error code for programmatic handling.
func (e *BaseError) Code() ErrorCode {
	return e.ErrorCode
}

// Category returns the error category for handling strategy.
func (e *BaseError) Category() ErrorCategory {
	return e.ErrorCategory
}

// Wrap wraps another error with additional context, creating a new BaseError.
func (e *BaseError) Wrap(err error) OCIError {
	return &BaseError{
		ErrorCode:      e.ErrorCode,
		ErrorCategory:  e.ErrorCategory,
		Message:        e.Message,
		Cause:          err,
		ErrorContext:   e.ErrorContext,
		ErrorTimestamp: time.Now().UTC(),
	}
}

// Unwrap returns the underlying error, supporting Go 1.13+ error unwrapping.
func (e *BaseError) Unwrap() error {
	return e.Cause
}

// Context returns the additional error context data.
func (e *BaseError) Context() map[string]interface{} {
	return e.ErrorContext
}

// Timestamp returns when the error occurred.
func (e *BaseError) Timestamp() time.Time {
	return e.ErrorTimestamp
}

// String returns a detailed string representation of the error.
func (e *BaseError) String() string {
	return fmt.Sprintf("OCIError{code=%d, category=%s, message=%q, timestamp=%s}",
		e.ErrorCode, e.ErrorCategory, e.Message, e.ErrorTimestamp.Format(time.RFC3339))
}

// ErrorContext provides structured context information for errors.
type ErrorContext struct {
	// Operation describes what operation was being performed
	Operation string `json:"operation,omitempty"`

	// Component identifies which component generated the error
	Component string `json:"component,omitempty"`

	// ResourceType specifies the type of resource involved
	ResourceType string `json:"resource_type,omitempty"`

	// ResourceID identifies the specific resource
	ResourceID string `json:"resource_id,omitempty"`

	// Registry identifies the registry being accessed
	Registry string `json:"registry,omitempty"`

	// Repository identifies the repository being accessed
	Repository string `json:"repository,omitempty"`

	// Tag identifies the image tag involved
	Tag string `json:"tag,omitempty"`

	// Additional allows for arbitrary additional context
	Additional map[string]interface{} `json:"additional,omitempty"`
}

// ErrorStack represents a chain of related errors, useful for tracking
// error propagation through different system layers.
type ErrorStack struct {
	// Errors contains the chain of errors from most recent to root cause
	Errors []OCIError `json:"errors"`

	// MaxDepth limits how deep the error stack can grow
	MaxDepth int `json:"max_depth"`
}

// NewErrorStack creates a new error stack with the given initial error.
func NewErrorStack(err OCIError, maxDepth int) *ErrorStack {
	return &ErrorStack{
		Errors:   []OCIError{err},
		MaxDepth: maxDepth,
	}
}

// Push adds a new error to the top of the stack.
func (es *ErrorStack) Push(err OCIError) {
	if len(es.Errors) >= es.MaxDepth {
		// Remove oldest error (root) to maintain max depth
		es.Errors = es.Errors[:len(es.Errors)-1]
	}
	es.Errors = append([]OCIError{err}, es.Errors...)
}

// Root returns the root cause error (last in the chain).
func (es *ErrorStack) Root() OCIError {
	if len(es.Errors) == 0 {
		return nil
	}
	return es.Errors[len(es.Errors)-1]
}

// Latest returns the most recent error (first in the chain).
func (es *ErrorStack) Latest() OCIError {
	if len(es.Errors) == 0 {
		return nil
	}
	return es.Errors[0]
}

// Depth returns the current depth of the error stack.
func (es *ErrorStack) Depth() int {
	return len(es.Errors)
}

// Error implements the error interface for the entire stack.
func (es *ErrorStack) Error() string {
	if len(es.Errors) == 0 {
		return "empty error stack"
	}
	return es.Latest().Error()
}