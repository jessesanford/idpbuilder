//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/cnoe-io/idpbuilder/cmd/push"
	"github.com/cnoe-io/idpbuilder/test/harness"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestProgressUpdatesReceived verifies progress callback invoked correctly
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

	imageName, err := env.BuildTestImage(ctx, harness.ImageConfig{
		Name:   "progresstest",
		Tag:    "v1.0",
		Layers: 3,
		SizeMB: 10,
		Arch:   "amd64",
	})
	require.NoError(t, err)

	// Track progress updates
	var progressUpdates []push.ProgressUpdate
	var callbackInvoked bool

	pushOpts := &push.PushOptions{
		ImageRef:     imageName,
		Registry:     env.RegistryURL,
		Username:     env.AdminUsername,
		Password:     env.AdminPassword,
		Insecure:     true,
		DockerClient: env.DockerClient,
		ProgressCallback: func(update push.ProgressUpdate) {
			callbackInvoked = true
			progressUpdates = append(progressUpdates, update)
		},
	}

	err = push.Execute(ctx, pushOpts)
	require.NoError(t, err)

	// Verify callback invoked
	assert.True(t, callbackInvoked, "Progress callback never invoked")
	assert.NotEmpty(t, progressUpdates, "No progress updates received")

	// Verify updates contain layer digests
	for _, update := range progressUpdates {
		assert.NotEmpty(t, update.LayerDigest, "Layer digest missing in update")
	}

	// Verify bytes pushed increments
	var lastBytes int64
	for _, update := range progressUpdates {
		if update.BytesPushed > 0 {
			assert.GreaterOrEqual(t, update.BytesPushed, lastBytes,
				"Bytes pushed should not decrease")
			lastBytes = update.BytesPushed
		}
	}
}

// TestProgressForAllLayers verifies progress for each layer individually
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
	imageName, err := env.BuildTestImage(ctx, harness.ImageConfig{
		Name:   "layertest",
		Tag:    "v1.0",
		Layers: 5,
		SizeMB: 25,
		Arch:   "amd64",
	})
	require.NoError(t, err)

	// Track unique layers
	layersSeen := make(map[string]bool)

	pushOpts := &push.PushOptions{
		ImageRef:     imageName,
		Registry:     env.RegistryURL,
		Username:     env.AdminUsername,
		Password:     env.AdminPassword,
		Insecure:     true,
		DockerClient: env.DockerClient,
		ProgressCallback: func(update push.ProgressUpdate) {
			if update.LayerDigest != "" {
				layersSeen[update.LayerDigest] = true
			}
		},
	}

	err = push.Execute(ctx, pushOpts)
	require.NoError(t, err)

	// Verify all 5 layers reported (may have base layer too)
	assert.GreaterOrEqual(t, len(layersSeen), 5,
		"Expected at least 5 layers, got %d", len(layersSeen))
}

// TestProgressMemoryEfficiency verifies memory stays reasonable during large push
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
	imageName, err := env.BuildTestImage(ctx, harness.ImageConfig{
		Name:   "memorytest",
		Tag:    "v1.0",
		Layers: 10,
		SizeMB: 500,
		Arch:   "amd64",
	})
	require.NoError(t, err)

	// Track memory usage (rough check)
	var updateCount int

	pushOpts := &push.PushOptions{
		ImageRef:     imageName,
		Registry:     env.RegistryURL,
		Username:     env.AdminUsername,
		Password:     env.AdminPassword,
		Insecure:     true,
		DockerClient: env.DockerClient,
		ProgressCallback: func(update push.ProgressUpdate) {
			updateCount++
		},
	}

	err = push.Execute(ctx, pushOpts)
	require.NoError(t, err)

	// Verify streaming worked (got progress updates)
	assert.Greater(t, updateCount, 0, "No progress updates received")

	// Note: Actual memory profiling would require runtime/pprof
	// This test primarily verifies the push completes without OOM
}
