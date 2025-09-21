package registry

import (
	"context"
	"io"
	"time"

	"github.com/google/go-containerregistry/pkg/v1"
)

// Registry interface defines methods for OCI registry operations
type Registry interface {
	// Push pushes an image from an io.Reader to the registry
	Push(ctx context.Context, imageRef string, content io.Reader) error

	// List returns all repositories in the registry
	List(ctx context.Context) ([]string, error)

	// Exists checks if a repository exists in the registry
	Exists(ctx context.Context, repository string) (bool, error)

	// Delete removes a repository from the registry
	Delete(ctx context.Context, repository string) error
}

// Client interface defines the underlying client operations
type Client interface {
	Push(ctx context.Context, image v1.Image, imageRef string, opts PushOptions) error
	Catalog(ctx context.Context) ([]string, error)
	Tags(ctx context.Context, repository string) ([]string, error)
}

// Options contains common configuration options
type Options struct {
	Insecure bool
	Timeout  time.Duration
}

// PushOptions contains push-specific options
type PushOptions struct {
	Options
}

// GiteaConfig holds configuration for Gitea registry adapter
type GiteaConfig struct {
	Insecure bool
	Timeout  time.Duration
}

// GiteaClient is defined in gitea_client.go

// GiteaRegistryAdapter adapts GiteaClient to implement Registry interface
type GiteaRegistryAdapter struct {
	client *GiteaClient
	config GiteaConfig
}

// NewGiteaRegistryAdapter creates a new adapter instance
func NewGiteaRegistryAdapter(client *GiteaClient, config GiteaConfig) *GiteaRegistryAdapter {
	return &GiteaRegistryAdapter{
		client: client,
		config: config,
	}
}

// Push implements Registry interface by delegating to GiteaClient
func (g *GiteaRegistryAdapter) Push(ctx context.Context, imageRef string, content io.Reader) error {
	// Delegate directly to the GiteaClient's Push method which now implements Registry interface
	return g.client.Push(ctx, imageRef, content)
}

// List implements Registry interface by delegating to GiteaClient
func (g *GiteaRegistryAdapter) List(ctx context.Context) ([]string, error) {
	// Delegate directly to the GiteaClient's List method which now implements Registry interface
	return g.client.List(ctx)
}

// Exists implements Registry interface by delegating to GiteaClient
func (g *GiteaRegistryAdapter) Exists(ctx context.Context, repository string) (bool, error) {
	// Delegate directly to the GiteaClient's Exists method which now implements Registry interface
	return g.client.Exists(ctx, repository)
}

// Delete implements Registry interface by delegating to GiteaClient
func (g *GiteaRegistryAdapter) Delete(ctx context.Context, repository string) error {
	// Delegate directly to the GiteaClient's Delete method which now implements Registry interface
	return g.client.Delete(ctx, repository)
}