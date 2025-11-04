//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cnoe-io/idpbuilder/cmd/push"
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
	imageName, err := env.BuildTestImage(ctx, harness.ImageConfig{
		Name:   "testapp",
		Tag:    "v1.0",
		Layers: 2,
		SizeMB: 5,
		Arch:   "amd64",
	})
	require.NoError(t, err, "Failed to build test image")

	// Step 3: Create push options with test environment credentials
	pushOpts := &push.PushOptions{
		ImageRef:     imageName,
		Registry:     env.RegistryURL,
		Username:     env.AdminUsername,
		Password:     env.AdminPassword,
		Insecure:     true, // Test environment uses self-signed certs
		Verbose:      true,
		DockerClient: env.DockerClient,
	}

	// Step 4: Capture progress updates
	var progressUpdates []push.ProgressUpdate
	pushOpts.ProgressCallback = func(update push.ProgressUpdate) {
		progressUpdates = append(progressUpdates, update)
	}

	// Step 5: Execute push
	err = push.Execute(ctx, pushOpts)
	require.NoError(t, err, "Push command failed")

	// Step 6: Verify image in registry
	pushedRef := fmt.Sprintf("%s/%s", env.RegistryURL, imageName)
	exists, err := env.VerifyImageInRegistry(ctx, pushedRef)
	require.NoError(t, err, "Failed to verify image in registry")
	assert.True(t, exists, "Image not found in registry after push")

	// Step 7: Verify progress updates received
	assert.NotEmpty(t, progressUpdates, "No progress updates received")
	assert.True(t, hasCompleteStatus(progressUpdates), "Missing complete status")
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
	imageName, err := env.BuildTestImage(ctx, harness.ImageConfig{
		Name:   "largeapp",
		Tag:    "latest",
		Layers: 10,
		SizeMB: 100,
		Arch:   "amd64",
	})
	require.NoError(t, err)

	// Track progress in detail
	var progressUpdates []push.ProgressUpdate
	var layersCompleted int

	pushOpts := &push.PushOptions{
		ImageRef:     imageName,
		Registry:     env.RegistryURL,
		Username:     env.AdminUsername,
		Password:     env.AdminPassword,
		Insecure:     true,
		Verbose:      true,
		DockerClient: env.DockerClient,
		ProgressCallback: func(update push.ProgressUpdate) {
			progressUpdates = append(progressUpdates, update)
			if update.Status == "Complete" {
				layersCompleted++
			}
		},
	}

	err = push.Execute(ctx, pushOpts)
	require.NoError(t, err)

	// Verify all 10 layers reported progress
	assert.Equal(t, 10, layersCompleted, "Not all layers completed")
	assert.True(t, len(progressUpdates) >= 10, "Expected at least 10 progress updates")

	// Verify total bytes processed
	var totalBytesProcessed int64
	for _, update := range progressUpdates {
		totalBytesProcessed += update.BytesPushed
	}
	assert.True(t, totalBytesProcessed > 100*1024*1024,
		"Expected >100MB processed, got %d bytes", totalBytesProcessed)
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

	imageName, err := env.BuildTestImage(ctx, harness.ImageConfig{
		Name:   "authtest",
		Tag:    "latest",
		Layers: 2,
		SizeMB: 5,
		Arch:   "amd64",
	})
	require.NoError(t, err)

	// Test with CORRECT credentials from environment
	pushOpts := &push.PushOptions{
		ImageRef:     imageName,
		Registry:     env.RegistryURL,
		Username:     env.AdminUsername,
		Password:     env.AdminPassword,
		Insecure:     true,
		DockerClient: env.DockerClient,
	}

	err = push.Execute(ctx, pushOpts)
	assert.NoError(t, err, "Push with correct credentials should succeed")

	// Verify image pushed successfully
	pushedRef := fmt.Sprintf("%s/%s", env.RegistryURL, imageName)
	exists, err := env.VerifyImageInRegistry(ctx, pushedRef)
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

	imageName, err := env.BuildTestImage(ctx, harness.ImageConfig{
		Name:   "customapp",
		Tag:    "v1.0",
		Layers: 2,
		SizeMB: 5,
		Arch:   "amd64",
	})
	require.NoError(t, err)

	// Use custom registry URL (Gitea dynamic port)
	customRegistryURL := env.RegistryURL

	pushOpts := &push.PushOptions{
		ImageRef:     imageName,
		Registry:     customRegistryURL, // Use custom URL
		Username:     env.AdminUsername,
		Password:     env.AdminPassword,
		Insecure:     true,
		DockerClient: env.DockerClient,
	}

	err = push.Execute(ctx, pushOpts)
	require.NoError(t, err, "Push to custom registry failed")

	// Verify image pushed to custom registry
	pushedRef := fmt.Sprintf("%s/%s", customRegistryURL, imageName)
	exists, err := env.VerifyImageInRegistry(ctx, pushedRef)
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

	// Build 3 different images
	imageNames := make([]string, 3)
	for i := 0; i < 3; i++ {
		imageName, err := env.BuildTestImage(ctx, harness.ImageConfig{
			Name:   fmt.Sprintf("app%d", i+1),
			Tag:    "v1.0",
			Layers: 2,
			SizeMB: 5,
			Arch:   "amd64",
		})
		require.NoError(t, err, "Failed to build image %d", i+1)
		imageNames[i] = imageName
	}

	// Push all 3 images sequentially
	for i, imageName := range imageNames {
		pushOpts := &push.PushOptions{
			ImageRef:     imageName,
			Registry:     env.RegistryURL,
			Username:     env.AdminUsername,
			Password:     env.AdminPassword,
			Insecure:     true,
			DockerClient: env.DockerClient,
		}

		err = push.Execute(ctx, pushOpts)
		require.NoError(t, err, "Push %d failed", i+1)
	}

	// Verify all 3 images exist in registry
	for i, imageName := range imageNames {
		pushedRef := fmt.Sprintf("%s/%s", env.RegistryURL, imageName)
		exists, err := env.VerifyImageInRegistry(ctx, pushedRef)
		require.NoError(t, err)
		assert.True(t, exists, "Image %d not found in registry", i+1)
	}
}

// hasCompleteStatus checks if any progress update has "Complete" status
func hasCompleteStatus(updates []push.ProgressUpdate) bool {
	for _, update := range updates {
		if update.Status == "Complete" {
			return true
		}
	}
	return false
}
