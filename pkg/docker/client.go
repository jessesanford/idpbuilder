// Package docker provides Docker daemon client functionality.
package docker

import (
	"context"
	"fmt"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
)

// Client provides Docker daemon operations
type Client interface {
	// GetImage retrieves an image from the Docker daemon
	GetImage(ctx context.Context, imageName string) (v1.Image, error)

	// Close releases Docker client resources
	Close() error
}

// NewClient creates a new Docker client
func NewClient() (Client, error) {
	return &dockerClient{}, nil
}

// dockerClient implements Client using go-containerregistry daemon package
type dockerClient struct{}

func (c *dockerClient) GetImage(ctx context.Context, imageName string) (v1.Image, error) {
	// Parse the image reference
	ref, err := name.ParseReference(imageName)
	if err != nil {
		return nil, fmt.Errorf("invalid image name %q: %w", imageName, err)
	}

	// Get the image from the Docker daemon
	img, err := daemon.Image(ref)
	if err != nil {
		return nil, fmt.Errorf("failed to get image from Docker daemon: %w", err)
	}

	return img, nil
}

func (c *dockerClient) Close() error {
	// go-containerregistry daemon client doesn't require explicit cleanup
	return nil
}
