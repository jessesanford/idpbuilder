package registry

import (
	"context"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/cnoe-io/idpbuilder/pkg/certs"
	"github.com/google/go-containerregistry/pkg/authn"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/random"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// mockTrustStore implements certs.TrustStoreManager for testing
type mockTrustStore struct {
	httpClient         *http.Client
	insecureRegistries map[string]bool
	certificates       map[string][]*x509.Certificate
	shouldFailCreate   bool
	shouldFailConfig   bool
	shouldFailAdd      bool
}

func newMockTrustStore() *mockTrustStore {
	return &mockTrustStore{
		httpClient:         &http.Client{Transport: http.DefaultTransport},
		insecureRegistries: make(map[string]bool),
		certificates:       make(map[string][]*x509.Certificate),
	}
}

// AddCertificate adds a certificate for a specific registry
func (m *mockTrustStore) AddCertificate(registry string, cert *x509.Certificate) error {
	if m.shouldFailAdd {
		return fmt.Errorf("mock add certificate failure")
	}
	m.certificates[registry] = append(m.certificates[registry], cert)
	return nil
}

// RemoveCertificate removes the certificate for a registry
func (m *mockTrustStore) RemoveCertificate(registry string) error {
	delete(m.certificates, registry)
	return nil
}

// SetInsecureRegistry marks a registry as insecure (skip TLS verification)
func (m *mockTrustStore) SetInsecureRegistry(registry string, insecure bool) error {
	m.insecureRegistries[registry] = insecure
	return nil
}

// GetTrustedCerts returns all trusted certificates for a registry
func (m *mockTrustStore) GetTrustedCerts(registry string) ([]*x509.Certificate, error) {
	return m.certificates[registry], nil
}

// GetCertPool returns a configured cert pool for a registry
func (m *mockTrustStore) GetCertPool(registry string) (*x509.CertPool, error) {
	pool := x509.NewCertPool()
	if certs, exists := m.certificates[registry]; exists {
		for _, cert := range certs {
			pool.AddCert(cert)
		}
	}
	return pool, nil
}

// IsInsecure checks if a registry is marked as insecure
func (m *mockTrustStore) IsInsecure(registry string) bool {
	return m.insecureRegistries[registry]
}

// LoadFromDisk loads all certificates from persistent storage
func (m *mockTrustStore) LoadFromDisk() error {
	return nil // Mock implementation - no disk operations
}

// SaveToDisk saves a certificate to persistent storage
func (m *mockTrustStore) SaveToDisk(registry string, cert *x509.Certificate) error {
	return nil // Mock implementation - no disk operations
}

// ConfigureTransport creates a remote.Option with proper TLS configuration
func (m *mockTrustStore) ConfigureTransport(registry string) (remote.Option, error) {
	if m.shouldFailConfig {
		return nil, fmt.Errorf("mock trust store config failure")
	}
	return remote.WithTransport(http.DefaultTransport), nil
}

// ConfigureTransportWithConfig creates a remote.Option with custom transport configuration
func (m *mockTrustStore) ConfigureTransportWithConfig(registry string, config *certs.TransportConfig) (remote.Option, error) {
	if m.shouldFailConfig {
		return nil, fmt.Errorf("mock trust store config failure")
	}
	return remote.WithTransport(http.DefaultTransport), nil
}

// CreateHTTPClient creates an HTTP client with proper TLS configuration
func (m *mockTrustStore) CreateHTTPClient(registry string) (*http.Client, error) {
	if m.shouldFailCreate {
		return nil, fmt.Errorf("mock trust store create failure")
	}
	return m.httpClient, nil
}

// CreateHTTPClientWithConfig creates an HTTP client with custom configuration
func (m *mockTrustStore) CreateHTTPClientWithConfig(registry string, config *certs.TransportConfig) (*http.Client, error) {
	if m.shouldFailCreate {
		return nil, fmt.Errorf("mock trust store create failure")
	}
	return m.httpClient, nil
}

// Test NewGiteaClient constructor
func TestNewGiteaClient(t *testing.T) {
	trustStore := newMockTrustStore()
	
	tests := []struct {
		name        string
		baseURL     string
		username    string
		password    string
		opts        []ClientOption
		wantErr     bool
		errContains string
	}{
		{
			name:     "valid client creation",
			baseURL:  "https://gitea.example.com",
			username: "testuser",
			password: "testpass",
			wantErr:  false,
		},
		{
			name:        "empty base URL",
			baseURL:     "",
			username:    "testuser", 
			password:    "testpass",
			wantErr:     true,
			errContains: "base URL cannot be empty",
		},
		{
			name:     "invalid URL",
			baseURL:  "://invalid-url",
			username: "testuser",
			password: "testpass", 
			wantErr:  true,
			errContains: "invalid base URL",
		},
		{
			name:     "URL without scheme gets https",
			baseURL:  "gitea.example.com",
			username: "testuser",
			password: "testpass",
			wantErr:  false,
		},
		{
			name:     "with timeout option",
			baseURL:  "https://gitea.example.com",
			username: "testuser",
			password: "testpass",
			opts:     []ClientOption{WithTimeout(10 * time.Second)},
			wantErr:  false,
		},
		{
			name:     "with invalid timeout option",
			baseURL:  "https://gitea.example.com",
			username: "testuser",
			password: "testpass",
			opts:     []ClientOption{WithTimeout(-1 * time.Second)},
			wantErr:  true,
			errContains: "timeout must be positive",
		},
		{
			name:     "with retry config option",
			baseURL:  "https://gitea.example.com",
			username: "testuser",
			password: "testpass",
			opts:     []ClientOption{WithRetryConfig(5, 2*time.Second)},
			wantErr:  false,
		},
		{
			name:     "with invalid retry config",
			baseURL:  "https://gitea.example.com",
			username: "testuser",
			password: "testpass",
			opts:     []ClientOption{WithRetryConfig(-1, time.Second)},
			wantErr:  true,
			errContains: "max retries cannot be negative",
		},
		{
			name:     "with insecure option",
			baseURL:  "https://gitea.example.com",
			username: "testuser",
			password: "testpass",
			opts:     []ClientOption{WithInsecure(true)},
			wantErr:  false,
		},
		{
			name:     "with user agent option", 
			baseURL:  "https://gitea.example.com",
			username: "testuser",
			password: "testpass",
			opts:     []ClientOption{WithUserAgent("test-agent/1.0")},
			wantErr:  false,
		},
		{
			name:     "with empty user agent",
			baseURL:  "https://gitea.example.com",
			username: "testuser",
			password: "testpass",
			opts:     []ClientOption{WithUserAgent("")},
			wantErr:  true,
			errContains: "user agent cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewGiteaClient(tt.baseURL, tt.username, tt.password, trustStore, tt.opts...)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewGiteaClient() expected error but got none")
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("NewGiteaClient() error = %v, want error containing %s", err, tt.errContains)
				}
				return
			}
			
			if err != nil {
				t.Errorf("NewGiteaClient() error = %v, want nil", err)
				return
			}
			
			if client == nil {
				t.Errorf("NewGiteaClient() returned nil client")
			}
		})
	}
}

