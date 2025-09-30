package retry

import (
	"errors"
	"strings"
	"testing"
)

func TestMaxRetriesExceededError_Error(t *testing.T) {
	baseErr := errors.New("connection failed")
	err := &MaxRetriesExceededError{
		Attempts: 3,
		LastErr:  baseErr,
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "max retries exceeded") {
		t.Errorf("expected error message to contain 'max retries exceeded', got: %s", errMsg)
	}
	if !strings.Contains(errMsg, "3 attempts") {
		t.Errorf("expected error message to contain '3 attempts', got: %s", errMsg)
	}
	if !strings.Contains(errMsg, "connection failed") {
		t.Errorf("expected error message to contain underlying error, got: %s", errMsg)
	}
}

func TestMaxRetriesExceededError_Unwrap(t *testing.T) {
	baseErr := errors.New("connection failed")
	err := &MaxRetriesExceededError{
		Attempts: 3,
		LastErr:  baseErr,
	}

	unwrapped := errors.Unwrap(err)
	if unwrapped != baseErr {
		t.Errorf("expected unwrapped error to be base error, got: %v", unwrapped)
	}
}

func TestMaxRetriesExceededError_Is(t *testing.T) {
	baseErr := errors.New("connection failed")
	err := &MaxRetriesExceededError{
		Attempts: 3,
		LastErr:  baseErr,
	}

	// Test that errors.Is works with the wrapped error
	if !errors.Is(err, baseErr) {
		t.Error("expected errors.Is to find base error in wrapped error")
	}
}