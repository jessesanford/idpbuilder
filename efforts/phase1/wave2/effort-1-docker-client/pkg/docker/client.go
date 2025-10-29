// Package docker implements the Docker client for interacting with the Docker daemon.
package docker

import (
	"context"
	"strings"

	"github.com/docker/docker/client"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
)

// dockerClient implements the Client interface for interacting with Docker daemon.
type dockerClient struct {
	cli *client.Client
}

// NewClient creates a new Docker client and verifies daemon connectivity.
//
// The client connects to the Docker daemon using environment variables
// (DOCKER_HOST, etc.) or defaults to the platform's standard socket.
//
// Returns:
//   - Client: Docker client interface implementation
//   - error: DaemonConnectionError if daemon is not reachable or not running
func NewClient() (Client, error) {
	// Create Docker client with environment-based configuration and API version negotiation
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, &DaemonConnectionError{Cause: err}
	}

	// Verify daemon is reachable by pinging it
	ctx := context.Background()
	_, err = cli.Ping(ctx)
	if err != nil {
		cli.Close() // Clean up client if ping fails
		return nil, &DaemonConnectionError{Cause: err}
	}

	return &dockerClient{cli: cli}, nil
}

// ImageExists checks if an image exists in the local Docker daemon.
//
// This method first validates the image name for security, then queries
// the Docker daemon. If the image is not found, it returns (false, nil)
// rather than an error, as a missing image is a valid query result.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - imageName: Image name in format "name:tag" (e.g., "myapp:latest")
//
// Returns:
//   - bool: true if image exists, false if not found
//   - error: DaemonConnectionError if cannot connect to daemon,
//            ValidationError if imageName is malformed
func (d *dockerClient) ImageExists(ctx context.Context, imageName string) (bool, error) {
	// Validate image name first to prevent command injection
	if err := d.ValidateImageName(imageName); err != nil {
		return false, err
	}

	// Inspect the image to check existence
	_, _, err := d.cli.ImageInspectWithRaw(ctx, imageName)
	if err != nil {
		// Image not found is NOT an error - it's a valid negative result
		if client.IsErrNotFound(err) {
			return false, nil
		}
		// Other errors (like connection issues) ARE errors
		return false, &DaemonConnectionError{Cause: err}
	}

	return true, nil
}

// GetImage retrieves an image from the Docker daemon and converts it to OCI v1.Image format.
//
// This method checks if the image exists, retrieves it from the daemon,
// and converts it to the OCI v1.Image format compatible with go-containerregistry.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - imageName: Image name in format "name:tag"
//
// Returns:
//   - v1.Image: OCI-compatible image object
//   - error: ImageNotFoundError if image doesn't exist,
//            DaemonConnectionError if cannot connect,
//            ImageConversionError if conversion fails
func (d *dockerClient) GetImage(ctx context.Context, imageName string) (v1.Image, error) {
	// Check if image exists first
	exists, err := d.ImageExists(ctx, imageName)
	if err != nil {
		// ImageExists already returns appropriate error types
		return nil, err
	}
	if !exists {
		return nil, &ImageNotFoundError{ImageName: imageName}
	}

	// Parse the image reference as a tag
	ref, err := name.ParseReference(imageName)
	if err != nil {
		return nil, &ImageConversionError{
			ImageName: imageName,
			Cause:     err,
		}
	}

	// Convert Docker image to OCI v1.Image format
	img, err := daemon.Image(ref, daemon.WithContext(ctx))
	if err != nil {
		return nil, &ImageConversionError{
			ImageName: imageName,
			Cause:     err,
		}
	}

	return img, nil
}

// ValidateImageName checks if an image name follows naming conventions and security requirements.
//
// This method performs the following validations:
//   - Checks for empty image name
//   - Prevents command injection by rejecting dangerous characters
//   - Ensures only valid characters are present (alphanumeric, dots, slashes, colons, hyphens, underscores)
//
// Parameters:
//   - imageName: Image name to validate
//
// Returns:
//   - error: ValidationError with details if invalid, nil if valid
func (d *dockerClient) ValidateImageName(imageName string) error {
	// Check for empty image name
	if imageName == "" {
		return &ValidationError{
			Field:   "imageName",
			Message: "cannot be empty",
		}
	}

	// Check for command injection attempts
	// These characters could be used for shell command injection
	dangerousChars := []string{";", "|", "&", "$", "`", "(", ")", "<", ">", "\\", "\n", "\r"}
	if containsAny(imageName, dangerousChars) {
		return &ValidationError{
			Field:   "imageName",
			Message: "contains invalid or potentially dangerous characters",
		}
	}

	// Additional validation: check for only allowed characters
	// Valid: alphanumeric, dots, slashes, colons, hyphens, underscores
	for _, char := range imageName {
		if !isValidImageChar(char) {
			return &ValidationError{
				Field:   "imageName",
				Message: "contains invalid characters (only alphanumeric, ., /, :, -, _ allowed)",
			}
		}
	}

	return nil
}

// Close cleans up Docker client resources and closes connections.
//
// This method should be called when the client is no longer needed
// to properly release system resources and network connections.
//
// Returns:
//   - error: Error if cleanup fails
func (d *dockerClient) Close() error {
	if d.cli == nil {
		return nil
	}
	return d.cli.Close()
}

// Helper functions

// containsAny checks if a string contains any of the specified substrings.
func containsAny(s string, substrs []string) bool {
	for _, substr := range substrs {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

// isValidImageChar checks if a character is valid for an image name.
// Valid characters: alphanumeric, dot, slash, colon, hyphen, underscore
func isValidImageChar(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '.' ||
		c == '/' ||
		c == ':' ||
		c == '-' ||
		c == '_'
}
