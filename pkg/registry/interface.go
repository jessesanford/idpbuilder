package registry

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
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

// Push implements Registry interface by converting io.Reader to v1.Image
func (g *GiteaRegistryAdapter) Push(ctx context.Context, imageRef string, content io.Reader) error {
	// The content io.Reader is already provided by the CLI after loading from Docker daemon
	// We need to convert it to v1.Image for the underlying GiteaClient.Push method

	// Parse the reference first
	_, err := name.ParseReference(imageRef)
	if err != nil {
		return fmt.Errorf("invalid reference %s: %w", imageRef, err)
	}

	// Try multiple approaches to get the v1.Image
	var image v1.Image

	// Approach 1: Try to load from Docker daemon (most reliable)
	image, err = crane.Pull(imageRef)
	if err != nil {
		// Approach 2: If crane.Pull fails, try using remote.Image with local Docker
		// This handles cases where the image tag might be different
		localRepo, _ := name.ParseReference(imageRef, name.WeakValidation)
		if localRepo != nil {
			image, _ = crane.Pull(localRepo.String())
		}

		// Approach 3: If still no image, check if it's already in the registry
		if image == nil {
			// As a fallback, we can try to use tarball format if the reader contains it
			// Note: This requires the io.Reader to be a tar stream
			// For now, return a clear error message
			return fmt.Errorf("image %s must be loaded in Docker daemon before pushing (docker load or docker pull required)", imageRef)
		}
	}

	// Now use the GiteaClient's Push method with the v1.Image
	pushOpts := PushOptions{
		Options: Options{
			Insecure: g.config.Insecure,
			Timeout:  5 * time.Minute,
		},
	}

	// Use the actual GiteaClient Push method
	if err := g.client.Push(ctx, image, imageRef, pushOpts); err != nil {
		return fmt.Errorf("failed to push image via GiteaClient: %w", err)
	}

	return nil
}

// List implements Registry interface by calling GiteaClient.Catalog
func (g *GiteaRegistryAdapter) List(ctx context.Context) ([]string, error) {
	// Use the GiteaClient's Catalog method which lists repositories
	repositories, err := g.client.Catalog(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list repositories: %w", err)
	}
	return repositories, nil
}

// Exists implements Registry interface by checking if repository has tags
func (g *GiteaRegistryAdapter) Exists(ctx context.Context, repository string) (bool, error) {
	// Use the GiteaClient's Tags method to check if repository exists
	// If we can list tags, the repository exists
	_, err := g.client.Tags(ctx, repository)
	if err != nil {
		// Check if it's a "not found" error
		if strings.Contains(strings.ToLower(err.Error()), "not found") ||
		   strings.Contains(err.Error(), "404") {
			return false, nil
		}
		// For other errors, return the error
		return false, fmt.Errorf("failed to check repository existence: %w", err)
	}
	// If we got tags (even empty list), repository exists
	return true, nil
}

// Delete implements Registry interface
func (g *GiteaRegistryAdapter) Delete(ctx context.Context, repository string) error {
	// Delete is a complex operation that requires:
	// 1. List all tags in the repository
	// 2. Delete each tag/manifest
	// 3. Run garbage collection if needed

	// First, check if repository exists
	exists, err := g.Exists(ctx, repository)
	if err != nil {
		return fmt.Errorf("failed to check repository existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("repository %s does not exist", repository)
	}

	// Get all tags for the repository
	tags, err := g.client.Tags(ctx, repository)
	if err != nil {
		return fmt.Errorf("failed to list tags for deletion: %w", err)
	}

	// Delete each tag using registry API
	// Note: This requires direct API calls as go-containerregistry doesn't have delete
	for _, tag := range tags {
		ref := fmt.Sprintf("%s:%s", repository, tag)
		// Parse the reference
		repo, err := name.ParseReference(ref)
		if err != nil {
			return fmt.Errorf("invalid reference %s: %w", ref, err)
		}

		// Delete using remote.Delete with our auth and transport
		deleteOpts := []remote.Option{
			remote.WithAuth(g.client.auth),
			remote.WithTransport(g.client.transport),
		}

		if err := remote.Delete(repo, deleteOpts...); err != nil {
			// Log but continue with other tags
			fmt.Fprintf(os.Stderr, "Warning: failed to delete %s: %v\n", ref, err)
		}
	}

	return nil
}