package errors

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"syscall"
	"time"
)

// StandardErrorHandler implements the ErrorHandler interface with standard classification rules
type StandardErrorHandler struct {
	retryHandler *RetryHandler
	context      map[string]interface{}
}

// NewStandardErrorHandler creates a new standard error handler
func NewStandardErrorHandler() *StandardErrorHandler {
	return &StandardErrorHandler{
		retryHandler: NewDefaultRetryHandler(),
		context:      make(map[string]interface{}),
	}
}

// NewStandardErrorHandlerWithRetry creates a standard error handler with custom retry configuration
func NewStandardErrorHandlerWithRetry(retryConfig RetryInfo) *StandardErrorHandler {
	return &StandardErrorHandler{
		retryHandler: NewRetryHandler(retryConfig),
		context:      make(map[string]interface{}),
	}
}

// ClassifyError categorizes an error based on its type and characteristics
func (h *StandardErrorHandler) ClassifyError(err error) ErrorCategory {
	return ClassifyError(err)
}

// ClassifyError is a standalone function to classify errors
func ClassifyError(err error) ErrorCategory {
	if err == nil {
		return ErrorCategoryUnknown
	}

	// Check for HTTP status codes
	if httpErr := extractHTTPError(err); httpErr != nil {
		return classifyHTTPError(httpErr.StatusCode)
	}

	// Check for network errors
	if isNetworkError(err) {
		return ErrorCategoryNetwork
	}

	// Check for authentication errors
	if isAuthError(err) {
		return ErrorCategoryAuth
	}

	// Check for format/validation errors
	if isFormatError(err) {
		return ErrorCategoryFormat
	}

	// Check for quota/resource errors
	if isQuotaError(err) {
		return ErrorCategoryQuota
	}

	// Check for temporary errors
	if IsTemporary(err) {
		return ErrorCategoryTransient
	}

	// Check error message content for classification hints
	errMsg := strings.ToLower(err.Error())
	switch {
	case strings.Contains(errMsg, "unauthorized") || strings.Contains(errMsg, "forbidden"):
		return ErrorCategoryAuth
	case strings.Contains(errMsg, "timeout") || strings.Contains(errMsg, "deadline"):
		return ErrorCategoryTransient
	case strings.Contains(errMsg, "connection") || strings.Contains(errMsg, "network"):
		return ErrorCategoryNetwork
	case strings.Contains(errMsg, "invalid") || strings.Contains(errMsg, "malformed"):
		return ErrorCategoryFormat
	case strings.Contains(errMsg, "quota") || strings.Contains(errMsg, "limit"):
		return ErrorCategoryQuota
	default:
		return ErrorCategoryUnknown
	}
}

// ShouldRetry determines if an operation should be retried
func (h *StandardErrorHandler) ShouldRetry(err error, attempt int) bool {
	if err == nil {
		return false
	}

	category := h.ClassifyError(err)
	if !category.ShouldRetry() {
		return false
	}

	return h.retryHandler.ShouldRetry(err, attempt)
}

// WrapError wraps an error with operation context
func (h *StandardErrorHandler) WrapError(operation, resource string, err error, attempt int) *OperationError {
	if err == nil {
		return nil
	}

	category := h.ClassifyError(err)

	return &OperationError{
		Operation: operation,
		Resource:  resource,
		Category:  category,
		Retryable: category.ShouldRetry() && IsRetryable(err),
		Attempt:   attempt,
		Timestamp: time.Now(),
		Cause:     err,
		Context:   h.copyContext(),
	}
}

// HandleError processes an error and returns appropriate action
func (h *StandardErrorHandler) HandleError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	// If it's already an OperationError, return as-is
	if opErr, ok := err.(*OperationError); ok {
		return opErr
	}

	// Wrap the error with classification
	category := h.ClassifyError(err)

	wrappedErr := &OperationError{
		Operation: "unknown",
		Resource:  "",
		Category:  category,
		Retryable: category.ShouldRetry() && IsRetryable(err),
		Attempt:   1,
		Timestamp: time.Now(),
		Cause:     err,
		Context:   h.copyContext(),
	}

	return wrappedErr
}

