package errors

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"
)

func TestErrorCategory_String(t *testing.T) {
	tests := []struct {
		category ErrorCategory
		expected string
	}{
		{ErrorCategoryUnknown, "unknown"},
		{ErrorCategoryTransient, "transient"},
		{ErrorCategoryPermanent, "permanent"},
		{ErrorCategoryAuth, "authentication"},
		{ErrorCategoryNetwork, "network"},
		{ErrorCategoryFormat, "format"},
		{ErrorCategoryQuota, "quota"},
	}

	for _, tt := range tests {
		if got := tt.category.String(); got != tt.expected {
			t.Errorf("ErrorCategory.String() = %v, want %v", got, tt.expected)
		}
	}
}

func TestErrorCategory_ShouldRetry(t *testing.T) {
	tests := []struct {
		category ErrorCategory
		expected bool
	}{
		{ErrorCategoryUnknown, false},
		{ErrorCategoryTransient, true},
		{ErrorCategoryPermanent, false},
		{ErrorCategoryAuth, false},
		{ErrorCategoryNetwork, true},
		{ErrorCategoryFormat, false},
		{ErrorCategoryQuota, false},
	}

	for _, tt := range tests {
		if got := tt.category.ShouldRetry(); got != tt.expected {
			t.Errorf("ErrorCategory.ShouldRetry() for %v = %v, want %v", tt.category, got, tt.expected)
		}
	}
}

func TestOperationError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *OperationError
		expected string
	}{
		{
			name: "with resource",
			err: &OperationError{
				Operation: "push",
				Resource:  "registry.example.com/repo:tag",
				Category:  ErrorCategoryNetwork,
				Attempt:   2,
				Cause:     errors.New("connection failed"),
			},
			expected: "push operation failed for registry.example.com/repo:tag (attempt 2, network): connection failed",
		},
		{
			name: "without resource",
			err: &OperationError{
				Operation: "list",
				Category:  ErrorCategoryAuth,
				Attempt:   1,
				Cause:     errors.New("unauthorized"),
			},
			expected: "list operation failed (attempt 1, authentication): unauthorized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("OperationError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestOperationError_Unwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := &OperationError{
		Operation: "test",
		Cause:     cause,
	}

	if got := err.Unwrap(); got != cause {
		t.Errorf("OperationError.Unwrap() = %v, want %v", got, cause)
	}
}

func TestOperationError_Is(t *testing.T) {
	baseErr := errors.New("base error")
	opErr := &OperationError{
		Operation: "test",
		Category:  ErrorCategoryNetwork,
		Cause:     baseErr,
	}

	tests := []struct {
		name     string
		target   error
		expected bool
	}{
		{
			name:     "same operation error",
			target:   &OperationError{Operation: "test", Category: ErrorCategoryNetwork, Cause: baseErr},
			expected: true,
		},
		{
			name:     "different operation",
			target:   &OperationError{Operation: "other", Category: ErrorCategoryNetwork, Cause: baseErr},
			expected: false,
		},
		{
			name:     "different category",
			target:   &OperationError{Operation: "test", Category: ErrorCategoryAuth, Cause: baseErr},
			expected: false,
		},
		{
			name:     "underlying error",
			target:   baseErr,
			expected: true,
		},
		{
			name:     "nil target",
			target:   nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := opErr.Is(tt.target); got != tt.expected {
				t.Errorf("OperationError.Is() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDefaultRetryInfo(t *testing.T) {
	info := DefaultRetryInfo()

	if info.MaxAttempts != 3 {
		t.Errorf("DefaultRetryInfo().MaxAttempts = %v, want 3", info.MaxAttempts)
	}

	if info.BaseDelay != 100*time.Millisecond {
		t.Errorf("DefaultRetryInfo().BaseDelay = %v, want 100ms", info.BaseDelay)
	}

	if info.MaxDelay != 10*time.Second {
		t.Errorf("DefaultRetryInfo().MaxDelay = %v, want 10s", info.MaxDelay)
	}

	if info.Multiplier != 2.0 {
		t.Errorf("DefaultRetryInfo().Multiplier = %v, want 2.0", info.Multiplier)
	}

	if !info.Jitter {
		t.Errorf("DefaultRetryInfo().Jitter = %v, want true", info.Jitter)
	}
}

func TestIsTemporary(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "connection refused",
			err:      ErrConnectionRefused,
			expected: true,
		},
		{
			name:     "timeout",
			err:      ErrTimeout,
			expected: true,
		},
		{
			name:     "rate limit",
			err:      ErrRateLimit,
			expected: true,
		},
		{
			name:     "unauthorized",
			err:      ErrUnauthorized,
			expected: false,
		},
		{
			name:     "context deadline exceeded",
			err:      context.DeadlineExceeded,
			expected: false,
		},
		{
			name:     "context canceled",
			err:      context.Canceled,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsTemporary(tt.err); got != tt.expected {
				t.Errorf("IsTemporary() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsRetryable(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name: "retryable operation error",
			err: &OperationError{
				Category:  ErrorCategoryTransient,
				Retryable: true,
			},
			expected: true,
		},
		{
			name: "non-retryable operation error",
			err: &OperationError{
				Category:  ErrorCategoryPermanent,
				Retryable: false,
			},
			expected: false,
		},
		{
			name:     "temporary error",
			err:      ErrTimeout,
			expected: true,
		},
		{
			name:     "permanent error",
			err:      ErrUnauthorized,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRetryable(tt.err); got != tt.expected {
				t.Errorf("IsRetryable() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Mock network error for testing
type mockNetError struct {
	message   string
	temporary bool
	timeout   bool
}

func (e mockNetError) Error() string   { return e.message }
func (e mockNetError) Temporary() bool { return e.temporary }
func (e mockNetError) Timeout() bool   { return e.timeout }

func TestIsTemporary_NetError(t *testing.T) {
	tests := []struct {
		name     string
		err      net.Error
		expected bool
	}{
		{
			name:     "temporary network error",
			err:      mockNetError{message: "temp error", temporary: true},
			expected: true,
		},
		{
			name:     "timeout network error",
			err:      mockNetError{message: "timeout", timeout: true},
			expected: true,
		},
		{
			name:     "permanent network error",
			err:      mockNetError{message: "permanent error", temporary: false, timeout: false},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsTemporary(tt.err); got != tt.expected {
				t.Errorf("IsTemporary() = %v, want %v", got, tt.expected)
			}
		})
	}
}