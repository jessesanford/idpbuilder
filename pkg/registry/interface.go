// Package registry provides interfaces and types for interacting with OCI registries.
package registry

import (
	"context"

	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// Client defines operations for pushing OCI images to container registries.
type Client interface {
	// Push pushes an OCI image to the specified registry with optional progress reporting.
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeout control
	//   - image: OCI v1.Image to push
	//   - targetRef: Fully qualified image reference
	//   - progressCallback: Optional callback for progress updates (can be nil)
	//
	// Returns:
	//   - error: AuthenticationError if credentials invalid,
	//            NetworkError if registry unreachable,
	//            PushFailedError if upload fails
	Push(ctx context.Context, image v1.Image, targetRef string, progressCallback ProgressCallback) error

	// BuildImageReference constructs a fully qualified registry image reference.
	//
	// Parameters:
	//   - registryURL: Base registry URL
	//   - imageName: Image name with optional tag
	//
	// Returns:
	//   - string: Fully qualified image reference
	//   - error: ValidationError if URL or name invalid
	BuildImageReference(registryURL, imageName string) (string, error)

	// ValidateRegistry checks if the registry is reachable by pinging the /v2/ endpoint.
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeout control
	//   - registryURL: Registry base URL to validate
	//
	// Returns:
	//   - error: NetworkError if unreachable,
	//            RegistryUnavailableError if invalid response
	ValidateRegistry(ctx context.Context, registryURL string) error
}

// ProgressCallback is a function type for receiving progress updates during image push.
type ProgressCallback func(update ProgressUpdate)

// ProgressUpdate contains progress information for a single layer upload.
type ProgressUpdate struct {
	// LayerDigest is the SHA256 digest of the layer being uploaded.
	LayerDigest string

	// LayerSize is the total size of the layer in bytes.
	LayerSize int64

	// BytesPushed is the number of bytes uploaded so far.
	BytesPushed int64

	// Status indicates the current state of the layer upload.
	// Values: "uploading", "complete", "exists"
	Status string
}

// NewClient creates a new registry client with authentication and TLS configuration.
//
// Parameters:
//   - authProvider: Authentication provider (from pkg/auth)
//   - tlsConfig: TLS configuration provider (from pkg/tls)
//
// Returns:
//   - Client: Registry client interface implementation
//   - error: ValidationError if providers are invalid
func NewClient(authProvider AuthProvider, tlsConfig TLSConfigProvider) (Client, error) {
	// Implementation will be provided in Wave 2 (pkg/registry/client.go)
	panic("not implemented - interface definition only")
}

// AuthProvider is imported from pkg/auth (forward reference for clarity)
type AuthProvider interface {
	GetAuthenticator() (interface{}, error) // Returns authn.Authenticator
	ValidateCredentials() error
}

// TLSConfigProvider is imported from pkg/tls (forward reference for clarity)
type TLSConfigProvider interface {
	GetTLSConfig() interface{} // Returns *tls.Config
	IsInsecure() bool
}