// SetContext adds context information to the error handler
func (h *StandardErrorHandler) SetContext(key string, value interface{}) {
	h.context[key] = value
}

// GetContext retrieves context information from the error handler
func (h *StandardErrorHandler) GetContext(key string) interface{} {
	return h.context[key]
}

// copyContext creates a copy of the current context
func (h *StandardErrorHandler) copyContext() map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range h.context {
		result[k] = v
	}
	return result
}

// Helper functions for error classification

// HTTPError represents an HTTP error with status code
type HTTPError struct {
	StatusCode int
	Message    string
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Message)
}

// extractHTTPError attempts to extract HTTP error information
func extractHTTPError(err error) *HTTPError {
	// Check for direct HTTP error
	if httpErr, ok := err.(*HTTPError); ok {
		return httpErr
	}

	// Try to extract from error message
	errMsg := err.Error()
	if strings.Contains(errMsg, "HTTP ") {
		// This is a simplified extraction - in real implementation,
		// you might want more sophisticated parsing
		return &HTTPError{StatusCode: 0, Message: errMsg}
	}

	return nil
}

// classifyHTTPError categorizes errors based on HTTP status codes
func classifyHTTPError(statusCode int) ErrorCategory {
	switch {
	case statusCode >= 400 && statusCode < 500:
		// 4xx errors
		switch statusCode {
		case http.StatusUnauthorized, http.StatusForbidden:
			return ErrorCategoryAuth
		case http.StatusBadRequest, http.StatusNotAcceptable, http.StatusUnsupportedMediaType:
			return ErrorCategoryFormat
		case http.StatusTooManyRequests:
			return ErrorCategoryQuota
		default:
			return ErrorCategoryPermanent
		}
	case statusCode >= 500 && statusCode < 600:
		// 5xx errors are generally transient
		return ErrorCategoryTransient
	default:
		return ErrorCategoryUnknown
	}
}

// isNetworkError checks if an error is network-related
func isNetworkError(err error) bool {
	// Check for net.Error
	if _, ok := err.(net.Error); ok {
		return true
	}

	// Check for specific network errors
	var netOpErr *net.OpError
	if errors.As(err, &netOpErr) {
		return true
	}

	// Check for syscall errors
	var syscallErr syscall.Errno
	if errors.As(err, &syscallErr) {
		switch syscallErr {
		case syscall.ECONNREFUSED, syscall.ECONNRESET, syscall.ETIMEDOUT, syscall.EHOSTUNREACH:
			return true
		}
	}

	return false
}

// isAuthError checks if an error is authentication-related
func isAuthError(err error) bool {
	errMsg := strings.ToLower(err.Error())
	keywords := []string{"unauthorized", "forbidden", "authentication", "invalid token", "access denied"}

	for _, keyword := range keywords {
		if strings.Contains(errMsg, keyword) {
			return true
		}
	}

	return false
}

// isFormatError checks if an error is format/validation-related
func isFormatError(err error) bool {
	errMsg := strings.ToLower(err.Error())
	keywords := []string{"invalid", "malformed", "parse", "unmarshal", "decode", "format", "syntax"}

	for _, keyword := range keywords {
		if strings.Contains(errMsg, keyword) {
			return true
		}
	}

	return false
}

// isQuotaError checks if an error is quota/resource-related
func isQuotaError(err error) bool {
	errMsg := strings.ToLower(err.Error())

	// Check for specific quota-related phrases first
	quotaPhrases := []string{"quota exceeded", "limit exceeded", "limit reached", "rate limit", "storage full"}
	for _, phrase := range quotaPhrases {
		if strings.Contains(errMsg, phrase) {
			return true
		}
	}

	// Check for other quota keywords, but exclude deadline-related errors
	if strings.Contains(errMsg, "deadline") {
		return false
	}

	keywords := []string{"quota", "throttle"}
	for _, keyword := range keywords {
		if strings.Contains(errMsg, keyword) {
			return true
		}
	}

	return false
}