// Test NewGiteaClient with insecure registry environment variable
func TestNewGiteaClient_InsecureRegistryFlag(t *testing.T) {
	trustStore := newMockTrustStore()
	
	// Save original env var
	origInsecure := os.Getenv("IDPBUILDER_INSECURE_REGISTRY")
	defer os.Setenv("IDPBUILDER_INSECURE_REGISTRY", origInsecure)
	
	// Set insecure flag
	os.Setenv("IDPBUILDER_INSECURE_REGISTRY", "true")
	
	client, err := NewGiteaClient("https://gitea.example.com", "user", "pass", trustStore)
	if err != nil {
		t.Fatalf("NewGiteaClient() error = %v, want nil", err)
	}
	
	if client == nil {
		t.Fatal("NewGiteaClient() returned nil client")
	}
	
	// Verify that insecure registry was set in trust store
	if !trustStore.insecureRegistries["https://gitea.example.com"] {
		t.Error("Expected registry to be marked as insecure")
	}
}

// Test NewGiteaClient with trust store failures
func TestNewGiteaClient_TrustStoreFailures(t *testing.T) {
	tests := []struct {
		name           string
		failCreate     bool
		failConfig     bool
		expectedError  string
	}{
		{
			name:          "trust store create failure",
			failCreate:    true,
			expectedError: "failed to configure transport",
		},
		{
			name:          "trust store config failure", 
			failConfig:    true,
			expectedError: "failed to configure transport",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trustStore := newMockTrustStore()
			trustStore.shouldFailCreate = tt.failCreate
			trustStore.shouldFailConfig = tt.failConfig
			
			_, err := NewGiteaClient("https://gitea.example.com", "user", "pass", trustStore)
			if err == nil {
				t.Errorf("NewGiteaClient() expected error but got none")
				return
			}
			
			if !strings.Contains(err.Error(), tt.expectedError) {
				t.Errorf("NewGiteaClient() error = %v, want error containing %s", err, tt.expectedError)
			}
		})
	}
}

