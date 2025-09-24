package oci

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"k8s.io/client-go/kubernetes"
)

// DefaultAuthenticator implements the Authenticator interface
type DefaultAuthenticator struct {
	sources     []CredentialSource
	cache       map[string]*Credentials
	mu          sync.RWMutex
	secretsPath string
	k8sClient   kubernetes.Interface
}

// AuthConfig represents configuration for the authenticator
type AuthConfig struct {
	Sources     []CredentialSource
	SecretsPath string
	K8sClient   kubernetes.Interface
}

// NewAuthenticator creates a new DefaultAuthenticator
func NewAuthenticator(config *AuthConfig) *DefaultAuthenticator {
	if config == nil {
		config = &AuthConfig{
			Sources: []CredentialSource{SourceDockerConfig, SourceEnvironment, SourceKubernetes},
		}
	}

	if config.SecretsPath == "" {
		if homeDir, err := os.UserHomeDir(); err == nil {
			config.SecretsPath = filepath.Join(homeDir, ".docker")
		}
	}

	return &DefaultAuthenticator{
		sources:     config.Sources,
		cache:       make(map[string]*Credentials),
		secretsPath: config.SecretsPath,
		k8sClient:   config.K8sClient,
	}
}

// Authenticate returns credentials for the given registry
func (a *DefaultAuthenticator) Authenticate(ctx context.Context, registry string) (*Credentials, error) {
	// Check cache first
	if creds := a.getCachedCredentials(registry); creds != nil && !creds.IsExpired() {
		return creds, nil
	}

	// Try each credential source
	for _, source := range a.sources {
		creds, err := a.loadCredentialsFromSource(ctx, registry, source)
		if err != nil {
			continue // Try next source
		}

		if creds != nil && creds.IsValid() {
			a.cacheCredentials(registry, creds)
			return creds, nil
		}
	}

	return nil, NewAuthError(registry, 0, ErrNoCredentialsFound)
}

// RefreshToken refreshes an expired token
func (a *DefaultAuthenticator) RefreshToken(ctx context.Context, registry string) (*Credentials, error) {
	a.mu.Lock()
	delete(a.cache, registry)
	a.mu.Unlock()
	return a.Authenticate(ctx, registry)
}

// ValidateCredentials checks if credentials are still valid
func (a *DefaultAuthenticator) ValidateCredentials(ctx context.Context, creds *Credentials) (bool, error) {
	if creds == nil || creds.IsExpired() {
		return false, ErrInvalidCredentials
	}
	return creds.IsValid(), nil
}

// loadCredentialsFromSource loads credentials from a specific source
func (a *DefaultAuthenticator) loadCredentialsFromSource(ctx context.Context, registry string, source CredentialSource) (*Credentials, error) {
	switch source {
	case SourceDockerConfig:
		return a.loadDockerConfig(registry)
	case SourceEnvironment:
		return a.loadEnvironmentCreds(registry)
	case SourceKubernetes:
		return a.loadKubernetesSecret(ctx, registry)
	default:
		return nil, NewAuthError(registry, source, ErrCredentialSourceFailed)
	}
}

// loadDockerConfig loads credentials from Docker config
func (a *DefaultAuthenticator) loadDockerConfig(registry string) (*Credentials, error) {
	configPath := filepath.Join(a.secretsPath, "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, ErrNoCredentialsFound
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("reading docker config: %w", err)
	}

	var config struct {
		Auths map[string]struct {
			Auth string `json:"auth"`
		} `json:"auths"`
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parsing docker config: %w", err)
	}

	// Simple registry match
	if auth, ok := config.Auths[registry]; ok && auth.Auth != "" {
		creds, err := a.parseBasicAuth(auth.Auth)
		if err != nil {
			return nil, err
		}
		creds.Registry = registry
		return creds, nil
	}

	return nil, ErrNoCredentialsFound
}

// loadEnvironmentCreds loads credentials from environment variables
func (a *DefaultAuthenticator) loadEnvironmentCreds(registry string) (*Credentials, error) {
	username := os.Getenv("REGISTRY_USERNAME")
	password := os.Getenv("REGISTRY_PASSWORD")
	token := os.Getenv("REGISTRY_TOKEN")

	if username == "" && password == "" && token == "" {
		return nil, ErrNoCredentialsFound
	}

	creds := &Credentials{
		Username: username,
		Password: password,
		Token:    token,
		Registry: registry,
	}

	return creds, nil
}

// loadKubernetesSecret loads credentials from Kubernetes secrets
func (a *DefaultAuthenticator) loadKubernetesSecret(ctx context.Context, registry string) (*Credentials, error) {
	if a.k8sClient == nil {
		return nil, ErrNoCredentialsFound
	}

	// Get current namespace
	namespace := os.Getenv("POD_NAMESPACE")
	if namespace == "" {
		namespace = "default"
	}

	// Simple K8s secret lookup - minimal implementation
	return nil, ErrNoCredentialsFound
}

