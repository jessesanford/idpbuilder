// Package docker provides an interface to the local Docker daemon for image operations.
package docker

import (
	"context"

	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// DockerClient provides access to the local Docker daemon for retrieving and validating container images.
// Implementations may use the Docker Engine API or fall back to subprocess commands.
type DockerClient interface {
	// ImageExists checks if an image exists in the local Docker daemon.
	// Returns true if the image is found, false otherwise.
	//
	// Example:
	//   exists, err := client.ImageExists(ctx, "myapp:latest")
	//   if err != nil {
	//       return fmt.Errorf("checking image: %w", err)
	//   }
	//   if !exists {
	//       return fmt.Errorf("image not found")
	//   }
	ImageExists(ctx context.Context, imageName string) (bool, error)

	// GetImage retrieves an image from the Docker daemon as a v1.Image.
	// The returned image can be pushed to an OCI registry.
	//
	// Example:
	//   image, err := client.GetImage(ctx, "myapp:latest")
	//   if err != nil {
	//       return fmt.Errorf("getting image: %w", err)
	//   }
	//   // Use image with RegistryClient.Push()
	GetImage(ctx context.Context, imageName string) (v1.Image, error)

	// ValidateImageName validates an image name against the OCI specification.
	// Returns an error if the name is invalid (wrong format, invalid characters, too long).
	//
	// Example:
	//   if err := client.ValidateImageName("My App:v1"); err != nil {
	//       // Handle invalid name (spaces not allowed)
	//   }
	ValidateImageName(imageName string) error

	// Close releases the Docker daemon connection and any associated resources.
	// Must be called when the client is no longer needed.
	//
	// Example:
	//   defer client.Close()
	Close() error
}

// NewDockerClient creates a new Docker client.
// Attempts to connect to the Docker Engine API first.
// Falls back to subprocess execution (docker save) if the API is unavailable.
//
// Returns an error if both the API and subprocess methods fail.
//
// Example:
//   client, err := docker.NewDockerClient()
//   if err != nil {
//       return fmt.Errorf("creating docker client: %w", err)
//   }
//   defer client.Close()
func NewDockerClient() (DockerClient, error) {
	// Implementation will be provided in Wave 2
	panic("not implemented")
}

// ImageNotFoundError indicates that the requested image does not exist in the Docker daemon.
type ImageNotFoundError struct {
	ImageName string
}

func (e *ImageNotFoundError) Error() string {
	return "image not found: " + e.ImageName
}

// DaemonConnectionError indicates that the Docker daemon is unreachable.
type DaemonConnectionError struct {
	Endpoint string
	Cause    error
}

func (e *DaemonConnectionError) Error() string {
	if e.Cause != nil {
		return "cannot connect to Docker daemon at " + e.Endpoint + ": " + e.Cause.Error()
	}
	return "cannot connect to Docker daemon at " + e.Endpoint
}

func (e *DaemonConnectionError) Unwrap() error {
	return e.Cause
}

// InvalidImageNameError indicates that an image name violates the OCI specification.
type InvalidImageNameError struct {
	ImageName string
	Reason    string
}

func (e *InvalidImageNameError) Error() string {
	return "invalid image name '" + e.ImageName + "': " + e.Reason
}
