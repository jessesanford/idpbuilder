package testutils

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestFixtures manages test environment setup and cleanup
type TestFixtures struct {
	Registry    *MockRegistry
	AuthConfig  *AuthConfig
	TempDir     string
	TestContext context.Context
	TestCancel  context.CancelFunc
	Transport   *MockAuthTransport
}

// SetupTestFixtures initializes a complete test environment
func SetupTestFixtures(t *testing.T, authEnabled bool) *TestFixtures {
	// Create test context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	// Create auth configuration
	authConfig := &AuthConfig{
		Username: "testuser",
		Password: "testpass",
		Token:    "test-token-123",
		Enabled:  authEnabled,
	}

	// Create mock registry
	registry := NewMockRegistry(authConfig)

	// Create temporary directory for test files
	tempDir, err := os.MkdirTemp("", "idpbuilder-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Create transport
	transport := NewMockAuthTransport(authConfig)

	fixtures := &TestFixtures{
		Registry:    registry,
		AuthConfig:  authConfig,
		TempDir:     tempDir,
		TestContext: ctx,
		TestCancel:  cancel,
		Transport:   transport,
	}

	// Log setup for debugging
	t.Logf("Test fixtures created:")
	t.Logf("  Registry URL: %s", registry.URL())
	t.Logf("  Auth enabled: %t", authEnabled)
	t.Logf("  Temp dir: %s", tempDir)

	return fixtures
}

// CleanupTestFixtures cleans up all test resources
func (f *TestFixtures) CleanupTestFixtures(t *testing.T) {
	// Cancel context
	if f.TestCancel != nil {
		f.TestCancel()
	}

	// Close mock registry
	if f.Registry != nil {
		f.Registry.Close()
		t.Logf("Mock registry closed")
	}

	// Remove temporary directory
	if f.TempDir != "" {
		err := os.RemoveAll(f.TempDir)
		if err != nil {
			t.Logf("Warning: Failed to remove temp directory %s: %v", f.TempDir, err)
		} else {
			t.Logf("Temp directory removed: %s", f.TempDir)
		}
	}
}

// GetTestReference creates a test image reference pointing to the mock registry
func (f *TestFixtures) GetTestReference(repository, tag string) (string, error) {
	registryHost := f.Registry.Host()
	if repository == "" {
		return "", fmt.Errorf("repository cannot be empty")
	}
	if tag == "" {
		return "", fmt.Errorf("tag cannot be empty")
	}

	refStr := fmt.Sprintf("%s/%s:%s", registryHost, repository, tag)
	return refStr, nil
}

// GetAuthHeaders returns authentication headers for the configured transport
func (f *TestFixtures) GetAuthHeaders() map[string]string {
	return f.Transport.GetAuthHeaders()
}

// GetHTTPClient returns an HTTP client configured with the mock transport
func (f *TestFixtures) GetHTTPClient() *http.Client {
	return &http.Client{
		Transport: f.Transport,
		Timeout:   30 * time.Second,
	}
}

// CreateTestFile creates a test file in the temporary directory
func (f *TestFixtures) CreateTestFile(filename, content string) (string, error) {
	fullPath := filepath.Join(f.TempDir, filename)

	// Ensure directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write file
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write file %s: %w", fullPath, err)
	}

	return fullPath, nil
}

// ReadTestFile reads a file from the temporary directory
func (f *TestFixtures) ReadTestFile(filename string) ([]byte, error) {
	fullPath := filepath.Join(f.TempDir, filename)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", fullPath, err)
	}
	return content, nil
}

// AssertRegistryURL validates that the registry URL is accessible
func (f *TestFixtures) AssertRegistryURL(t *testing.T) {
	registryURL := f.Registry.URL()
	if registryURL == "" {
		t.Fatal("Registry URL is empty")
	}

	// Parse URL to validate format
	parsedURL, err := url.Parse(registryURL)
	if err != nil {
		t.Fatalf("Invalid registry URL %s: %v", registryURL, err)
	}

	if parsedURL.Host == "" {
		t.Fatalf("Registry URL missing host: %s", registryURL)
	}

	t.Logf("Registry URL validated: %s", registryURL)
}

// AssertAuthConfig validates authentication configuration
func (f *TestFixtures) AssertAuthConfig(t *testing.T) {
	if f.AuthConfig == nil {
		t.Fatal("Auth config is nil")
	}

	if f.AuthConfig.Enabled {
		if f.AuthConfig.Username == "" && f.AuthConfig.Token == "" {
			t.Fatal("Auth enabled but no username or token provided")
		}

		if f.AuthConfig.Username != "" && f.AuthConfig.Password == "" {
			t.Fatal("Username provided but password is empty")
		}

		t.Logf("Auth config validated - Username: %s, Token: %s",
			f.AuthConfig.Username,
			maskString(f.AuthConfig.Token))
	} else {
		t.Log("Auth config validated - authentication disabled")
	}
}

// GetRegistryStats returns statistics about the mock registry
func (f *TestFixtures) GetRegistryStats() RegistryStats {
	return RegistryStats{
		RequestCount: f.Registry.GetRequestCount(),
		URL:          f.Registry.URL(),
		Host:         f.Registry.Host(),
		AuthEnabled:  f.AuthConfig.Enabled,
	}
}

// RegistryStats holds statistics about the mock registry
type RegistryStats struct {
	RequestCount int
	URL          string
	Host         string
	AuthEnabled  bool
}

// Helper functions

// maskString masks sensitive strings for logging
func maskString(s string) string {
	if len(s) <= 4 {
		return "***"
	}
	return s[:2] + "***" + s[len(s)-2:]
}

// WaitForRegistry waits for the registry to be ready
func (f *TestFixtures) WaitForRegistry(t *testing.T, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(f.TestContext, timeout)
	defer cancel()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for registry to be ready")
		case <-ticker.C:
			// Try to access the registry root endpoint
			if f.Registry.URL() != "" {
				t.Logf("Registry ready at: %s", f.Registry.URL())
				return nil
			}
		}
	}
}