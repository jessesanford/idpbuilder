package testutils

import (
	"net/http"
	"strings"
	"testing"
	"time"
)

// TestMockRegistryCreation validates that mock registry can be created and configured
func TestMockRegistryCreation(t *testing.T) {
	tests := []struct {
		name        string
		authEnabled bool
		username    string
		password    string
		token       string
	}{
		{
			name:        "registry without auth",
			authEnabled: false,
		},
		{
			name:        "registry with basic auth",
			authEnabled: true,
			username:    "testuser",
			password:    "testpass",
		},
		{
			name:        "registry with token auth",
			authEnabled: true,
			token:       "test-token-123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create auth config
			authConfig := &AuthConfig{
				Username: tt.username,
				Password: tt.password,
				Token:    tt.token,
				Enabled:  tt.authEnabled,
			}

			// Create mock registry
			registry := NewMockRegistry(authConfig)
			defer registry.Close()

			// Validate registry properties
			if registry.URL() == "" {
				t.Error("Registry URL should not be empty")
			}

			if registry.Host() == "" {
				t.Error("Registry host should not be empty")
			}

			// Test registry endpoint
			resp, err := http.Get(registry.URL() + "/v2/")
			if err != nil {
				if tt.authEnabled {
					// Auth-enabled registry should reject unauthorized requests
					t.Logf("Expected auth failure for auth-enabled registry: %v", err)
				} else {
					t.Errorf("Failed to connect to registry: %v", err)
				}
			} else {
				defer resp.Body.Close()

				if tt.authEnabled && resp.StatusCode != http.StatusUnauthorized {
					t.Errorf("Expected 401 Unauthorized for auth-enabled registry, got %d", resp.StatusCode)
				} else if !tt.authEnabled && resp.StatusCode != http.StatusOK {
					t.Errorf("Expected 200 OK for non-auth registry, got %d", resp.StatusCode)
				}
			}

			t.Logf("Registry created successfully at %s", registry.URL())
		})
	}
}

// TestAuthTransport validates the mock authentication transport
func TestAuthTransport(t *testing.T) {
	tests := []struct {
		name     string
		config   *AuthConfig
		wantAuth bool
	}{
		{
			name: "no auth configured",
			config: &AuthConfig{
				Enabled: false,
			},
			wantAuth: false,
		},
		{
			name: "basic auth configured",
			config: &AuthConfig{
				Username: "testuser",
				Password: "testpass",
				Enabled:  true,
			},
			wantAuth: true,
		},
		{
			name: "token auth configured",
			config: &AuthConfig{
				Token:   "test-token-123",
				Enabled: true,
			},
			wantAuth: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create auth transport
			transport := NewMockAuthTransport(tt.config)

			// Test GetAuthHeaders method
			headers := transport.GetAuthHeaders()

			if tt.wantAuth {
				if tt.config.Token != "" {
					expectedAuth := "Bearer " + tt.config.Token
					if headers["Authorization"] != expectedAuth {
						t.Errorf("Expected Authorization header %s, got %s", expectedAuth, headers["Authorization"])
					}
				}
				if tt.config.Username != "" {
					if headers["X-Username"] != tt.config.Username {
						t.Errorf("Expected X-Username header %s, got %s", tt.config.Username, headers["X-Username"])
					}
					if headers["X-Password"] != tt.config.Password {
						t.Errorf("Expected X-Password header %s, got %s", tt.config.Password, headers["X-Password"])
					}
				}
			} else {
				if len(headers) > 0 {
					t.Error("Expected empty auth headers for disabled auth")
				}
			}

			// Verify the transport was created successfully
			if transport == nil {
				t.Error("Transport should not be nil")
			}

			t.Logf("Auth transport created successfully for config: %+v", tt.config)
		})
	}
}

