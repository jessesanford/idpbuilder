// Package docker provides interfaces and types for interacting with Docker daemon.
package docker

import (
	"context"

	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// Client defines operations for interacting with the Docker daemon
// to retrieve and validate OCI images stored locally.
type Client interface {
	// ImageExists checks if an image exists in the local Docker daemon.
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeout control
	//   - imageName: Image name in format "name:tag" (e.g., "myapp:latest")
	//
	// Returns:
	//   - bool: true if image exists, false otherwise
	//   - error: DaemonConnectionError if cannot connect to daemon,
	//            ValidationError if imageName is malformed
	ImageExists(ctx context.Context, imageName string) (bool, error)

	// GetImage retrieves an image from the Docker daemon and converts it
	// to an OCI v1.Image format compatible with go-containerregistry.
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeout control
	//   - imageName: Image name in format "name:tag"
	//
	// Returns:
	//   - v1.Image: OCI-compatible image object
	//   - error: ImageNotFoundError if image doesn't exist,
	//            DaemonConnectionError if cannot connect,
	//            ImageConversionError if tar conversion fails
	GetImage(ctx context.Context, imageName string) (v1.Image, error)

	// ValidateImageName checks if an image name follows the OCI naming specification.
	//
	// Parameters:
	//   - imageName: Image name to validate
	//
	// Returns:
	//   - error: ValidationError with details if invalid, nil if valid
	ValidateImageName(imageName string) error

	// Close cleans up Docker client resources and closes connections.
	//
	// Returns:
	//   - error: Error if cleanup fails
	Close() error
}
