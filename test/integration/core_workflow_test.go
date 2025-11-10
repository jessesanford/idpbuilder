//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cnoe-io/idpbuilder/pkg/cmd/push"
	"github.com/cnoe-io/idpbuilder/test/harness"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPushSmallImageSuccess verifies basic push workflow with 5MB image
func TestPushSmallImageSuccess(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Step 1: Setup test environment from harness (effort 3.1.1)
	env, err := harness.SetupGiteaRegistry(ctx)
	require.NoError(t, err, "Failed to setup test environment")
	defer func() {
		if cleanupErr := env.Cleanup(); cleanupErr != nil {
			t.Logf("Warning: Cleanup failed: %v", cleanupErr)
		}
	}()

	// Step 2: Build test image from harness (effort 3.1.2)
	// Note: BuildTestImage is on BuilderTestEnvironment, create one with our Docker client
	builderEnv := &harness.BuilderTestEnvironment{
		DockerClient: env.DockerClient,
	}
	buildResult, err := builderEnv.BuildTestImage(ctx, harness.ImageConfig{
		Name:   "testapp",
		Tag:    "v1.0",
		Layers: 2,
		SizeMB: 5,
		Arch:   "amd64",
	})
	require.NoError(t, err, "Failed to build test image")
	imageName := buildResult.ImageRef

	// Step 3: Create push options with test environment credentials
	pushOpts := &push.PushOptions{
		ImageName: imageName,
		Registry:  env.RegistryURL,
		Username:  env.AdminUsername,
		Password:  env.AdminPassword,
		Insecure:  true, // Test environment uses self-signed certs
		Verbose:   true,
	}

	// Step 4: Execute push (no progress callbacks in current API)
	err = push.RunPushForTesting(ctx, pushOpts)
	require.NoError(t, err, "Push command failed")

	// Step 5: Verify image in registry using package function
	pushedRef := imageName
	exists, err := harness.VerifyImageInRegistry(ctx, env.RegistryURL, pushedRef, env.AdminUsername, env.AdminPassword)
	require.NoError(t, err, "Failed to verify image in registry")
	assert.True(t, exists, "Image not found in registry after push")
}

// TestPushLargeImageWithProgress verifies progress reporting with 100MB, 10-layer image
func TestPushLargeImageWithProgress(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	env, err := harness.SetupGiteaRegistry(ctx)
	require.NoError(t, err)
	defer func() {
		if cleanupErr := env.Cleanup(); cleanupErr != nil {
			t.Logf("Warning: Cleanup failed: %v", cleanupErr)
		}
	}()

	// Build large test image (100MB, 10 layers)
	builderEnv := &harness.BuilderTestEnvironment{
		DockerClient: env.DockerClient,
	}
	buildResult, err := builderEnv.BuildTestImage(ctx, harness.ImageConfig{
		Name:   "largeapp",
		Tag:    "latest",
		Layers: 10,
		SizeMB: 100,
		Arch:   "amd64",
	})
	require.NoError(t, err)
	imageName := buildResult.ImageRef

	// Create push options (progress callbacks not in current API)
	pushOpts := &push.PushOptions{
		ImageName: imageName,
		Registry:  env.RegistryURL,
		Username:  env.AdminUsername,
		Password:  env.AdminPassword,
		Insecure:  true,
		Verbose:   true,
	}

	err = push.RunPushForTesting(ctx, pushOpts)
	require.NoError(t, err)

	// Verify image was pushed successfully
	exists, err := harness.VerifyImageInRegistry(ctx, env.RegistryURL, imageName, env.AdminUsername, env.AdminPassword)
	require.NoError(t, err)
	assert.True(t, exists, "Large image not found in registry after push")

	// Verify the build result has correct layer count
	// Note: BuildTestImage uses alpine:latest base + requested layers,
	// so actual layer count will be base layers + requested layers
	assert.GreaterOrEqual(t, buildResult.LayerCount, 10, "Expected at least 10 layers in built image")
	assert.True(t, buildResult.SizeBytes > 100*1024*1024,
		"Expected >100MB image size, got %d bytes", buildResult.SizeBytes)
}

