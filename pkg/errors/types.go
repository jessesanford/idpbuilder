package errors

import (
	"fmt"
	"time"
)

// OCIError defines the standard interface for all OCI management errors.
// This interface extends the standard error interface with additional
// contextual information and error handling capabilities.
type OCIError interface {
	error
	
	// Code returns the specific error code
	Code() ErrorCode
	
	// Context returns structured context information
	Context() ErrorContext
	
	// Severity returns the error severity level
	Severity() string
	
	// IsTransient returns true if the error may resolve with retry
	IsTransient() bool
	
	// IsPermanent returns true if the error will not resolve with retry
	IsPermanent() bool
	
	// Unwrap returns the wrapped error if any
	Unwrap() error
	
	// Stack returns the error propagation chain
	Stack() ErrorStack
	
	// WithContext adds context information to the error
	WithContext(key, value string) OCIError
	
	// WithOperation sets the operation context
	WithOperation(operation string) OCIError
	
	// WithResource sets the resource context
	WithResource(resource string) OCIError
}

// ErrorContext provides structured context information for errors
type ErrorContext struct {
	// Operation identifies the operation that failed
	Operation string `json:"operation,omitempty"`
	
	// Resource identifies the resource involved
	Resource string `json:"resource,omitempty"`
	
	// Namespace identifies the Kubernetes namespace
	Namespace string `json:"namespace,omitempty"`
	
	// Registry identifies the OCI registry
	Registry string `json:"registry,omitempty"`
	
	// Image identifies the OCI image
	Image string `json:"image,omitempty"`
	
	// Tag identifies the image tag
	Tag string `json:"tag,omitempty"`
	
	// Repository identifies the repository
	Repository string `json:"repository,omitempty"`
	
	// URL identifies a URL involved in the operation
	URL string `json:"url,omitempty"`
	
	// Path identifies a file or directory path
	Path string `json:"path,omitempty"`
	
	// Timeout identifies timeout duration
	Timeout time.Duration `json:"timeout,omitempty"`
	
	// Custom holds additional context key-value pairs
	Custom map[string]string `json:"custom,omitempty"`
	
	// Timestamp records when the error occurred
	Timestamp time.Time `json:"timestamp"`
}

// ErrorStack represents a chain of related errors for tracking propagation
type ErrorStack struct {
	// Errors contains the chain of errors from root cause to current
	Errors []StackFrame `json:"errors"`
	
	// Depth indicates the current depth in the error chain
	Depth int `json:"depth"`
}

// StackFrame represents a single frame in the error stack
type StackFrame struct {
	// Error is the error at this frame
	Error string `json:"error"`
	
	// Code is the error code at this frame
	Code ErrorCode `json:"code"`
	
	// Context is the error context at this frame
	Context ErrorContext `json:"context"`
	
	// Timestamp records when this frame was added
	Timestamp time.Time `json:"timestamp"`
}

// BaseError provides a concrete implementation of the OCIError interface
type BaseError struct {
	// code identifies the specific error condition
	code ErrorCode
	
	// message is the human-readable error message
	message string
	
	// context provides structured error information
	context ErrorContext
	
	// wrapped holds any wrapped error
	wrapped error
	
	// stack tracks error propagation
	stack ErrorStack
	
	// timestamp records when the error was created
	timestamp time.Time
}

// New creates a new OCIError with the specified code and message
func New(code ErrorCode, message string) OCIError {
	return &BaseError{
		code:      code,
		message:   message,
		context:   ErrorContext{Timestamp: time.Now()},
		stack:     ErrorStack{Errors: []StackFrame{}, Depth: 0},
		timestamp: time.Now(),
	}
}

// Newf creates a new OCIError with the specified code and formatted message
func Newf(code ErrorCode, format string, args ...interface{}) OCIError {
	return New(code, fmt.Sprintf(format, args...))
}

// Wrap creates a new OCIError that wraps another error
func Wrap(code ErrorCode, message string, err error) OCIError {
	baseErr := &BaseError{
		code:      code,
		message:   message,
		context:   ErrorContext{Timestamp: time.Now()},
		wrapped:   err,
		stack:     ErrorStack{Errors: []StackFrame{}, Depth: 0},
		timestamp: time.Now(),
	}
	
	// If wrapping another OCIError, inherit its stack
	if ociErr, ok := err.(OCIError); ok {
		baseErr.stack = ociErr.Stack()
		baseErr.stack.Depth++
	}
	
	// Add current frame to stack
	frame := StackFrame{
		Error:     message,
		Code:      code,
		Context:   baseErr.context,
		Timestamp: time.Now(),
	}
	baseErr.stack.Errors = append(baseErr.stack.Errors, frame)
	
	return baseErr
}

// Wrapf creates a new OCIError that wraps another error with formatted message
func Wrapf(code ErrorCode, err error, format string, args ...interface{}) OCIError {
	return Wrap(code, fmt.Sprintf(format, args...), err)
}

// Error implements the standard error interface
func (e *BaseError) Error() string {
	if e.wrapped != nil {
		return fmt.Sprintf("%s: %v", e.message, e.wrapped)
	}
	return e.message
}

// Code returns the error code
func (e *BaseError) Code() ErrorCode {
	return e.code
}

// Context returns the error context
func (e *BaseError) Context() ErrorContext {
	return e.context
}

// Severity returns the error severity level
func (e *BaseError) Severity() string {
	return e.code.Severity
}

// IsTransient returns true if the error may resolve with retry
func (e *BaseError) IsTransient() bool {
	return IsTransient(e.code)
}

// IsPermanent returns true if the error will not resolve with retry
func (e *BaseError) IsPermanent() bool {
	return IsPermanent(e.code)
}

// Unwrap returns the wrapped error if any
func (e *BaseError) Unwrap() error {
	return e.wrapped
}

// Stack returns the error propagation chain
func (e *BaseError) Stack() ErrorStack {
	return e.stack
}

// WithContext adds context information to the error
func (e *BaseError) WithContext(key, value string) OCIError {
	if e.context.Custom == nil {
		e.context.Custom = make(map[string]string)
	}
	e.context.Custom[key] = value
	return e
}

// WithOperation sets the operation context
func (e *BaseError) WithOperation(operation string) OCIError {
	e.context.Operation = operation
	return e
}

// WithResource sets the resource context
func (e *BaseError) WithResource(resource string) OCIError {
	e.context.Resource = resource
	return e
}

// Utility functions for error handling

// IsOCIError checks if an error implements the OCIError interface
func IsOCIError(err error) bool {
	_, ok := err.(OCIError)
	return ok
}

// AsOCIError attempts to cast an error to OCIError
func AsOCIError(err error) (OCIError, bool) {
	ociErr, ok := err.(OCIError)
	return ociErr, ok
}

// GetRootCause traverses the error chain to find the root cause
func GetRootCause(err error) error {
	for err != nil {
		if wrapped := Unwrap(err); wrapped != nil {
			err = wrapped
		} else {
			break
		}
	}
	return err
}

// Unwrap extracts the wrapped error, supporting both OCIError and standard errors
func Unwrap(err error) error {
	if ociErr, ok := err.(OCIError); ok {
		return ociErr.Unwrap()
	}
	
	// Support standard library error unwrapping
	if unwrapper, ok := err.(interface{ Unwrap() error }); ok {
		return unwrapper.Unwrap()
	}
	
	return nil
}