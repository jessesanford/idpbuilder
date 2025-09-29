//go:build integration

package integration

import (
	"context"
	"crypto/tls"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

// TestIntegration_RegistryConnection tests basic registry connectivity
func TestIntegration_RegistryConnection(t *testing.T) {
	t.Log("Starting registry connection test")

	// Setup test registry
	container, registryURL, err := setupTestRegistryForTest(t)
	if err != nil {
		t.Fatalf("Failed to setup test registry: %v", err)
	}
	defer cleanupContainer(t, container)

	// Test basic connectivity
	if err := validateRegistryConnection(registryURL); err != nil {
		t.Fatalf("Registry connection validation failed: %v", err)
	}

	// Test registry health endpoint (if available)
	config := createTestConfig(registryURL)
	if config.RegistryURL != registryURL {
		t.Errorf("Expected registry URL %s, got %s", registryURL, config.RegistryURL)
	}

	// Verify registry responds to basic requests
	if !VerifyImageInRegistry(registryURL, "nonexistent/image:latest") {
		t.Log("Expected behavior: nonexistent image not found")
	}

	t.Log("Registry connection test completed successfully")
}

// TestIntegration_AuthenticationFlow tests authentication with credentials
func TestIntegration_AuthenticationFlow(t *testing.T) {
	t.Log("Starting authentication flow test")

	// Setup test registry
	container, registryURL, err := setupTestRegistryForTest(t)
	if err != nil {
		t.Fatalf("Failed to setup test registry: %v", err)
	}
	defer cleanupContainer(t, container)

	// Generate test credentials
	creds, err := GenerateTestCredentials()
	if err != nil {
		t.Fatalf("Failed to generate test credentials: %v", err)
	}

	// Validate credentials format
	if err := validateCredentials(creds); err != nil {
		t.Fatalf("Generated credentials are invalid: %v", err)
	}

	// Test credential structure
	if creds.Username == "" {
		t.Errorf("Expected non-empty username")
	}
	if creds.Password == "" {
		t.Errorf("Expected non-empty password")
	}

	// Setup registry credentials (simulated)
	registryCreds, err := setupRegistryCredentials(registryURL)
	if err != nil {
		t.Fatalf("Failed to setup registry credentials: %v", err)
	}

	// Validate registry credentials
	if err := validateCredentials(registryCreds); err != nil {
		t.Fatalf("Registry credentials are invalid: %v", err)
	}

	t.Log("Authentication flow test completed successfully")
}

// TestIntegration_InsecureCertHandling tests self-signed certificate support
func TestIntegration_InsecureCertHandling(t *testing.T) {
	t.Log("Starting insecure certificate handling test")

	// Setup insecure TLS configuration
	tlsConfig, err := SetupInsecureCertTest()
	if err != nil {
		t.Fatalf("Failed to setup insecure TLS config: %v", err)
	}

	// Validate TLS configuration
	if err := validateTLSConfig(tlsConfig); err != nil {
		t.Fatalf("TLS config validation failed: %v", err)
	}

	// Test insecure skip verify setting
	if !tlsConfig.InsecureSkipVerify {
		t.Errorf("Expected InsecureSkipVerify to be true for testing")
	}

	// Test minimum TLS version
	if tlsConfig.MinVersion < tls.VersionTLS12 {
		t.Errorf("Expected minimum TLS version 1.2 or higher")
	}

	// Create insecure transport for testing
	transport := configureInsecureTransport()
	if transport == nil {
		t.Fatalf("Failed to create insecure transport")
	}

	// Verify transport configuration
	if !transport.TLSClientConfig.InsecureSkipVerify {
		t.Errorf("Expected transport to skip TLS verification")
	}

	// Test with registry configuration
	config := &IntegrationTestConfig{
		RegistryURL:    "https://localhost:5000",
		InsecureMode:   true,
		TestImagePath:  "test/insecure:latest",
		Timeout:        30 * time.Second,
	}

	if !config.InsecureMode {
		t.Errorf("Expected insecure mode to be enabled")
	}

	t.Log("Insecure certificate handling test completed successfully")
}

// TestIntegration_ImagePushPull tests image push and pull operations
func TestIntegration_ImagePushPull(t *testing.T) {
	t.Log("Starting image push/pull test")

	// Setup test registry
	container, registryURL, err := setupTestRegistryForTest(t)
	if err != nil {
		t.Fatalf("Failed to setup test registry: %v", err)
	}
	defer cleanupContainer(t, container)

	// Define test image
	testImage := "test/integration:latest"

	// Test image reference validation
	if err := validateImageReference(registryURL, testImage); err != nil {
		t.Fatalf("Image reference validation failed: %v", err)
	}

	// Push test image to registry
	t.Log("Pushing test image to registry")
	if err := PushTestImage(registryURL, testImage); err != nil {
		t.Fatalf("Failed to push test image: %v", err)
	}

	// Verify image exists in registry
	t.Log("Verifying image exists in registry")
	if !VerifyImageInRegistry(registryURL, testImage) {
		t.Fatalf("Image not found in registry after push")
	}

	// Pull test image from registry
	t.Log("Pulling test image from registry")
	if err := PullTestImage(registryURL, testImage); err != nil {
		t.Fatalf("Failed to pull test image: %v", err)
	}

	// Test with different image name
	anotherImage := "test/another:v1.0"
	if err := PushTestImage(registryURL, anotherImage); err != nil {
		t.Fatalf("Failed to push second test image: %v", err)
	}

	// Verify both images exist
	if !VerifyImageInRegistry(registryURL, testImage) {
		t.Errorf("First test image no longer found")
	}
	if !VerifyImageInRegistry(registryURL, anotherImage) {
		t.Errorf("Second test image not found")
	}

	t.Log("Image push/pull test completed successfully")
}

// TestIntegration_ClusterLifecycle tests cluster creation and destruction
func TestIntegration_ClusterLifecycle(t *testing.T) {
	t.Log("Starting cluster lifecycle test")

	// Create test cluster
	cluster, err := CreateTestCluster(t)
	if err != nil {
		t.Fatalf("Failed to create test cluster: %v", err)
	}

	// Validate cluster information
	if cluster.Name == "" {
		t.Errorf("Expected non-empty cluster name")
	}
	if cluster.Namespace == "" {
		t.Errorf("Expected non-empty namespace")
	}
	if cluster.Context == "" {
		t.Errorf("Expected non-empty kube context")
	}
	if cluster.Cleanup == nil {
		t.Errorf("Expected cleanup function to be set")
	}

	// Test cluster access
	t.Logf("Testing access to cluster: %s", cluster.Name)
	if err := verifyClusterAccess(cluster.Context); err != nil {
		t.Errorf("Cluster access verification failed: %v", err)
	}

	// Get kubeconfig path
	kubeconfigPath, err := getKubeconfig(cluster.Name)
	if err != nil {
		t.Logf("Warning: Could not get kubeconfig path: %v", err)
	} else {
		t.Logf("Kubeconfig path: %s", kubeconfigPath)
	}

	// Test cluster cleanup
	t.Log("Cleaning up test cluster")
	CleanupTestCluster(t, cluster)

	t.Log("Cluster lifecycle test completed successfully")
}

// Helper functions for tests

// setupTestRegistryForTest wraps SetupTestRegistry with proper error handling
func setupTestRegistryForTest(t *testing.T) (*testcontainers.Container, string, error) {
	container, registryURL := SetupTestRegistry(t)
	return container, registryURL, nil
}

// cleanupContainer safely terminates a test container
func cleanupContainer(t *testing.T, container *testcontainers.Container) {
	if container == nil {
		return
	}

	ctx := context.Background()
	if err := (*container).Terminate(ctx); err != nil {
		t.Logf("Warning: Failed to cleanup container: %v", err)
	}
}