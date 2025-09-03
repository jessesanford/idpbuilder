package registry

import (
	"context"
	"time"

	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// Client defines the interface for registry operations
type Client interface {
	// Push pushes an image to the registry
	Push(ctx context.Context, image v1.Image, ref string, opts PushOptions) error
	
	// Pull pulls an image from the registry
	Pull(ctx context.Context, ref string, opts PullOptions) (v1.Image, error)
	
	// Catalog lists repositories in the registry
	Catalog(ctx context.Context) ([]string, error)
	
	// Tags lists tags for a repository
	Tags(ctx context.Context, repository string) ([]string, error)
	
	// Close cleans up resources used by the client
	Close() error
}

// Options contains common options for registry operations
type Options struct {
	Insecure bool
	Timeout  time.Duration
	Platform *v1.Platform
}

// PushOptions contains options for pushing images
type PushOptions struct {
	Options
}

// PullOptions contains options for pulling images
type PullOptions struct {
	Options
}

// RegistryInfo contains information about a registry
type RegistryInfo struct {
	// URL is the base URL of the registry
	URL string
	
	// Version is the registry API version
	Version string
	
	// Features contains supported registry features
	Features []string
	
	// TLSEnabled indicates if TLS is configured
	TLSEnabled bool
	
	// AuthRequired indicates if authentication is required
	AuthRequired bool
}

// ImageInfo contains metadata about an image
type ImageInfo struct {
	// Repository is the image repository name
	Repository string
	
	// Tag is the image tag
	Tag string
	
	// Digest is the image digest
	Digest string
	
	// Size is the total size in bytes
	Size int64
	
	// CreatedAt is when the image was created
	CreatedAt time.Time
	
	// Architecture is the target architecture
	Architecture string
	
	// OS is the target operating system
	OS string
}

// ClientError represents a registry client error with simplified structure
type ClientError struct {
	Code    string
	Message string
	Details map[string]interface{}
}

// Error implements the error interface
func (e *ClientError) Error() string {
	return e.Message
}