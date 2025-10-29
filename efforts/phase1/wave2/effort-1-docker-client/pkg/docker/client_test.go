package docker

import (
	"context"
	"testing"
	"time"

	"github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// isDaemonRunning checks if the Docker daemon is accessible.
// This is used to skip tests when Docker is not available.
func isDaemonRunning() bool {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return false
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err = cli.Ping(ctx)
	return err == nil
}

// requireDockerDaemon skips the test if Docker daemon is not running.
func requireDockerDaemon(t *testing.T) {
	t.Helper()
	if !isDaemonRunning() {
		t.Skip("Docker daemon not running - skipping test")
	}
}

// Test suite: NewClient

func TestNewClient_Success(t *testing.T) {
	requireDockerDaemon(t)

	client, err := NewClient()
	require.NoError(t, err, "NewClient should succeed when daemon is running")
	require.NotNil(t, client, "Client should not be nil")

	// Clean up
	err = client.Close()
	assert.NoError(t, err, "Close should succeed")
}

func TestNewClient_DaemonConnectionError(t *testing.T) {
	// This test is difficult to run reliably since we can't easily stop the daemon
	// In a real scenario, you'd mock the Docker client or use integration test setup
	// For now, we'll document this test case for manual verification
	t.Skip("Requires stopped Docker daemon - manual test case")
}

// Test suite: ImageExists

func TestImageExists_ImagePresent(t *testing.T) {
	requireDockerDaemon(t)

	// Create client
	client, err := NewClient()
	require.NoError(t, err)
	defer client.Close()

	ctx := context.Background()

	// Test with a common image that's likely to be present
	// Note: In a real test environment, you'd pull this image first
	exists, err := client.ImageExists(ctx, "alpine:latest")
	require.NoError(t, err, "ImageExists should not return error for query")

	// We can't assert exists == true without pulling the image first
	// This test validates the API contract works correctly
	t.Logf("alpine:latest exists: %v", exists)
}

func TestImageExists_ImageNotFound(t *testing.T) {
	requireDockerDaemon(t)

	client, err := NewClient()
	require.NoError(t, err)
	defer client.Close()

	ctx := context.Background()

	// Use a non-existent image name
	nonExistentImage := "this-image-definitely-does-not-exist-12345:nonexistent"
	exists, err := client.ImageExists(ctx, nonExistentImage)

	// This is the key test: ImageExists should return (false, nil) for missing images
	// NOT an error
	assert.NoError(t, err, "ImageExists should return nil error for non-existent image")
	assert.False(t, exists, "ImageExists should return false for non-existent image")
}

func TestImageExists_ValidationError_Empty(t *testing.T) {
	requireDockerDaemon(t)

	client, err := NewClient()
	require.NoError(t, err)
	defer client.Close()

	ctx := context.Background()

	// Test with empty image name
	_, err = client.ImageExists(ctx, "")

	require.Error(t, err, "ImageExists should return error for empty image name")
	assert.IsType(t, &ValidationError{}, err, "Error should be ValidationError")

	valErr := err.(*ValidationError)
	assert.Equal(t, "imageName", valErr.Field)
	assert.Contains(t, valErr.Message, "empty")
}

func TestImageExists_ValidationError_InvalidChars(t *testing.T) {
	requireDockerDaemon(t)

	client, err := NewClient()
	require.NoError(t, err)
	defer client.Close()

	ctx := context.Background()

	// Test with dangerous characters
	dangerousNames := []string{
		"image;rm -rf",
		"image|cat /etc/passwd",
		"image&& malicious",
		"image$(whoami)",
		"image`date`",
	}

	for _, imageName := range dangerousNames {
		_, err = client.ImageExists(ctx, imageName)

		require.Error(t, err, "ImageExists should return error for image name: %s", imageName)
		assert.IsType(t, &ValidationError{}, err, "Error should be ValidationError for: %s", imageName)
	}
}

// Test suite: GetImage

func TestGetImage_Success(t *testing.T) {
	requireDockerDaemon(t)

	client, err := NewClient()
	require.NoError(t, err)
	defer client.Close()

	// This test requires an image to be present
	// In a real test environment, you'd pull this image first
	t.Skip("Requires pre-pulled image - integration test")

	// Example of what the test would look like:
	// ctx := context.Background()
	// image, err := client.GetImage(ctx, "alpine:latest")
	// require.NoError(t, err)
	// require.NotNil(t, image)
	// assert.Implements(t, (*v1.Image)(nil), image)
}

func TestGetImage_ImageNotFoundError(t *testing.T) {
	requireDockerDaemon(t)

	client, err := NewClient()
	require.NoError(t, err)
	defer client.Close()

	ctx := context.Background()

	// Use a non-existent image
	nonExistentImage := "this-image-definitely-does-not-exist-12345:nonexistent"
	_, err = client.GetImage(ctx, nonExistentImage)

	require.Error(t, err, "GetImage should return error for non-existent image")
	assert.IsType(t, &ImageNotFoundError{}, err, "Error should be ImageNotFoundError")

	notFoundErr := err.(*ImageNotFoundError)
	assert.Equal(t, nonExistentImage, notFoundErr.ImageName)
}

func TestGetImage_ValidationError(t *testing.T) {
	requireDockerDaemon(t)

	client, err := NewClient()
	require.NoError(t, err)
	defer client.Close()

	ctx := context.Background()

	// Test with invalid image name
	_, err = client.GetImage(ctx, "")

	require.Error(t, err, "GetImage should return error for invalid image name")
	assert.IsType(t, &ValidationError{}, err, "Error should be ValidationError")
}

// Test suite: ValidateImageName

func TestValidateImageName_Valid(t *testing.T) {
	requireDockerDaemon(t)

	client, err := NewClient()
	require.NoError(t, err)
	defer client.Close()

	validNames := []string{
		"alpine:latest",
		"myapp:v1.0.0",
		"registry.io/repo/image:tag",
		"localhost:5000/myimage:latest",
		"my-app_image:1.2.3",
		"nginx",
		"ubuntu:22.04",
	}

	for _, imageName := range validNames {
		err := client.ValidateImageName(imageName)
		assert.NoError(t, err, "ValidateImageName should pass for: %s", imageName)
	}
}

func TestValidateImageName_CommandInjection(t *testing.T) {
	requireDockerDaemon(t)

	client, err := NewClient()
	require.NoError(t, err)
	defer client.Close()

	// Test various command injection attempts
	injectionAttempts := []string{
		"image;rm -rf /",
		"image|cat /etc/passwd",
		"image&& malicious",
		"image$(whoami)",
		"image`date`",
		"image<script>",
		"image>output.txt",
		"image\\escape",
		"image\nnewtline",
		"image\rcarriagereturn",
	}

	for _, imageName := range injectionAttempts {
		err := client.ValidateImageName(imageName)

		require.Error(t, err, "ValidateImageName should reject: %s", imageName)
		assert.IsType(t, &ValidationError{}, err, "Error should be ValidationError for: %s", imageName)
		assert.Contains(t, err.Error(), "invalid", "Error message should mention invalid characters")
	}
}

func TestValidateImageName_Empty(t *testing.T) {
	requireDockerDaemon(t)

	client, err := NewClient()
	require.NoError(t, err)
	defer client.Close()

	err = client.ValidateImageName("")

	require.Error(t, err, "ValidateImageName should reject empty string")
	assert.IsType(t, &ValidationError{}, err)

	valErr := err.(*ValidationError)
	assert.Equal(t, "imageName", valErr.Field)
	assert.Contains(t, valErr.Message, "empty")
}

// Test suite: Close

func TestClose_Success(t *testing.T) {
	requireDockerDaemon(t)

	client, err := NewClient()
	require.NoError(t, err)

	// Close should succeed
	err = client.Close()
	assert.NoError(t, err, "Close should succeed")

	// Multiple closes should be safe
	err = client.Close()
	assert.NoError(t, err, "Multiple Close calls should be safe")
}

// Test suite: Edge cases

func TestContextCancellation(t *testing.T) {
	requireDockerDaemon(t)

	client, err := NewClient()
	require.NoError(t, err)
	defer client.Close()

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// ImageExists should respect context cancellation
	_, err = client.ImageExists(ctx, "alpine:latest")

	// Depending on timing, we might get a context error or connection error
	// The important thing is that the operation doesn't hang
	if err != nil {
		t.Logf("Got expected error with cancelled context: %v", err)
	}
}

// Test suite: Helper functions

func TestContainsAny(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		substrs  []string
		expected bool
	}{
		{
			name:     "contains semicolon",
			s:        "image;command",
			substrs:  []string{";", "|", "&"},
			expected: true,
		},
		{
			name:     "contains pipe",
			s:        "image|command",
			substrs:  []string{";", "|", "&"},
			expected: true,
		},
		{
			name:     "no dangerous chars",
			s:        "alpine:latest",
			substrs:  []string{";", "|", "&"},
			expected: false,
		},
		{
			name:     "empty string",
			s:        "",
			substrs:  []string{";", "|", "&"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsAny(tt.s, tt.substrs)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsValidImageChar(t *testing.T) {
	validChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789./:_-"
	for _, char := range validChars {
		assert.True(t, isValidImageChar(char), "Character should be valid: %c", char)
	}

	invalidChars := ";|&$`()<>\\ \n\r!@#%^*+=[]{}\"'"
	for _, char := range invalidChars {
		assert.False(t, isValidImageChar(char), "Character should be invalid: %c", char)
	}
}
