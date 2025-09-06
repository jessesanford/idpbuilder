package registry

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/cnoe-io/idpbuilder/pkg/certs"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// GiteaClient implements Client for Gitea registries with Phase 1 certificate integration
type GiteaClient struct {
	// Configuration
	baseURL   string
	username  string
	password  string
	userAgent string

	// Phase 1 certificate integration
	trustStore certs.TrustStoreManager

	// Transport and authentication
	transport http.RoundTripper
	auth      authn.Authenticator

	// Client configuration
	timeout    time.Duration
	maxRetries int
	retryDelay time.Duration
	insecure   bool

	// Synchronization
	mu sync.RWMutex

	// Connection tracking
	lastUsed time.Time

	// Feature flags (R307)
	insecureRegistryFlag bool
}

// NewGiteaClient creates a new Gitea registry client with Phase 1 certificate integration
func NewGiteaClient(baseURL, username, password string, trustStore certs.TrustStoreManager, opts ...ClientOption) (*GiteaClient, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("base URL cannot be empty")
	}

	// Parse and validate URL
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL %s: %w", baseURL, err)
	}

	// Ensure we have a scheme
	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "https"
		baseURL = parsedURL.String()
	}

	client := &GiteaClient{
		baseURL:    baseURL,
		username:   username,
		password:   password,
		userAgent:  "idpbuilder-oci/1.0",
		timeout:    30 * time.Second,
		maxRetries: 3,
		retryDelay: 1 * time.Second,
		insecure:   false,
		trustStore: trustStore,
		lastUsed:   time.Now(),
	}

	// Apply options
	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, fmt.Errorf("failed to apply client option: %w", err)
		}
	}

	// Check for insecure registry feature flag (R307)
	if os.Getenv("IDPBUILDER_INSECURE_REGISTRY") == "true" {
		client.insecureRegistryFlag = true
		client.insecure = true
	}

	// Configure authentication
	client.auth = configureAuth(username, password)

	// Configure transport with Phase 1 certificate integration
	if err := client.configureTransport(); err != nil {
		return nil, fmt.Errorf("failed to configure transport: %w", err)
	}

	return client, nil
}

// ClientOption allows customization of the GiteaClient
type ClientOption func(*GiteaClient) error

// WithTimeout sets the client timeout
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *GiteaClient) error {
		if timeout <= 0 {
			return fmt.Errorf("timeout must be positive")
		}
		c.timeout = timeout
		return nil
	}
}

// WithRetryConfig sets retry configuration
func WithRetryConfig(maxRetries int, retryDelay time.Duration) ClientOption {
	return func(c *GiteaClient) error {
		if maxRetries < 0 {
			return fmt.Errorf("max retries cannot be negative")
		}
		c.maxRetries = maxRetries
		c.retryDelay = retryDelay
		return nil
	}
}

// WithInsecure enables insecure mode (R307 feature flag)
func WithInsecure(insecure bool) ClientOption {
	return func(c *GiteaClient) error {
		c.insecure = insecure
		return nil
	}
}

// WithUserAgent sets custom user agent
func WithUserAgent(userAgent string) ClientOption {
	return func(c *GiteaClient) error {
		if userAgent == "" {
			return fmt.Errorf("user agent cannot be empty")
		}
		c.userAgent = userAgent
		return nil
	}
}

// Push pushes an image to the Gitea registry
func (c *GiteaClient) Push(ctx context.Context, image v1.Image, ref string, opts PushOptions) error {
	c.mu.Lock()
	c.lastUsed = time.Now()
	c.mu.Unlock()

	// Parse the reference
	repo, err := name.ParseReference(ref)
	if err != nil {
		return &ClientError{
			Code:    "invalid_reference",
			Message: fmt.Sprintf("invalid reference %s", ref),
			Details: map[string]interface{}{"registry": c.baseURL, "operation": "push"},
		}
	}

	// Apply timeout if specified
	if opts.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	// Configure remote options
	remoteOpts, err := c.buildRemoteOptions(opts.Insecure, opts.Platform)
	if err != nil {
		return fmt.Errorf("failed to build remote options: %w", err)
	}

	// Perform push with basic retry
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(c.retryDelay):
			}
		}

		err := remote.Write(repo, image, remoteOpts...)
		if err == nil {
			return nil
		}

		// For auth/access errors, don't retry
		if strings.Contains(err.Error(), "401") || strings.Contains(err.Error(), "403") {
			return c.wrapError("push", err)
		}
	}

	return c.wrapError("push", fmt.Errorf("push failed after %d attempts", c.maxRetries+1))
}

