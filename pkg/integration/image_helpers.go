//go:build integration

package integration

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

// PushTestImage pushes a test image to the specified registry
func PushTestImage(registry string, image string) error {
	// Create a minimal test image
	baseImage := empty.Image

	// Add a simple layer to make it a valid image
	layer, err := tarball.LayerFromOpener(func() (io.ReadCloser, error) {
		return io.NopCloser(strings.NewReader("test content")), nil
	})
	if err != nil {
		return fmt.Errorf("failed to create test layer: %w", err)
	}

	testImage, err := mutate.AppendLayers(baseImage, layer)
	if err != nil {
		return fmt.Errorf("failed to create test image: %w", err)
	}

	// Construct full image reference
	imageRef := fmt.Sprintf("%s/%s", registry, image)

	// Parse image reference
	ref, err := name.ParseReference(imageRef)
	if err != nil {
		return fmt.Errorf("invalid image reference %s: %w", imageRef, err)
	}

	// Configure remote options for insecure registry
	options := []remote.Option{
		remote.WithTransport(&http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}),
	}

	// Push image to registry
	if err := remote.Write(ref, testImage, options...); err != nil {
		return fmt.Errorf("failed to push image %s: %w", imageRef, err)
	}

	return nil
}

// PullTestImage pulls a test image from the specified registry
func PullTestImage(registry string, image string) error {
	// Construct full image reference
	imageRef := fmt.Sprintf("%s/%s", registry, image)

	// Parse image reference
	ref, err := name.ParseReference(imageRef)
	if err != nil {
		return fmt.Errorf("invalid image reference %s: %w", imageRef, err)
	}

	// Configure remote options for insecure registry
	options := []remote.Option{
		remote.WithTransport(&http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}),
	}

	// Pull image from registry
	_, err = remote.Image(ref, options...)
	if err != nil {
		return fmt.Errorf("failed to pull image %s: %w", imageRef, err)
	}

	return nil
}

// VerifyImageInRegistry checks if an image exists in the registry
func VerifyImageInRegistry(registry string, image string) bool {
	// Construct full image reference
	imageRef := fmt.Sprintf("%s/%s", registry, image)

	// Parse image reference
	ref, err := name.ParseReference(imageRef)
	if err != nil {
		return false
	}

	// Configure remote options for insecure registry
	options := []remote.Option{
		remote.WithTransport(&http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}),
	}

	// Check if image exists by trying to get its manifest
	_, err = remote.Head(ref, options...)
	return err == nil
}

// validateImageReference checks if an image reference is valid
func validateImageReference(registry, image string) error {
	if strings.TrimSpace(registry) == "" {
		return fmt.Errorf("registry cannot be empty")
	}
	if strings.TrimSpace(image) == "" {
		return fmt.Errorf("image name cannot be empty")
	}
	imageRef := fmt.Sprintf("%s/%s", registry, image)
	_, err := name.ParseReference(imageRef)
	if err != nil {
		return fmt.Errorf("invalid image reference %s: %w", imageRef, err)
	}
	return nil
}