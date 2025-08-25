package errors

import (
	"fmt"
	"strings"
	"time"
)

type ErrorCategory string

const (
	CategoryBuild          ErrorCategory = "build"
	CategoryRegistry       ErrorCategory = "registry"
	CategoryConfiguration  ErrorCategory = "configuration"
	CategoryStack          ErrorCategory = "stack"
	CategoryAuthentication ErrorCategory = "authentication"
	CategorySystem         ErrorCategory = "system"
)

type OCIError struct {
	Code       string                 `json:"code"`
	Message    string                 `json:"message"`
	Component  string                 `json:"component,omitempty"`
	Operation  string                 `json:"operation,omitempty"`
	Timestamp  time.Time              `json:"timestamp"`
	Cause      error                  `json:"cause,omitempty"`
	Details    map[string]interface{} `json:"details,omitempty"`
	RequestID  string                 `json:"request_id,omitempty"`
	Retry      bool                   `json:"retry"`
	RetryAfter time.Duration          `json:"retry_after,omitempty"`
}

func (e *OCIError) Error() string {
	var parts []string
	if e.Code != "" {
		parts = append(parts, fmt.Sprintf("code=%s", e.Code))
	}
	if e.Component != "" {
		parts = append(parts, fmt.Sprintf("component=%s", e.Component))
	}
	if e.Operation != "" {
		parts = append(parts, fmt.Sprintf("operation=%s", e.Operation))
	}
	
	prefix := ""
	if len(parts) > 0 {
		prefix = fmt.Sprintf("[%s] ", strings.Join(parts, ","))
	}
	
	message := e.Message
	if message == "" {
		message = "OCI operation failed"
	}
	return fmt.Sprintf("%s%s", prefix, message)
}

func (e *OCIError) Is(target error) bool {
	if t, ok := target.(*OCIError); ok {
		return e.Code == t.Code
	}
	return false
}

func (e *OCIError) Unwrap() error {
	return e.Cause
}

func (e *OCIError) WithCause(cause error) *OCIError {
	e.Cause = cause
	return e
}

func (e *OCIError) WithDetails(details map[string]interface{}) *OCIError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	for k, v := range details {
		e.Details[k] = v
	}
	return e
}

func (e *OCIError) WithRequestID(requestID string) *OCIError {
	e.RequestID = requestID
	return e
}

func (e *OCIError) WithRetryAfter(duration time.Duration) *OCIError {
	e.Retry = true
	e.RetryAfter = duration
	return e
}

func NewOCIError(code, message string) *OCIError {
	return &OCIError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
		Details:   make(map[string]interface{}),
	}
}

func NewOCIErrorWithComponent(code, message, component, operation string) *OCIError {
	return &OCIError{
		Code:      code,
		Message:   message,
		Component: component,
		Operation: operation,
		Timestamp: time.Now(),
		Details:   make(map[string]interface{}),
	}
}

func GetCategory(code string) ErrorCategory {
	if code == "" {
		return CategorySystem
	}
	switch {
	case strings.HasPrefix(code, "1"):
		return CategoryBuild
	case strings.HasPrefix(code, "2"):
		return CategoryRegistry
	case strings.HasPrefix(code, "3"):
		return CategoryConfiguration
	case strings.HasPrefix(code, "4"):
		return CategoryStack
	case strings.HasPrefix(code, "5"):
		return CategoryAuthentication
	case strings.HasPrefix(code, "6"):
		return CategorySystem
	default:
		return CategorySystem
	}
}

func IsRetryable(err error) bool {
	if ociErr, ok := err.(*OCIError); ok {
		return ociErr.Retry
	}
	return false
}
