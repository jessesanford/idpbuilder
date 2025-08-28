package certs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCertificateError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *CertificateError
		expected string
	}{
		{
			name: "basic error",
			err: &CertificateError{
				Operation: "test_op",
				Cause:     fmt.Errorf("test cause"),
				Context:   "",
			},
			expected: "certificate error in test_op: test cause",
		},
		{
			name: "error with context",
			err: &CertificateError{
				Operation: "test_op",
				Cause:     fmt.Errorf("test cause"),
				Context:   "additional context",
			},
			expected: "certificate error in test_op: test cause (additional context)",
		},
		{
			name: "error with empty context",
			err: &CertificateError{
				Operation: "validation",
				Cause:     fmt.Errorf("invalid certificate"),
				Context:   "",
			},
			expected: "certificate error in validation: invalid certificate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCertificateError_Unwrap(t *testing.T) {
	originalErr := fmt.Errorf("original error")
	certErr := &CertificateError{
		Operation: "test",
		Cause:     originalErr,
	}

	unwrapped := certErr.Unwrap()
	assert.Equal(t, originalErr, unwrapped)
}

func TestNewCertificateError(t *testing.T) {
	operation := "test_operation"
	cause := fmt.Errorf("test cause")
	context := "test context"
	suggestions := []string{"suggestion 1", "suggestion 2"}

	err := NewCertificateError(operation, cause, context, suggestions)

	assert.NotNil(t, err)
	assert.Equal(t, operation, err.Operation)
	assert.Equal(t, cause, err.Cause)
	assert.Equal(t, context, err.Context)
	assert.Equal(t, suggestions, err.Suggestions)
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name        string
		err         *CertificateError
		operation   string
		hasCause    bool
		hasContext  bool
		hasSuggestions bool
	}{
		{
			name:           "ErrClusterNotFound",
			err:            ErrClusterNotFound,
			operation:      "cluster_detection",
			hasCause:       true,
			hasContext:     true,
			hasSuggestions: true,
		},
		{
			name:           "ErrGiteaPodNotFound",
			err:            ErrGiteaPodNotFound,
			operation:      "pod_discovery",
			hasCause:       true,
			hasContext:     true,
			hasSuggestions: true,
		},
		{
			name:           "ErrCertificateNotFound",
			err:            ErrCertificateNotFound,
			operation:      "certificate_extraction",
			hasCause:       true,
			hasContext:     true,
			hasSuggestions: true,
		},
		{
			name:           "ErrCertificateExpired",
			err:            ErrCertificateExpired,
			operation:      "certificate_validation",
			hasCause:       true,
			hasContext:     true,
			hasSuggestions: true,
		},
		{
			name:           "ErrCertificateInvalid",
			err:            ErrCertificateInvalid,
			operation:      "certificate_validation",
			hasCause:       true,
			hasContext:     true,
			hasSuggestions: true,
		},
		{
			name:           "ErrStoragePermission",
			err:            ErrStoragePermission,
			operation:      "certificate_storage",
			hasCause:       true,
			hasContext:     true,
			hasSuggestions: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.operation, tt.err.Operation)
			
			if tt.hasCause {
				assert.NotNil(t, tt.err.Cause)
				assert.NotEmpty(t, tt.err.Cause.Error())
			}
			
			if tt.hasContext {
				assert.NotEmpty(t, tt.err.Context)
			}
			
			if tt.hasSuggestions {
				assert.NotEmpty(t, tt.err.Suggestions)
				assert.True(t, len(tt.err.Suggestions) > 0)
				for _, suggestion := range tt.err.Suggestions {
					assert.NotEmpty(t, suggestion)
				}
			}
		})
	}
}

func TestErrorChaining(t *testing.T) {
	originalErr := fmt.Errorf("root cause")
	certErr := NewCertificateError("test_op", originalErr, "test context", []string{"fix it"})

	// Test that error can be unwrapped
	assert.Equal(t, originalErr, certErr.Unwrap())
	
	// Test error message
	expectedMsg := "certificate error in test_op: root cause (test context)"
	assert.Equal(t, expectedMsg, certErr.Error())
}

func TestErrorSuggestions(t *testing.T) {
	suggestions := []string{
		"Check cluster connectivity",
		"Verify pod is running",
		"Try again later",
	}
	
	err := NewCertificateError("test", fmt.Errorf("test"), "context", suggestions)
	
	assert.Equal(t, suggestions, err.Suggestions)
	assert.Len(t, err.Suggestions, 3)
}

func TestPredefineErrorMessages(t *testing.T) {
	// Test that predefined errors have meaningful error messages
	errorMessages := map[string]string{
		"ErrClusterNotFound":      ErrClusterNotFound.Error(),
		"ErrGiteaPodNotFound":     ErrGiteaPodNotFound.Error(),
		"ErrCertificateNotFound":  ErrCertificateNotFound.Error(),
		"ErrCertificateExpired":   ErrCertificateExpired.Error(),
		"ErrCertificateInvalid":   ErrCertificateInvalid.Error(),
		"ErrStoragePermission":    ErrStoragePermission.Error(),
	}
	
	for errorName, message := range errorMessages {
		t.Run(errorName, func(t *testing.T) {
			assert.NotEmpty(t, message)
			assert.Contains(t, message, "certificate error in")
		})
	}
}