package builder

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// DefaultRegistryClient provides a default implementation of RegistryClient.
type DefaultRegistryClient struct {
	config    RegistryConfig
	transport http.RoundTripper
	options   []remote.Option
}

// NewDefaultRegistryClient creates a new default registry client.
func NewDefaultRegistryClient(config RegistryConfig) *DefaultRegistryClient {
	client := &DefaultRegistryClient{
		config:  config,
		options: make([]remote.Option, 0),
	}

	// Configure authentication
	client.configureAuth()

	// Configure transport
	client.configureTransport()

	return client
}

// configureAuth sets up authentication options.
func (drc *DefaultRegistryClient) configureAuth() {
	var auth authn.Authenticator

	// Determine authentication method
	if drc.config.Token != "" {
		auth = &authn.Bearer{Token: drc.config.Token}
	} else if drc.config.RegistryToken != "" {
		auth = &authn.Bearer{Token: drc.config.RegistryToken}
	} else if drc.config.Username != "" && drc.config.Password != "" {
		auth = &authn.Basic{
			Username: drc.config.Username,
			Password: drc.config.Password,
		}
	} else {
		auth = authn.Anonymous
	}

	drc.options = append(drc.options, remote.WithAuth(auth))
}

// configureTransport sets up HTTP transport options.
func (drc *DefaultRegistryClient) configureTransport() {
	transport := http.DefaultTransport.(*http.Transport).Clone()

	// Configure TLS settings
	transport.TLSClientConfig.InsecureSkipVerify = drc.config.SkipTLSVerify

	// Configure timeouts
	transport.ResponseHeaderTimeout = drc.config.Timeout
	transport.IdleConnTimeout = drc.config.Timeout

	drc.transport = &retryTransport{
		base:       transport,
		retryCount: drc.config.RetryCount,
		retryDelay: drc.config.RetryDelay,
	}

	drc.options = append(drc.options, remote.WithTransport(drc.transport))
}

// GetImage retrieves an image from the registry.
func (drc *DefaultRegistryClient) GetImage(ref string) (v1.Image, error) {
	nameRef, err := parseReference(ref)
	if err != nil {
		return nil, fmt.Errorf("failed to parse reference %s: %w", ref, err)
	}

	image, err := remote.Image(nameRef, drc.options...)
	if err != nil {
		return nil, fmt.Errorf("failed to get image %s: %w", ref, err)
	}

	return image, nil
}

// PushImage pushes an image to the registry.
func (drc *DefaultRegistryClient) PushImage(ref string, image v1.Image) error {
	nameRef, err := parseReference(ref)
	if err != nil {
		return fmt.Errorf("failed to parse reference %s: %w", ref, err)
	}

	err = remote.Write(nameRef, image, drc.options...)
	if err != nil {
		return fmt.Errorf("failed to push image %s: %w", ref, err)
	}

	return nil
}

// CheckImageExists checks if an image exists in the registry.
func (drc *DefaultRegistryClient) CheckImageExists(ref string) (bool, error) {
	nameRef, err := parseReference(ref)
	if err != nil {
		return false, fmt.Errorf("failed to parse reference %s: %w", ref, err)
	}

	_, err = remote.Head(nameRef, drc.options...)
	if err != nil {
		// If it's a 404, the image doesn't exist
		if isNotFoundError(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if image exists %s: %w", ref, err)
	}

	return true, nil
}

// retryTransport wraps an HTTP transport with retry logic.
type retryTransport struct {
	base       http.RoundTripper
	retryCount int
	retryDelay time.Duration
}

// RoundTrip implements http.RoundTripper with retry logic.
func (rt *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var lastErr error

	for attempt := 0; attempt <= rt.retryCount; attempt++ {
		resp, err := rt.base.RoundTrip(req)
		if err == nil {
			return resp, nil
		}

		lastErr = err

		// Don't retry on the last attempt
		if attempt == rt.retryCount {
			break
		}

		// Only retry on certain errors
		if !shouldRetry(err, resp) {
			break
		}

		// Wait before retrying
		if rt.retryDelay > 0 {
			time.Sleep(rt.retryDelay)
		}
	}

	return nil, lastErr
}

// shouldRetry determines if a request should be retried based on the error or response.
func shouldRetry(err error, resp *http.Response) bool {
	// Always retry on network errors
	if err != nil {
		return true
	}

	// Retry on server errors (5xx) but not client errors (4xx)
	if resp != nil && resp.StatusCode >= 500 {
		return true
	}

	return false
}

// parseReference parses a reference string and returns a name.Reference.
func parseReference(ref string) (name.Reference, error) {
	// Parse as a tag first
	if tag, err := name.NewTag(ref); err == nil {
		return tag, nil
	}
	
	// Try as a digest
	if digest, err := name.NewDigest(ref); err == nil {
		return digest, nil
	}
	
	// If neither works, add a default tag and parse as tag
	if !strings.Contains(ref, ":") && !strings.Contains(ref, "@") {
		ref = ref + ":latest"
	}
	
	return name.NewTag(ref)
}

// isNotFoundError checks if an error indicates that a resource was not found.
func isNotFoundError(err error) bool {
	// This is a simplified check. In practice, you'd check for specific error types
	// or HTTP status codes
	return err != nil && (err.Error() == "not found" || err.Error() == "404")
}

// MockRegistryClient provides a mock implementation for testing.
type MockRegistryClient struct {
	images map[string]v1.Image
	pushes []string
	errors map[string]error
}

// NewMockRegistryClient creates a new mock registry client.
func NewMockRegistryClient() *MockRegistryClient {
	return &MockRegistryClient{
		images: make(map[string]v1.Image),
		pushes: make([]string, 0),
		errors: make(map[string]error),
	}
}

// SetImage configures the mock to return a specific image for a reference.
func (mrc *MockRegistryClient) SetImage(ref string, image v1.Image) {
	mrc.images[ref] = image
}

// SetError configures the mock to return an error for a specific operation.
func (mrc *MockRegistryClient) SetError(operation, ref string, err error) {
	mrc.errors[operation+":"+ref] = err
}

// GetImage returns a mock image or error.
func (mrc *MockRegistryClient) GetImage(ref string) (v1.Image, error) {
	if err, exists := mrc.errors["get:"+ref]; exists {
		return nil, err
	}

	image, exists := mrc.images[ref]
	if !exists {
		return nil, fmt.Errorf("image not found: %s", ref)
	}

	return image, nil
}

// PushImage records the push operation.
func (mrc *MockRegistryClient) PushImage(ref string, image v1.Image) error {
	if err, exists := mrc.errors["push:"+ref]; exists {
		return err
	}

	mrc.pushes = append(mrc.pushes, ref)
	mrc.images[ref] = image
	return nil
}

// CheckImageExists checks if an image was configured in the mock.
func (mrc *MockRegistryClient) CheckImageExists(ref string) (bool, error) {
	if err, exists := mrc.errors["exists:"+ref]; exists {
		return false, err
	}

	_, exists := mrc.images[ref]
	return exists, nil
}

// GetPushes returns all the pushes that were made.
func (mrc *MockRegistryClient) GetPushes() []string {
	pushes := make([]string, len(mrc.pushes))
	copy(pushes, mrc.pushes)
	return pushes
}

// Clear clears all mock data.
func (mrc *MockRegistryClient) Clear() {
	mrc.images = make(map[string]v1.Image)
	mrc.pushes = make([]string, 0)
	mrc.errors = make(map[string]error)
}