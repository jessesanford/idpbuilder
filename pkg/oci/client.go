package oci

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// OCIClient implements the RegistryClient interface using go-containerregistry
type OCIClient struct {
	registry  string
	insecure  bool
	transport http.RoundTripper
	auth      Authenticator
	mu        sync.RWMutex
}

// NewRegistryClient creates a new RegistryClient instance
func NewRegistryClient() RegistryClient {
	return &OCIClient{
		auth: NewAuthenticator(nil), // Use existing DefaultAuthenticator with defaults
	}
}

// Connect establishes a connection to the specified OCI registry
func (c *OCIClient) Connect(ctx context.Context, registry string) error {
	if registry == "" {
		return errors.New("registry cannot be empty")
	}

	// Validate the registry URL format
	if err := c.validateURL(registry); err != nil {
		return fmt.Errorf("invalid registry URL: %w", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Store registry for future use
	c.registry = registry

	// Configure HTTP transport
	c.transport = c.configureTransport(c.insecure)

	// Test connection by making a simple ping request with context timeout
	timeout := 30 * time.Second
	if deadline, ok := ctx.Deadline(); ok {
		timeout = time.Until(deadline)
	}

	// Create a context with timeout for connection test
	testCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Parse the registry URL to extract just the host for go-containerregistry
	parsedURL, err := url.Parse(registry)
	if err != nil {
		return fmt.Errorf("parsing registry URL %q: %w", registry, err)
	}

	// go-containerregistry expects just the host, not the full URL
	registryHost := parsedURL.Host
	reg, err := name.NewRegistry(registryHost)
	if err != nil {
		return fmt.Errorf("parsing registry host %q: %w", registryHost, err)
	}

	// Use go-containerregistry to ping the registry
	_, err = remote.Catalog(testCtx, reg, remote.WithTransport(c.transport))
	if err != nil {
		// Don't fail on catalog errors - some registries don't support it
		// Just try a basic ping instead
		pingURL := fmt.Sprintf("%s/v2/", strings.TrimSuffix(registry, "/"))
		req, err := http.NewRequestWithContext(testCtx, "GET", pingURL, nil)
		if err != nil {
			return fmt.Errorf("creating ping request: %w", err)
		}

		client := &http.Client{Transport: c.transport}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("connecting to registry %q: %w", registry, err)
		}
		defer resp.Body.Close()

		// Accept any 2xx, 401 (auth required), or 403 (forbidden) as valid responses
		// These indicate the registry is responding
		if resp.StatusCode >= 500 {
			return fmt.Errorf("registry %q returned server error: %d", registry, resp.StatusCode)
		}
	}

	return nil
}

// Authenticate configures authentication credentials for registry access
func (c *OCIClient) Authenticate(credentials *ClientCredentials) error {
	if credentials == nil {
		// Allow nil credentials for anonymous access
		return nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Convert ClientCredentials to internal Credentials type
	creds := c.convertCredentials(credentials)

	// Cache the credentials in the authenticator for this registry
	if c.auth != nil {
		// Store credentials in the authenticator's cache
		// Since DefaultAuthenticator has a private cache, we'll create a new one
		// that can be used with the existing interface
		if defaultAuth, ok := c.auth.(*DefaultAuthenticator); ok {
			defaultAuth.mu.Lock()
			defaultAuth.cache[creds.Registry] = creds
			defaultAuth.mu.Unlock()
		}
	}

	return nil
}

// SetInsecure configures whether to allow insecure (HTTP) connections
func (c *OCIClient) SetInsecure(insecure bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.insecure = insecure

	// Reconfigure transport if we already have one
	if c.transport != nil {
		c.transport = c.configureTransport(insecure)
	}
}

// GetTransport returns the configured HTTP transport
func (c *OCIClient) GetTransport() http.RoundTripper {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.transport
}

// Close closes the connection and cleans up resources
func (c *OCIClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Clean up transport if it has a close method
	if transport, ok := c.transport.(*http.Transport); ok {
		transport.CloseIdleConnections()
	}

	// Reset state
	c.registry = ""
	c.transport = nil

	return nil
}

// Helper Functions

// convertCredentials converts ClientCredentials to internal Credentials
func (c *OCIClient) convertCredentials(cc *ClientCredentials) *Credentials {
	if cc == nil {
		return nil
	}

	return &Credentials{
		Username: cc.Username,
		Password: cc.Password,
		Token:    cc.Token,
		Registry: cc.Registry,
		// ExpiresAt is not set as ClientCredentials doesn't have expiry info
	}
}

// configureTransport creates and configures an HTTP transport
func (c *OCIClient) configureTransport(insecure bool) *http.Transport {
	transport := &http.Transport{
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		MaxIdleConnsPerHost: 10,
	}

	if insecure {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	return transport
}

// validateURL validates the registry URL format
func (c *OCIClient) validateURL(registry string) error {
	if registry == "" {
		return errors.New("registry URL cannot be empty")
	}

	// Parse URL
	parsedURL, err := url.Parse(registry)
	if err != nil {
		return fmt.Errorf("malformed URL: %w", err)
	}

	// Check scheme
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("unsupported scheme %q, must be http or https", parsedURL.Scheme)
	}

	// Check host
	if parsedURL.Host == "" {
		return errors.New("missing host in URL")
	}

	return nil
}
