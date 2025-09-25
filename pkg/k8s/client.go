package k8s

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
)

// Client provides a Kubernetes client with scheme management
type Client struct {
	restConfig *rest.Config
	scheme     *runtime.Scheme
	helper     *ResourceHelper
}

// NewClient creates a new Client with the provided rest config
func NewClient(config *rest.Config) (*Client, error) {
	if config == nil {
		return nil, fmt.Errorf("rest config cannot be nil")
	}

	// Validate config has required fields
	if config.Host == "" {
		return nil, fmt.Errorf("rest config must have a host")
	}

	// Build the scheme
	schemeBuilder := NewSchemeBuilder()
	scheme := schemeBuilder.Build()

	// Create the helper
	helper := NewResourceHelper()
	helper.SetScheme(scheme)

	return &Client{
		restConfig: config,
		scheme:     scheme,
		helper:     helper,
	}, nil
}

// GetScheme returns the runtime scheme used by this client
func (c *Client) GetScheme() *runtime.Scheme {
	return c.scheme
}

// GetHelper returns the resource helper used by this client
func (c *Client) GetHelper() *ResourceHelper {
	return c.helper
}

// GetConfig returns the rest config used by this client
func (c *Client) GetConfig() *rest.Config {
	return c.restConfig
}

// IsHealthy performs a basic health check on the client configuration
func (c *Client) IsHealthy() bool {
	return c.restConfig != nil && c.scheme != nil && c.helper != nil
}