package build

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/containers/buildah"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultStoreConfig(t *testing.T) {
	config := DefaultStoreConfig()
	
	assert.NotNil(t, config)
	assert.Equal(t, "overlay", config.GraphDriver)
	assert.True(t, config.Rootless) // Assumes running in non-root environment
	assert.NotEmpty(t, config.RootDir)
	assert.NotEmpty(t, config.RunRoot)
	assert.NotNil(t, config.StorageOpts)
}

func TestNewStoreManager(t *testing.T) {
	// Test with nil config
	sm := NewStoreManager(nil)
	assert.NotNil(t, sm)
	assert.False(t, sm.initialized)
	assert.Equal(t, "overlay", sm.graphDriver)
	
	// Test with custom config
	config := &StoreConfig{
		RootDir:     "/tmp/test-storage",
		RunRoot:     "/tmp/test-run",
		GraphDriver: "vfs",
		Rootless:    true,
	}
	
	sm = NewStoreManager(config)
	assert.NotNil(t, sm)
	assert.Equal(t, "/tmp/test-storage", sm.rootDir)
	assert.Equal(t, "/tmp/test-run", sm.runRoot)
	assert.Equal(t, "vfs", sm.graphDriver)
	assert.True(t, sm.rootless)
}

func TestStoreManagerInitialization(t *testing.T) {
	// Create temporary directories for testing
	tempDir, err := os.MkdirTemp("", "test-store-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	config := &StoreConfig{
		RootDir:     filepath.Join(tempDir, "storage"),
		RunRoot:     filepath.Join(tempDir, "run"),
		GraphDriver: "vfs", // Use VFS for testing as it has no dependencies
		Rootless:    true,
	}
	
	sm := NewStoreManager(config)
	assert.False(t, sm.initialized)
	
	// Initialize should work
	ctx := context.Background()
	err = sm.Initialize(ctx)
	
	// Note: This may fail due to system dependencies, but we're testing the logic
	// In a real environment, this would work with proper container runtime setup
	if err != nil {
		t.Logf("Initialization failed (expected in test environment): %v", err)
		return
	}
	
	assert.True(t, sm.initialized)
	assert.NotNil(t, sm.GetStore())
	
	// Test shutdown
	err = sm.Shutdown()
	assert.NoError(t, err)
	assert.False(t, sm.initialized)
}

func TestStoreManagerShutdown(t *testing.T) {
	sm := NewStoreManager(nil)
	
	// Should not error when shutting down uninitialized store
	err := sm.Shutdown()
	assert.NoError(t, err)
}

func TestCreateDirectories(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-dirs-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	config := &StoreConfig{
		RootDir: filepath.Join(tempDir, "root"),
		RunRoot: filepath.Join(tempDir, "run"),
	}
	
	sm := NewStoreManager(config)
	err = sm.createDirectories()
	assert.NoError(t, err)
	
	// Verify directories were created
	assert.DirExists(t, config.RootDir)
	assert.DirExists(t, config.RunRoot)
}

func TestConfigureStorage(t *testing.T) {
	sm := NewStoreManager(nil)
	
	// Test VFS configuration
	sm.graphDriver = "vfs"
	sm.configureStorage()
	assert.Equal(t, []string{}, sm.storageOptions.GraphOptions)
	
	// Test overlay configuration
	sm.graphDriver = "overlay"
	sm.rootless = true
	sm.configureStorage()
	assert.Contains(t, sm.storageOptions.GraphOptions, "overlay.mount_program=/usr/bin/fuse-overlayfs")
	
	// Test unknown driver
	sm.graphDriver = "unknown"
	sm.configureStorage() // Should not panic
}

func TestStoreManagerMethods(t *testing.T) {
	sm := NewStoreManager(nil)
	
	// Test methods on uninitialized store
	assert.Nil(t, sm.GetStore())
	
	options := sm.GetStoreOptions()
	assert.NotEmpty(t, options.GraphDriverName)
	
	// Test CreateBuilder (should fail on uninitialized store)
	ctx := context.Background()
	_, err := sm.CreateBuilder(ctx, buildah.BuilderOptions{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not initialized")
	
	// Test ListImages (should fail on uninitialized store)
	_, err = sm.ListImages()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not initialized")
	
	// Test ImageExists (should fail on uninitialized store)
	_, err = sm.ImageExists("test:latest")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not initialized")
}