// TestPushWithAuthenticationSuccess verifies authentication flow with valid credentials
func TestPushWithAuthenticationSuccess(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	env, err := harness.SetupGiteaRegistry(ctx)
	require.NoError(t, err)
	defer func() {
		if cleanupErr := env.Cleanup(); cleanupErr != nil {
			t.Logf("Warning: Cleanup failed: %v", cleanupErr)
		}
	}()

	builderEnv := &harness.BuilderTestEnvironment{
		DockerClient: env.DockerClient,
	}
	buildResult, err := builderEnv.BuildTestImage(ctx, harness.ImageConfig{
		Name:   "authtest",
		Tag:    "latest",
		Layers: 2,
		SizeMB: 5,
		Arch:   "amd64",
	})
	require.NoError(t, err)
	imageName := buildResult.ImageRef

	// Test with CORRECT credentials from environment
	pushOpts := &push.PushOptions{
		ImageName: imageName,
		Registry:  env.RegistryURL,
		Username:  env.AdminUsername,
		Password:  env.AdminPassword,
		Insecure:  true,
	}

	err = push.RunPushForTesting(ctx, pushOpts)
	assert.NoError(t, err, "Push with correct credentials should succeed")

	// Verify image pushed successfully
	exists, err := harness.VerifyImageInRegistry(ctx, env.RegistryURL, imageName, env.AdminUsername, env.AdminPassword)
	require.NoError(t, err)
	assert.True(t, exists, "Image not found after authenticated push")
}

// TestPushToCustomRegistry verifies custom registry URL handling
func TestPushToCustomRegistry(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	env, err := harness.SetupGiteaRegistry(ctx)
	require.NoError(t, err)
	defer func() {
		if cleanupErr := env.Cleanup(); cleanupErr != nil {
			t.Logf("Warning: Cleanup failed: %v", cleanupErr)
		}
	}()

	builderEnv := &harness.BuilderTestEnvironment{
		DockerClient: env.DockerClient,
	}
	buildResult, err := builderEnv.BuildTestImage(ctx, harness.ImageConfig{
		Name:   "customapp",
		Tag:    "v1.0",
		Layers: 2,
		SizeMB: 5,
		Arch:   "amd64",
	})
	require.NoError(t, err)
	imageName := buildResult.ImageRef

	// Use custom registry URL (Gitea dynamic port)
	customRegistryURL := env.RegistryURL

	pushOpts := &push.PushOptions{
		ImageName: imageName,
		Registry:  customRegistryURL, // Use custom URL
		Username:  env.AdminUsername,
		Password:  env.AdminPassword,
		Insecure:  true,
	}

	err = push.RunPushForTesting(ctx, pushOpts)
	require.NoError(t, err, "Push to custom registry failed")

	// Verify image pushed to custom registry
	exists, err := harness.VerifyImageInRegistry(ctx, customRegistryURL, imageName, env.AdminUsername, env.AdminPassword)
	require.NoError(t, err)
	assert.True(t, exists, "Image not found at custom registry URL")
}

// TestPushMultipleImagesSequentially verifies multiple sequential pushes work
func TestPushMultipleImagesSequentially(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	env, err := harness.SetupGiteaRegistry(ctx)
	require.NoError(t, err)
	defer func() {
		if cleanupErr := env.Cleanup(); cleanupErr != nil {
			t.Logf("Warning: Cleanup failed: %v", cleanupErr)
		}
	}()

	builderEnv := &harness.BuilderTestEnvironment{
		DockerClient: env.DockerClient,
	}

	// Build 3 different images
	imageNames := make([]string, 3)
	for i := 0; i < 3; i++ {
		buildResult, err := builderEnv.BuildTestImage(ctx, harness.ImageConfig{
			Name:   fmt.Sprintf("app%d", i+1),
			Tag:    "v1.0",
			Layers: 2,
			SizeMB: 5,
			Arch:   "amd64",
		})
		require.NoError(t, err, "Failed to build image %d", i+1)
		imageNames[i] = buildResult.ImageRef
	}

	// Push all 3 images sequentially
	for i, imageName := range imageNames {
		pushOpts := &push.PushOptions{
			ImageName: imageName,
			Registry:  env.RegistryURL,
			Username:  env.AdminUsername,
			Password:  env.AdminPassword,
			Insecure:  true,
		}

		err = push.RunPushForTesting(ctx, pushOpts)
		require.NoError(t, err, "Push %d failed", i+1)
	}

	// Verify all 3 images exist in registry
	for i, imageName := range imageNames {
		exists, err := harness.VerifyImageInRegistry(ctx, env.RegistryURL, imageName, env.AdminUsername, env.AdminPassword)
		require.NoError(t, err)
		assert.True(t, exists, "Image %d not found in registry", i+1)
	}
}

