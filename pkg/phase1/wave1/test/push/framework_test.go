package push

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMockRegistryCreation validates mock registry creation and basic functionality
func TestMockRegistryCreation(t *testing.T) {
	// Create mock registry
	registry := NewMockRegistry()
	require.NotNil(t, registry, "Mock registry should be created successfully")

	// Ensure cleanup happens
	defer registry.Cleanup()

	// Verify server is running
	assert.NotNil(t, registry.Server, "HTTP server should be initialized")
	assert.NotEmpty(t, registry.GetURL(), "Registry URL should not be empty")

	// Verify images storage is initialized
	assert.NotNil(t, registry.Images, "Images storage should be initialized")
	assert.Empty(t, registry.Images, "Images storage should start empty")

	// Test registry health check
	resp, err := http.Get(registry.GetURL() + "/v2/")
	require.NoError(t, err, "Registry health check should succeed")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Registry should return 200 OK")
	assert.Equal(t, "registry/2.0", resp.Header.Get("Docker-Distribution-API-Version"),
		"Registry should return correct API version header")

	// Test authentication configuration
	registry.SetAuth("testuser", "testpass", true)
	assert.NotNil(t, registry.AuthConfig, "Auth config should be set")
	assert.True(t, registry.AuthConfig.Required, "Auth should be required when configured")
	assert.Equal(t, "testuser", registry.AuthConfig.Username, "Username should be set correctly")
	assert.Equal(t, "testpass", registry.AuthConfig.Password, "Password should be set correctly")
}

// TestAuthTransport validates the mock authentication transport functionality
func TestAuthTransport(t *testing.T) {
	// Create mock auth transport
	transport := NewMockAuthTransport("testuser", "testpass")
	require.NotNil(t, transport, "Mock auth transport should be created successfully")

	// Verify transport configuration
	assert.Equal(t, "testuser", transport.Username, "Username should be set correctly")
	assert.Equal(t, "testpass", transport.Password, "Password should be set correctly")
	assert.NotNil(t, transport.Base, "Base transport should be set")

	// Create HTTP client with auth transport
	client := &http.Client{
		Transport: transport,
	}

	// Create registry with auth required
	registry := NewMockRegistry()
	defer registry.Cleanup()

	registry.SetAuth("testuser", "testpass", true)

	// Test authenticated request succeeds
	resp, err := client.Get(registry.GetURL() + "/v2/")
	require.NoError(t, err, "Authenticated request should succeed")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode,
		"Authenticated request should return 200 OK")

	// Test unauthenticated client fails
	unauthClient := &http.Client{}
	resp2, err := unauthClient.Get(registry.GetURL() + "/v2/")
	require.NoError(t, err, "Unauthenticated request should not error")
	defer resp2.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp2.StatusCode,
		"Unauthenticated request should return 401 Unauthorized")
}

// TestImageCreation validates test image creation and structure
func TestImageCreation(t *testing.T) {
	// Create test image
	img := CreateTestImage("test-app", "v1.0.0")
	require.NotNil(t, img, "Test image should be created successfully")

	// Verify image digest
	digest, err := img.Digest()
	require.NoError(t, err, "Image should have valid digest")
	assert.Equal(t, "sha256", digest.Algorithm, "Digest should use SHA256")
	assert.NotEmpty(t, digest.Hex, "Digest should have hex value")

	// Verify manifest
	manifest, err := img.Manifest()
	require.NoError(t, err, "Image should have valid manifest")
	assert.Equal(t, int64(2), manifest.SchemaVersion, "Manifest should be schema version 2")
	assert.NotEmpty(t, manifest.MediaType, "Manifest should have media type")
	assert.NotEmpty(t, manifest.Config.Digest.Hex, "Config should have digest")
	assert.True(t, len(manifest.Layers) > 0, "Image should have at least one layer")

	// Verify config file
	config, err := img.ConfigFile()
	require.NoError(t, err, "Image should have valid config file")
	assert.Equal(t, "amd64", config.Architecture, "Config should specify architecture")
	assert.Equal(t, "linux", config.OS, "Config should specify OS")
	assert.Equal(t, "layers", config.RootFS.Type, "RootFS should specify layers type")

	// Verify layers
	layers, err := img.Layers()
	require.NoError(t, err, "Image should have accessible layers")
	assert.True(t, len(layers) > 0, "Image should have at least one layer")

	// Test layer functionality
	layer := layers[0]
	layerDigest, err := layer.Digest()
	require.NoError(t, err, "Layer should have valid digest")
	assert.NotEmpty(t, layerDigest.Hex, "Layer digest should not be empty")

	size, err := layer.Size()
	require.NoError(t, err, "Layer should have valid size")
	assert.True(t, size > 0, "Layer size should be positive")
}

