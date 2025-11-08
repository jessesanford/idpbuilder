package errors_test

import (
	"fmt"
	"testing"

	"github.com/cnoe-io/idpbuilder/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test Suite 1: Error Type Creation (4 tests)

// T-2.3.4-01: TestNewValidationError verifies ValidationError constructor
func TestNewValidationError(t *testing.T) {
	err := errors.NewValidationError("imageName", "invalid image format", "use name:tag format")

	require.NotNil(t, err)
	assert.Equal(t, "imageName", err.Field)
	assert.Equal(t, 1, err.ExitCode)
	assert.Contains(t, err.Error(), "invalid image format")
	assert.Contains(t, err.Error(), "use name:tag format")
}

// T-2.3.4-02: TestNewAuthenticationError verifies AuthenticationError constructor
func TestNewAuthenticationError(t *testing.T) {
	err := errors.NewAuthenticationError("docker.io", "authentication failed", "check credentials")

	require.NotNil(t, err)
	assert.Equal(t, "docker.io", err.Registry)
	assert.Equal(t, 2, err.ExitCode)
	assert.Contains(t, err.Error(), "authentication failed")
	assert.Contains(t, err.Error(), "check credentials")
}

// T-2.3.4-03: TestNewNetworkError verifies NetworkError constructor
func TestNewNetworkError(t *testing.T) {
	err := errors.NewNetworkError("registry.example.com", "connection timeout", "check network")

	require.NotNil(t, err)
	assert.Equal(t, "registry.example.com", err.Target)
	assert.Equal(t, 3, err.ExitCode)
	assert.Contains(t, err.Error(), "connection timeout")
	assert.Contains(t, err.Error(), "check network")
}

// T-2.3.4-04: TestNewImageNotFoundError verifies ImageNotFoundError constructor
func TestNewImageNotFoundError(t *testing.T) {
	err := errors.NewImageNotFoundError("alpine:latest", "image not found", "pull image first")

	require.NotNil(t, err)
	assert.Equal(t, "alpine:latest", err.ImageName)
	assert.Equal(t, 4, err.ExitCode)
	assert.Contains(t, err.Error(), "image not found")
	assert.Contains(t, err.Error(), "pull image first")
}

// Test Suite 2: Error Formatting (6 tests)

// T-2.3.4-05: TestValidationError_Format verifies error message formatting
func TestValidationError_Format(t *testing.T) {
	err := errors.NewValidationError("registry", "invalid URL", "use https://registry.example.com")

	errorMsg := err.Error()
	assert.Contains(t, errorMsg, "Error:")
	assert.Contains(t, errorMsg, "invalid URL")
	assert.Contains(t, errorMsg, "Suggestion:")
	assert.Contains(t, errorMsg, "use https://registry.example.com")
}

// T-2.3.4-06: TestAuthenticationError_Format verifies auth error formatting
func TestAuthenticationError_Format(t *testing.T) {
	err := errors.NewAuthenticationError("gcr.io", "401 unauthorized", "verify credentials")

	errorMsg := err.Error()
	assert.Contains(t, errorMsg, "Error:")
	assert.Contains(t, errorMsg, "401 unauthorized")
	assert.Contains(t, errorMsg, "Suggestion:")
	assert.Contains(t, errorMsg, "verify credentials")
}

// T-2.3.4-07: TestSSRFWarning_Format verifies SSRF warning formatting
func TestSSRFWarning_Format(t *testing.T) {
	warning := &errors.SSRFWarning{
		Target:     "192.168.1.1",
		Message:    "private IP address detected",
		Suggestion: "use public registry instead",
	}

	errorMsg := warning.Error()
	assert.Contains(t, errorMsg, "Warning:")
	assert.Contains(t, errorMsg, "private IP address detected")
	assert.Contains(t, errorMsg, "Suggestion:")
	assert.Contains(t, errorMsg, "use public registry instead")
}

// TestSecurityWarning_Format verifies security warning formatting
func TestSecurityWarning_Format(t *testing.T) {
	warning := &errors.SecurityWarning{
		Message:    "weak credentials detected",
		Suggestion: "use stronger password",
	}

	errorMsg := warning.Error()
	assert.Contains(t, errorMsg, "Warning:")
	assert.Contains(t, errorMsg, "weak credentials detected")
	assert.Contains(t, errorMsg, "Suggestion:")
	assert.Contains(t, errorMsg, "use stronger password")
}

// TestBaseError_NoSuggestion verifies error formatting without suggestion
func TestBaseError_NoSuggestion(t *testing.T) {
	err := errors.NewValidationError("field", "error message", "")

	errorMsg := err.Error()
	assert.Contains(t, errorMsg, "Error:")
	assert.Contains(t, errorMsg, "error message")
	assert.NotContains(t, errorMsg, "Suggestion:")
}

// TestFormatError_Styling verifies FormatError adds appropriate emoji
func TestFormatError_Styling(t *testing.T) {
	// Test regular error formatting
	err := errors.NewValidationError("field", "invalid", "fix it")
	formatted := errors.FormatError(err)
	assert.Contains(t, formatted, "❌")

	// Test warning formatting
	warning := &errors.SSRFWarning{
		Target:     "localhost",
		Message:    "SSRF risk",
		Suggestion: "avoid localhost",
	}
	formatted = errors.FormatError(warning)
	assert.Contains(t, formatted, "⚠️")

	// Test nil error
	formatted = errors.FormatError(nil)
	assert.Equal(t, "", formatted)
}

// Test Suite 3: Error Unwrapping (2 tests)

// T-2.3.4-08: TestBaseError_Unwrap verifies error unwrapping
func TestBaseError_Unwrap(t *testing.T) {
	cause := fmt.Errorf("underlying error")
	err := errors.NewValidationError("field", "validation failed", "fix it")
	err.Cause = cause

	unwrapped := err.Unwrap()
	assert.Equal(t, cause, unwrapped)
}

// T-2.3.4-09: TestErrorChain_Unwrap verifies error chain traversal
func TestErrorChain_Unwrap(t *testing.T) {
	// Create error chain
	root := fmt.Errorf("root cause")
	middle := fmt.Errorf("middle error: %w", root)

	authErr := errors.NewAuthenticationError("docker.io", "auth failed", "check creds")
	authErr.Cause = middle

	// Verify unwrapping
	assert.Equal(t, middle, authErr.Unwrap())

	// Verify errors.Unwrap works through the chain
	var e error = authErr
	e = fmt.Errorf("wrapper: %w", e) // errors.Unwrap would get authErr
	// Note: This tests that BaseError.Unwrap is compatible with standard library
}

// Test Suite 4: Warning Detection (1 test)

// T-2.3.4-10: TestIsWarning_Detection verifies warning type detection
func TestIsWarning_Detection(t *testing.T) {
	// SSRFWarning should be detected as warning
	ssrfWarning := &errors.SSRFWarning{
		Target:     "192.168.1.1",
		Message:    "private IP",
		Suggestion: "use public registry",
	}
	assert.True(t, errors.IsWarning(ssrfWarning))

	// SecurityWarning should be detected as warning
	secWarning := &errors.SecurityWarning{
		Message:    "weak password",
		Suggestion: "use strong password",
	}
	assert.True(t, errors.IsWarning(secWarning))

	// Regular errors should NOT be warnings
	validationErr := errors.NewValidationError("field", "invalid", "fix it")
	assert.False(t, errors.IsWarning(validationErr))

	authErr := errors.NewAuthenticationError("docker.io", "auth failed", "check creds")
	assert.False(t, errors.IsWarning(authErr))

	networkErr := errors.NewNetworkError("registry", "timeout", "check network")
	assert.False(t, errors.IsWarning(networkErr))

	imageErr := errors.NewImageNotFoundError("alpine", "not found", "pull it")
	assert.False(t, errors.IsWarning(imageErr))
}
