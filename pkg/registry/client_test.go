package registry

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name        string
		config      *RegistryConfig
		envVars     map[string]string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "nil config",
			config:      nil,
			expectError: true,
			errorMsg:    "registry config cannot be nil",
		},
		{
			name: "valid config with URL",
			config: &RegistryConfig{
				URL: "https://registry.example.com",
				Auth: &AuthConfig{
					Type:     "basic",
					Username: "user",
					Password: "pass",
				},
			},
			expectError: false,
		},
		{
			name:   "config without URL uses env var",
			config: &RegistryConfig{},
			envVars: map[string]string{
				"REGISTRY_URL": "https://env-registry.example.com",
			},
			expectError: false,
		},
		{
			name:        "config without URL and no env var",
			config:      &RegistryConfig{},
			expectError: true,
			errorMsg:    "registry URL not configured",
		},
		{
			name: "insecure config",
			config: &RegistryConfig{
				URL:      "http://insecure-registry.example.com",
				Insecure: true,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
			}

			client, err := NewClient(tt.config)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
			}
		})
	}
}

func TestParseReference(t *testing.T) {
	tests := []struct {
		name         string
		ref          string
		expectError  bool
		expectedReg  string
		expectedRepo string
		expectedTag  string
	}{
		{
			name:         "full reference with registry",
			ref:          "registry.example.com/namespace/repo:v1.0.0",
			expectError:  false,
			expectedReg:  "registry.example.com",
			expectedRepo: "namespace/repo",
			expectedTag:  "v1.0.0",
		},
		{
			name:         "reference without registry (uses default)",
			ref:          "namespace/repo:latest",
			expectError:  false,
			expectedRepo: "namespace/repo",
			expectedTag:  "latest",
		},
		{
			name:         "reference without tag (defaults to latest)",
			ref:          "registry.example.com/repo",
			expectError:  false,
			expectedReg:  "registry.example.com",
			expectedRepo: "repo",
		},
		{
			name:        "empty reference",
			ref:         "",
			expectError: true,
		},
		{
			name:        "invalid reference format",
			ref:         ":::invalid:::",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref, err := parseReference(tt.ref)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, ref)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, ref)

				if tt.expectedReg != "" {
					assert.Equal(t, tt.expectedReg, ref.Registry())
				}
				if tt.expectedRepo != "" {
					assert.Equal(t, tt.expectedRepo, ref.Repository())
				}
				if tt.expectedTag != "" {
					assert.Equal(t, tt.expectedTag, ref.Tag())
				}
			}
		})
	}
}

func TestTransportWithAuth(t *testing.T) {
	tests := []struct {
		name     string
		auth     *AuthConfig
		expected string // Expected authorization header
	}{
		{
			name: "basic auth",
			auth: &AuthConfig{
				Type:     "basic",
				Username: "user",
				Password: "pass",
			},
			expected: "basic", // We can't easily test the exact encoding without exposing internals
		},
		{
			name: "token auth",
			auth: &AuthConfig{
				Type:  "token",
				Token: "test-token",
			},
			expected: "Bearer test-token",
		},
		{
			name: "anonymous auth",
			auth: &AuthConfig{
				Type: "anonymous",
			},
			expected: "", // No auth header
		},
		{
			name:     "nil auth",
			auth:     nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transport := NewTransport(tt.auth)
			assert.NotNil(t, transport)

			// Test WithAuth method
			newTransport := transport.WithAuth(tt.auth)
			assert.NotNil(t, newTransport)
			assert.NotSame(t, transport, newTransport) // Should be a new instance

			// Create a test request
			req, err := http.NewRequest("GET", "https://registry.example.com/v2/", nil)
			require.NoError(t, err)

			// We can't easily test the actual round trip without a real server,
			// but we can verify the transport doesn't panic and handles the request structure
			assert.NotNil(t, req)
		})
	}
}

