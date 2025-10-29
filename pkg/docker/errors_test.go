package docker

import (
	"errors"
	"strings"
	"testing"
)

// TestDaemonConnectionError verifies DaemonConnectionError formatting and unwrapping.
func TestDaemonConnectionError(t *testing.T) {
	cause := errors.New("connection refused")
	err := &DaemonConnectionError{Cause: cause}

	// Test Error() formatting
	errMsg := err.Error()
	if !strings.Contains(errMsg, "Docker daemon connection error") {
		t.Errorf("Error message should contain 'Docker daemon connection error', got: %s", errMsg)
	}
	if !strings.Contains(errMsg, "connection refused") {
		t.Errorf("Error message should contain cause, got: %s", errMsg)
	}

	// Test Unwrap()
	unwrapped := errors.Unwrap(err)
	if unwrapped != cause {
		t.Errorf("Unwrap should return cause, got: %v", unwrapped)
	}

	t.Log("✅ DaemonConnectionError works correctly")
}

// TestImageNotFoundError verifies ImageNotFoundError formatting.
func TestImageNotFoundError(t *testing.T) {
	err := &ImageNotFoundError{ImageName: "myapp:latest"}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "not found") {
		t.Errorf("Error message should contain 'not found', got: %s", errMsg)
	}
	if !strings.Contains(errMsg, "myapp:latest") {
		t.Errorf("Error message should contain image name, got: %s", errMsg)
	}

	t.Log("✅ ImageNotFoundError works correctly")
}

// TestImageConversionError verifies ImageConversionError formatting and unwrapping.
func TestImageConversionError(t *testing.T) {
	cause := errors.New("tar export failed")
	err := &ImageConversionError{
		ImageName: "myapp:latest",
		Cause:     cause,
	}

	// Test Error() formatting
	errMsg := err.Error()
	if !strings.Contains(errMsg, "failed to convert") {
		t.Errorf("Error message should contain 'failed to convert', got: %s", errMsg)
	}
	if !strings.Contains(errMsg, "myapp:latest") {
		t.Errorf("Error message should contain image name, got: %s", errMsg)
	}
	if !strings.Contains(errMsg, "tar export failed") {
		t.Errorf("Error message should contain cause, got: %s", errMsg)
	}

	// Test Unwrap()
	unwrapped := errors.Unwrap(err)
	if unwrapped != cause {
		t.Errorf("Unwrap should return cause, got: %v", unwrapped)
	}

	t.Log("✅ ImageConversionError works correctly")
}

// TestValidationError verifies ValidationError formatting.
func TestValidationError(t *testing.T) {
	err := &ValidationError{
		Field:   "imageName",
		Message: "contains invalid characters",
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "validation error") {
		t.Errorf("Error message should contain 'validation error', got: %s", errMsg)
	}
	if !strings.Contains(errMsg, "imageName") {
		t.Errorf("Error message should contain field name, got: %s", errMsg)
	}
	if !strings.Contains(errMsg, "invalid characters") {
		t.Errorf("Error message should contain message, got: %s", errMsg)
	}

	t.Log("✅ ValidationError works correctly")
}

// TestErrorTypesImplementError verifies all error types satisfy error interface.
func TestErrorTypesImplementError(t *testing.T) {
	var _ error = &DaemonConnectionError{}
	var _ error = &ImageNotFoundError{}
	var _ error = &ImageConversionError{}
	var _ error = &ValidationError{}

	t.Log("✅ All error types implement error interface")
}
