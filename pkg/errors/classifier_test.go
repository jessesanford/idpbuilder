package errors

import (
	"context"
	"errors"
	"net"
	"net/http"
	"syscall"
	"testing"
	"time"
)

func TestNewStandardErrorHandler(t *testing.T) {
	handler := NewStandardErrorHandler()

	if handler == nil {
		t.Fatal("NewStandardErrorHandler() returned nil")
	}

	if handler.retryHandler == nil {
		t.Error("NewStandardErrorHandler() retryHandler is nil")
	}

	if handler.context == nil {
		t.Error("NewStandardErrorHandler() context is nil")
	}
}

func TestNewStandardErrorHandlerWithRetry(t *testing.T) {
	config := RetryInfo{
		MaxAttempts: 5,
		BaseDelay:   50 * time.Millisecond,
		MaxDelay:    5 * time.Second,
		Multiplier:  1.5,
		Jitter:      false,
	}

	handler := NewStandardErrorHandlerWithRetry(config)

	if handler == nil {
		t.Fatal("NewStandardErrorHandlerWithRetry() returned nil")
	}

	if handler.retryHandler.GetConfig() != config {
		t.Errorf("NewStandardErrorHandlerWithRetry() config = %v, want %v",
			handler.retryHandler.GetConfig(), config)
	}
}