// TestTestFixturesSetup validates the test fixtures setup and cleanup
func TestTestFixturesSetup(t *testing.T) {
	tests := []struct {
		name        string
		authEnabled bool
	}{
		{
			name:        "fixtures without auth",
			authEnabled: false,
		},
		{
			name:        "fixtures with auth",
			authEnabled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup fixtures
			fixtures := SetupTestFixtures(t, tt.authEnabled)

			// Validate fixtures
			if fixtures == nil {
				t.Fatal("Fixtures should not be nil")
			}

			if fixtures.Registry == nil {
				t.Error("Registry should not be nil")
			}

			if fixtures.AuthConfig == nil {
				t.Error("AuthConfig should not be nil")
			}

			if fixtures.TempDir == "" {
				t.Error("TempDir should not be empty")
			}

			if fixtures.TestContext == nil {
				t.Error("TestContext should not be nil")
			}

			if fixtures.TestCancel == nil {
				t.Error("TestCancel should not be nil")
			}

			if fixtures.Transport == nil {
				t.Error("Transport should not be nil")
			}

			// Validate auth configuration
			if fixtures.AuthConfig.Enabled != tt.authEnabled {
				t.Errorf("Expected auth enabled = %t, got %t", tt.authEnabled, fixtures.AuthConfig.Enabled)
			}

			// Test registry access
			fixtures.AssertRegistryURL(t)
			fixtures.AssertAuthConfig(t)

			// Test file operations
			testContent := "test file content"
			filename, err := fixtures.CreateTestFile("test.txt", testContent)
			if err != nil {
				t.Errorf("Failed to create test file: %v", err)
			} else {
				content, err := fixtures.ReadTestFile("test.txt")
				if err != nil {
					t.Errorf("Failed to read test file: %v", err)
				} else if string(content) != testContent {
					t.Errorf("Expected content %s, got %s", testContent, string(content))
				}
				t.Logf("Test file operations successful: %s", filename)
			}

			// Test registry statistics
			stats := fixtures.GetRegistryStats()
			if stats.URL != fixtures.Registry.URL() {
				t.Errorf("Stats URL mismatch: expected %s, got %s", fixtures.Registry.URL(), stats.URL)
			}

			if stats.AuthEnabled != tt.authEnabled {
				t.Errorf("Stats auth enabled mismatch: expected %t, got %t", tt.authEnabled, stats.AuthEnabled)
			}

			// Test wait for registry
			err = fixtures.WaitForRegistry(t, 5*time.Second)
			if err != nil {
				t.Errorf("Failed to wait for registry: %v", err)
			}

			// Cleanup fixtures
			fixtures.CleanupTestFixtures(t)

			t.Logf("Test fixtures validated successfully")
		})
	}
}

// TestGetTestReference validates reference creation
func TestGetTestReference(t *testing.T) {
	fixtures := SetupTestFixtures(t, false)
	defer fixtures.CleanupTestFixtures(t)

	tests := []struct {
		name       string
		repository string
		tag        string
		wantError  bool
	}{
		{
			name:       "valid reference",
			repository: "test/repo",
			tag:        "latest",
			wantError:  false,
		},
		{
			name:       "valid reference with org",
			repository: "myorg/myrepo",
			tag:        "v1.0.0",
			wantError:  false,
		},
		{
			name:       "empty repository",
			repository: "",
			tag:        "latest",
			wantError:  true,
		},
		{
			name:       "empty tag",
			repository: "test/repo",
			tag:        "",
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref, err := fixtures.GetTestReference(tt.repository, tt.tag)

			if tt.wantError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}

				if ref == "" {
					t.Error("Reference should not be empty")
					return
				}

				// Validate reference contains registry host
				registryHost := fixtures.Registry.Host()
				if !strings.Contains(ref, registryHost) {
					t.Errorf("Reference %s should contain registry host %s", ref, registryHost)
				}

				t.Logf("Reference created successfully: %s", ref)
			}
		})
	}
}

// TestHTTPClient validates HTTP client configuration
func TestHTTPClient(t *testing.T) {
	fixtures := SetupTestFixtures(t, true)
	defer fixtures.CleanupTestFixtures(t)

	client := fixtures.GetHTTPClient()

	if client == nil {
		t.Error("HTTP client should not be nil")
	}

	if client.Transport == nil {
		t.Error("HTTP client transport should not be nil")
	}

	if client.Timeout == 0 {
		t.Error("HTTP client timeout should be set")
	}

	// Test auth headers
	headers := fixtures.GetAuthHeaders()
	if len(headers) == 0 {
		t.Error("Auth headers should not be empty for auth-enabled fixtures")
	}

	t.Logf("HTTP client created successfully with timeout: %v", client.Timeout)
	t.Logf("Auth headers: %d configured", len(headers))
}