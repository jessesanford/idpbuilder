package docker

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A. Constructor Tests

func TestNewClient_Success(t *testing.T) {
	// Verify client creation with running daemon
	client, err := NewClient()
	require.NoError(t, err, "should create client when daemon is running")
	require.NotNil(t, client, "client should not be nil")

	// Clean up
	if client != nil {
		defer client.Close()
	}

	// Verify client is functional by checking existence of test image
	ctx := context.Background()
	exists, err := client.ImageExists(ctx, "alpine:latest")
	require.NoError(t, err, "should be able to check image existence")
	assert.True(t, exists, "alpine:latest should exist (test prerequisite)")
}

func TestNewClient_DaemonNotRunning(t *testing.T) {
	t.Skip("Manual test - requires stopping Docker daemon")
	// This test would verify DaemonConnectionError when daemon is stopped
	// Manual test procedure:
	// 1. Stop Docker daemon: sudo systemctl stop docker
	// 2. Run this test
	// 3. Should get DaemonConnectionError
	// 4. Start daemon again: sudo systemctl start docker
}

// B. ImageExists Tests

func TestImageExists_ImagePresent(t *testing.T) {
	// Test with alpine:latest which should exist (test prerequisite)
	client, err := NewClient()
	require.NoError(t, err)
	defer client.Close()

	ctx := context.Background()
	exists, err := client.ImageExists(ctx, "alpine:latest")

	require.NoError(t, err, "should not return error for valid image check")
	assert.True(t, exists, "alpine:latest should exist")
}

func TestImageExists_ImageNotPresent(t *testing.T) {
	// Test with non-existent image
	// CRITICAL: This should return (false, nil) NOT an error
	client, err := NewClient()
	require.NoError(t, err)
	defer client.Close()

	ctx := context.Background()
	exists, err := client.ImageExists(ctx, "nonexistent-image-for-testing:v999")

	// Image not found is NOT an error condition
	require.NoError(t, err, "image not found should NOT return error")
	assert.False(t, exists, "non-existent image should return false")
}

func TestImageExists_InvalidImageName(t *testing.T) {
	client, err := NewClient()
	require.NoError(t, err)
	defer client.Close()

	ctx := context.Background()

	testCases := []struct {
		name      string
		imageName string
	}{
		{"empty string", ""},
		{"command injection semicolon", "alpine;rm -rf /"},
		{"command injection pipe", "alpine|whoami"},
		{"command injection ampersand", "alpine&whoami"},
		{"command injection dollar", "alpine$USER"},
		{"command injection backtick", "alpine`whoami`"},
		{"command injection paren", "alpine()"},
		{"command injection redirect", "alpine>file"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			exists, err := client.ImageExists(ctx, tc.imageName)
			assert.Error(t, err, "should return ValidationError for: %s", tc.imageName)
			assert.False(t, exists, "should return false for invalid image name")

			// Verify it's a ValidationError
			var validationErr *ValidationError
			assert.ErrorAs(t, err, &validationErr, "should be ValidationError")
		})
	}
}

// C. GetImage Tests

func TestGetImage_Success(t *testing.T) {
	// Retrieve alpine:latest and convert to OCI v1.Image
	client, err := NewClient()
	require.NoError(t, err)
	defer client.Close()

	ctx := context.Background()
	image, err := client.GetImage(ctx, "alpine:latest")

	require.NoError(t, err, "should successfully retrieve image")
	require.NotNil(t, image, "image should not be nil")

	// Verify it's a valid OCI v1.Image by checking basic methods
	configName, err := image.ConfigName()
	require.NoError(t, err, "should be able to get config name")
	assert.NotEmpty(t, configName.String(), "config name should not be empty")

	digest, err := image.Digest()
	require.NoError(t, err, "should be able to get digest")
	assert.NotEmpty(t, digest.String(), "digest should not be empty")
}

func TestGetImage_ImageNotFound(t *testing.T) {
	// Test with non-existent image
	// Should return ImageNotFoundError
	client, err := NewClient()
	require.NoError(t, err)
	defer client.Close()

	ctx := context.Background()
	image, err := client.GetImage(ctx, "nonexistent-image-for-testing:v999")

	require.Error(t, err, "should return error for non-existent image")
	assert.Nil(t, image, "image should be nil")

	// Verify it's an ImageNotFoundError
	var imageNotFoundErr *ImageNotFoundError
	assert.ErrorAs(t, err, &imageNotFoundErr, "should be ImageNotFoundError")
	assert.Contains(t, imageNotFoundErr.ImageName, "nonexistent-image-for-testing:v999")
}

