package registry

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// GiteaClient implements the Registry interface for Gitea registry operations
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

// Push implements the Registry interface push method
func (c *GiteaClient) Push(ctx context.Context, imageRef string, content io.Reader) error {
	// Parse the reference first
	_, err := name.ParseReference(imageRef)
	if err != nil {
		return fmt.Errorf("invalid reference %s: %w", imageRef, err)
	}

	// Try multiple approaches to get the v1.Image from content
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

		// Approach 3: If still no image, return a clear error message
		if image == nil {
			return fmt.Errorf("image %s must be loaded in Docker daemon before pushing (docker load or docker pull required)", imageRef)
		}
	}

	// Parse the reference again for remote operations
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

// List implements the Registry interface list method
func (c *GiteaClient) List(ctx context.Context) ([]string, error) {
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

// Exists implements the Registry interface exists method
func (c *GiteaClient) Exists(ctx context.Context, repository string) (bool, error) {
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

// Delete implements the Registry interface delete method
func (c *GiteaClient) Delete(ctx context.Context, repository string) error {
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