// Test Push method
func TestGiteaClient_Push(t *testing.T) {
	trustStore := newMockTrustStore()
	client, err := NewGiteaClient("https://gitea.example.com", "user", "pass", trustStore)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	ctx := context.Background()
	
	// Create a test image
	img, err := random.Image(1024, 1)
	if err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}
	
	// Explicit cast to ensure v1 import is used
	_ = v1.Image(img)
	
	tests := []struct {
		name        string
		ref         string
		opts        PushOptions
		wantErr     bool
		errContains string
	}{
		{
			name:        "invalid reference",
			ref:         "invalid::reference",
			wantErr:     true,
			errContains: "invalid_reference",
		},
		{
			name:    "valid reference - will fail with network error", 
			ref:     "gitea.example.com/test/image:latest",
			wantErr: true, // This will fail due to network/auth, but we're testing the flow
		},
		{
			name: "with timeout",
			ref:  "gitea.example.com/test/image:latest", 
			opts: PushOptions{
				Options: Options{
					Timeout: 1 * time.Second,
				},
			},
			wantErr: true, // This will fail due to network/auth
		},
		{
			name: "with insecure option",
			ref:  "gitea.example.com/test/image:latest",
			opts: PushOptions{
				Options: Options{
					Insecure: true,
				},
			},
			wantErr: true, // This will fail due to network/auth
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Push(ctx, img, tt.ref, tt.opts)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("Push() expected error but got none")
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Push() error = %v, want error containing %s", err, tt.errContains)
				}
				return
			}
			
			if err != nil {
				t.Errorf("Push() error = %v, want nil", err)
			}
		})
	}
}

// Test Pull method
func TestGiteaClient_Pull(t *testing.T) {
	trustStore := newMockTrustStore()
	client, err := NewGiteaClient("https://gitea.example.com", "user", "pass", trustStore)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	ctx := context.Background()
	
	tests := []struct {
		name        string
		ref         string
		opts        PullOptions
		wantErr     bool
		errContains string
	}{
		{
			name:        "invalid reference",
			ref:         "invalid::reference",
			wantErr:     true,
			errContains: "invalid_reference",
		},
		{
			name:    "valid reference - will fail with network error",
			ref:     "gitea.example.com/test/image:latest",
			wantErr: true, // This will fail due to network/auth
		},
		{
			name: "with timeout",
			ref:  "gitea.example.com/test/image:latest",
			opts: PullOptions{
				Options: Options{
					Timeout: 1 * time.Second,
				},
			},
			wantErr: true, // This will fail due to network/auth
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Pull(ctx, tt.ref, tt.opts)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("Pull() expected error but got none")
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Pull() error = %v, want error containing %s", err, tt.errContains)
				}
				return
			}
			
			if err != nil {
				t.Errorf("Pull() error = %v, want nil", err)
			}
		})
	}
}

// Test Catalog method
func TestGiteaClient_Catalog(t *testing.T) {
	trustStore := newMockTrustStore()
	client, err := NewGiteaClient("https://gitea.example.com", "user", "pass", trustStore)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	ctx := context.Background()
	
	// This will fail with network error but tests the flow
	_, err = client.Catalog(ctx)
	if err == nil {
		t.Errorf("Catalog() expected error but got none (should fail with network error)")
	}
}