func TestClassifyError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected ErrorCategory
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: ErrorCategoryUnknown,
		},
		{
			name:     "connection refused",
			err:      &net.OpError{Op: "dial", Err: syscall.ECONNREFUSED},
			expected: ErrorCategoryNetwork,
		},
		{
			name:     "timeout error",
			err:      errors.New("operation timeout"),
			expected: ErrorCategoryTransient,
		},
		{
			name:     "unauthorized error",
			err:      errors.New("unauthorized access"),
			expected: ErrorCategoryAuth,
		},
		{
			name:     "forbidden error",
			err:      errors.New("forbidden resource"),
			expected: ErrorCategoryAuth,
		},
		{
			name:     "invalid format error",
			err:      errors.New("invalid manifest format"),
			expected: ErrorCategoryFormat,
		},
		{
			name:     "quota exceeded error",
			err:      errors.New("quota exceeded for user"),
			expected: ErrorCategoryQuota,
		},
		{
			name:     "rate limit error",
			err:      errors.New("rate limit exceeded"),
			expected: ErrorCategoryQuota,
		},
		{
			name:     "network connection error",
			err:      errors.New("network connection failed"),
			expected: ErrorCategoryNetwork,
		},
		{
			name:     "deadline exceeded",
			err:      errors.New("deadline exceeded"),
			expected: ErrorCategoryTransient,
		},
		{
			name:     "unknown error",
			err:      errors.New("some unknown error"),
			expected: ErrorCategoryUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ClassifyError(tt.err); got != tt.expected {
				t.Errorf("ClassifyError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestClassifyHTTPError(t *testing.T) {
	tests := []struct {
		statusCode int
		expected   ErrorCategory
	}{
		{http.StatusBadRequest, ErrorCategoryFormat},
		{http.StatusUnauthorized, ErrorCategoryAuth},
		{http.StatusForbidden, ErrorCategoryAuth},
		{http.StatusNotAcceptable, ErrorCategoryFormat},
		{http.StatusTooManyRequests, ErrorCategoryQuota},
		{http.StatusUnsupportedMediaType, ErrorCategoryFormat},
		{http.StatusInternalServerError, ErrorCategoryTransient},
		{http.StatusBadGateway, ErrorCategoryTransient},
		{http.StatusServiceUnavailable, ErrorCategoryTransient},
		{200, ErrorCategoryUnknown},
	}

	for _, tt := range tests {
		t.Run(http.StatusText(tt.statusCode), func(t *testing.T) {
			if got := classifyHTTPError(tt.statusCode); got != tt.expected {
				t.Errorf("classifyHTTPError(%d) = %v, want %v", tt.statusCode, got, tt.expected)
			}
		})
	}
}

func TestIsNetworkError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "net.OpError",
			err:      &net.OpError{Op: "dial", Err: syscall.ECONNREFUSED},
			expected: true,
		},
		{
			name:     "syscall ECONNREFUSED",
			err:      syscall.ECONNREFUSED,
			expected: true,
		},
		{
			name:     "syscall ECONNRESET",
			err:      syscall.ECONNRESET,
			expected: true,
		},
		{
			name:     "syscall ETIMEDOUT",
			err:      syscall.ETIMEDOUT,
			expected: true,
		},
		{
			name:     "syscall EHOSTUNREACH",
			err:      syscall.EHOSTUNREACH,
			expected: true,
		},
		{
			name:     "regular error",
			err:      errors.New("not a network error"),
			expected: false,
		},
		{
			name:     "mock net error",
			err:      mockNetError{message: "network error"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNetworkError(tt.err); got != tt.expected {
				t.Errorf("isNetworkError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsAuthError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "unauthorized",
			err:      errors.New("unauthorized access"),
			expected: true,
		},
		{
			name:     "forbidden",
			err:      errors.New("forbidden resource"),
			expected: true,
		},
		{
			name:     "authentication failed",
			err:      errors.New("authentication failed"),
			expected: true,
		},
		{
			name:     "invalid token",
			err:      errors.New("invalid token provided"),
			expected: true,
		},
		{
			name:     "access denied",
			err:      errors.New("access denied to resource"),
			expected: true,
		},
		{
			name:     "not an auth error",
			err:      errors.New("some other error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isAuthError(tt.err); got != tt.expected {
				t.Errorf("isAuthError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsFormatError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "invalid format",
			err:      errors.New("invalid manifest format"),
			expected: true,
		},
		{
			name:     "malformed data",
			err:      errors.New("malformed JSON data"),
			expected: true,
		},
		{
			name:     "parse error",
			err:      errors.New("failed to parse content"),
			expected: true,
		},
		{
			name:     "unmarshal error",
			err:      errors.New("failed to unmarshal JSON"),
			expected: true,
		},
		{
			name:     "decode error",
			err:      errors.New("failed to decode base64"),
			expected: true,
		},
		{
			name:     "syntax error",
			err:      errors.New("syntax error in YAML"),
			expected: true,
		},
		{
			name:     "not a format error",
			err:      errors.New("network connection failed"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isFormatError(tt.err); got != tt.expected {
				t.Errorf("isFormatError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsQuotaError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "quota exceeded",
			err:      errors.New("quota exceeded for user"),
			expected: true,
		},
		{
			name:     "limit reached",
			err:      errors.New("storage limit reached"),
			expected: true,
		},
		{
			name:     "rate limit",
			err:      errors.New("rate limit exceeded"),
			expected: true,
		},
		{
			name:     "storage full",
			err:      errors.New("storage full on device"),
			expected: true,
		},
		{
			name:     "throttled",
			err:      errors.New("request throttled"),
			expected: true,
		},
		{
			name:     "not a quota error",
			err:      errors.New("authentication failed"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isQuotaError(tt.err); got != tt.expected {
				t.Errorf("isQuotaError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestStandardErrorHandler_ClassifyError(t *testing.T) {
	handler := NewStandardErrorHandler()

	err := errors.New("unauthorized access")
	category := handler.ClassifyError(err)

	if category != ErrorCategoryAuth {
		t.Errorf("ClassifyError() = %v, want %v", category, ErrorCategoryAuth)
	}
}

func TestStandardErrorHandler_ShouldRetry(t *testing.T) {
	handler := NewStandardErrorHandler()

	tests := []struct {
		name     string
		err      error
		attempt  int
		expected bool
	}{
		{
			name:     "retryable error within limit",
			err:      ErrTimeout,
			attempt:  1,
			expected: true,
		},
		{
			name:     "non-retryable error",
			err:      ErrUnauthorized,
			attempt:  1,
			expected: false,
		},
		{
			name:     "retryable error at limit",
			err:      ErrTimeout,
			attempt:  3,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handler.ShouldRetry(tt.err, tt.attempt); got != tt.expected {
				t.Errorf("ShouldRetry() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestStandardErrorHandler_WrapError(t *testing.T) {
	handler := NewStandardErrorHandler()
	handler.SetContext("test_key", "test_value")

	baseErr := errors.New("connection failed")
	wrappedErr := handler.WrapError("push", "registry.example.com/repo:tag", baseErr, 2)

	if wrappedErr == nil {
		t.Fatal("WrapError() returned nil")
	}

	if wrappedErr.Operation != "push" {
		t.Errorf("WrapError().Operation = %v, want push", wrappedErr.Operation)
	}

	if wrappedErr.Resource != "registry.example.com/repo:tag" {
		t.Errorf("WrapError().Resource = %v, want registry.example.com/repo:tag", wrappedErr.Resource)
	}

	if wrappedErr.Attempt != 2 {
		t.Errorf("WrapError().Attempt = %v, want 2", wrappedErr.Attempt)
	}

	if wrappedErr.Cause != baseErr {
		t.Errorf("WrapError().Cause = %v, want %v", wrappedErr.Cause, baseErr)
	}

	if wrappedErr.Context["test_key"] != "test_value" {
		t.Errorf("WrapError().Context[test_key] = %v, want test_value", wrappedErr.Context["test_key"])
	}
}

func TestStandardErrorHandler_HandleError(t *testing.T) {
	handler := NewStandardErrorHandler()
	ctx := context.Background()

	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: true, // Should return nil
		},
		{
			name:     "operation error",
			err:      &OperationError{Operation: "test"},
			expected: false, // Should return as-is
		},
		{
			name:     "regular error",
			err:      errors.New("test error"),
			expected: false, // Should wrap in OperationError
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.HandleError(ctx, tt.err)

			if tt.expected {
				if result != nil {
					t.Errorf("HandleError() = %v, want nil", result)
				}
			} else {
				if result == nil {
					t.Error("HandleError() = nil, want non-nil")
				}
			}
		})
	}
}

func TestStandardErrorHandler_Context(t *testing.T) {
	handler := NewStandardErrorHandler()

	// Test setting and getting context
	handler.SetContext("key1", "value1")
	handler.SetContext("key2", 42)

	if got := handler.GetContext("key1"); got != "value1" {
		t.Errorf("GetContext(key1) = %v, want value1", got)
	}

	if got := handler.GetContext("key2"); got != 42 {
		t.Errorf("GetContext(key2) = %v, want 42", got)
	}

	if got := handler.GetContext("nonexistent"); got != nil {
		t.Errorf("GetContext(nonexistent) = %v, want nil", got)
	}
}

func TestHTTPError(t *testing.T) {
	err := HTTPError{StatusCode: 404, Message: "Not Found"}
	expected := "HTTP 404: Not Found"

	if got := err.Error(); got != expected {
		t.Errorf("HTTPError.Error() = %v, want %v", got, expected)
	}
}