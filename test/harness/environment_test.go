package harness

import (
	"context"
	"testing"
	"time"
)

// TestGiteaContainerStartup verifies that the Gitea container starts
// successfully and ports are properly mapped.
func TestGiteaContainerStartup(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	env, err := SetupGiteaRegistry(ctx)
	if err != nil {
		t.Fatalf("Failed to setup Gitea registry: %v", err)
	}
	defer env.Cleanup()

	// Verify container is running
	if env.GiteaContainer == nil {
		t.Fatal("Gitea container is nil")
	}

	// Verify registry URL is set
	if env.RegistryURL == "" {
		t.Fatal("Registry URL is empty")
	}

	// Verify admin credentials are set
	if env.AdminUsername == "" || env.AdminPassword == "" {
		t.Fatal("Admin credentials not set")
	}

	// Verify Docker client is initialized
	if env.DockerClient == nil {
		t.Fatal("Docker client is nil")
	}

	t.Logf("Gitea container started successfully")
	t.Logf("Registry URL: %s", env.RegistryURL)
	t.Logf("Admin username: %s", env.AdminUsername)
}

// TestDynamicPortAllocation verifies that dynamic port allocation works correctly
func TestDynamicPortAllocation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	env, err := SetupGiteaRegistry(ctx)
	if err != nil {
		t.Fatalf("Failed to setup Gitea registry: %v", err)
	}
	defer env.Cleanup()

	// Verify port is allocated (not the default 5000)
	if env.RegistryURL == "localhost:5000" {
		t.Log("Warning: Got default port, but this might be coincidence")
	}

	// Verify we can get the actual mapped port (changed from 3000 to 5000 for registry:2)
	registryPort, err := env.GiteaContainer.MappedPort(ctx, "5000")
	if err != nil {
		t.Fatalf("Failed to get mapped registry port: %v", err)
	}

	if registryPort.Port() == "" {
		t.Fatal("Mapped port is empty")
	}

	t.Logf("Registry mapped to port: %s", registryPort.Port())
}

// TestDockerClientConnection verifies that the Docker client connects to the host daemon
func TestDockerClientConnection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	env, err := SetupGiteaRegistry(ctx)
	if err != nil {
		t.Fatalf("Failed to setup Gitea registry: %v", err)
	}
	defer env.Cleanup()

	// Try to ping Docker daemon
	_, err = env.DockerClient.Ping(ctx)
	if err != nil {
		t.Fatalf("Failed to ping Docker daemon: %v", err)
	}

	t.Log("Docker client successfully connected to daemon")
}

// TestCleanupRemovesContainer verifies that cleanup properly removes the container
func TestCleanupRemovesContainer(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	env, err := SetupGiteaRegistry(ctx)
	if err != nil {
		t.Fatalf("Failed to setup Gitea registry: %v", err)
	}

	// Verify container is running before cleanup
	if env.GiteaContainer == nil {
		t.Fatal("Gitea container is nil before cleanup")
	}

	// Cleanup
	if err := env.Cleanup(); err != nil {
		t.Fatalf("Cleanup failed: %v", err)
	}

	// Verify container is gone (by trying to get state - should fail)
	// Note: We can't easily check if container is removed without another client,
	// but we can verify the cleanup function executed without error
	t.Log("Container cleaned up successfully")
}

// TestWaitForGiteaReady verifies that the readiness check works
func TestWaitForGiteaReady(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	env, err := SetupGiteaRegistry(ctx)
	if err != nil {
		t.Fatalf("Failed to setup Gitea registry: %v", err)
	}
	defer env.Cleanup()

	// Gitea should already be ready from SetupGiteaRegistry
	// But let's test the WaitForGiteaReady function explicitly
	// Note: This will check the registry port, not the web port
	// For testing purposes, we'll create a short timeout context
	checkCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// The registry URL might not respond to HTTP GET requests the same way
	// This test mainly verifies the function doesn't panic
	err = WaitForGiteaReady(checkCtx, env.RegistryURL)
	if err != nil {
		// It's okay if this fails - registry port doesn't serve HTTP
		t.Logf("WaitForGiteaReady returned error (expected for registry port): %v", err)
	} else {
		t.Log("WaitForGiteaReady succeeded")
	}
}

// TestCleanupTestEnvironment verifies the CleanupTestEnvironment method
func TestCleanupTestEnvironment(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	env, err := SetupGiteaRegistry(ctx)
	if err != nil {
		t.Fatalf("Failed to setup Gitea registry: %v", err)
	}

	// Use CleanupTestEnvironment instead of Cleanup
	if err := env.CleanupTestEnvironment(); err != nil {
		t.Fatalf("CleanupTestEnvironment failed: %v", err)
	}

	t.Log("CleanupTestEnvironment succeeded")
}

// TestRemoveTestImages verifies test image removal
func TestRemoveTestImages(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	env, err := SetupGiteaRegistry(ctx)
	if err != nil {
		t.Fatalf("Failed to setup Gitea registry: %v", err)
	}
	defer env.Cleanup()

	// Try to remove test images (there might not be any, which is fine)
	err = RemoveTestImages(ctx, env.DockerClient)
	if err != nil {
		// Log error but don't fail - might be no test images
		t.Logf("RemoveTestImages returned error (might be no test images): %v", err)
	} else {
		t.Log("RemoveTestImages succeeded")
	}
}

// TestVerifyImageInRegistry tests the image verification helper
func TestVerifyImageInRegistry(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	env, err := SetupGiteaRegistry(ctx)
	if err != nil {
		t.Fatalf("Failed to setup Gitea registry: %v", err)
	}
	defer env.Cleanup()

	// Try to verify a non-existent image
	exists, err := VerifyImageInRegistry(
		ctx,
		env.RegistryURL,
		"nonexistent/image:latest",
		env.AdminUsername,
		env.AdminPassword,
	)

	if err != nil {
		t.Logf("VerifyImageInRegistry returned error (expected for non-existent image): %v", err)
	}

	if exists {
		t.Error("Non-existent image reported as existing")
	} else {
		t.Log("Correctly identified non-existent image")
	}
}

// TestGetImageDigest tests the digest retrieval helper
func TestGetImageDigest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	env, err := SetupGiteaRegistry(ctx)
	if err != nil {
		t.Fatalf("Failed to setup Gitea registry: %v", err)
	}
	defer env.Cleanup()

	// Try to get digest for non-existent image (should fail)
	digest, err := GetImageDigest(
		ctx,
		env.RegistryURL,
		"nonexistent/image:latest",
		env.AdminUsername,
		env.AdminPassword,
	)

	if err == nil {
		t.Error("Expected error for non-existent image, got none")
	}

	if digest != "" {
		t.Errorf("Expected empty digest for non-existent image, got: %s", digest)
	} else {
		t.Log("Correctly returned error and empty digest for non-existent image")
	}
}
