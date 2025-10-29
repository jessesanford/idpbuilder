package registry

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAuthenticationError verifies AuthenticationError formatting.
func TestAuthenticationError(t *testing.T) {
	cause := errors.New("invalid credentials")
	err := &AuthenticationError{
		Registry: "registry.io",
		Cause:    cause,
	}

	assert.Equal(t, "registry authentication failed for registry.io: invalid credentials", err.Error())
	assert.Equal(t, cause, err.Unwrap())
}

// TestNetworkError verifies NetworkError formatting.
func TestNetworkError(t *testing.T) {
	cause := errors.New("connection timeout")
	err := &NetworkError{
		Registry: "registry.io",
		Cause:    cause,
	}

	assert.Equal(t, "network error connecting to registry registry.io: connection timeout", err.Error())
	assert.Equal(t, cause, err.Unwrap())
}

// TestRegistryUnavailableError verifies RegistryUnavailableError formatting.
func TestRegistryUnavailableError(t *testing.T) {
	err := &RegistryUnavailableError{
		Registry:   "registry.io",
		StatusCode: 503,
	}

	assert.Equal(t, "registry registry.io unavailable (status: 503)", err.Error())
}

// TestPushFailedError verifies PushFailedError formatting.
func TestPushFailedError(t *testing.T) {
	cause := errors.New("layer upload failed")
	err := &PushFailedError{
		TargetRef: "registry.io/myapp:latest",
		Cause:     cause,
	}

	assert.Equal(t, "push to registry.io/myapp:latest failed: layer upload failed", err.Error())
	assert.Equal(t, cause, err.Unwrap())
}

// TestErrorTypesImplementError verifies all error types satisfy error interface.
func TestErrorTypesImplementError(t *testing.T) {
	var _ error = &AuthenticationError{}
	var _ error = &NetworkError{}
	var _ error = &RegistryUnavailableError{}
	var _ error = &PushFailedError{}

	assert.True(t, true, "All error types implement error interface")
}

// TestErrorUnwrapping verifies error unwrapping support.
func TestErrorUnwrapping(t *testing.T) {
	rootCause := errors.New("root cause")

	// Test AuthenticationError unwrapping
	authErr := &AuthenticationError{Registry: "test", Cause: rootCause}
	assert.Equal(t, rootCause, errors.Unwrap(authErr))

	// Test NetworkError unwrapping
	netErr := &NetworkError{Registry: "test", Cause: rootCause}
	assert.Equal(t, rootCause, errors.Unwrap(netErr))

	// Test PushFailedError unwrapping
	pushErr := &PushFailedError{TargetRef: "test", Cause: rootCause}
	assert.Equal(t, rootCause, errors.Unwrap(pushErr))

	// RegistryUnavailableError does not support unwrapping (no Cause field)
	unavailErr := &RegistryUnavailableError{Registry: "test", StatusCode: 503}
	assert.Nil(t, errors.Unwrap(unavailErr))
}
