package build

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/containers/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockStore implements a mock storage.Store for testing
type MockStore struct {
	mock.Mock
}

func (m *MockStore) Shutdown(force bool) ([]string, error) {
	args := m.Called(force)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockStore) GraphDriverName() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockStore) GraphRoot() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockStore) RunRoot() string {
	args := m.Called()
	return args.String(0)
}

// MockLogger implements a mock Logger for testing
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Infof(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Warnf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Errorf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Debugf(format string, args ...interface{}) {
	m.Called(format, args)
}

func TestExtendedConfigManager_New(t *testing.T) {
	store := &MockStore{}
	logger := &MockLogger{}

	manager := NewExtendedConfigManager(store, logger)

	assert.NotNil(t, manager)
	assert.Equal(t, store, manager.store)
	assert.Equal(t, logger, manager.logger)
	assert.NotNil(t, manager.cache)
	assert.Equal(t, 5*time.Minute, manager.cacheTTL)
}

func TestExtendedConfigManager_SaveAndLoad(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "extended_config_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	store := &MockStore{}
	store.On("GraphDriverName").Return("overlay")
	store.On("GraphRoot").Return(filepath.Join(tempDir, "graph"))
	store.On("RunRoot").Return(filepath.Join(tempDir, "run"))

	logger := &MockLogger{}
	logger.On("Infof", mock.Anything, mock.Anything).Return()
	logger.On("Debugf", mock.Anything, mock.Anything).Return()

	manager := NewExtendedConfigManager(store, logger)
	configPath := filepath.Join(tempDir, "test_config.json")

	// Test configuration
	config := &ExtendedBuildConfig{
		Name:    "test-config",
		Version: "1.0.0",
		Build: BuildOptions{
			Dockerfile: "Dockerfile.test",
			Context:    tempDir,
		},
		Runtime: RuntimeOptions{
			Runtime: "runc",
			Memory:  "512m",
		},
		Storage: StorageOptions{
			Driver: "overlay",
		},
		Features: FeatureFlags{
			LayerCaching: true,
		},
	}

	// Test save
	err = manager.SaveConfig(configPath, config)
	require.NoError(t, err)

	// Verify file exists
	assert.FileExists(t, configPath)

	// Test load
	loadedConfig, err := manager.LoadConfig(configPath)
	require.NoError(t, err)

	assert.Equal(t, config.Name, loadedConfig.Name)
	assert.Equal(t, config.Version, loadedConfig.Version)
	assert.Equal(t, config.Build.Dockerfile, loadedConfig.Build.Dockerfile)
	assert.Equal(t, config.Runtime.Runtime, loadedConfig.Runtime.Runtime)
	assert.Equal(t, config.Storage.Driver, loadedConfig.Storage.Driver)
}

func TestExtendedConfigManager_Inheritance(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "inheritance_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	store := &MockStore{}
	logger := &MockLogger{}
	logger.On("Infof", mock.Anything, mock.Anything).Return()
	logger.On("Debugf", mock.Anything, mock.Anything).Return()

	manager := NewExtendedConfigManager(store, logger)

	// Create parent config
	parentConfig := &ExtendedBuildConfig{
		Name:        "parent-config",
		Description: "Parent configuration",
		Environment: map[string]string{
			"PARENT_VAR": "parent_value",
		},
		Build: BuildOptions{
			Dockerfile: "Dockerfile.parent",
			Context:    tempDir,
			Args: map[string]string{
				"PARENT_ARG": "parent_arg_value",
			},
		},
	}

	parentPath := filepath.Join(tempDir, "parent.json")
	err = manager.SaveConfig(parentPath, parentConfig)
	require.NoError(t, err)

	// Create child config with inheritance
	childConfig := &ExtendedBuildConfig{
		Name: "child-config",
		Inherits: []string{"parent.json"},
		Environment: map[string]string{
			"CHILD_VAR": "child_value",
		},
		Build: BuildOptions{
			Args: map[string]string{
				"CHILD_ARG": "child_arg_value",
			},
		},
	}

	childPath := filepath.Join(tempDir, "child.json")
	err = manager.SaveConfig(childPath, childConfig)
	require.NoError(t, err)

	// Load child config (should resolve inheritance)
	loadedChild, err := manager.LoadConfig(childPath)
	require.NoError(t, err)

	// Check inherited values
	assert.Equal(t, "child-config", loadedChild.Name)
	assert.Equal(t, "Parent configuration", loadedChild.Description) // Inherited
	assert.Equal(t, "Dockerfile.parent", loadedChild.Build.Dockerfile) // Inherited
	assert.Equal(t, tempDir, loadedChild.Build.Context) // Inherited

	// Check merged environment variables
	assert.Equal(t, "parent_value", loadedChild.Environment["PARENT_VAR"]) // Inherited
	assert.Equal(t, "child_value", loadedChild.Environment["CHILD_VAR"])   // Child's own

	// Check merged build args
	assert.Equal(t, "parent_arg_value", loadedChild.Build.Args["PARENT_ARG"]) // Inherited
	assert.Equal(t, "child_arg_value", loadedChild.Build.Args["CHILD_ARG"])   // Child's own
}

func TestExtendedConfigManager_Profiles(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "profile_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	store := &MockStore{}
	logger := &MockLogger{}
	logger.On("Infof", mock.Anything, mock.Anything).Return()
	logger.On("Debugf", mock.Anything, mock.Anything).Return()

	manager := NewExtendedConfigManager(store, logger)

	// Create config with profiles
	config := &ExtendedBuildConfig{
		Name:    "profile-config",
		Profile: "development",
		Profiles: map[string]Profile{
			"development": {
				Name:        "Development Profile",
				Description: "Development environment settings",
				Settings: map[string]interface{}{
					"build.dockerfile": "Dockerfile.dev",
					"runtime.memory":   "1g",
					"cache.enabled":    true,
				},
			},
		},
	}

	configPath := filepath.Join(tempDir, "profile_config.json")
	err = manager.SaveConfig(configPath, config)
	require.NoError(t, err)

	// Load config (should apply profile)
	loadedConfig, err := manager.LoadConfig(configPath)
	require.NoError(t, err)

	// Check profile was applied
	assert.Equal(t, "Dockerfile.dev", loadedConfig.Build.Dockerfile)
	assert.Equal(t, "1g", loadedConfig.Runtime.Memory)
	assert.True(t, loadedConfig.Cache.Enabled)
}

func TestExtendedConfigManager_EnvironmentOverrides(t *testing.T) {
	// Set environment variables
	os.Setenv("BUILDAH_DOCKERFILE", "Dockerfile.env")
	os.Setenv("BUILDAH_RUNTIME", "crun")
	defer os.Unsetenv("BUILDAH_DOCKERFILE")
	defer os.Unsetenv("BUILDAH_RUNTIME")

	tempDir, err := os.MkdirTemp("", "env_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	store := &MockStore{}
	logger := &MockLogger{}
	logger.On("Infof", mock.Anything, mock.Anything).Return()
	logger.On("Debugf", mock.Anything, mock.Anything).Return()

	manager := NewExtendedConfigManager(store, logger)

	config := &ExtendedBuildConfig{
		Name: "env-config",
	}

	configPath := filepath.Join(tempDir, "env_config.json")
	err = manager.SaveConfig(configPath, config)
	require.NoError(t, err)

	// Load config (should apply environment overrides)
	loadedConfig, err := manager.LoadConfig(configPath)
	require.NoError(t, err)

	// Check environment overrides were applied
	assert.Equal(t, "Dockerfile.env", loadedConfig.Build.Dockerfile)
	assert.Equal(t, "crun", loadedConfig.Runtime.Runtime)
}

func TestExtendedConfigManager_Cache(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "cache_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	store := &MockStore{}
	logger := &MockLogger{}
	logger.On("Infof", mock.Anything, mock.Anything).Return()
	logger.On("Debugf", mock.Anything, mock.Anything).Return()

	manager := NewExtendedConfigManager(store, logger)

	config := &ExtendedBuildConfig{
		Name: "cache-test",
	}

	configPath := filepath.Join(tempDir, "cache_config.json")
	err = manager.SaveConfig(configPath, config)
	require.NoError(t, err)

	// First load
	_, err = manager.LoadConfig(configPath)
	require.NoError(t, err)
	assert.Equal(t, 1, len(manager.cache))

	// Second load (should use cache)
	_, err = manager.LoadConfig(configPath)
	require.NoError(t, err)
	assert.Equal(t, 1, len(manager.cache))

	// Clear cache
	manager.ClearCache()
	assert.Equal(t, 0, len(manager.cache))

	// Test cache disable
	manager.EnableCache(false)
	_, err = manager.LoadConfig(configPath)
	require.NoError(t, err)
	assert.Equal(t, 0, len(manager.cache)) // Should not be cached

	// Re-enable cache
	manager.EnableCache(true)
	_, err = manager.LoadConfig(configPath)
	require.NoError(t, err)
	assert.Equal(t, 1, len(manager.cache)) // Now cached
}

func TestExtendedConfigManager_Summary(t *testing.T) {
	store := &MockStore{}
	logger := &MockLogger{}

	manager := NewExtendedConfigManager(store, logger)

	config := &ExtendedBuildConfig{
		Name:    "summary-test",
		Version: "2.0.0",
		Profile: "production",
		Build: BuildOptions{
			Dockerfile: "Dockerfile.prod",
			Context:    "/app",
			Pull:       true,
		},
		Runtime: RuntimeOptions{
			Runtime: "runc",
			Memory:  "2g",
			CPUs:    "2",
		},
		Storage: StorageOptions{
			Driver: "overlay2",
			Root:   "/var/lib/storage",
		},
		Features: FeatureFlags{
			ParallelBuilds: true,
			LayerCaching:   true,
		},
	}

	summary := manager.GetConfigSummary(config)

	assert.Equal(t, "summary-test", summary["name"])
	assert.Equal(t, "2.0.0", summary["version"])
	assert.Equal(t, "production", summary["profile"])

	buildInfo := summary["build"].(map[string]interface{})
	assert.Equal(t, "Dockerfile.prod", buildInfo["dockerfile"])
	assert.Equal(t, "/app", buildInfo["context"])
	assert.True(t, buildInfo["pull"].(bool))

	runtimeInfo := summary["runtime"].(map[string]interface{})
	assert.Equal(t, "runc", runtimeInfo["runtime"])
	assert.Equal(t, "2g", runtimeInfo["memory"])
	assert.Equal(t, "2", runtimeInfo["cpus"])

	storageInfo := summary["storage"].(map[string]interface{})
	assert.Equal(t, "overlay2", storageInfo["driver"])
	assert.Equal(t, "/var/lib/storage", storageInfo["root"])

	featuresInfo := summary["features"].(map[string]interface{})
	assert.True(t, featuresInfo["parallel_builds"].(bool))
	assert.True(t, featuresInfo["layer_caching"].(bool))
}

func TestConfigValidation(t *testing.T) {
	store := &MockStore{}
	logger := &MockLogger{}
	logger.On("Infof", mock.Anything, mock.Anything).Return()
	logger.On("Debugf", mock.Anything, mock.Anything).Return()

	manager := NewExtendedConfigManager(store, logger)

	tests := []struct {
		name        string
		config      *ExtendedBuildConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid config",
			config: &ExtendedBuildConfig{
				Name: "valid-config",
			},
			expectError: false,
		},
		{
			name: "empty name",
			config: &ExtendedBuildConfig{
				Name: "",
			},
			expectError: true,
			errorMsg:    "name is required",
		},
		{
			name: "name too long",
			config: &ExtendedBuildConfig{
				Name: string(make([]byte, 254)), // 254 characters
			},
			expectError: true,
			errorMsg:    "name too long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.validateConfig(tt.config)
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}