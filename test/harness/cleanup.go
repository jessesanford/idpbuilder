package harness

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

// CleanupTestEnvironment ensures all resources are cleaned up properly.
// This includes terminating containers, closing connections, and removing
// test images. Errors are logged but don't fail the cleanup process.
func (env *TestEnvironment) CleanupTestEnvironment() error {
	if env == nil {
		return nil
	}

	var errs []error

	// Close Docker client
	if env.DockerClient != nil {
		if err := env.DockerClient.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close Docker client: %w", err))
		}
	}

	// Terminate Gitea container
	if env.GiteaContainer != nil {
		ctx := context.Background()
		if err := env.GiteaContainer.Terminate(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to terminate Gitea container: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("cleanup encountered %d error(s): %v", len(errs), errs)
	}

	return nil
}

// RemoveTestImages removes all test images from the Docker daemon.
// This helps prevent accumulation of test artifacts. Images are identified
// by having "test" in their repository name.
func RemoveTestImages(ctx context.Context, dockerClient *client.Client) error {
	if dockerClient == nil {
		return fmt.Errorf("Docker client is nil")
	}

	// List all images
	images, err := dockerClient.ImageList(ctx, image.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list images: %w", err)
	}

	// Track errors but continue removing other images
	var errs []error

	// Remove images with "test" in their repository tags
	for _, image := range images {
		for _, tag := range image.RepoTags {
			// Only remove test images (containing "test" in name)
			if contains(tag, "test") || contains(tag, "Test") {
				_, err := dockerClient.ImageRemove(ctx, image.ID, types.ImageRemoveOptions{
					Force:         true,
					PruneChildren: true,
				})
				if err != nil {
					// Don't fail on "image not found" errors
					if !contains(err.Error(), "No such image") {
						errs = append(errs, fmt.Errorf("failed to remove image %s: %w", tag, err))
					}
				}
				break // Only need to remove once per image
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("encountered %d error(s) removing images: %v", len(errs), errs)
	}

	return nil
}

// WaitForGiteaReady polls the Gitea instance until it responds successfully
// or the timeout is reached. This ensures Gitea is fully initialized before
// running tests that depend on it.
func WaitForGiteaReady(ctx context.Context, registryURL string) error {
	// Try to connect to Gitea web interface
	// Use the web port (typically exposed alongside registry)
	// For simplicity, we'll check if we can make an HTTP request
	maxRetries := 30
	retryInterval := 1 * time.Second

	var lastErr error
	for i := 0; i < maxRetries; i++ {
		// Create a request with timeout
		checkCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		req, err := http.NewRequestWithContext(checkCtx, "GET", "http://"+registryURL, nil)
		if err != nil {
			cancel()
			lastErr = fmt.Errorf("failed to create request: %w", err)
			time.Sleep(retryInterval)
			continue
		}

		// Try to connect
		resp, err := http.DefaultClient.Do(req)
		cancel()

		if err == nil {
			resp.Body.Close()
			// Any response means Gitea is up
			return nil
		}

		lastErr = err
		time.Sleep(retryInterval)
	}

	return fmt.Errorf("Gitea not ready after %d retries: %w", maxRetries, lastErr)
}

// contains checks if a string contains a substring (case-sensitive helper)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && hasSubstring(s, substr))
}

// hasSubstring performs the actual substring search
func hasSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