// Pull pulls an image from the Gitea registry
func (c *GiteaClient) Pull(ctx context.Context, ref string, opts PullOptions) (v1.Image, error) {
	c.mu.Lock()
	c.lastUsed = time.Now()
	c.mu.Unlock()

	// Parse the reference
	repo, err := name.ParseReference(ref)
	if err != nil {
		return nil, &ClientError{
			Code:    "invalid_reference",
			Message: fmt.Sprintf("invalid reference %s", ref),
			Details: map[string]interface{}{"registry": c.baseURL, "operation": "pull"},
		}
	}

	// Apply timeout if specified
	if opts.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	// Configure remote options
	remoteOpts, err := c.buildRemoteOptions(opts.Insecure, opts.Platform)
	if err != nil {
		return nil, fmt.Errorf("failed to build remote options: %w", err)
	}

	// Perform pull with basic retry
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(c.retryDelay):
			}
		}

		image, err := remote.Image(repo, remoteOpts...)
		if err == nil {
			return image, nil
		}

		// For auth/access errors, don't retry
		if strings.Contains(err.Error(), "401") || strings.Contains(err.Error(), "403") {
			return nil, c.wrapError("pull", err)
		}
	}

	return nil, c.wrapError("pull", fmt.Errorf("pull failed after %d attempts", c.maxRetries+1))
}

// Catalog lists repositories in the Gitea registry
func (c *GiteaClient) Catalog(ctx context.Context) ([]string, error) {
	c.mu.Lock()
	c.lastUsed = time.Now()
	c.mu.Unlock()

	// Parse registry URL
	registry, err := name.NewRegistry(c.baseURL)
	if err != nil {
		return nil, &ClientError{
			Code:    "invalid_registry",
			Message: "invalid registry URL",
			Details: map[string]interface{}{"registry": c.baseURL, "operation": "catalog"},
		}
	}

	// Configure remote options
	remoteOpts, err := c.buildRemoteOptions(false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build remote options: %w", err)
	}

	// List repositories
	repos, err := remote.Catalog(ctx, registry, remoteOpts...)
	if err != nil {
		return nil, c.wrapError("catalog", err)
	}

	return repos, nil
}

// Tags lists tags for a repository in the Gitea registry
func (c *GiteaClient) Tags(ctx context.Context, repository string) ([]string, error) {
	c.mu.Lock()
	c.lastUsed = time.Now()
	c.mu.Unlock()

	// Parse repository reference
	repo, err := name.NewRepository(repository)
	if err != nil {
		return nil, &ClientError{
			Code:    "invalid_repository",
			Message: fmt.Sprintf("invalid repository %s", repository),
			Details: map[string]interface{}{"registry": c.baseURL, "operation": "tags"},
		}
	}

	// Configure remote options
	remoteOpts, err := c.buildRemoteOptions(false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build remote options: %w", err)
	}

	// List tags
	tags, err := remote.List(repo, remoteOpts...)
	if err != nil {
		return nil, c.wrapError("tags", err)
	}

	return tags, nil
}

// Close cleans up resources used by the client
func (c *GiteaClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Currently no resources to clean up explicitly
	// In future implementations, this might close connection pools
	return nil
}

// configureTransport sets up HTTP transport with Phase 1 certificate integration
func (c *GiteaClient) configureTransport() error {
	if c.trustStore == nil {
		return fmt.Errorf("trust store manager is required")
	}

	// Check if registry should use insecure mode
	if c.insecure || c.insecureRegistryFlag {
		// Mark registry as insecure in trust store
		if err := c.trustStore.SetInsecureRegistry(c.baseURL, true); err != nil {
			return fmt.Errorf("failed to set registry as insecure: %w", err)
		}
	}

	// Use Phase 1's TrustStoreManager to create HTTP client
	httpClient, err := c.trustStore.CreateHTTPClient(c.baseURL)
	if err != nil {
		return fmt.Errorf("failed to create HTTP client with Phase 1 certificates: %w", err)
	}

	// Extract the transport
	c.transport = httpClient.Transport

	return nil
}

// buildRemoteOptions creates remote options for go-containerregistry operations
func (c *GiteaClient) buildRemoteOptions(insecureOverride bool, platform *v1.Platform) ([]remote.Option, error) {
	var opts []remote.Option

	// Add authentication
	opts = append(opts, remote.WithAuth(c.auth))

	// Add user agent
	opts = append(opts, remote.WithUserAgent(c.userAgent))

	// Configure transport using Phase 1's trust store
	transportOpt, err := c.trustStore.ConfigureTransport(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to configure transport: %w", err)
	}
	opts = append(opts, transportOpt)

	// Add platform if specified
	if platform != nil {
		opts = append(opts, remote.WithPlatform(*platform))
	}

	return opts, nil
}

// wrapError wraps an error in a ClientError with simplified categorization
func (c *GiteaClient) wrapError(operation string, err error) error {
	if err == nil {
		return nil
	}

	errStr := strings.ToLower(err.Error())
	var code string

	switch {
	case strings.Contains(errStr, "unauthorized") || strings.Contains(errStr, "401"):
		code = "auth_failed"
	case strings.Contains(errStr, "forbidden") || strings.Contains(errStr, "403"):
		code = "access_denied"
	case strings.Contains(errStr, "not found") || strings.Contains(errStr, "404"):
		code = "not_found"
	default:
		code = "registry_error"
	}

	return &ClientError{
		Code:    code,
		Message: fmt.Sprintf("%s operation failed: %v", operation, err),
		Details: map[string]interface{}{"registry": c.baseURL, "operation": operation},
	}
}
