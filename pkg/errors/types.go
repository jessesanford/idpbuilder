package errors

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"
)

// ErrorCategory defines the category of errors for proper handling
type ErrorCategory int

const (
	// ErrorCategoryUnknown represents an unclassified error
	ErrorCategoryUnknown ErrorCategory = iota

	// ErrorCategoryTransient represents errors that may succeed on retry
	ErrorCategoryTransient

	// ErrorCategoryPermanent represents errors that will not succeed on retry
	ErrorCategoryPermanent

	// ErrorCategoryAuth represents authentication/authorization errors
	ErrorCategoryAuth

	// ErrorCategoryNetwork represents network connectivity errors
	ErrorCategoryNetwork

	// ErrorCategoryFormat represents data format/validation errors
	ErrorCategoryFormat

	// ErrorCategoryQuota represents resource quota/limit errors
	ErrorCategoryQuota
)

// String returns a string representation of the error category
func (c ErrorCategory) String() string {
	switch c {
	case ErrorCategoryTransient:
		return "transient"
	case ErrorCategoryPermanent:
		return "permanent"
	case ErrorCategoryAuth:
		return "authentication"
	case ErrorCategoryNetwork:
		return "network"
	case ErrorCategoryFormat:
		return "format"
	case ErrorCategoryQuota:
		return "quota"
	default:
		return "unknown"
	}
}

// ShouldRetry returns whether errors in this category should be retried
func (c ErrorCategory) ShouldRetry() bool {
	switch c {
	case ErrorCategoryTransient, ErrorCategoryNetwork:
		return true
	default:
		return false
	}
}

// OperationError wraps errors with operation context and retry information
type OperationError struct {
	// Operation is the name of the operation that failed
	Operation string

	// Resource is the resource being operated on (e.g., registry/repo:tag)
	Resource string

	// Category categorizes the error for handling decisions
	Category ErrorCategory

	// Retryable indicates if this specific error instance should be retried
	Retryable bool

	// Attempt tracks which attempt this error occurred on
	Attempt int

	// Timestamp records when the error occurred
	Timestamp time.Time

	// Cause is the underlying error
	Cause error

	// Context provides additional metadata about the error
	Context map[string]interface{}
}

// Error implements the error interface
func (e *OperationError) Error() string {
	if e.Resource != "" {
		return fmt.Sprintf("%s operation failed for %s (attempt %d, %s): %v",
			e.Operation, e.Resource, e.Attempt, e.Category, e.Cause)
	}
	return fmt.Sprintf("%s operation failed (attempt %d, %s): %v",
		e.Operation, e.Attempt, e.Category, e.Cause)
}

// Unwrap returns the underlying error for errors.Is and errors.As support
func (e *OperationError) Unwrap() error {
	return e.Cause
}

// Is implements error comparison for errors.Is
func (e *OperationError) Is(target error) bool {
	if target == nil {
		return false
	}

	if other, ok := target.(*OperationError); ok {
		return e.Operation == other.Operation &&
			e.Category == other.Category &&
			errors.Is(e.Cause, other.Cause)
	}

	return errors.Is(e.Cause, target)
}

// RetryInfo contains information about retry configuration and state
type RetryInfo struct {
	// MaxAttempts is the maximum number of retry attempts
	MaxAttempts int

	// BaseDelay is the base delay between retries
	BaseDelay time.Duration

	// MaxDelay is the maximum delay between retries
	MaxDelay time.Duration

	// Multiplier is the backoff multiplier for each retry
	Multiplier float64

	// Jitter adds randomness to retry delays to avoid thundering herd
	Jitter bool
}

// DefaultRetryInfo returns sensible default retry configuration
func DefaultRetryInfo() RetryInfo {
	return RetryInfo{
		MaxAttempts: 3,
		BaseDelay:   100 * time.Millisecond,
		MaxDelay:    10 * time.Second,
		Multiplier:  2.0,
		Jitter:      true,
	}
}

// ErrorHandler provides a standardized way to handle and classify errors
type ErrorHandler interface {
	// ClassifyError categorizes an error for appropriate handling
	ClassifyError(err error) ErrorCategory

	// ShouldRetry determines if an operation should be retried based on the error
	ShouldRetry(err error, attempt int) bool

	// WrapError wraps an error with operation context
	WrapError(operation, resource string, err error, attempt int) *OperationError

	// HandleError processes an error and returns appropriate action
	HandleError(ctx context.Context, err error) error
}

// Registry of common error types for classification
var (
	// Common network errors
	ErrConnectionRefused = errors.New("connection refused")
	ErrTimeout           = errors.New("operation timeout")
	ErrDNSFailure       = errors.New("DNS resolution failed")

	// Common authentication errors
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrInvalidToken     = errors.New("invalid token")

	// Common format errors
	ErrInvalidManifest  = errors.New("invalid manifest")
	ErrInvalidDigest    = errors.New("invalid digest")
	ErrUnsupportedType  = errors.New("unsupported media type")

	// Common quota errors
	ErrQuotaExceeded    = errors.New("quota exceeded")
	ErrRateLimit        = errors.New("rate limit exceeded")
	ErrStorageFull      = errors.New("storage full")
)

// IsTemporary checks if an error is temporary and may succeed on retry
func IsTemporary(err error) bool {
	if err == nil {
		return false
	}

	// Check for context errors first (these are never temporary)
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return false
	}

	// Check for net.Error interface
	if netErr, ok := err.(net.Error); ok {
		return netErr.Temporary() || netErr.Timeout()
	}

	// Check for known temporary errors
	switch {
	case errors.Is(err, ErrConnectionRefused):
		return true
	case errors.Is(err, ErrTimeout):
		return true
	case errors.Is(err, ErrRateLimit):
		return true
	default:
		return false
	}
}

// IsRetryable determines if an error should trigger a retry
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	// Check if it's an OperationError with retry information
	if opErr, ok := err.(*OperationError); ok {
		return opErr.Retryable && opErr.Category.ShouldRetry()
	}

	// Fall back to temporary check
	return IsTemporary(err)
}