// TestCleanup validates test fixtures cleanup functionality
func TestCleanup(t *testing.T) {
	// Setup test fixtures
	fixtures := SetupTestFixtures(t)
	require.NotNil(t, fixtures, "Test fixtures should be created successfully")

	// Verify fixtures are properly initialized
	assert.NotNil(t, fixtures.Registry, "Registry should be initialized")
	assert.NotNil(t, fixtures.Client, "HTTP client should be initialized")
	assert.NotEmpty(t, fixtures.TempDir, "Temp directory should be created")

	// Verify registry is running
	registryURL := fixtures.Registry.GetURL()
	assert.NotEmpty(t, registryURL, "Registry URL should not be empty")

	resp, err := http.Get(registryURL + "/v2/")
	require.NoError(t, err, "Registry should be accessible before cleanup")
	resp.Body.Close()

	// Verify temp directory exists
	assert.DirExists(t, fixtures.TempDir, "Temp directory should exist before cleanup")

	// Test manual cleanup (automatic cleanup happens via t.Cleanup)
	oldTempDir := fixtures.TempDir
	CleanupTestFixtures(fixtures)

	// Verify temp directory is removed
	assert.NoDirExists(t, oldTempDir, "Temp directory should be removed after cleanup")

	// Test that cleanup handles nil fixtures gracefully
	CleanupTestFixtures(nil)
	// Should not panic or error

	// Test cleanup with nil registry
	emptyFixtures := &TestFixtures{
		Registry: nil,
		Client:   &http.Client{},
		TempDir:  "",
	}
	CleanupTestFixtures(emptyFixtures)
	// Should not panic or error
}

// TestAuthenticatedFixtures validates authenticated test setup
func TestAuthenticatedFixtures(t *testing.T) {
	// Create authenticated fixtures
	fixtures := CreateAuthenticatedFixtures(t, "authuser", "authpass")
	require.NotNil(t, fixtures, "Authenticated fixtures should be created")

	// Verify auth is configured on registry
	require.NotNil(t, fixtures.Registry.AuthConfig, "Auth config should be set")
	assert.True(t, fixtures.Registry.AuthConfig.Required, "Auth should be required")
	assert.Equal(t, "authuser", fixtures.Registry.AuthConfig.Username, "Username should match")
	assert.Equal(t, "authpass", fixtures.Registry.AuthConfig.Password, "Password should match")

	// Verify client has auth transport
	transport, ok := fixtures.Client.Transport.(*MockAuthTransport)
	require.True(t, ok, "Client should use MockAuthTransport")
	assert.Equal(t, "authuser", transport.Username, "Transport username should match")
	assert.Equal(t, "authpass", transport.Password, "Transport password should match")

	// Test that authentication works
	resp, err := fixtures.Client.Get(fixtures.Registry.GetURL() + "/v2/")
	require.NoError(t, err, "Authenticated request should succeed")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode,
		"Authenticated request should return 200 OK")
}

// TestInsecureTransport validates the insecure transport mock
func TestInsecureTransport(t *testing.T) {
	// Create insecure transport
	transport := MockInsecureTransport()
	require.NotNil(t, transport, "Insecure transport should be created")

	// Create client with insecure transport
	client := &http.Client{
		Transport: transport,
	}

	// Create registry for testing
	registry := NewMockRegistry()
	defer registry.Cleanup()

	// Test that insecure transport can make requests
	resp, err := client.Get(registry.GetURL() + "/v2/")
	require.NoError(t, err, "Insecure transport should handle requests")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode,
		"Request via insecure transport should succeed")
}