// Test Tags method
func TestGiteaClient_Tags(t *testing.T) {
	trustStore := newMockTrustStore()
	client, err := NewGiteaClient("https://gitea.example.com", "user", "pass", trustStore)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	ctx := context.Background()
	
	tests := []struct {
		name        string
		repository  string
		wantErr     bool
		errContains string
	}{
		{
			name:        "invalid repository",
			repository:  "invalid::repo",
			wantErr:     true,
			errContains: "invalid_repository",
		},
		{
			name:       "valid repository - will fail with network error",
			repository: "gitea.example.com/test/repo",
			wantErr:    true, // This will fail due to network
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Tags(ctx, tt.repository)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("Tags() expected error but got none")
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Tags() error = %v, want error containing %s", err, tt.errContains)
				}
				return
			}
			
			if err != nil {
				t.Errorf("Tags() error = %v, want nil", err)
			}
		})
	}
}

// Test Close method
func TestGiteaClient_Close(t *testing.T) {
	trustStore := newMockTrustStore()
	client, err := NewGiteaClient("https://gitea.example.com", "user", "pass", trustStore)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	err = client.Close()
	if err != nil {
		t.Errorf("Close() error = %v, want nil", err)
	}
}

// Test configureAuth function
func TestConfigureAuth(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		wantType string
	}{
		{
			name:     "empty credentials",
			username: "",
			password: "",
			wantType: "*authn.Anonymous",
		},
		{
			name:     "valid credentials",
			username: "testuser",
			password: "testpass",
			wantType: "*authn.Basic",
		},
		{
			name:     "invalid credentials - short username",
			username: "a",
			password: "testpass",
			wantType: "*authn.Anonymous", // Falls back to anonymous
		},
		{
			name:     "invalid credentials - short password",
			username: "testuser",
			password: "abc",
			wantType: "*authn.Anonymous", // Falls back to anonymous
		},
		{
			name:     "invalid credentials - colon in username",
			username: "test:user",
			password: "testpass",
			wantType: "*authn.Anonymous", // Falls back to anonymous
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We can't directly test configureAuth since it's not exported,
			// but we can test it indirectly through NewGiteaClient
			trustStore := newMockTrustStore()
			client, err := NewGiteaClient("https://gitea.example.com", tt.username, tt.password, trustStore)
			if err != nil {
				t.Errorf("NewGiteaClient() error = %v, want nil", err)
				return
			}
			if client == nil {
				t.Error("NewGiteaClient() returned nil client")
			}
		})
	}
}