func TestGetImage_InvalidImageName(t *testing.T) {
	// Test with invalid image name
	// Should return ValidationError
	client, err := NewClient()
	require.NoError(t, err)
	defer client.Close()

	ctx := context.Background()
	image, err := client.GetImage(ctx, "")

	require.Error(t, err, "should return error for empty image name")
	assert.Nil(t, image, "image should be nil")

	// Verify it's a ValidationError
	var validationErr *ValidationError
	assert.ErrorAs(t, err, &validationErr, "should be ValidationError")
}

// D. ValidateImageName Tests

func TestValidateImageName_Valid(t *testing.T) {
	client, err := NewClient()
	require.NoError(t, err)
	defer client.Close()

	validNames := []string{
		"alpine:latest",
		"myapp:v1.2.3",
		"registry.io/repo/app:tag",
		"docker.io/library/ubuntu:20.04",
		"localhost:5000/myimage:dev",
		"my_app-image:1.0",
	}

	for _, name := range validNames {
		t.Run(name, func(t *testing.T) {
			err := client.ValidateImageName(name)
			assert.NoError(t, err, "should accept valid image name: %s", name)
		})
	}
}

func TestValidateImageName_Invalid(t *testing.T) {
	client, err := NewClient()
	require.NoError(t, err)
	defer client.Close()

	invalidCases := []struct {
		name      string
		imageName string
		reason    string
	}{
		{"empty", "", "empty string"},
		{"semicolon", "alpine;whoami", "command injection - semicolon"},
		{"pipe", "alpine|whoami", "command injection - pipe"},
		{"ampersand", "alpine&whoami", "command injection - ampersand"},
		{"dollar", "alpine$USER", "command injection - dollar"},
		{"backtick", "alpine`whoami`", "command injection - backtick"},
		{"open paren", "alpine(test)", "command injection - open paren"},
		{"close paren", "alpine)test", "command injection - close paren"},
		{"less than", "alpine<file", "command injection - less than"},
		{"greater than", "alpine>file", "command injection - greater than"},
		{"backslash", "alpine\\test", "command injection - backslash"},
	}

	for _, tc := range invalidCases {
		t.Run(tc.name, func(t *testing.T) {
			err := client.ValidateImageName(tc.imageName)
			assert.Error(t, err, "should reject %s: %s", tc.reason, tc.imageName)

			// Verify it's a ValidationError
			var validationErr *ValidationError
			assert.ErrorAs(t, err, &validationErr, "should be ValidationError")
			assert.Equal(t, "imageName", validationErr.Field)
		})
	}
}

// E. Close Tests

func TestClose_Success(t *testing.T) {
	client, err := NewClient()
	require.NoError(t, err)
	require.NotNil(t, client)

	err = client.Close()
	assert.NoError(t, err, "should clean up resources successfully")
}

// F. Edge Case Tests

func TestImageExists_ContextCancellation(t *testing.T) {
	client, err := NewClient()
	require.NoError(t, err)
	defer client.Close()

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	exists, err := client.ImageExists(ctx, "alpine:latest")

	// Should respect context cancellation
	assert.Error(t, err, "should return error when context is cancelled")
	assert.False(t, exists, "should return false when context is cancelled")
}

// G. Error Type Tests (for coverage)

func TestErrorTypes_DaemonConnectionError(t *testing.T) {
	baseErr := assert.AnError
	err := &DaemonConnectionError{Cause: baseErr}

	// Test Error() method
	assert.Contains(t, err.Error(), "Docker daemon connection error")
	assert.Contains(t, err.Error(), baseErr.Error())

	// Test Unwrap() method
	assert.Equal(t, baseErr, err.Unwrap())
}

func TestErrorTypes_ImageNotFoundError(t *testing.T) {
	imageName := "test-image:v1.0"
	err := &ImageNotFoundError{ImageName: imageName}

	// Test Error() method
	assert.Contains(t, err.Error(), imageName)
	assert.Contains(t, err.Error(), "not found")
}

func TestErrorTypes_ImageConversionError(t *testing.T) {
	imageName := "test-image:v1.0"
	baseErr := assert.AnError
	err := &ImageConversionError{
		ImageName: imageName,
		Cause:     baseErr,
	}

	// Test Error() method
	assert.Contains(t, err.Error(), imageName)
	assert.Contains(t, err.Error(), "OCI format")
	assert.Contains(t, err.Error(), baseErr.Error())

	// Test Unwrap() method
	assert.Equal(t, baseErr, err.Unwrap())
}

func TestErrorTypes_ValidationError(t *testing.T) {
	field := "imageName"
	message := "cannot be empty"
	err := &ValidationError{
		Field:   field,
		Message: message,
	}

	// Test Error() method
	assert.Contains(t, err.Error(), field)
	assert.Contains(t, err.Error(), message)
	assert.Contains(t, err.Error(), "validation error")
}
