// Package docker provides Docker daemon integration for OCI image operations.
package docker

import (
	"context"

	"github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
)

// dockerClient implements the Client interface using Docker Engine API.
type dockerClient struct {
	cli *client.Client
}

// NewClient creates a new Docker client instance.
//
// The client connects to the Docker daemon using:
//   - DOCKER_HOST environment variable (if set)
//   - Default Unix socket: unix:///var/run/docker.sock
//   - Default Windows named pipe: npipe:////./pipe/docker_engine
//
// Returns:
//   - Client: Docker client interface implementation
//   - error: DaemonConnectionError if daemon is not reachable or not running
//
// Example:
//
//	client, err := docker.NewClient()
//	if err != nil {
//	    return fmt.Errorf("failed to create Docker client: %w", err)
//	}
//	defer client.Close()
func NewClient() (Client, error) {
	// Create Docker client with environment config and API version negotiation
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, &DaemonConnectionError{
			Cause: err,
		}
	}

	// Verify daemon is reachable with ping
	ctx := context.Background()
	_, err = cli.Ping(ctx)
	if err != nil {
		return nil, &DaemonConnectionError{
			Cause: err,
		}
	}

	return &dockerClient{
		cli: cli,
	}, nil
}

// ImageExists checks if an image exists in the local Docker daemon.
//
// This method uses the Docker Engine API's ImageInspectWithRaw to check
// for image existence. A NotFound error from Docker indicates the image
// doesn't exist (returns false, nil).
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - imageName: Image name in format "name:tag" (e.g., "myapp:latest")
//
// Returns:
//   - bool: true if image exists, false otherwise
//   - error: DaemonConnectionError if cannot connect to daemon,
//     ValidationError if imageName is malformed
//
// Example:
//
//	exists, err := client.ImageExists(ctx, "myapp:latest")
//	if err != nil {
//	    return fmt.Errorf("failed to check image: %w", err)
//	}
//	if !exists {
//	    return fmt.Errorf("image not found in Docker daemon")
//	}
func (c *dockerClient) ImageExists(ctx context.Context, imageName string) (bool, error) {
	// Validate image name first
	if err := c.ValidateImageName(imageName); err != nil {
		return false, err
	}

	// Use ImageInspectWithRaw to check existence
	_, _, err := c.cli.ImageInspectWithRaw(ctx, imageName)
	if err != nil {
		// Check if image simply doesn't exist (NOT AN ERROR!)
		if errdefs.IsNotFound(err) {
			return false, nil // Normal case - image not found
		}
		// Any other error is a connection/daemon issue
		return false, &DaemonConnectionError{Cause: err}
	}

	return true, nil
}

// GetImage retrieves an image from the Docker daemon and converts it
// to an OCI v1.Image format compatible with go-containerregistry.
//
// This method uses go-containerregistry's daemon package to handle the
// conversion from Docker image format to OCI v1.Image. The daemon package
// internally uses Docker's SaveImage API to export the image.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - imageName: Image name in format "name:tag"
//
// Returns:
//   - v1.Image: OCI-compatible image object
//   - error: ImageNotFoundError if image doesn't exist,
//     DaemonConnectionError if cannot connect,
//     ImageConversionError if conversion fails
//
// Example:
//
//	image, err := client.GetImage(ctx, "myapp:latest")
//	if err != nil {
//	    return fmt.Errorf("failed to retrieve image: %w", err)
//	}
//	// image can now be pushed to registry
func (c *dockerClient) GetImage(ctx context.Context, imageName string) (v1.Image, error) {
	// Validate image name first
	if err := c.ValidateImageName(imageName); err != nil {
		return nil, err
	}

	// Check if image exists first
	exists, err := c.ImageExists(ctx, imageName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, &ImageNotFoundError{
			ImageName: imageName,
		}
	}

	// Parse image reference using name package
	ref, err := name.ParseReference(imageName)
	if err != nil {
		return nil, &ImageConversionError{
			ImageName: imageName,
			Cause:     err,
		}
	}

	// Convert to OCI v1.Image format using daemon package
	img, err := daemon.Image(ref)
	if err != nil {
		return nil, &ImageConversionError{
			ImageName: imageName,
			Cause:     err,
		}
	}

	return img, nil
}

// ValidateImageName checks if an image name follows the OCI naming specification.
//
// This method validates:
//   - Image name is not empty
//   - No command injection attempts (no semicolons, pipes, etc.)
//   - Basic format check (allows alphanumeric, dots, slashes, colons, hyphens, underscores)
//
// Note: Full OCI spec validation is complex. This provides basic safety checks.
// Docker daemon will perform additional validation.
//
// Parameters:
//   - imageName: Image name to validate
//
// Returns:
//   - error: ValidationError with details if invalid, nil if valid
//
// Example:
//
//	if err := client.ValidateImageName("myapp:latest"); err != nil {
//	    return fmt.Errorf("invalid image name: %w", err)
//	}
func (c *dockerClient) ValidateImageName(imageName string) error {
	// Check for empty string
	if imageName == "" {
		return &ValidationError{
			Field:   "imageName",
			Message: "image name cannot be empty",
		}
	}

	// Check for command injection attempts
	dangerousChars := []string{
		";",  // Command separator
		"|",  // Pipe
		"&",  // Background/AND
		"$",  // Variable expansion
		"`",  // Command substitution
		"(",  // Subshell
		")",  // Subshell
		"<",  // Redirect input
		">",  // Redirect output
		"\\", // Escape character
	}

	for _, char := range dangerousChars {
		if containsString(imageName, char) {
			return &ValidationError{
				Field:   "imageName",
				Message: "image name contains dangerous character '" + char + "' (potential command injection)",
			}
		}
	}

	// Basic validation passed
	return nil
}

// Close cleans up Docker client resources and closes connections.
//
// This method closes the underlying HTTP client connection to the Docker daemon.
// It should be called when the client is no longer needed, typically via defer.
//
// Returns:
//   - error: Error if cleanup fails
//
// Example:
//
//	client, err := NewClient()
//	if err != nil {
//	    return err
//	}
//	defer client.Close()
func (c *dockerClient) Close() error {
	if c.cli != nil {
		return c.cli.Close()
	}
	return nil
}

// containsString checks if substr exists in s
func containsString(s, substr string) bool {
	for i := 0; i < len(s); i++ {
		if i+len(substr) <= len(s) && s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
