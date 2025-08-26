package build

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/containers/buildah/define"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultBuildConfig(t *testing.T) {
	config := DefaultBuildConfig()
	
	assert.NotNil(t, config)
	assert.Equal(t, "scratch", config.From)
	assert.Equal(t, ".", config.Context)
	assert.Equal(t, runtime.GOARCH, config.Platform)
	assert.Equal(t, "runc", config.Runtime)
	assert.Equal(t, define.IsolationDefault, config.Isolation)
	assert.Equal(t, "bridge", config.NetworkMode)
	assert.Equal(t, "oci", config.Format)
	assert.False(t, config.Push)
	assert.NotNil(t, config.Args)
	assert.NotNil(t, config.Labels)
	assert.NotNil(t, config.Storage)
	assert.NotNil(t, config.SystemContext)
}

func TestNewConfigManager(t *testing.T) {
	// Test with nil config
	cm := NewConfigManager(nil)
	assert.NotNil(t, cm)
	assert.NotNil(t, cm.config)
	assert.False(t, cm.validated)
	
	// Test with custom config
	config := &BuildConfig{
		From:    "alpine:latest",
		Context: "/tmp/build",
		Format:  "docker",
	}
	
	cm = NewConfigManager(config)
	assert.NotNil(t, cm)
	assert.Equal(t, "alpine:latest", cm.config.From)
	assert.Equal(t, "/tmp/build", cm.config.Context)
	assert.Equal(t, "docker", cm.config.Format)
}

func TestConfigManagerValidation(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test-config-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	config := &BuildConfig{
		From:    "scratch",
		Context: tempDir, // Use existing directory
		Storage: DefaultStoreConfig(),
	}
	
	cm := NewConfigManager(config)
	ctx := context.Background()
	
	// Validation should succeed with valid config
	err = cm.Validate(ctx)
	assert.NoError(t, err)
	assert.True(t, cm.validated)
	
	// Second validation should return immediately
	err = cm.Validate(ctx)
	assert.NoError(t, err)
}

func TestValidateBuildParams(t *testing.T) {
	cm := NewConfigManager(nil)
	
	// Test with non-existent context directory
	cm.config.Context = "/non/existent/path"
	err := cm.validateBuildParams()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
	
	// Test with valid relative path
	tempDir, err := os.MkdirTemp("", "test-context-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Change to temp directory to test relative paths
	oldWd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(oldWd)
	
	cm.config.Context = "."
	err = cm.validateBuildParams()
	assert.NoError(t, err)
	assert.True(t, filepath.IsAbs(cm.config.Context))
	
	// Test with valid platform
	cm.config.Platform = "amd64"
	err = cm.validateBuildParams()
	assert.NoError(t, err)
	
	// Test with invalid platform
	cm.config.Platform = "invalid-arch"
	err = cm.validateBuildParams()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported platform")
}

func TestValidateStorageConfig(t *testing.T) {
	cm := NewConfigManager(nil)
	
	// Test with nil storage config
	cm.config.Storage = nil
	err := cm.validateStorageConfig()
	assert.NoError(t, err)
	
	// Test with valid storage config
	tempDir, err := os.MkdirTemp("", "test-storage-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	cm.config.Storage = &StoreConfig{
		RootDir:     filepath.Join(tempDir, "root"),
		RunRoot:     filepath.Join(tempDir, "run"),
		GraphDriver: "vfs",
	}
	
	err = cm.validateStorageConfig()
	assert.NoError(t, err)
	
	// Directories should have been created
	assert.DirExists(t, cm.config.Storage.RootDir)
	assert.DirExists(t, cm.config.Storage.RunRoot)
	
	// Test with invalid driver
	cm.config.Storage.GraphDriver = "invalid-driver"
	err = cm.validateStorageConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported storage driver")
}

func TestGettersAndSetters(t *testing.T) {
	config := &BuildConfig{From: "alpine:latest"}
	cm := NewConfigManager(config)
	
	// Test GetConfig
	retrievedConfig := cm.GetConfig()
	assert.Equal(t, config, retrievedConfig)
	assert.Equal(t, "alpine:latest", retrievedConfig.From)
	
	// Test GetSystemContext
	sysCtx := cm.GetSystemContext()
	assert.NotNil(t, sysCtx)
	
	// Multiple calls should return same instance
	sysCtx2 := cm.GetSystemContext()
	assert.Equal(t, sysCtx, sysCtx2)
}

func TestUpdateConfig(t *testing.T) {
	cm := NewConfigManager(nil)
	
	// Mark as validated first
	cm.validated = true
	
	updates := map[string]interface{}{
		"from":     "ubuntu:20.04",
		"tag":      "myimage:latest", 
		"platform": "arm64",
		"unknown":  "should be ignored",
	}
	
	err := cm.UpdateConfig(updates)
	assert.NoError(t, err)
	
	// Should have updated known fields
	assert.Equal(t, "ubuntu:20.04", cm.config.From)
	assert.Equal(t, "myimage:latest", cm.config.Tag)
	assert.Equal(t, "arm64", cm.config.Platform)
	
	// Should have marked as not validated
	assert.False(t, cm.validated)
}

func TestApplyEnvironmentVariables(t *testing.T) {
	cm := NewConfigManager(nil)
	
	// Set proxy values
	cm.config.HTTPProxy = "http://proxy.example.com:8080"
	cm.config.HTTPSProxy = "https://proxy.example.com:8080"
	
	// Store original values
	origHTTP := os.Getenv("HTTP_PROXY")
	origHTTPS := os.Getenv("HTTPS_PROXY")
	
	// Apply environment variables
	cm.ApplyEnvironmentVariables()
	
	// Check that environment variables were set
	assert.Equal(t, "http://proxy.example.com:8080", os.Getenv("HTTP_PROXY"))
	assert.Equal(t, "http://proxy.example.com:8080", os.Getenv("http_proxy"))
	assert.Equal(t, "https://proxy.example.com:8080", os.Getenv("HTTPS_PROXY"))
	assert.Equal(t, "https://proxy.example.com:8080", os.Getenv("https_proxy"))
	
	// Restore original values
	os.Setenv("HTTP_PROXY", origHTTP)
	os.Setenv("HTTPS_PROXY", origHTTPS)
}

func TestCreateBuildahOptions(t *testing.T) {
	cm := NewConfigManager(nil)
	ctx := context.Background()
	
	// Should fail with unvalidated config
	_, err := cm.CreateBuildahOptions(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not validated")
	
	// Mark as validated
	cm.validated = true
	
	// Should succeed with validated config
	options, err := cm.CreateBuildahOptions(ctx, nil)
	assert.NoError(t, err)
	
	assert.Equal(t, cm.config.From, options.FromImage)
	assert.Equal(t, cm.config.Isolation, options.Isolation)
	assert.Equal(t, cm.config.Format, options.Format)
	assert.Equal(t, cm.config.Args, options.Args)
	assert.Equal(t, cm.config.Capabilities, options.Capabilities)
	
	// Test network configuration
	cm.config.NetworkMode = "host"
	options, err = cm.CreateBuildahOptions(ctx, nil)
	assert.NoError(t, err)
	assert.Equal(t, define.NetworkDisabled, options.ConfigureNetwork)
	
	cm.config.NetworkMode = "none"
	options, err = cm.CreateBuildahOptions(ctx, nil)
	assert.NoError(t, err)
	assert.Equal(t, define.NetworkDisabled, options.ConfigureNetwork)
}