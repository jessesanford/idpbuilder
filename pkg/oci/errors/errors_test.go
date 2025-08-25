package errors

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestNewOCIError(t *testing.T) {
	err := NewOCIError(ErrCodeBuildFailed, "TestComponent", "TestOperation", "Test message")

	if err.Code != ErrCodeBuildFailed {
		t.Errorf("Expected code %s, got %s", ErrCodeBuildFailed, err.Code)
	}
	if err.Component != "TestComponent" {
		t.Errorf("Expected component TestComponent, got %s", err.Component)
	}
	if err.Operation != "TestOperation" {
		t.Errorf("Expected operation TestOperation, got %s", err.Operation)
	}
	if err.Message != "Test message" {
		t.Errorf("Expected message 'Test message', got %s", err.Message)
	}
	if err.Category != ErrorCategoryBuild {
		t.Errorf("Expected category Build, got %s", err.Category.String())
	}
	if len(err.Details) != 0 {
		t.Errorf("Expected empty details map")
	}
}

func TestOCIError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *OCIError
		expected string
	}{
		{
			name: "Basic error",
			err: &OCIError{
				Code:      ErrCodeBuildFailed,
				Component: "TestComponent",
				Operation: "TestOperation",
				Message:   "Test message",
			},
			expected: "[1000] - TestComponent.TestOperation - Test message",
		},
		{
			name: "Error with cause",
			err: &OCIError{
				Code:      ErrCodeBuildFailed,
				Component: "TestComponent",
				Operation: "TestOperation",
				Message:   "Test message",
				Cause:     fmt.Errorf("underlying error"),
			},
			expected: "[1000] - TestComponent.TestOperation - Test message - caused by: underlying error",
		},
		{
			name: "Error without operation",
			err: &OCIError{
				Code:      ErrCodeBuildFailed,
				Component: "TestComponent",
				Message:   "Test message",
			},
			expected: "[1000] - TestComponent - Test message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestOCIError_WithCause(t *testing.T) {
	originalErr := fmt.Errorf("original error")
	err := NewOCIError(ErrCodeBuildFailed, "TestComponent", "TestOperation", "Test message").
		WithCause(originalErr)

	if err.Cause != originalErr {
		t.Errorf("Expected cause to be set")
	}

	if err.Unwrap() != originalErr {
		t.Errorf("Expected Unwrap() to return original error")
	}
}

func TestOCIError_WithDetails(t *testing.T) {
	details := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}

	err := NewOCIError(ErrCodeBuildFailed, "TestComponent", "TestOperation", "Test message").
		WithDetails(details)

	if len(err.Details) != 2 {
		t.Errorf("Expected 2 details, got %d", len(err.Details))
	}

	if err.Details["key1"] != "value1" {
		t.Errorf("Expected details[key1] = value1, got %v", err.Details["key1"])
	}

	if err.Details["key2"] != 42 {
		t.Errorf("Expected details[key2] = 42, got %v", err.Details["key2"])
	}
}

func TestOCIError_WithDetail(t *testing.T) {
	err := NewOCIError(ErrCodeBuildFailed, "TestComponent", "TestOperation", "Test message").
		WithDetail("testKey", "testValue")

	if err.Details["testKey"] != "testValue" {
		t.Errorf("Expected details[testKey] = testValue, got %v", err.Details["testKey"])
	}
}

func TestOCIError_WithRequestID(t *testing.T) {
	requestID := "test-request-123"
	err := NewOCIError(ErrCodeBuildFailed, "TestComponent", "TestOperation", "Test message").
		WithRequestID(requestID)

	if err.RequestID != requestID {
		t.Errorf("Expected RequestID %s, got %s", requestID, err.RequestID)
	}
}

