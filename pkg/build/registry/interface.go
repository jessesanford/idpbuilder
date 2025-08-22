// Package registry defines the interface for container registry operations  
package registry

import (
	"context"
)

// Registry defines the interface for container registry operations
type Registry interface {
	// Push pushes an image to the registry
	Push(ctx context.Context, imageRef string) error
	
	// Tag tags an image in the registry
	Tag(ctx context.Context, imageID string, tag string) error
	
	// Exists checks if an image exists in the registry  
	Exists(ctx context.Context, imageRef string) (bool, error)
}

// RegistryConfig holds configuration for registry instances
type RegistryConfig struct {
	// URL is the registry base URL
	URL string
	
	// Namespace is the registry namespace
	Namespace string
	
	// Username for authentication
	Username string
	
	// Password for authentication
	Password string
	
	// InsecureSkipTLSVerify skips TLS verification
	InsecureSkipTLSVerify bool
}