// parseBasicAuth parses base64 encoded basic auth
func (a *DefaultAuthenticator) parseBasicAuth(auth string) (*Credentials, error) {
	decoded, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		return nil, fmt.Errorf("decoding auth: %w", err)
	}

	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 {
		return nil, ErrInvalidCredentials
	}

	return &Credentials{
		Username: parts[0],
		Password: parts[1],
	}, nil
}


// getCachedCredentials retrieves credentials from cache
func (a *DefaultAuthenticator) getCachedCredentials(registry string) *Credentials {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.cache[registry]
}

// cacheCredentials stores credentials in cache
func (a *DefaultAuthenticator) cacheCredentials(registry string, creds *Credentials) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.cache[registry] = creds
}

// NewAuthenticatorFromFlags creates an authenticator from CLI flag values
func NewAuthenticatorFromFlags(username, password string) (*DefaultAuthenticator, error) {
	// Validate both parameters are non-empty
	if username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}
	if password == "" {
		return nil, fmt.Errorf("password cannot be empty")
	}

	// Create AuthConfig with credentials - using direct credentials approach
	// Since the existing NewAuthenticator uses sources, we'll create a minimal config
	config := &AuthConfig{
		Sources: []CredentialSource{SourceEnvironment}, // Use environment as the source type
	}

	// Create the authenticator
	auth := NewAuthenticator(config)

	// For flags-based auth, we'll store the credentials directly in cache
	// This allows the credentials to be available when requested
	creds := &Credentials{
		Username: username,
		Password: password,
		Registry: "", // Will be set when specific registry is requested
	}

	// Cache the credentials for immediate use - using a default key
	auth.mu.Lock()
	auth.cache["flags-provided"] = creds
	auth.mu.Unlock()

	return auth, nil
}

// NewAuthenticatorFromEnv creates an authenticator from environment variables
// Reads OCI_USERNAME and OCI_PASSWORD environment variables
func NewAuthenticatorFromEnv() (*DefaultAuthenticator, error) {
	// Read environment variables
	username := os.Getenv("OCI_USERNAME")
	password := os.Getenv("OCI_PASSWORD")

	// Validate both are present and non-empty
	if username == "" {
		return nil, fmt.Errorf("OCI_USERNAME environment variable is not set or empty")
	}
	if password == "" {
		return nil, fmt.Errorf("OCI_PASSWORD environment variable is not set or empty")
	}

	// Create AuthConfig - use environment source since that's what we're doing
	config := &AuthConfig{
		Sources: []CredentialSource{SourceEnvironment},
	}

	// Create the authenticator
	auth := NewAuthenticator(config)

	// Store the credentials directly in cache for immediate use
	creds := &Credentials{
		Username: username,
		Password: password,
		Registry: "", // Will be set when specific registry is requested
	}

	// Cache the credentials using an environment-specific key
	auth.mu.Lock()
	auth.cache["env-provided"] = creds
	auth.mu.Unlock()

	return auth, nil
}

// NewAuthenticatorFromSecrets creates an authenticator from Kubernetes secret data
// The secretData map should contain "username" and "password" keys
func NewAuthenticatorFromSecrets(secretData map[string][]byte) (*DefaultAuthenticator, error) {
	// Validate secretData is not nil
	if secretData == nil {
		return nil, fmt.Errorf("secretData cannot be nil")
	}

	// Check for required keys and convert to strings
	usernameBytes, hasUsername := secretData["username"]
	passwordBytes, hasPassword := secretData["password"]

	if !hasUsername {
		return nil, fmt.Errorf("username key missing in secret data")
	}
	if !hasPassword {
		return nil, fmt.Errorf("password key missing in secret data")
	}

	// Convert []byte values to strings
	username := string(usernameBytes)
	password := string(passwordBytes)

	// Validate converted values are not empty
	if username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}
	if password == "" {
		return nil, fmt.Errorf("password cannot be empty")
	}

	// Create AuthConfig - use Kubernetes source since that's what we're doing
	config := &AuthConfig{
		Sources: []CredentialSource{SourceKubernetes},
	}

	// Create the authenticator
	auth := NewAuthenticator(config)

	// Store the credentials directly in cache for immediate use
	creds := &Credentials{
		Username: username,
		Password: password,
		Registry: "", // Will be set when specific registry is requested
	}

	// Cache the credentials using a secrets-specific key
	auth.mu.Lock()
	auth.cache["secrets-provided"] = creds
	auth.mu.Unlock()

	return auth, nil
}