func TestOCIError_WithRetryAfter(t *testing.T) {
	retryAfter := 5 * time.Second
	err := NewOCIError(ErrCodeBuildFailed, "TestComponent", "TestOperation", "Test message").
		WithRetryAfter(retryAfter)

	if err.RetryAfter == nil {
		t.Errorf("Expected RetryAfter to be set")
	}

	if *err.RetryAfter != retryAfter {
		t.Errorf("Expected RetryAfter %v, got %v", retryAfter, *err.RetryAfter)
	}
}

func TestOCIError_Is(t *testing.T) {
	err1 := NewOCIError(ErrCodeBuildFailed, "TestComponent", "TestOperation", "Test message")
	err2 := NewOCIError(ErrCodeBuildFailed, "OtherComponent", "OtherOperation", "Other message")
	err3 := NewOCIError(ErrCodeRegistryUnreachable, "TestComponent", "TestOperation", "Test message")
	regularErr := fmt.Errorf("regular error")

	// Same error code should match
	if !errors.Is(err1, err2) {
		t.Errorf("Expected err1 to match err2 (same error code)")
	}

	// Different error codes should not match
	if errors.Is(err1, err3) {
		t.Errorf("Expected err1 not to match err3 (different error codes)")
	}

	// OCIError should not match regular error
	if errors.Is(err1, regularErr) {
		t.Errorf("Expected OCIError not to match regular error")
	}

	// Regular error should not match OCIError
	if errors.Is(regularErr, err1) {
		t.Errorf("Expected regular error not to match OCIError")
	}
}

func TestGetCategoryFromCode(t *testing.T) {
	tests := []struct {
		code     string
		expected ErrorCategory
	}{
		{ErrCodeBuildFailed, ErrorCategoryBuild},
		{ErrCodeRegistryUnreachable, ErrorCategoryRegistry},
		{ErrCodeConfigMissing, ErrorCategoryConfiguration},
		{ErrCodeStackNotFound, ErrorCategoryStack},
		{ErrCodeAuthTokenInvalid, ErrorCategoryAuthentication},
		{ErrCodeSystemResourceExhausted, ErrorCategorySystem},
		{"invalid", ErrorCategorySystem}, // Default case
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			if got := GetCategoryFromCode(tt.code); got != tt.expected {
				t.Errorf("GetCategoryFromCode(%s) = %v, want %v", tt.code, got, tt.expected)
			}
		})
	}
}

func TestIsRetryable(t *testing.T) {
	tests := []struct {
		code     string
		expected bool
	}{
		{ErrCodeRegistryTimeout, true},
		{ErrCodeRegistryRateLimit, true},
		{ErrCodeSystemDiskFull, true},
		{ErrCodeSystemNetworkError, true},
		{ErrCodeBuildTempFailure, true},
		{ErrCodeBuildFailed, false},
		{ErrCodeAuthTokenInvalid, false},
		{ErrCodeConfigMissing, false},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			if got := IsRetryable(tt.code); got != tt.expected {
				t.Errorf("IsRetryable(%s) = %v, want %v", tt.code, got, tt.expected)
			}
		})
	}
}

func TestOCIError_IsRetryable(t *testing.T) {
	retryableErr := NewOCIError(ErrCodeRegistryTimeout, "TestComponent", "TestOperation", "Test message")
	nonRetryableErr := NewOCIError(ErrCodeBuildFailed, "TestComponent", "TestOperation", "Test message")

	if !retryableErr.IsRetryable() {
		t.Errorf("Expected retryable error to be retryable")
	}

	if nonRetryableErr.IsRetryable() {
		t.Errorf("Expected non-retryable error to not be retryable")
	}
}

func TestErrorCategory_String(t *testing.T) {
	tests := []struct {
		category ErrorCategory
		expected string
	}{
		{ErrorCategoryBuild, "Build"},
		{ErrorCategoryRegistry, "Registry"},
		{ErrorCategoryConfiguration, "Configuration"},
		{ErrorCategoryStack, "Stack"},
		{ErrorCategoryAuthentication, "Authentication"},
		{ErrorCategorySystem, "System"},
		{ErrorCategory(999), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.category.String(); got != tt.expected {
				t.Errorf("ErrorCategory.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}