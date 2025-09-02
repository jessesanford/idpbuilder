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

// PushOptions contains options for pushing images
type PushOptions struct {
	// Insecure skips TLS verification (R307 - feature flag)
	Insecure bool
	
	// Progress callback for operation progress
	Progress ProgressFunc
	
	// Timeout for the operation
	Timeout time.Duration
	
	// MaxRetries for failed operations
	MaxRetries int
	
	// RetryDelay between retry attempts
	RetryDelay time.Duration
	
	// Platform specifies the target platform for multi-arch images
	Platform *v1.Platform
}

// PullOptions contains options for pulling images
type PullOptions struct {
	// Platform specifies the target platform
	Platform *v1.Platform
	
	// Insecure skips TLS verification (R307 - feature flag)
	Insecure bool
	
	// Timeout for the operation
	Timeout time.Duration
	
	// MaxRetries for failed operations
	MaxRetries int
	
	// RetryDelay between retry attempts
	RetryDelay time.Duration
}

// CatalogOptions contains options for listing repositories
type CatalogOptions struct {
	// MaxEntries limits the number of repositories to return
	MaxEntries int
	
	// LastEntry for pagination support
	LastEntry string
	
	// Timeout for the operation
	Timeout time.Duration
}

// TagsOptions contains options for listing tags
type TagsOptions struct {
	// MaxEntries limits the number of tags to return
	MaxEntries int
	
	// LastEntry for pagination support
	LastEntry string
	
	// Timeout for the operation
	Timeout time.Duration
}

// ProgressFunc reports progress during operations
// current: bytes processed, total: total bytes (-1 if unknown)
type ProgressFunc func(current, total int64)

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

// ErrorType represents different categories of registry errors
type ErrorType int

const (
	// ErrorUnknown represents an unspecified error
	ErrorUnknown ErrorType = iota
	
	// ErrorNetwork represents network-related errors
	ErrorNetwork
	
	// ErrorAuthentication represents authentication failures
	ErrorAuthentication
	
	// ErrorAuthorization represents authorization failures
	ErrorAuthorization
	
	// ErrorNotFound represents resource not found errors
	ErrorNotFound
	
	// ErrorConflict represents resource conflict errors
	ErrorConflict
	
	// ErrorInvalidReference represents invalid image reference errors
	ErrorInvalidReference
	
	// ErrorCertificate represents TLS certificate errors
	ErrorCertificate
	
	// ErrorTimeout represents operation timeout errors
	ErrorTimeout
)

// RegistryError represents a registry-specific error
type RegistryError struct {
	// Type categorizes the error
	Type ErrorType
	
	// Registry is the registry URL where the error occurred
	Registry string
	
	// Operation is the operation that failed
	Operation string
	
	// Message is the error message
	Message string
	
	// Underlying is the original error
	Underlying error
	
	// Retryable indicates if the operation can be retried
	Retryable bool
	
	// StatusCode is the HTTP status code if applicable
	StatusCode int
}

// Error implements the error interface
func (e *RegistryError) Error() string {
	if e.Underlying != nil {
		return e.Message + ": " + e.Underlying.Error()
	}
	return e.Message
}

// Unwrap returns the underlying error for error unwrapping
func (e *RegistryError) Unwrap() error {
	return e.Underlying
}

// IsRetryable returns whether the error represents a retryable condition
func (e *RegistryError) IsRetryable() bool {
	return e.Retryable
}