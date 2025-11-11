// Package registry provides an interface to OCI-compliant container registries.
package registry

import (
	"context"

	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// LayerStatus represents the current status of a layer during push.
type LayerStatus int

const (
	// LayerWaiting indicates the layer is queued but not yet uploading.
	LayerWaiting LayerStatus = iota
	// LayerUploading indicates the layer is currently being uploaded.
	LayerUploading
	// LayerComplete indicates the layer upload completed successfully.
	LayerComplete
	// LayerFailed indicates the layer upload failed.
	LayerFailed
)

// String returns a string representation of the LayerStatus.
func (s LayerStatus) String() string {
	switch s {
	case LayerWaiting:
		return "Waiting"
	case LayerUploading:
		return "Uploading"
	case LayerComplete:
		return "Complete"
	case LayerFailed:
		return "Failed"
	default:
		return "Unknown"
	}
}

// ProgressUpdate represents the progress of a single layer upload.
type ProgressUpdate struct {
	// LayerDigest is the SHA256 digest of the layer (e.g., "sha256:abc123...")
	LayerDigest string

	// LayerSize is the total size of the layer in bytes.
	LayerSize int64

	// BytesUploaded is the number of bytes uploaded so far.
	BytesUploaded int64

	// Status is the current status of this layer.
	Status LayerStatus
}

// ProgressCallback is invoked during image push to report upload progress.
// The callback is called once for each layer as its status changes.
// Implementations should not block; long-running operations should be asynchronous.
type ProgressCallback func(update ProgressUpdate)

// RegistryClient handles operations with OCI-compliant container registries.
type RegistryClient interface {
	// Push uploads an image to a registry with progress reporting.
	//
	// The targetRef must be a full registry reference in the format:
	//   "registry.example.com:port/namespace/repository:tag"
	//
	// The progress callback is invoked for each layer as it uploads.
	// Pass nil for progress if callbacks are not needed.
	//
	// Example:
	//   err := client.Push(ctx, image, "gitea.example.com:8443/giteaadmin/myapp:latest", func(update ProgressUpdate) {
	//       fmt.Printf("Layer %s: %d/%d bytes (%s)\n",
	//           update.LayerDigest[:12], update.BytesUploaded, update.LayerSize, update.Status)
	//   })
	Push(ctx context.Context, image v1.Image, targetRef string, progress ProgressCallback) error

	// BuildImageReference constructs a full registry reference from a registry URL and image name.
	//
	// The registryURL should be the base URL of the registry (e.g., "https://gitea.example.com:8443").
	// The imageName should be in the format "repository:tag" or "namespace/repository:tag".
	//
	// Example:
	//   ref, err := client.BuildImageReference("https://gitea.example.com:8443", "myapp:v1")
	//   // Returns: "gitea.example.com:8443/giteaadmin/myapp:v1", nil
	BuildImageReference(registryURL, imageName string) (string, error)

	// ValidateRegistry performs connectivity and OCI compliance checks on a registry.
	// Returns an error if the registry is unreachable or does not support the OCI Distribution spec.
	//
	// Example:
	//   if err := client.ValidateRegistry(ctx, "https://gitea.example.com:8443"); err != nil {
	//       return fmt.Errorf("registry validation failed: %w", err)
	//   }
	ValidateRegistry(ctx context.Context, registryURL string) error
}

// NewRegistryClient creates a new OCI registry client.
// The auth and tls parameters customize the connection behavior.
//
// Example:
//   auth, _ := auth.NewAuthProvider("username", "password")
//   tls, _ := tls.NewTLSProvider(true) // insecure mode for local dev
//   client, err := registry.NewRegistryClient(auth, tls)
func NewRegistryClient(auth AuthProvider, tls TLSProvider) (RegistryClient, error) {
	// Implementation will be provided in Wave 2
	panic("not implemented")
}

// AuthProvider is imported from pkg/auth (defined below)
type AuthProvider interface{}

// TLSProvider is imported from pkg/tls (defined below)
type TLSProvider interface{}

// RegistryAuthError indicates that authentication to the registry failed.
type RegistryAuthError struct {
	Registry string
	Cause    error
}

func (e *RegistryAuthError) Error() string {
	if e.Cause != nil {
		return "authentication failed for registry " + e.Registry + ": " + e.Cause.Error()
	}
	return "authentication failed for registry " + e.Registry
}

func (e *RegistryAuthError) Unwrap() error {
	return e.Cause
}

// RegistryConnectionError indicates that the registry is unreachable.
type RegistryConnectionError struct {
	Registry string
	Cause    error
}

func (e *RegistryConnectionError) Error() string {
	if e.Cause != nil {
		return "cannot connect to registry " + e.Registry + ": " + e.Cause.Error()
	}
	return "cannot connect to registry " + e.Registry
}

func (e *RegistryConnectionError) Unwrap() error {
	return e.Cause
}

// LayerPushError indicates that uploading a specific layer failed.
type LayerPushError struct {
	LayerDigest string
	Cause       error
}

func (e *LayerPushError) Error() string {
	if e.Cause != nil {
		return "failed to push layer " + e.LayerDigest + ": " + e.Cause.Error()
	}
	return "failed to push layer " + e.LayerDigest
}

func (e *LayerPushError) Unwrap() error {
	return e.Cause
}
