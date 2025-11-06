package push

import (
	goerrors "errors"
	"fmt"
	"testing"

	"github.com/cnoe-io/idpbuilder/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test Suite 3: Error Wrapping Integration (12 tests)

// T-2.3.6-01: TestWrapDockerError_ImageNotFound verifies image not found detection
func TestWrapDockerError_ImageNotFound(t *testing.T) {
	originalErr := fmt.Errorf("No such image: alpine:latest")

	wrappedErr := WrapDockerError(originalErr, "alpine:latest")

	require.NotNil(t, wrappedErr)

	// Should be an ImageNotFoundError
	var imageErr *errors.ImageNotFoundError
	assert.ErrorAs(t, wrappedErr, &imageErr)
	assert.Equal(t, "alpine:latest", imageErr.ImageName)
	assert.Contains(t, wrappedErr.Error(), "alpine:latest")
	assert.Contains(t, wrappedErr.Error(), "docker pull")
}

// T-2.3.6-02: TestWrapDockerError_ConnectionRefused verifies connection refused detection
func TestWrapDockerError_ConnectionRefused(t *testing.T) {
	originalErr := fmt.Errorf("connection refused")

	wrappedErr := WrapDockerError(originalErr, "alpine:latest")

	require.NotNil(t, wrappedErr)

	// Should be a NetworkError
	var networkErr *errors.NetworkError
	assert.ErrorAs(t, wrappedErr, &networkErr)
	assert.Equal(t, "docker daemon", networkErr.Target)
	assert.Contains(t, wrappedErr.Error(), "Docker daemon")
	assert.Contains(t, wrappedErr.Error(), "systemctl start docker")
}

// T-2.3.6-03: TestWrapDockerError_CannotConnect verifies "Cannot connect" detection
func TestWrapDockerError_CannotConnect(t *testing.T) {
	originalErr := fmt.Errorf("Cannot connect to the Docker daemon")

	wrappedErr := WrapDockerError(originalErr, "myimage:v1")

	require.NotNil(t, wrappedErr)

	// Should be a NetworkError
	var networkErr *errors.NetworkError
	assert.ErrorAs(t, wrappedErr, &networkErr)
	assert.Equal(t, "docker daemon", networkErr.Target)
	assert.Contains(t, wrappedErr.Error(), "Docker daemon")
}

// T-2.3.6-04: TestWrapDockerError_GenericError verifies fallback for unknown errors
func TestWrapDockerError_GenericError(t *testing.T) {
	originalErr := fmt.Errorf("unknown Docker error")

	wrappedErr := WrapDockerError(originalErr, "alpine:latest")

	require.NotNil(t, wrappedErr)

	// Should be a generic wrapped error
	assert.Contains(t, wrappedErr.Error(), "Docker error")
	assert.Contains(t, wrappedErr.Error(), "unknown Docker error")

	// Should NOT be a typed error
	var imageErr *errors.ImageNotFoundError
	assert.False(t, goerrors.As(wrappedErr, &imageErr))

	var networkErr *errors.NetworkError
	assert.False(t, goerrors.As(wrappedErr, &networkErr))
}

// T-2.3.6-05: TestWrapRegistryError_Unauthorized verifies 401 detection
func TestWrapRegistryError_Unauthorized(t *testing.T) {
	originalErr := fmt.Errorf("401 Unauthorized")

	wrappedErr := WrapRegistryError(originalErr, "docker.io")

	require.NotNil(t, wrappedErr)

	// Should be an AuthenticationError
	var authErr *errors.AuthenticationError
	assert.ErrorAs(t, wrappedErr, &authErr)
	assert.Equal(t, "docker.io", authErr.Registry)
	assert.Contains(t, wrappedErr.Error(), "authentication failed")
	assert.Contains(t, wrappedErr.Error(), "docker.io")
}

// T-2.3.6-06: TestWrapRegistryError_ConnectionRefused verifies connection refused
func TestWrapRegistryError_ConnectionRefused(t *testing.T) {
	originalErr := fmt.Errorf("connection refused")

	wrappedErr := WrapRegistryError(originalErr, "registry.example.com")

	require.NotNil(t, wrappedErr)

	// Should be a NetworkError
	var networkErr *errors.NetworkError
	assert.ErrorAs(t, wrappedErr, &networkErr)
	assert.Equal(t, "registry.example.com", networkErr.Target)
	assert.Contains(t, wrappedErr.Error(), "cannot connect")
}

// T-2.3.6-07: TestWrapRegistryError_Timeout verifies timeout detection
func TestWrapRegistryError_Timeout(t *testing.T) {
	originalErr := fmt.Errorf("request timeout")

	wrappedErr := WrapRegistryError(originalErr, "gcr.io")

	require.NotNil(t, wrappedErr)

	// Should be a NetworkError
	var networkErr *errors.NetworkError
	assert.ErrorAs(t, wrappedErr, &networkErr)
	assert.Equal(t, "gcr.io", networkErr.Target)
	assert.Contains(t, wrappedErr.Error(), "cannot connect")
}

// T-2.3.6-08: TestWrapRegistryError_TLSError verifies x509 certificate error detection
func TestWrapRegistryError_TLSError(t *testing.T) {
	originalErr := fmt.Errorf("x509: certificate signed by unknown authority")

	wrappedErr := WrapRegistryError(originalErr, "self-signed.example.com")

	require.NotNil(t, wrappedErr)

	// Should be a NetworkError
	var networkErr *errors.NetworkError
	assert.ErrorAs(t, wrappedErr, &networkErr)
	assert.Equal(t, "self-signed.example.com", networkErr.Target)
	assert.Contains(t, wrappedErr.Error(), "TLS certificate verification failed")
	assert.Contains(t, wrappedErr.Error(), "--insecure")
}

// T-2.3.6-09: TestWrapRegistryError_GenericError verifies fallback for unknown errors
func TestWrapRegistryError_GenericError(t *testing.T) {
	originalErr := fmt.Errorf("unknown registry error")

	wrappedErr := WrapRegistryError(originalErr, "docker.io")

	require.NotNil(t, wrappedErr)

	// Should be a generic wrapped error
	assert.Contains(t, wrappedErr.Error(), "registry error")
	assert.Contains(t, wrappedErr.Error(), "unknown registry error")

	// Should NOT be a typed error
	var authErr *errors.AuthenticationError
	assert.False(t, goerrors.As(wrappedErr, &authErr))

	var networkErr *errors.NetworkError
	assert.False(t, goerrors.As(wrappedErr, &networkErr))
}

// T-2.3.6-10: TestWrapDockerError_PreservesImageName verifies image name context
func TestWrapDockerError_PreservesImageName(t *testing.T) {
	tests := []struct {
		name      string
		imageName string
		errMsg    string
	}{
		{"simple", "alpine:latest", "No such image: alpine:latest"},
		{"with_registry", "gcr.io/project/image:v1", "No such image: gcr.io/project/image:v1"},
		{"sha256", "alpine@sha256:abc123", "No such image: alpine@sha256:abc123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalErr := fmt.Errorf("%s", tt.errMsg)
			wrappedErr := WrapDockerError(originalErr, tt.imageName)

			var imageErr *errors.ImageNotFoundError
			require.True(t, assert.ErrorAs(t, wrappedErr, &imageErr))
			assert.Equal(t, tt.imageName, imageErr.ImageName)
			assert.Contains(t, wrappedErr.Error(), tt.imageName)
		})
	}
}

