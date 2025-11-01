// Package registry provides OCI registry client functionality.
// This is a Phase 1 stub interface for Phase 2 development.
package registry

import (
	"context"

	"github.com/cnoe-io/idpbuilder/pkg/auth"
	"github.com/cnoe-io/idpbuilder/pkg/tls"
	v1 "github.com/google/go-containerregistry/pkg/v1"
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

// NewClient creates a new registry client (stub for Phase 1)
func NewClient(authProvider auth.Provider, tlsProvider tls.ConfigProvider) (Client, error) {
	// Phase 1 would implement actual registry client
	return &stubClient{}, nil
}

// stubClient is a minimal stub for planning purposes
type stubClient struct{}

func (c *stubClient) Push(ctx context.Context, image v1.Image, targetRef string, callback ProgressCallback) error {
	// Phase 1 would implement actual push logic
	if callback != nil {
		callback(ProgressUpdate{
			LayerDigest:  "sha256:stub",
			BytesPushed:  100,
			LayerSize:    100,
			Status:       "complete",
		})
	}
	return nil
}

func (c *stubClient) BuildImageReference(registryURL, imageName string) (string, error) {
	// Phase 1 would implement reference building logic
	return registryURL + "/" + imageName, nil
}
