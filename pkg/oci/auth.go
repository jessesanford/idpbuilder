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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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
			Sources: []CredentialSource{
				SourceDockerConfig,
				SourceEnvironment,
				SourceKubernetes,
			},
		}
	}

	// Set default secrets path if not provided
	if config.SecretsPath == "" {
		homeDir, err := os.UserHomeDir()
		if err == nil {
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
	// Clear cache for this registry
	a.mu.Lock()
	delete(a.cache, registry)
	a.mu.Unlock()

	// Re-authenticate
	return a.Authenticate(ctx, registry)
}

// ValidateCredentials checks if credentials are still valid
func (a *DefaultAuthenticator) ValidateCredentials(ctx context.Context, creds *Credentials) (bool, error) {
	if creds == nil {
		return false, ErrInvalidCredentials
	}

	if creds.IsExpired() {
		return false, ErrTokenExpired
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

	// Try exact match first, then partial matches
	for configRegistry, auth := range config.Auths {
		if configRegistry == registry || strings.Contains(registry, configRegistry) {
			if auth.Auth == "" {
				continue
			}

			creds, err := a.parseBasicAuth(auth.Auth)
			if err != nil {
				continue
			}

			creds.Registry = registry
			return creds, nil
		}
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
		// Try to create in-cluster client
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, ErrNoCredentialsFound
		}

		client, err := kubernetes.NewForConfig(config)
		if err != nil {
			return nil, fmt.Errorf("creating k8s client: %w", err)
		}
		a.k8sClient = client
	}

	// Get current namespace
	namespace := os.Getenv("POD_NAMESPACE")
	if namespace == "" {
		namespace = "default"
	}

	// List secrets and look for dockercfg secrets
	secrets, err := a.k8sClient.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{
		TypeMeta: metav1.TypeMeta{
			Kind: "Secret",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("listing secrets: %w", err)
	}

	for _, secret := range secrets.Items {
		if secret.Type == "kubernetes.io/dockercfg" || secret.Type == "kubernetes.io/dockerconfigjson" {
			creds, err := a.parseK8sSecretData(secret.Data, registry)
			if err != nil {
				continue
			}
			if creds != nil {
				return creds, nil
			}
		}
	}

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

// parseK8sSecretData parses Kubernetes secret data
func (a *DefaultAuthenticator) parseK8sSecretData(data map[string][]byte, registry string) (*Credentials, error) {
	var configData []byte
	var exists bool

	if configData, exists = data[".dockerconfigjson"]; exists {
		// Parse dockerconfigjson format
	} else if configData, exists = data[".dockercfg"]; exists {
		// Parse dockercfg format
	} else {
		return nil, ErrNoCredentialsFound
	}

	var config struct {
		Auths map[string]struct {
			Auth string `json:"auth"`
		} `json:"auths"`
	}

	if err := json.Unmarshal(configData, &config); err != nil {
		return nil, fmt.Errorf("parsing k8s secret config: %w", err)
	}

	// Find matching registry
	for secretRegistry, auth := range config.Auths {
		if secretRegistry == registry || strings.Contains(registry, secretRegistry) {
			creds, err := a.parseBasicAuth(auth.Auth)
			if err != nil {
				continue
			}
			creds.Registry = registry
			return creds, nil
		}
	}

	return nil, ErrNoCredentialsFound
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