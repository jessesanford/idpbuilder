package errors

import (
	"errors"
	"testing"
	"time"
)

func TestNewOCIError(t *testing.T) {
	err := NewOCIError(CodeBuildFailed, "Build failed")
	if err.Code != CodeBuildFailed {
		t.Errorf("Expected code %s, got %s", CodeBuildFailed, err.Code)
	}
	if err.Details == nil {
		t.Error("Expected Details to be initialized")
	}
}

func TestOCIErrorWithCause(t *testing.T) {
	originalErr := errors.New("underlying error")
	ociErr := NewOCIError(CodeBuildFailed, "Build failed").WithCause(originalErr)
	if ociErr.Unwrap() != originalErr {
		t.Error("Expected Unwrap to return the original error")
	}
}

func TestOCIErrorWithDetails(t *testing.T) {
	details := map[string]interface{}{"dockerfile": "Dockerfile"}
	err := NewOCIError(CodeBuildFailed, "Build failed").WithDetails(details)
	if err.Details["dockerfile"] != "Dockerfile" {
		t.Error("Expected dockerfile detail to be set")
	}
}

func TestOCIErrorWithRetryAfter(t *testing.T) {
	duration := 5 * time.Second
	err := NewOCIError(CodeBuildTimeout, "Build timed out").WithRetryAfter(duration)
	if !err.Retry {
		t.Error("Expected error to be marked as retryable")
	}
}

func TestOCIErrorError(t *testing.T) {
	err := NewOCIErrorWithComponent(CodeBuildFailed, "Build failed", "builder", "docker-build")
	expected := "[code=1001,component=builder,operation=docker-build] Build failed"
	if err.Error() != expected {
		t.Errorf("Expected %s, got %s", expected, err.Error())
	}
}

func TestOCIErrorIs(t *testing.T) {
	err1 := NewOCIError(CodeBuildFailed, "Build failed")
	err2 := NewOCIError(CodeBuildFailed, "Different message")
	if !err1.Is(err2) {
		t.Error("Expected errors with same code to match")
	}
}

func TestGetCategory(t *testing.T) {
	tests := []struct {
		code     string
		expected ErrorCategory
	}{
		{CodeBuildFailed, CategoryBuild},
		{CodeRegistryAuthFailed, CategoryRegistry},
		{CodeConfigInvalid, CategoryConfiguration},
		{CodeStackDeployFailed, CategoryStack},
		{CodeAuthTokenInvalid, CategoryAuthentication},
		{CodeSystemStorageError, CategorySystem},
	}
	
	for _, test := range tests {
		category := GetCategory(test.code)
		if category != test.expected {
			t.Errorf("For code %s, expected %s, got %s", test.code, test.expected, category)
		}
	}
}

func TestIsRetryable(t *testing.T) {
	retryableErr := NewOCIError(CodeBuildTimeout, "Timeout").WithRetryAfter(5 * time.Second)
	if !IsRetryable(retryableErr) {
		t.Error("Expected retryable error to be retryable")
	}
}
