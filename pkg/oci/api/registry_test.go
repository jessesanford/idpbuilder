package api

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// MockRegistryClient implements RegistryClient for testing
type MockRegistryClient struct {
	PushFunc             func(context.Context, string, AuthConfig) error
	PullFunc             func(context.Context, string, AuthConfig) (*Image, error)
	GetManifestFunc      func(context.Context, string) (*Manifest, error)
	ListTagsFunc         func(context.Context, string) ([]string, error)
	DeleteFunc           func(context.Context, string, AuthConfig) error
	PingFunc             func(context.Context) error
	GetRegistryInfoFunc  func(context.Context) (*RegistryInfo, error)
	ListRepositoriesFunc func(context.Context, AuthConfig) ([]string, error)
	CopyImageFunc        func(context.Context, string, string, AuthConfig) error
	GetImageHistoryFunc  func(context.Context, string) ([]*HistoryEntry, error)
}

func (m *MockRegistryClient) Push(ctx context.Context, image string, auth AuthConfig) error {
	if m.PushFunc != nil {
		return m.PushFunc(ctx, image, auth)
	}
	return nil
}

func (m *MockRegistryClient) Pull(ctx context.Context, image string, auth AuthConfig) (*Image, error) {
	if m.PullFunc != nil {
		return m.PullFunc(ctx, image, auth)
	}
	return &Image{}, nil
}

func (m *MockRegistryClient) GetManifest(ctx context.Context, image string) (*Manifest, error) {
	if m.GetManifestFunc != nil {
		return m.GetManifestFunc(ctx, image)
	}
	return &Manifest{}, nil
}

func (m *MockRegistryClient) ListTags(ctx context.Context, repository string) ([]string, error) {
	if m.ListTagsFunc != nil {
		return m.ListTagsFunc(ctx, repository)
	}
	return []string{"latest", "v1.0.0"}, nil
}

func (m *MockRegistryClient) Delete(ctx context.Context, image string, auth AuthConfig) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, image, auth)
	}
	return nil
}

func (m *MockRegistryClient) Ping(ctx context.Context) error {
	if m.PingFunc != nil {
		return m.PingFunc(ctx)
	}
	return nil
}

func (m *MockRegistryClient) GetRegistryInfo(ctx context.Context) (*RegistryInfo, error) {
	if m.GetRegistryInfoFunc != nil {
		return m.GetRegistryInfoFunc(ctx)
	}
	return &RegistryInfo{}, nil
}

func (m *MockRegistryClient) ListRepositories(ctx context.Context, auth AuthConfig) ([]string, error) {
	if m.ListRepositoriesFunc != nil {
		return m.ListRepositoriesFunc(ctx, auth)
	}
	return []string{"myapp", "myservice"}, nil
}

func (m *MockRegistryClient) CopyImage(ctx context.Context, source, destination string, auth AuthConfig) error {
	if m.CopyImageFunc != nil {
		return m.CopyImageFunc(ctx, source, destination, auth)
	}
	return nil
}

func (m *MockRegistryClient) GetImageHistory(ctx context.Context, image string) ([]*HistoryEntry, error) {
	if m.GetImageHistoryFunc != nil {
		return m.GetImageHistoryFunc(ctx, image)
	}
	return []*HistoryEntry{}, nil
}

// Test interface compliance
func TestRegistryClientInterface(t *testing.T) {
	var _ RegistryClient = &MockRegistryClient{}
}

func TestAuthConfigValidate(t *testing.T) {
	tests := []struct {
		name      string
		auth      AuthConfig
		expectErr bool
	}{
		{
			name: "valid basic auth",
			auth: AuthConfig{
				Username:      "user",
				Password:      "pass",
				ServerAddress: "https://registry.example.com",
			},
			expectErr: false,
		},
		{
			name: "valid token auth",
			auth: AuthConfig{
				IdentityToken: "token123",
				ServerAddress: "https://registry.example.com",
			},
			expectErr: false,
		},
		{
			name: "missing server address",
			auth: AuthConfig{
				Username: "user",
				Password: "pass",
			},
			expectErr: true,
		},
		{
			name: "no auth method",
			auth: AuthConfig{
				ServerAddress: "https://registry.example.com",
			},
			expectErr: true,
		},
		{
			name: "invalid server address",
			auth: AuthConfig{
				Username:      "user",
				Password:      "pass",
				ServerAddress: "registry.example.com",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.auth.Validate()
			if (err != nil) != tt.expectErr {
				t.Errorf("Expected error: %v, got: %v", tt.expectErr, err)
			}
		})
	}
}

func TestAuthConfigIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		auth     AuthConfig
		expected bool
	}{
		{
			name:     "empty auth",
			auth:     AuthConfig{},
			expected: true,
		},
		{
			name: "has username",
			auth: AuthConfig{
				Username: "user",
			},
			expected: false,
		},
		{
			name: "has token",
			auth: AuthConfig{
				IdentityToken: "token",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.auth.IsEmpty()
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}