// Test GetAuthToken function
func TestGetAuthToken(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		want     string
	}{
		{
			name:     "empty credentials",
			username: "",
			password: "",
			want:     "",
		},
		{
			name:     "valid credentials",
			username: "testuser",
			password: "testpass",
			want:     "dGVzdHVzZXI6dGVzdHBhc3M=", // base64 of "testuser:testpass"
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetAuthToken(tt.username, tt.password)
			if got != tt.want {
				t.Errorf("GetAuthToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test LoadCredentialsFromEnv function
func TestLoadCredentialsFromEnv(t *testing.T) {
	// Save original environment variables
	origVars := map[string]string{
		"REGISTRY_USERNAME":                os.Getenv("REGISTRY_USERNAME"),
		"REGISTRY_PASSWORD":                os.Getenv("REGISTRY_PASSWORD"),
		"GITEA_USERNAME":                   os.Getenv("GITEA_USERNAME"),
		"GITEA_PASSWORD":                   os.Getenv("GITEA_PASSWORD"),
		"IDPBUILDER_REGISTRY_USERNAME":     os.Getenv("IDPBUILDER_REGISTRY_USERNAME"),
		"IDPBUILDER_REGISTRY_PASSWORD":     os.Getenv("IDPBUILDER_REGISTRY_PASSWORD"),
	}
	
	// Restore environment variables after test
	defer func() {
		for k, v := range origVars {
			os.Setenv(k, v)
		}
	}()
	
	tests := []struct {
		name         string
		envVars      map[string]string
		wantUsername string
		wantPassword string
	}{
		{
			name:         "no environment variables",
			envVars:      map[string]string{},
			wantUsername: "",
			wantPassword: "",
		},
		{
			name: "registry environment variables",
			envVars: map[string]string{
				"REGISTRY_USERNAME": "reguser",
				"REGISTRY_PASSWORD": "regpass",
			},
			wantUsername: "reguser",
			wantPassword: "regpass",
		},
		{
			name: "gitea environment variables",
			envVars: map[string]string{
				"GITEA_USERNAME": "giteauser",
				"GITEA_PASSWORD": "giteapass",
			},
			wantUsername: "giteauser",
			wantPassword: "giteapass",
		},
		{
			name: "idpbuilder environment variables",
			envVars: map[string]string{
				"IDPBUILDER_REGISTRY_USERNAME": "idpuser",
				"IDPBUILDER_REGISTRY_PASSWORD": "idppass",
			},
			wantUsername: "idpuser",
			wantPassword: "idppass",
		},
		{
			name: "priority order - registry wins",
			envVars: map[string]string{
				"REGISTRY_USERNAME":                "reguser",
				"REGISTRY_PASSWORD":                "regpass",
				"GITEA_USERNAME":                   "giteauser",
				"GITEA_PASSWORD":                   "giteapass",
				"IDPBUILDER_REGISTRY_USERNAME":     "idpuser",
				"IDPBUILDER_REGISTRY_PASSWORD":     "idppass",
			},
			wantUsername: "reguser",
			wantPassword: "regpass",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all environment variables
			for k := range origVars {
				os.Unsetenv(k)
			}
			
			// Set test environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}
			
			username, password := LoadCredentialsFromEnv()
			if username != tt.wantUsername {
				t.Errorf("LoadCredentialsFromEnv() username = %v, want %v", username, tt.wantUsername)
			}
			if password != tt.wantPassword {
				t.Errorf("LoadCredentialsFromEnv() password = %v, want %v", password, tt.wantPassword)
			}
		})
	}
}

// Test NewAuthConfig function
func TestNewAuthConfig(t *testing.T) {
	tests := []struct {
		name         string
		username     string
		password     string
		wantAnonymous bool
		wantToken    string
	}{
		{
			name:         "empty credentials",
			username:     "",
			password:     "",
			wantAnonymous: true,
			wantToken:    "",
		},
		{
			name:         "valid credentials",
			username:     "testuser",
			password:     "testpass", 
			wantAnonymous: false,
			wantToken:    "dGVzdHVzZXI6dGVzdHBhc3M=",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewAuthConfig(tt.username, tt.password)
			if config == nil {
				t.Fatal("NewAuthConfig() returned nil")
			}
			
			if config.Username != tt.username {
				t.Errorf("NewAuthConfig() Username = %v, want %v", config.Username, tt.username)
			}
			if config.Password != tt.password {
				t.Errorf("NewAuthConfig() Password = %v, want %v", config.Password, tt.password)
			}
			if config.Anonymous != tt.wantAnonymous {
				t.Errorf("NewAuthConfig() Anonymous = %v, want %v", config.Anonymous, tt.wantAnonymous)
			}
			if config.Token != tt.wantToken {
				t.Errorf("NewAuthConfig() Token = %v, want %v", config.Token, tt.wantToken)
			}
		})
	}
}

// Test AuthConfig.ToAuthenticator method
func TestAuthConfig_ToAuthenticator(t *testing.T) {
	tests := []struct {
		name     string
		config   *AuthConfig
		wantType string
	}{
		{
			name:     "anonymous config",
			config:   &AuthConfig{Anonymous: true},
			wantType: "*authn.Anonymous",
		},
		{
			name: "basic auth config",
			config: &AuthConfig{
				Username:  "testuser",
				Password:  "testpass",
				Anonymous: false,
			},
			wantType: "*authn.Basic",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := tt.config.ToAuthenticator()
			if auth == nil {
				t.Fatal("ToAuthenticator() returned nil")
			}
			
			// We can't directly compare types easily, but we can verify the authenticator works
			// by checking if it's anonymous or not
			if tt.config.Anonymous {
				if auth != authn.Anonymous {
					t.Error("ToAuthenticator() should return authn.Anonymous for anonymous config")
				}
			}
		})
	}
}

// Test AuthConfig.Validate method
func TestAuthConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *AuthConfig
		wantErr bool
	}{
		{
			name:    "anonymous config is valid",
			config:  &AuthConfig{Anonymous: true},
			wantErr: false,
		},
		{
			name: "valid basic auth config",
			config: &AuthConfig{
				Username:  "testuser",
				Password:  "testpass",
				Anonymous: false,
			},
			wantErr: false,
		},
		{
			name: "invalid - empty username",
			config: &AuthConfig{
				Username:  "",
				Password:  "testpass",
				Anonymous: false,
			},
			wantErr: true,
		},
		{
			name: "invalid - empty password",
			config: &AuthConfig{
				Username:  "testuser",
				Password:  "",
				Anonymous: false,
			},
			wantErr: true,
		},
		{
			name: "invalid - short username",
			config: &AuthConfig{
				Username:  "a",
				Password:  "testpass",
				Anonymous: false,
			},
			wantErr: true,
		},
		{
			name: "invalid - short password",
			config: &AuthConfig{
				Username:  "testuser",
				Password:  "abc",
				Anonymous: false,
			},
			wantErr: true,
		},
		{
			name: "invalid - colon in username",
			config: &AuthConfig{
				Username:  "test:user",
				Password:  "testpass",
				Anonymous: false,
			},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test AuthConfig.Clone method
func TestAuthConfig_Clone(t *testing.T) {
	original := &AuthConfig{
		Username:  "testuser",
		Password:  "testpass",
		Token:     "testtoken",
		Anonymous: false,
	}
	
	clone := original.Clone()
	
	if clone == nil {
		t.Fatal("Clone() returned nil")
	}
	
	// Verify values are copied
	if clone.Username != original.Username {
		t.Errorf("Clone() Username = %v, want %v", clone.Username, original.Username)
	}
	if clone.Password != original.Password {
		t.Errorf("Clone() Password = %v, want %v", clone.Password, original.Password)
	}
	if clone.Token != original.Token {
		t.Errorf("Clone() Token = %v, want %v", clone.Token, original.Token)
	}
	if clone.Anonymous != original.Anonymous {
		t.Errorf("Clone() Anonymous = %v, want %v", clone.Anonymous, original.Anonymous)
	}
	
	// Verify it's a different instance
	if clone == original {
		t.Error("Clone() returned same instance, should be different")
	}
	
	// Verify modifying clone doesn't affect original
	clone.Username = "modified"
	if original.Username == "modified" {
		t.Error("Modifying clone affected original")
	}
}

// Test ClientOptions validation
func TestClientOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		options ClientOptions
		wantErr bool
		errContains string
	}{
		{
			name: "valid options",
			options: ClientOptions{
				BaseURL: "https://gitea.example.com",
				Timeout: 30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "empty base URL",
			options: ClientOptions{
				BaseURL: "",
				Timeout: 30 * time.Second,
			},
			wantErr: true,
			errContains: "base URL cannot be empty",
		},
		{
			name: "invalid timeout",
			options: ClientOptions{
				BaseURL: "https://gitea.example.com",
				Timeout: -1 * time.Second,
			},
			wantErr: true,
			errContains: "timeout must be positive",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.options.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ClientOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if tt.wantErr && tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("ClientOptions.Validate() error = %v, want error containing %s", err, tt.errContains)
			}
		})
	}
}