// T-2.3.6-11: TestWrapRegistryError_PreservesRegistry verifies registry context
func TestWrapRegistryError_PreservesRegistry(t *testing.T) {
	tests := []struct {
		name     string
		registry string
		errMsg   string
	}{
		{"docker_hub", "docker.io", "401 Unauthorized"},
		{"gcr", "gcr.io", "unauthorized: authentication required"},
		{"custom_port", "registry.example.com:5000", "401"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalErr := fmt.Errorf("%s", tt.errMsg)
			wrappedErr := WrapRegistryError(originalErr, tt.registry)

			var authErr *errors.AuthenticationError
			require.True(t, assert.ErrorAs(t, wrappedErr, &authErr))
			assert.Equal(t, tt.registry, authErr.Registry)
			assert.Contains(t, wrappedErr.Error(), tt.registry)
		})
	}
}

// T-2.3.6-12: TestWrapErrors_ChainUnwraps verifies error chain traversal
func TestWrapErrors_ChainUnwraps(t *testing.T) {
	// Create a chain: Docker error -> wrapped -> wrapped again
	originalErr := fmt.Errorf("No such image: alpine:latest")
	wrappedOnce := WrapDockerError(originalErr, "alpine:latest")
	wrappedTwice := fmt.Errorf("failed to retrieve image: %w", wrappedOnce)

	// Verify we can still detect the ImageNotFoundError through the chain
	var imageErr *errors.ImageNotFoundError
	assert.ErrorAs(t, wrappedTwice, &imageErr)
	assert.Equal(t, "alpine:latest", imageErr.ImageName)

	// Verify exit code mapping works through the chain
	exitCode := errors.GetExitCode(wrappedTwice)
	assert.Equal(t, errors.ExitImageNotFound, exitCode)
	assert.Equal(t, 4, exitCode)
}

// TestWrapDockerError_NilError verifies nil handling
func TestWrapDockerError_NilError(t *testing.T) {
	wrappedErr := WrapDockerError(nil, "alpine:latest")
	assert.Nil(t, wrappedErr)
}

// TestWrapRegistryError_NilError verifies nil handling
func TestWrapRegistryError_NilError(t *testing.T) {
	wrappedErr := WrapRegistryError(nil, "docker.io")
	assert.Nil(t, wrappedErr)
}