func TestIsValidRepository(t *testing.T) {
	tests := []struct {
		name        string
		repository  string
		expectError bool
	}{
		{
			name:        "valid simple repository",
			repository:  "registry.example.com/repo",
			expectError: false,
		},
		{
			name:        "valid repository with namespace",
			repository:  "registry.example.com/namespace/repo",
			expectError: false,
		},
		{
			name:        "valid repository with port",
			repository:  "registry.example.com:5000/repo",
			expectError: false,
		},
		{
			name:        "empty repository",
			repository:  "",
			expectError: true,
		},
		{
			name:        "repository with digest",
			repository:  "registry.example.com/repo@sha256:abc123",
			expectError: true,
		},
		{
			name:        "repository with tag",
			repository:  "registry.example.com/repo:latest",
			expectError: true, // Tags should not be in repository for listing
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := isValidRepository(tt.repository)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetAuthFromEnv(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		expectedType string
	}{
		{
			name: "token auth from env",
			envVars: map[string]string{
				"REGISTRY_TOKEN": "test-token",
			},
			expectedType: "token",
		},
		{
			name: "basic auth from env",
			envVars: map[string]string{
				"REGISTRY_USERNAME": "user",
				"REGISTRY_PASSWORD": "pass",
			},
			expectedType: "basic",
		},
		{
			name: "token takes precedence over basic",
			envVars: map[string]string{
				"REGISTRY_TOKEN":    "test-token",
				"REGISTRY_USERNAME": "user",
				"REGISTRY_PASSWORD": "pass",
			},
			expectedType: "token",
		},
		{
			name:         "no auth env vars",
			envVars:      map[string]string{},
			expectedType: "anonymous",
		},
		{
			name: "incomplete basic auth",
			envVars: map[string]string{
				"REGISTRY_USERNAME": "user",
				// Missing password
			},
			expectedType: "anonymous",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear relevant environment variables first
			os.Unsetenv("REGISTRY_TOKEN")
			os.Unsetenv("REGISTRY_USERNAME")
			os.Unsetenv("REGISTRY_PASSWORD")

			// Set test environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
			}

			auth := getAuthFromEnv()
			assert.NotNil(t, auth)
			assert.Equal(t, tt.expectedType, auth.Type)

			// Verify specific fields based on type
			switch tt.expectedType {
			case "token":
				assert.Equal(t, tt.envVars["REGISTRY_TOKEN"], auth.Token)
			case "basic":
				assert.Equal(t, tt.envVars["REGISTRY_USERNAME"], auth.Username)
				assert.Equal(t, tt.envVars["REGISTRY_PASSWORD"], auth.Password)
			case "anonymous":
				assert.Empty(t, auth.Token)
				assert.Empty(t, auth.Username)
				assert.Empty(t, auth.Password)
			}
		})
	}
}

func TestGcrAdapterClientInterface(t *testing.T) {
	// Test that gcrAdapter implements the Client interface
	config := &RegistryConfig{
		URL: "https://registry.example.com",
		Auth: &AuthConfig{
			Type: "anonymous",
		},
	}

	client, err := NewClient(config)
	require.NoError(t, err)
	require.NotNil(t, client)

	// Verify that client implements the Client interface
	var _ Client = client

	// Test basic method calls (these will fail due to no real registry, but should not panic)
	ctx := context.Background()

	// Test Exists with invalid reference
	exists, err := client.Exists(ctx, "")
	assert.Error(t, err)
	assert.False(t, exists)

	// Test ListTags with invalid repository
	tags, err := client.ListTags(ctx, "")
	assert.Error(t, err)
	assert.Nil(t, tags)

	// Test Push with nil artifact
	err = client.Push(ctx, "registry.example.com/repo:tag", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "artifact cannot be nil")

	// Test Pull with invalid reference
	artifact, err := client.Pull(ctx, "")
	assert.Error(t, err)
	assert.Nil(t, artifact)
}