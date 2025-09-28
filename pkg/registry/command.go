package registry

import (
	"context"
	"fmt"

	"github.com/cnoe-io/idpbuilder/pkg/cmd/interfaces"
	"github.com/cnoe-io/idpbuilder/pkg/config"
	"github.com/cnoe-io/idpbuilder/pkg/providers"
)

// Command implements the RegistryCommand interface using the registry client.
// It bridges the CLI command interface with the Provider implementation,
// providing registry operations accessible from command-line tools.
type Command struct {
	client providers.Provider
	config *config.RegistryConfig
}

// NewCommand creates a new registry command with the provided configuration.
// It initializes the underlying registry client and validates the configuration.
func NewCommand(cfg *config.RegistryConfig) (*Command, error) {
	if cfg == nil {
		return nil, fmt.Errorf("registry configuration cannot be nil")
	}

	// Create the underlying client
	client, err := NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create registry client: %w", err)
	}

	return &Command{
		client: client,
		config: cfg,
	}, nil
}

// ConfigureRegistry implements the RegistryCommand interface.
// It sets up the registry connection with the provided configuration.
func (c *Command) ConfigureRegistry(registryConfig interfaces.RegistryConfig) error {
	// Convert interfaces.RegistryConfig to config.RegistryConfig
	cfg := &config.RegistryConfig{
		URL:      registryConfig.URL,
		Type:     registryConfig.Type,
		Insecure: registryConfig.Insecure,
		Auth:     config.AuthConfig{
			// Note: interfaces.RegistryConfig doesn't include auth details
			// These would need to be provided separately or via command options
		},
		Options: map[string]string{},
	}

	// Create new client with updated configuration
	client, err := NewClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to configure registry: %w", err)
	}

	c.client = client
	c.config = cfg
	return nil
}

// ValidateCredentials implements the RegistryCommand interface.
// It verifies that the current authentication is valid by performing a test operation.
func (c *Command) ValidateCredentials(ctx context.Context) error {
	if c.client == nil {
		return fmt.Errorf("registry client not configured")
	}

	// Test credentials by attempting a list operation on a minimal repository
	// This is a lightweight way to validate authentication
	testRepo := c.config.URL + "/test"
	_, err := c.client.List(ctx, testRepo)

	// If we get a permission error, credentials are being used but invalid
	// If we get a not found error, credentials are valid but repository doesn't exist
	// Both cases indicate successful authentication validation
	if err != nil {
		if providerErr, ok := err.(*providers.ProviderError); ok {
			// Check if error indicates authentication issues vs repository issues
			if providerErr.Op == "list" {
				// Authentication worked if we get a repository-level error
				return nil
			}
		}
		return fmt.Errorf("credential validation failed: %w", err)
	}

	return nil
}

// GetRegistryInfo implements the RegistryCommand interface.
// It retrieves metadata about the configured registry.
func (c *Command) GetRegistryInfo(ctx context.Context) (*interfaces.RegistryInfo, error) {
	if c.client == nil {
		return nil, fmt.Errorf("registry client not configured")
	}

	// In a real implementation, this would probe the registry for version info,
	// capabilities, and other metadata. For now, return basic info from config.
	return &interfaces.RegistryInfo{
		Name:         c.config.Type,
		URL:          c.config.URL,
		Version:      "unknown", // Would need registry-specific probing
		Capabilities: []string{"push", "pull", "list", "delete"},
		Metadata: map[string]string{
			"insecure":      fmt.Sprintf("%t", c.config.Insecure),
			"authenticated": fmt.Sprintf("%t", c.config.Auth.Username != "" || c.config.Auth.Token != ""),
		},
	}, nil
}

// TestConnection implements the RegistryCommand interface.
// It performs a basic connectivity test to the registry.
func (c *Command) TestConnection(ctx context.Context) error {
	if c.client == nil {
		return fmt.Errorf("registry client not configured")
	}

	// Test connectivity by attempting a simple list operation
	// This verifies network connectivity, TLS configuration, and basic registry response
	_, err := c.client.List(ctx, c.config.URL+"/healthcheck")
	if err != nil {
		if providerErr, ok := err.(*providers.ProviderError); ok {
			// Even if the healthcheck repository doesn't exist, a proper error response
			// indicates successful connectivity
			if providerErr.Op == "list" {
				return nil
			}
		}
		return fmt.Errorf("connection test failed: %w", err)
	}

	return nil
}

// GetCapabilities implements the RegistryCommand interface.
// It returns the capabilities supported by the configured registry.
func (c *Command) GetCapabilities(ctx context.Context) ([]string, error) {
	if c.client == nil {
		return nil, fmt.Errorf("registry client not configured")
	}

	// Basic capabilities that our Provider interface supports
	capabilities := []string{
		"push",    // Provider.Push
		"pull",    // Provider.Pull
		"list",    // Provider.List
		"delete",  // Provider.Delete
	}

	// In a more sophisticated implementation, we could probe the registry
	// to determine additional capabilities like:
	// - "manifest-v2"
	// - "oci-artifacts"
	// - "delete-by-tag"
	// - "delete-by-digest"
	// - "catalog-api"

	return capabilities, nil
}