// Test NewDefaultOptions function
func TestNewDefaultOptions(t *testing.T) {
	options := NewDefaultOptions()
	
	if options.BaseURL == "" {
		t.Error("NewDefaultOptions() BaseURL should not be empty")
	}
	if options.UserAgent == "" {
		t.Error("NewDefaultOptions() UserAgent should not be empty")
	}
	if options.Timeout <= 0 {
		t.Error("NewDefaultOptions() Timeout should be positive")
	}
	
	// Validate the options
	if err := options.Validate(); err != nil {
		t.Errorf("NewDefaultOptions() created invalid options: %v", err)
	}
}

// Test ClientError error interface implementation
func TestClientError_Error(t *testing.T) {
	err := &ClientError{
		Code:    "test_error",
		Message: "test error message",
		Details: map[string]interface{}{"key": "value"},
	}
	
	if err.Error() != "test error message" {
		t.Errorf("ClientError.Error() = %v, want %v", err.Error(), "test error message")
	}
}

// Benchmark tests for performance validation
func BenchmarkNewGiteaClient(b *testing.B) {
	trustStore := newMockTrustStore()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewGiteaClient("https://gitea.example.com", "user", "pass", trustStore)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGetAuthToken(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetAuthToken("testuser", "testpass")
	}
}

func BenchmarkNewAuthConfig(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewAuthConfig("testuser", "testpass")
	}
}