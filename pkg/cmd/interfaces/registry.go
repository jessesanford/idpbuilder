// Package interfaces defines registry management interfaces for CLI operations.
// These interfaces provide contracts for managing registry connections,
// authentication, and configuration across different registry types.
package interfaces

import "context"

// RegistryCommand defines the interface for registry-specific operations.
// This interface handles direct interaction with container registries including
// configuration, authentication, and metadata retrieval.
type RegistryCommand interface {
	// ConfigureRegistry sets up the registry connection with the provided configuration.
	// This should validate the configuration and establish any necessary connections.
	ConfigureRegistry(config RegistryConfig) error

	// ValidateCredentials verifies that the current authentication is valid.
	// This can be used to test connectivity and auth before performing operations.
	ValidateCredentials(ctx context.Context) error

	// GetRegistryInfo retrieves metadata about the configured registry.
	// This includes version information, capabilities, and other registry details.
	GetRegistryInfo(ctx context.Context) (*RegistryInfo, error)

	// TestConnection performs a basic connectivity test to the registry.
	// This is useful for troubleshooting configuration issues.
	TestConnection(ctx context.Context) error

	// GetCapabilities returns the capabilities supported by the configured registry.
	// This helps determine what operations are available.
	GetCapabilities(ctx context.Context) ([]string, error)
}

// RegistryManager handles configuration and management of multiple registries.
// This interface enables users to work with multiple registries simultaneously
// and switch between them as needed.
type RegistryManager interface {
	// AddRegistry adds a new registry configuration with the given name.
	// The name should be unique across all configured registries.
	AddRegistry(name string, config RegistryConfig) error

	// RemoveRegistry removes a registry configuration by name.
	// This will fail if the registry is currently in use.
	RemoveRegistry(name string) error

	// UpdateRegistry updates an existing registry configuration.
	// This allows modifying URLs, credentials, and other settings.
	UpdateRegistry(name string, config RegistryConfig) error

	// GetRegistry retrieves a registry configuration by name.
	// Returns an error if the registry is not found.
	GetRegistry(name string) (*RegistryConfig, error)

	// ListRegistries returns all configured registry configurations.
	// This provides a way to enumerate available registries.
	ListRegistries() ([]RegistryConfig, error)

	// SetDefault sets the default registry to use when none is specified.
	// This affects commands that don't explicitly specify a registry.
	SetDefault(name string) error

	// GetDefault returns the currently configured default registry.
	GetDefault() (*RegistryConfig, error)

	// ValidateRegistry performs validation on a registry configuration.
	// This checks connectivity, authentication, and configuration validity.
	ValidateRegistry(ctx context.Context, name string) error
}

// RegistryAuthenticator handles authentication with different registry types.
// This interface abstracts the various authentication mechanisms used by
// different registry implementations.
type RegistryAuthenticator interface {
	// Authenticate performs authentication with the registry using the provided config.
	// This should handle different auth types (basic, token, oauth) appropriately.
	Authenticate(ctx context.Context, config RegistryConfig) error

	// RefreshCredentials refreshes authentication tokens/credentials if supported.
	// This is useful for long-running operations or expired credentials.
	RefreshCredentials(ctx context.Context) error

	// GetAuthToken returns the current authentication token for API calls.
	// This may trigger authentication if not already performed.
	GetAuthToken(ctx context.Context) (string, error)

	// IsAuthenticated returns true if currently authenticated with the registry.
	IsAuthenticated() bool

	// GetSupportedAuthTypes returns the authentication types this authenticator supports.
	GetSupportedAuthTypes() []string

	// Logout clears any stored authentication state.
	Logout() error
}

// RegistryDiscovery helps discover and identify registry types and capabilities.
// This interface enables automatic detection of registry features and optimal
// configuration for different registry implementations.
type RegistryDiscovery interface {
	// DiscoverRegistry attempts to identify the type and capabilities of a registry.
	// This can be used to auto-configure registry settings.
	DiscoverRegistry(ctx context.Context, url string) (*RegistryInfo, error)

	// ProbeCapabilities tests specific capabilities on a registry.
	// This returns which features are supported by the target registry.
	ProbeCapabilities(ctx context.Context, url string) ([]string, error)

	// DetectRegistryType attempts to identify the specific registry implementation.
	// This helps optimize interactions for specific registry types.
	DetectRegistryType(ctx context.Context, url string) (string, error)

	// GetRecommendedConfig returns optimized configuration for the detected registry.
	GetRecommendedConfig(ctx context.Context, url string) (*RegistryConfig, error)
}

// RegistryHealthChecker provides health monitoring for registry connections.
// This interface enables proactive monitoring of registry availability and
// performance for better error handling and user experience.
type RegistryHealthChecker interface {
	// CheckHealth performs a comprehensive health check on the registry.
	// This includes connectivity, authentication, and basic functionality tests.
	CheckHealth(ctx context.Context, config RegistryConfig) (*HealthStatus, error)

	// CheckConnectivity tests basic network connectivity to the registry.
	CheckConnectivity(ctx context.Context, url string) error

	// CheckAuthentication verifies that authentication is working properly.
	CheckAuthentication(ctx context.Context, config RegistryConfig) error

	// GetHealthHistory returns historical health check results if available.
	GetHealthHistory(registryName string) ([]HealthStatus, error)
}

// HealthStatus represents the result of a registry health check.
type HealthStatus struct {
	RegistryName   string            // Name of the checked registry
	URL            string            // Registry URL
	Status         string            // Overall status: "healthy", "degraded", "unhealthy"
	ResponseTime   int64             // Response time in milliseconds
	LastChecked    int64             // Timestamp of last check
	Errors         []string          // Any errors encountered
	Capabilities   []string          // Working capabilities
	Metadata       map[string]string // Additional health information
}

// RegistryCatalogBrowser provides functionality for browsing registry contents.
// This interface enables exploration of repositories, tags, and other registry
// artifacts without requiring specific knowledge of registry APIs.
type RegistryCatalogBrowser interface {
	// ListRepositories returns a list of repositories in the registry.
	// This supports pagination through the provided options.
	ListRepositories(ctx context.Context, opts ListOptions) ([]string, error)

	// ListTags returns tags for a specific repository.
	ListTags(ctx context.Context, repository string, opts ListOptions) ([]string, error)

	// GetRepositoryInfo retrieves metadata about a specific repository.
	GetRepositoryInfo(ctx context.Context, repository string) (*RepositoryInfo, error)

	// SearchRepositories searches for repositories matching the given query.
	SearchRepositories(ctx context.Context, query string, opts ListOptions) ([]RepositoryInfo, error)

	// GetCatalogInfo returns overall catalog information for the registry.
	GetCatalogInfo(ctx context.Context) (*CatalogInfo, error)
}

// RepositoryInfo contains metadata about a repository in a registry.
type RepositoryInfo struct {
	Name        string            // Repository name
	Description string            // Repository description
	TagCount    int               // Number of tags
	LastPush    int64             // Timestamp of last push
	Size        int64             // Total repository size in bytes
	Metadata    map[string]string // Additional repository metadata
}

// CatalogInfo contains overall information about a registry's catalog.
type CatalogInfo struct {
	RepositoryCount int               // Total number of repositories
	TotalSize       int64             // Total size of all repositories
	LastUpdate      int64             // Timestamp of last catalog update
	Capabilities    []string          // Catalog-related capabilities
	Metadata        map[string]string // Additional catalog metadata
}