// Package docker provides Docker daemon client functionality.
// This is a Phase 1 stub interface for Phase 2 development.
package docker

import (
	"context"

	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// Client provides Docker daemon operations
type Client interface {
	// GetImage retrieves an image from the Docker daemon
	GetImage(ctx context.Context, imageName string) (v1.Image, error)

	// Close releases Docker client resources
	Close() error
}

// NewClient creates a new Docker client (stub for Phase 1)
func NewClient() (Client, error) {
	// Phase 1 implementation would return actual Docker client
	return &stubClient{}, nil
}

// stubClient is a minimal stub for planning purposes
type stubClient struct{}

func (c *stubClient) GetImage(ctx context.Context, imageName string) (v1.Image, error) {
	// Phase 1 would implement actual Docker daemon communication
	return nil, nil
}

func (c *stubClient) Close() error {
	return nil
}
