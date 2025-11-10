//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/cnoe-io/idpbuilder/pkg/cmd/push"
	"github.com/cnoe-io/idpbuilder/test/harness"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestProgressUpdatesReceived verifies basic push works (progress callbacks removed from API)
func TestProgressUpdatesReceived(t *testing.T) {
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
		Name:   "progresstest",
		Tag:    "v1.0",
		Layers: 3,
		SizeMB: 10,
		Arch:   "amd64",
	})
	require.NoError(t, err)
	imageName := buildResult.ImageRef

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

	// Verify image pushed successfully
	exists, err := harness.VerifyImageInRegistry(ctx, env.RegistryURL, imageName, env.AdminUsername, env.AdminPassword)
	require.NoError(t, err)
	assert.True(t, exists, "Image not found in registry after push")

	// Verify image has expected layer count
	assert.Equal(t, 3, buildResult.LayerCount, "Expected 3 layers in built image")
}

// TestProgressForAllLayers verifies multi-layer image push works
func TestProgressForAllLayers(t *testing.T) {
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

	// Build 5-layer image
	builderEnv := &harness.BuilderTestEnvironment{
		DockerClient: env.DockerClient,
	}
	buildResult, err := builderEnv.BuildTestImage(ctx, harness.ImageConfig{
		Name:   "layertest",
		Tag:    "v1.0",
		Layers: 5,
		SizeMB: 25,
		Arch:   "amd64",
	})
	require.NoError(t, err)
	imageName := buildResult.ImageRef

	pushOpts := &push.PushOptions{
		ImageName: imageName,
		Registry:  env.RegistryURL,
		Username:  env.AdminUsername,
		Password:  env.AdminPassword,
		Insecure:  true,
	}

	err = push.RunPushForTesting(ctx, pushOpts)
	require.NoError(t, err)

	// Verify image pushed successfully
	exists, err := harness.VerifyImageInRegistry(ctx, env.RegistryURL, imageName, env.AdminUsername, env.AdminPassword)
	require.NoError(t, err)
	assert.True(t, exists, "Multi-layer image not found in registry")

	// Verify all 5 layers built correctly
	assert.GreaterOrEqual(t, buildResult.LayerCount, 5,
		"Expected at least 5 layers, got %d", buildResult.LayerCount)
}

// TestProgressMemoryEfficiency verifies large image push completes successfully
func TestProgressMemoryEfficiency(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	env, err := harness.SetupGiteaRegistry(ctx)
	require.NoError(t, err)
	defer func() {
		if cleanupErr := env.Cleanup(); cleanupErr != nil {
			t.Logf("Warning: Cleanup failed: %v", cleanupErr)
		}
	}()

	// Build large 500MB image
	builderEnv := &harness.BuilderTestEnvironment{
		DockerClient: env.DockerClient,
	}
	buildResult, err := builderEnv.BuildTestImage(ctx, harness.ImageConfig{
		Name:   "memorytest",
		Tag:    "v1.0",
		Layers: 10,
		SizeMB: 500,
		Arch:   "amd64",
	})
	require.NoError(t, err)
	imageName := buildResult.ImageRef

	pushOpts := &push.PushOptions{
		ImageName: imageName,
		Registry:  env.RegistryURL,
		Username:  env.AdminUsername,
		Password:  env.AdminPassword,
		Insecure:  true,
	}

	err = push.RunPushForTesting(ctx, pushOpts)
	require.NoError(t, err)

	// Verify large image pushed successfully
	exists, err := harness.VerifyImageInRegistry(ctx, env.RegistryURL, imageName, env.AdminUsername, env.AdminPassword)
	require.NoError(t, err)
	assert.True(t, exists, "Large image not found in registry")

	// Verify size is reasonable
	assert.Greater(t, buildResult.SizeBytes, int64(500*1024*1024),
		"Image should be >500MB, got %d bytes", buildResult.SizeBytes)

	// Note: This test primarily verifies the push completes without OOM
}
