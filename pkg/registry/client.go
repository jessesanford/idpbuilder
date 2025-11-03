// Package registry provides OCI registry client functionality.
package registry

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/cnoe-io/idpbuilder/pkg/auth"
	"github.com/cnoe-io/idpbuilder/pkg/tls"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// Client provides registry operations
type Client interface {
	// Push pushes an image to the registry
	Push(ctx context.Context, image v1.Image, targetRef string, callback ProgressCallback) error

	// BuildImageReference constructs a full image reference
	BuildImageReference(registryURL, imageName string) (string, error)
}

// ProgressUpdate contains progress information for a layer push
type ProgressUpdate struct {
	LayerDigest  string
	BytesPushed  int64
	LayerSize    int64
	Status       string
}

// ProgressCallback is called with progress updates during push
type ProgressCallback func(update ProgressUpdate)

// NewClient creates a new registry client
func NewClient(authProvider auth.Provider, tlsProvider tls.ConfigProvider) (Client, error) {
	return &registryClient{
		authProvider: authProvider,
		tlsProvider:  tlsProvider,
	}, nil
}

// registryClient implements Client using go-containerregistry remote package
type registryClient struct {
	authProvider auth.Provider
	tlsProvider  tls.ConfigProvider
}

func (c *registryClient) Push(ctx context.Context, image v1.Image, targetRef string, callback ProgressCallback) error {
	// Parse the target reference
	ref, err := name.ParseReference(targetRef)
	if err != nil {
		return fmt.Errorf("invalid target reference %q: %w", targetRef, err)
	}

	// Get authenticator
	authenticator, err := c.authProvider.GetAuthenticator()
	if err != nil {
		return fmt.Errorf("failed to get authenticator: %w", err)
	}

	// Build push options
	opts := []remote.Option{
		remote.WithAuth(authenticator),
		remote.WithContext(ctx),
	}

	// Add TLS configuration if available
	if c.tlsProvider != nil {
		tlsConfig := c.tlsProvider.GetTLSConfig()
		if tlsConfig != nil {
			opts = append(opts, remote.WithTransport(&http.Transport{
				TLSClientConfig: tlsConfig,
			}))
		}
	}

	// Add progress callback if provided
	if callback != nil {
		// Get layers for progress tracking
		layers, err := image.Layers()
		if err == nil {
			for _, layer := range layers {
				digest, _ := layer.Digest()
				size, _ := layer.Size()
				callback(ProgressUpdate{
					LayerDigest:  digest.String(),
					BytesPushed:  0,
					LayerSize:    size,
					Status:       "pushing",
				})
			}
		}
	}

	// Push the image
	if err := remote.Write(ref, image, opts...); err != nil {
		return fmt.Errorf("failed to push image: %w", err)
	}

	// Report completion
	if callback != nil {
		layers, err := image.Layers()
		if err == nil {
			for _, layer := range layers {
				digest, _ := layer.Digest()
				size, _ := layer.Size()
				callback(ProgressUpdate{
					LayerDigest:  digest.String(),
					BytesPushed:  size,
					LayerSize:    size,
					Status:       "complete",
				})
			}
		}
	}

	return nil
}

func (c *registryClient) BuildImageReference(registryURL, imageName string) (string, error) {
	// Remove any protocol prefix
	registryURL = strings.TrimPrefix(registryURL, "http://")
	registryURL = strings.TrimPrefix(registryURL, "https://")

	// Split image name into repository and tag
	var repo, tag string
	if parts := strings.Split(imageName, ":"); len(parts) == 2 {
		repo = parts[0]
		tag = parts[1]
	} else {
		repo = imageName
		tag = "latest"
	}

	// Build the full reference
	fullRef := fmt.Sprintf("%s/%s:%s", registryURL, repo, tag)

	// Validate the reference
	if _, err := name.ParseReference(fullRef); err != nil {
		return "", fmt.Errorf("invalid image reference %q: %w", fullRef, err)
	}

	return fullRef, nil
}
