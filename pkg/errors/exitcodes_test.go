package errors_test

import (
	"fmt"
	"testing"

	"github.com/cnoe-io/idpbuilder/pkg/errors"
	"github.com/stretchr/testify/assert"
)

// Test Suite 2: Exit Code Mapping (8 tests)

// T-2.3.5-01: TestGetExitCode_ValidationError verifies exit code 1 for validation errors
func TestGetExitCode_ValidationError(t *testing.T) {
	err := errors.NewValidationError("imageName", "invalid format", "use name:tag")

	exitCode := errors.GetExitCode(err)
	assert.Equal(t, errors.ExitValidationError, exitCode)
	assert.Equal(t, 1, exitCode)
}

// T-2.3.5-02: TestGetExitCode_AuthenticationError verifies exit code 2 for auth errors
func TestGetExitCode_AuthenticationError(t *testing.T) {
	err := errors.NewAuthenticationError("docker.io", "auth failed", "check credentials")

	exitCode := errors.GetExitCode(err)
	assert.Equal(t, errors.ExitAuthError, exitCode)
	assert.Equal(t, 2, exitCode)
}

// T-2.3.5-03: TestGetExitCode_NetworkError verifies exit code 3 for network errors
func TestGetExitCode_NetworkError(t *testing.T) {
	err := errors.NewNetworkError("registry.example.com", "connection timeout", "check network")

	exitCode := errors.GetExitCode(err)
	assert.Equal(t, errors.ExitNetworkError, exitCode)
	assert.Equal(t, 3, exitCode)
}

// T-2.3.5-04: TestGetExitCode_ImageNotFoundError verifies exit code 4 for image errors
func TestGetExitCode_ImageNotFoundError(t *testing.T) {
	err := errors.NewImageNotFoundError("alpine:latest", "image not found", "pull image")

	exitCode := errors.GetExitCode(err)
	assert.Equal(t, errors.ExitImageNotFound, exitCode)
	assert.Equal(t, 4, exitCode)
}

// T-2.3.5-05: TestGetExitCode_GenericError verifies exit code 1 for untyped errors
func TestGetExitCode_GenericError(t *testing.T) {
	err := fmt.Errorf("generic error without type")

	exitCode := errors.GetExitCode(err)
	assert.Equal(t, errors.ExitGenericError, exitCode)
	assert.Equal(t, 1, exitCode)
}

// T-2.3.5-06: TestGetExitCode_NilError verifies exit code 0 for nil
func TestGetExitCode_NilError(t *testing.T) {
	exitCode := errors.GetExitCode(nil)
	assert.Equal(t, errors.ExitSuccess, exitCode)
	assert.Equal(t, 0, exitCode)
}

// T-2.3.5-07: TestGetExitCode_WrappedValidationError verifies unwrapping works
func TestGetExitCode_WrappedValidationError(t *testing.T) {
	// Create validation error and wrap it
	validationErr := errors.NewValidationError("field", "invalid", "fix it")
	wrappedErr := fmt.Errorf("wrapped error: %w", validationErr)

	// GetExitCode should unwrap and find the ValidationError
	exitCode := errors.GetExitCode(wrappedErr)
	assert.Equal(t, errors.ExitValidationError, exitCode)
	assert.Equal(t, 1, exitCode)
}

// T-2.3.5-08: TestGetExitCode_WrappedAuthError verifies unwrapping for auth errors
func TestGetExitCode_WrappedAuthError(t *testing.T) {
	// Create auth error and wrap it multiple times
	authErr := errors.NewAuthenticationError("docker.io", "401", "check creds")
	wrappedOnce := fmt.Errorf("registry error: %w", authErr)
	wrappedTwice := fmt.Errorf("push failed: %w", wrappedOnce)

	// GetExitCode should unwrap through the chain
	exitCode := errors.GetExitCode(wrappedTwice)
	assert.Equal(t, errors.ExitAuthError, exitCode)
	assert.Equal(t, 2, exitCode)
}

// TestGetExitCode_AllErrorTypes verifies all error types return correct codes
func TestGetExitCode_AllErrorTypes(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected int
	}{
		{
			name:     "ValidationError",
			err:      errors.NewValidationError("field", "invalid", "fix"),
			expected: 1,
		},
		{
			name:     "AuthenticationError",
			err:      errors.NewAuthenticationError("docker.io", "auth failed", "check"),
			expected: 2,
		},
		{
			name:     "NetworkError",
			err:      errors.NewNetworkError("registry", "timeout", "retry"),
			expected: 3,
		},
		{
			name:     "ImageNotFoundError",
			err:      errors.NewImageNotFoundError("alpine", "not found", "pull"),
			expected: 4,
		},
		{
			name:     "GenericError",
			err:      fmt.Errorf("generic error"),
			expected: 1,
		},
		{
			name:     "NilError",
			err:      nil,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exitCode := errors.GetExitCode(tt.err)
			assert.Equal(t, tt.expected, exitCode)
		})
	}
}
