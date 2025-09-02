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
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// GiteaClient implements Client for Gitea registries with Phase 1 certificate integration
type GiteaClient struct {
	// Configuration
	baseURL    string
	username   string
	password   string
	userAgent  string
	
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
		baseURL:     baseURL,
		username:    username,
		password:    password,
		userAgent:   "idpbuilder-oci/1.0",
		timeout:     30 * time.Second,
		maxRetries:  3,
		retryDelay:  1 * time.Second,
		insecure:    false,
		trustStore:  trustStore,
		lastUsed:    time.Now(),
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
		return &RegistryError{
			Type:       ErrorInvalidReference,
			Registry:   c.baseURL,
			Operation:  "push",
			Message:    fmt.Sprintf("invalid reference %s", ref),
			Underlying: err,
			Retryable:  false,
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
	
	// Add progress reporting if provided
	if opts.Progress != nil {
		// Note: go-containerregistry doesn't directly support progress callbacks
		// This would need custom implementation for detailed progress tracking
		opts.Progress(0, -1) // Signal start
	}
	
	// Perform push with retries
	maxRetries := c.maxRetries
	if opts.MaxRetries > 0 {
		maxRetries = opts.MaxRetries
	}
	
	retryDelay := c.retryDelay
	if opts.RetryDelay > 0 {
		retryDelay = opts.RetryDelay
	}
	
	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(retryDelay):
			}
		}
		
		err := remote.Write(repo, image, remoteOpts...)
		if err == nil {
			if opts.Progress != nil {
				opts.Progress(-1, -1) // Signal completion
			}
			return nil
		}
		
		lastErr = c.wrapError("push", err)
		
		// Check if error is retryable
		if regErr, ok := lastErr.(*RegistryError); ok && !regErr.IsRetryable() {
			break
		}
	}
	
	return lastErr
}

// Pull pulls an image from the Gitea registry
func (c *GiteaClient) Pull(ctx context.Context, ref string, opts PullOptions) (v1.Image, error) {
	c.mu.Lock()
	c.lastUsed = time.Now()
	c.mu.Unlock()
	
	// Parse the reference
	repo, err := name.ParseReference(ref)
	if err != nil {
		return nil, &RegistryError{
			Type:       ErrorInvalidReference,
			Registry:   c.baseURL,
			Operation:  "pull",
			Message:    fmt.Sprintf("invalid reference %s", ref),
			Underlying: err,
			Retryable:  false,
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
	
	// Perform pull with retries
	maxRetries := c.maxRetries
	if opts.MaxRetries > 0 {
		maxRetries = opts.MaxRetries
	}
	
	retryDelay := c.retryDelay
	if opts.RetryDelay > 0 {
		retryDelay = opts.RetryDelay
	}
	
	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(retryDelay):
			}
		}
		
		image, err := remote.Image(repo, remoteOpts...)
		if err == nil {
			return image, nil
		}
		
		lastErr = c.wrapError("pull", err)
		
		// Check if error is retryable
		if regErr, ok := lastErr.(*RegistryError); ok && !regErr.IsRetryable() {
			break
		}
	}
	
	return nil, lastErr
}

// Catalog lists repositories in the Gitea registry
func (c *GiteaClient) Catalog(ctx context.Context) ([]string, error) {
	c.mu.Lock()
	c.lastUsed = time.Now()
	c.mu.Unlock()
	
	// Parse registry URL
	registry, err := name.NewRegistry(c.baseURL)
	if err != nil {
		return nil, &RegistryError{
			Type:       ErrorInvalidReference,
			Registry:   c.baseURL,
			Operation:  "catalog",
			Message:    "invalid registry URL",
			Underlying: err,
			Retryable:  false,
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
		return nil, &RegistryError{
			Type:       ErrorInvalidReference,
			Registry:   c.baseURL,
			Operation:  "tags",
			Message:    fmt.Sprintf("invalid repository %s", repository),
			Underlying: err,
			Retryable:  false,
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
	
	// Configure transport based on insecure settings
	useInsecure := c.insecure || insecureOverride || c.insecureRegistryFlag
	
	if useInsecure {
		// Use Phase 1's transport configuration for insecure mode
		transportOpt, err := c.trustStore.ConfigureTransport(c.baseURL)
		if err != nil {
			return nil, fmt.Errorf("failed to configure insecure transport: %w", err)
		}
		opts = append(opts, transportOpt)
	} else {
		// Use Phase 1's transport configuration with certificates
		transportOpt, err := c.trustStore.ConfigureTransport(c.baseURL)
		if err != nil {
			return nil, fmt.Errorf("failed to configure secure transport: %w", err)
		}
		opts = append(opts, transportOpt)
	}
	
	// Add platform if specified
	if platform != nil {
		opts = append(opts, remote.WithPlatform(*platform))
	}
	
	return opts, nil
}

// wrapError wraps an error in a RegistryError with appropriate categorization
func (c *GiteaClient) wrapError(operation string, err error) error {
	if err == nil {
		return nil
	}
	
	// Determine error type and retryability
	errorType := ErrorUnknown
	retryable := true
	statusCode := 0
	
	errStr := strings.ToLower(err.Error())
	
	switch {
	case strings.Contains(errStr, "unauthorized") || strings.Contains(errStr, "401"):
		errorType = ErrorAuthentication
		retryable = false
		statusCode = 401
	case strings.Contains(errStr, "forbidden") || strings.Contains(errStr, "403"):
		errorType = ErrorAuthorization
		retryable = false
		statusCode = 403
	case strings.Contains(errStr, "not found") || strings.Contains(errStr, "404"):
		errorType = ErrorNotFound
		retryable = false
		statusCode = 404
	case strings.Contains(errStr, "conflict") || strings.Contains(errStr, "409"):
		errorType = ErrorConflict
		retryable = false
		statusCode = 409
	case strings.Contains(errStr, "timeout"):
		errorType = ErrorTimeout
		retryable = true
	case strings.Contains(errStr, "certificate") || strings.Contains(errStr, "tls") || strings.Contains(errStr, "x509"):
		errorType = ErrorCertificate
		retryable = false
	case strings.Contains(errStr, "network") || strings.Contains(errStr, "connection"):
		errorType = ErrorNetwork
		retryable = true
	}
	
	return &RegistryError{
		Type:       errorType,
		Registry:   c.baseURL,
		Operation:  operation,
		Message:    fmt.Sprintf("%s operation failed on %s", operation, c.baseURL),
		Underlying: err,
		Retryable:  retryable,
		StatusCode: statusCode,
	}
}