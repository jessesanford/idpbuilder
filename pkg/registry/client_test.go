package registry

import (
	"context"
	"testing"

	"github.com/google/go-containerregistry/pkg/v1/remote"

	"github.com/cnoe-io/idpbuilder/pkg/config"
	"github.com/cnoe-io/idpbuilder/pkg/providers"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.RegistryConfig
		expectError bool
		errorMsg    string
	}{
		{
			name:        "nil config",
			config:      nil,
			expectError: true,
			errorMsg:    "registry configuration cannot be nil",
		},
		{
			name: "empty URL",
			config: &config.RegistryConfig{
				URL: "",
			},
			expectError: true,
			errorMsg:    "registry URL cannot be empty",
		},
		{
			name: "valid config with URL only",
			config: &config.RegistryConfig{
				URL: "https://registry.example.com",
			},
			expectError: false,
		},
		{
			name: "valid config with basic auth",
			config: &config.RegistryConfig{
				URL: "https://registry.example.com",
				Auth: config.AuthConfig{
					Username: "testuser",
					Password: "testpass",
				},
			},
			expectError: false,
		},
		{
			name: "valid config with token auth",
			config: &config.RegistryConfig{
				URL: "https://registry.example.com",
				Auth: config.AuthConfig{
					Token: "test-token",
				},
			},
			expectError: false,
		},
		{
			name: "valid config with insecure",
			config: &config.RegistryConfig{
				URL:      "http://registry.example.com",
				Insecure: true,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("expected error %q, got %q", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if client == nil {
				t.Error("expected non-nil client")
				return
			}

			// Verify client configuration
			if client.config != tt.config {
				t.Error("client config doesn't match input config")
			}

			if client.httpClient == nil {
				t.Error("HTTP client should be initialized")
			}

			if client.options == nil {
				t.Error("options should be initialized")
			}
		})
	}
}

func TestClient_ConfigureAuth(t *testing.T) {
	tests := []struct {
		name   string
		config *config.RegistryConfig
	}{
		{
			name: "no auth",
			config: &config.RegistryConfig{
				URL: "https://registry.example.com",
			},
		},
		{
			name: "basic auth",
			config: &config.RegistryConfig{
				URL: "https://registry.example.com",
				Auth: config.AuthConfig{
					Username: "user",
					Password: "pass",
				},
			},
		},
		{
			name: "token auth",
			config: &config.RegistryConfig{
				URL: "https://registry.example.com",
				Auth: config.AuthConfig{
					Token: "token123",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				config:  tt.config,
				options: []remote.Option{},
			}

			err := client.configureAuth()
			if err != nil {
				t.Errorf("configureAuth() error = %v", err)
				return
			}

			// Verify auth is configured
			if client.auth == nil {
				t.Error("auth should be configured")
			}

			// Verify options include auth
			if len(client.options) == 0 {
				t.Error("options should include auth configuration")
			}
		})
	}
}

// TestClient_ProviderInterface verifies that Client implements the Provider interface
func TestClient_ProviderInterface(t *testing.T) {
	var _ providers.Provider = (*Client)(nil)
}

func TestClient_PushErrors(t *testing.T) {
	config := &config.RegistryConfig{
		URL: "https://registry.example.com",
	}

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()
	artifact := providers.Artifact{
		MediaType: "application/vnd.oci.image.manifest.v1+json",
		Manifest:  []byte("{}"),
	}

	tests := []struct {
		name string
		ref  string
	}{
		{
			name: "invalid reference",
			ref:  "invalid ref with spaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Push(ctx, tt.ref, artifact)
			if err == nil {
				t.Error("expected error for invalid reference")
				return
			}

			// Verify it's a ProviderError
			if providerErr, ok := err.(*providers.ProviderError); ok {
				if providerErr.Op != "push" {
					t.Errorf("expected op 'push', got %q", providerErr.Op)
				}
				if providerErr.Ref != tt.ref {
					t.Errorf("expected ref %q, got %q", tt.ref, providerErr.Ref)
				}
			} else {
				t.Errorf("expected ProviderError, got %T", err)
			}
		})
	}
}

func TestClient_PullErrors(t *testing.T) {
	config := &config.RegistryConfig{
		URL: "https://registry.example.com",
	}

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	tests := []struct {
		name string
		ref  string
	}{
		{
			name: "invalid reference",
			ref:  "invalid ref with spaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Pull(ctx, tt.ref)
			if err == nil {
				t.Error("expected error for invalid reference")
				return
			}

			// Verify it's a ProviderError
			if providerErr, ok := err.(*providers.ProviderError); ok {
				if providerErr.Op != "pull" {
					t.Errorf("expected op 'pull', got %q", providerErr.Op)
				}
				if providerErr.Ref != tt.ref {
					t.Errorf("expected ref %q, got %q", tt.ref, providerErr.Ref)
				}
			} else {
				t.Errorf("expected ProviderError, got %T", err)
			}
		})
	}
}

func TestClient_ListErrors(t *testing.T) {
	config := &config.RegistryConfig{
		URL: "https://registry.example.com",
	}

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	tests := []struct {
		name       string
		repository string
	}{
		{
			name:       "invalid repository",
			repository: "invalid repo with spaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.List(ctx, tt.repository)
			if err == nil {
				t.Error("expected error for invalid repository")
				return
			}

			// Verify it's a ProviderError
			if providerErr, ok := err.(*providers.ProviderError); ok {
				if providerErr.Op != "list" {
					t.Errorf("expected op 'list', got %q", providerErr.Op)
				}
				if providerErr.Ref != tt.repository {
					t.Errorf("expected ref %q, got %q", tt.repository, providerErr.Ref)
				}
			} else {
				t.Errorf("expected ProviderError, got %T", err)
			}
		})
	}
}

func TestClient_DeleteErrors(t *testing.T) {
	config := &config.RegistryConfig{
		URL: "https://registry.example.com",
	}

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	tests := []struct {
		name string
		ref  string
	}{
		{
			name: "invalid reference",
			ref:  "invalid ref with spaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Delete(ctx, tt.ref)
			if err == nil {
				t.Error("expected error for invalid reference")
				return
			}

			// Verify it's a ProviderError
			if providerErr, ok := err.(*providers.ProviderError); ok {
				if providerErr.Op != "delete" {
					t.Errorf("expected op 'delete', got %q", providerErr.Op)
				}
				if providerErr.Ref != tt.ref {
					t.Errorf("expected ref %q, got %q", tt.ref, providerErr.Ref)
				}
			} else {
				t.Errorf("expected ProviderError, got %T", err)
			}
		})
	}
}