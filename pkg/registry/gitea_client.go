package registry

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// GiteaClient implements the Client interface for Gitea registry operations
type GiteaClient struct {
	auth      authn.Authenticator
	transport http.RoundTripper
	insecure  bool
	timeout   time.Duration
}

// NewGiteaClient creates a new GiteaClient
func NewGiteaClient(auth authn.Authenticator, transport http.RoundTripper, insecure bool, timeout time.Duration) *GiteaClient {
	return &GiteaClient{
		auth:      auth,
		transport: transport,
		insecure:  insecure,
		timeout:   timeout,
	}
}

// Push implements the Client interface push method
func (c *GiteaClient) Push(ctx context.Context, image v1.Image, imageRef string, opts PushOptions) error {
	// Parse the reference
	repo, err := name.ParseReference(imageRef)
	if err != nil {
		return fmt.Errorf("invalid reference %s: %w", imageRef, err)
	}

	// Configure remote options with authentication and transport
	remoteOpts := []remote.Option{
		remote.WithAuth(c.auth),
		remote.WithTransport(c.transport),
		remote.WithContext(ctx),
	}

	// Perform the push using go-containerregistry's remote.Write
	if err := remote.Write(repo, image, remoteOpts...); err != nil {
		return fmt.Errorf("failed to push image: %w", err)
	}

	return nil
}

// Catalog implements the Client interface catalog method
func (c *GiteaClient) Catalog(ctx context.Context) ([]string, error) {
	// This would typically make a registry API call to list repositories
	// For now, we'll return a placeholder implementation
	// TODO: Implement actual catalog API call to Gitea registry
	return []string{}, nil
}

// Tags implements the Client interface tags method
func (c *GiteaClient) Tags(ctx context.Context, repository string) ([]string, error) {
	// Parse the repository reference
	repo, err := name.NewRepository(repository)
	if err != nil {
		return nil, fmt.Errorf("invalid repository %s: %w", repository, err)
	}

	// Configure remote options
	remoteOpts := []remote.Option{
		remote.WithAuth(c.auth),
		remote.WithTransport(c.transport),
		remote.WithContext(ctx),
	}

	// List tags using go-containerregistry
	tags, err := remote.List(repo, remoteOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags for %s: %w", repository, err)
	}

	return tags, nil
}

// RegistryPush provides a Registry-compatible push method
// This is a helper for GiteaRegistryAdapter to use
func (c *GiteaClient) RegistryPush(ctx context.Context, imageRef string, image v1.Image) error {
	// This method wraps the existing Push with Registry-style parameters
	pushOpts := PushOptions{
		Options: Options{
			Insecure: c.insecure,
			Timeout:  c.timeout,
		},
	}
	return c.Push(ctx, image, imageRef, pushOpts)
}

// RegistryList provides a Registry-compatible list method
// Alias for Catalog to match Registry interface naming
func (c *GiteaClient) RegistryList(ctx context.Context) ([]string, error) {
	return c.Catalog(ctx)
}

// RegistryExists checks if a repository exists in the registry
func (c *GiteaClient) RegistryExists(ctx context.Context, repository string) (bool, error) {
	// Try to get tags - if successful, repository exists
	_, err := c.Tags(ctx, repository)
	if err != nil {
		// Check for not found errors
		errStr := strings.ToLower(err.Error())
		if strings.Contains(errStr, "not found") || strings.Contains(errStr, "404") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// RegistryDelete deletes a repository from the registry
func (c *GiteaClient) RegistryDelete(ctx context.Context, repository string) error {
	// Get all tags
	tags, err := c.Tags(ctx, repository)
	if err != nil {
		return fmt.Errorf("failed to list tags for deletion: %w", err)
	}

	// Delete each tag
	for _, tag := range tags {
		ref := fmt.Sprintf("%s:%s", repository, tag)
		repo, err := name.ParseReference(ref)
		if err != nil {
			continue // Skip invalid refs
		}

		deleteOpts := []remote.Option{
			remote.WithAuth(c.auth),
			remote.WithTransport(c.transport),
		}

		// Ignore individual delete errors, try to delete all
		_ = remote.Delete(repo, deleteOpts...)
	}

	return nil
}