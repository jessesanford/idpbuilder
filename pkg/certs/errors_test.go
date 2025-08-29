package certs

import (
	"errors"
	"testing"
)

func TestCertificateError_Error(t *testing.T) {
	// Test basic error formatting
	err := &CertificateError{Code: "TEST", Message: "Test message"}
	if result := err.Error(); !containsString(result, "[TEST]") || !containsString(result, "Test message") {
		t.Errorf("Expected error to contain '[TEST]' and 'Test message', got '%s'", result)
	}
	
	// Test error with context
	err.WithContext("key", "value")
	if result := err.Error(); !containsString(result, "context:") {
		t.Error("Error string should contain context")
	}
	
	// Test error with underlying
	underlying := errors.New("original")
	err.Wrap(underlying)
	if result := err.Error(); !containsString(result, "underlying:") {
		t.Error("Error string should contain underlying error")
	}
}

func TestCertificateError_Unwrap(t *testing.T) {
	underlying := errors.New("original error")
	err := &CertificateError{Code: "TEST", Message: "Wrapped", Underlying: underlying}
	
	if unwrapped := err.Unwrap(); unwrapped != underlying {
		t.Error("Unwrap should return underlying error")
	}
	
	// Test nil underlying
	errNoUnderlying := &CertificateError{Code: "TEST", Message: "test"}
	if unwrapped := errNoUnderlying.Unwrap(); unwrapped != nil {
		t.Error("Unwrap should return nil when no underlying error")
	}
}

func TestCertificateError_WithContext(t *testing.T) {
	err := NewCertificateError("TEST", "Test message")
	result := err.WithContext("key1", "value1").WithContext("key2", 42)
	
	if result != err {
		t.Error("WithContext should return same instance")
	}
	if len(err.Context) != 2 {
		t.Errorf("Expected 2 context items, got %d", len(err.Context))
	}
}

func TestCertificateError_WithSuggestion(t *testing.T) {
	err := NewCertificateError("TEST", "Test message")
	result := err.WithSuggestion("Try this").WithSuggestion("Or that")
	
	if result != err || len(err.Suggestions) != 2 {
		t.Error("WithSuggestion should return same instance with 2 suggestions")
	}
}

func TestCertificateError_Wrap(t *testing.T) {
	original := errors.New("original")
	err := NewCertificateError("TEST", "Wrapper").Wrap(original)
	
	if err.Underlying != original || !errors.Is(err, original) {
		t.Error("Wrap should set underlying error and be detectable with errors.Is")
	}
}

func TestNewCertificateError(t *testing.T) {
	err := NewCertificateError("CODE", "Message")
	
	if err.Code != "CODE" || err.Message != "Message" || err.Context == nil {
		t.Error("NewCertificateError should initialize all fields correctly")
	}
}

func TestPredefinedErrors(t *testing.T) {
	predefinedErrors := []*CertificateError{
		ErrClusterNotFound, ErrClusterConnection, ErrGiteaPodNotFound, ErrMultipleGiteaPods,
		ErrCertificateNotFound, ErrCertificateRead, ErrCertificateParse, ErrCertificateStore,
	}
	
	for _, err := range predefinedErrors {
		if err.Code == "" || err.Message == "" || len(err.Suggestions) == 0 {
			t.Errorf("Predefined error %s should have code, message, and suggestions", err.Code)
		}
	}
}

func TestWrapError(t *testing.T) {
	original := errors.New("original")
	wrapped := WrapError(original, "CODE", "Message")
	
	if wrapped.Code != "CODE" || wrapped.Message != "Message" || wrapped.Underlying != original {
		t.Error("WrapError should set all fields correctly")
	}
}

func TestIsErrorCode(t *testing.T) {
	certErr := NewCertificateError("TEST_CODE", "Message")
	
	if !IsErrorCode(certErr, "TEST_CODE") {
		t.Error("IsErrorCode should return true for matching code")
	}
	if IsErrorCode(certErr, "OTHER") {
		t.Error("IsErrorCode should return false for non-matching code")
	}
	if IsErrorCode(errors.New("standard"), "TEST") {
		t.Error("IsErrorCode should return false for non-CertificateError")
	}
}

// Helper function
func containsString(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}