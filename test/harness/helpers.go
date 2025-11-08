package harness

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// VerifyImageInRegistry checks if an image exists in the Gitea registry.
// Returns true if the image manifest can be retrieved, false if not found,
// and an error for other failures (authentication, network, etc.).
func VerifyImageInRegistry(ctx context.Context, registryURL, imageRef, username, password string) (bool, error) {
	// Parse the image reference
	ref, err := name.ParseReference(fmt.Sprintf("%s/%s", registryURL, imageRef))
	if err != nil {
		return false, fmt.Errorf("failed to parse image reference: %w", err)
	}

	// Create authentication configuration
	auth := &authn.Basic{
		Username: username,
		Password: password,
	}

	// Create transport with insecure TLS (for testing)
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	// Try to get the image descriptor
	_, err = remote.Get(ref,
		remote.WithAuth(auth),
		remote.WithTransport(transport),
		remote.WithContext(ctx),
	)

	if err != nil {
		// Check if it's a "not found" error
		if contains(err.Error(), "404") || contains(err.Error(), "NAME_UNKNOWN") || contains(err.Error(), "MANIFEST_UNKNOWN") {
			return false, nil
		}
		// Other errors (authentication, network, etc.)
		return false, fmt.Errorf("failed to verify image: %w", err)
	}

	// Image exists
	return true, nil
}

// GetImageDigest retrieves the digest of a pushed image from the registry.
// Returns the digest string (e.g., "sha256:abc123...") or an error if the
// image doesn't exist or cannot be accessed.
func GetImageDigest(ctx context.Context, registryURL, imageRef, username, password string) (string, error) {
	// Parse the image reference
	ref, err := name.ParseReference(fmt.Sprintf("%s/%s", registryURL, imageRef))
	if err != nil {
		return "", fmt.Errorf("failed to parse image reference: %w", err)
	}

	// Create authentication configuration
	auth := &authn.Basic{
		Username: username,
		Password: password,
	}

	// Create transport with insecure TLS (for testing)
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	// Get the image descriptor
	desc, err := remote.Get(ref,
		remote.WithAuth(auth),
		remote.WithTransport(transport),
		remote.WithContext(ctx),
	)
	if err != nil {
		return "", fmt.Errorf("failed to get image descriptor: %w", err)
	}

	// Extract and return the digest
	return desc.Digest.String(), nil
}
