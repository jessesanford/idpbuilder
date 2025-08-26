package build

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/containers/buildah/define"
	"github.com/containers/storage"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
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

// ExtendedConfigTestSuite provides a test suite for extended configuration functionality
type ExtendedConfigTestSuite struct {
	suite.Suite
	manager    *ExtendedConfigManager
	store      *MockStore
	logger     *logrus.Logger
	tempDir    string
	configPath string
}

func (suite *ExtendedConfigTestSuite) SetupTest() {
	// Create temporary directory for test files
	tempDir, err := os.MkdirTemp("", "extended_config_test")
	suite.Require().NoError(err)
	suite.tempDir = tempDir

	// Create mock store
	suite.store = &MockStore{}
	suite.store.On("GraphDriverName").Return("overlay")
	suite.store.On("GraphRoot").Return(filepath.Join(tempDir, "graph"))
	suite.store.On("RunRoot").Return(filepath.Join(tempDir, "run"))

	// Create logger
	suite.logger = logrus.New()
	suite.logger.SetLevel(logrus.DebugLevel)

	// Create manager
	suite.manager = NewExtendedConfigManager(suite.store, suite.logger)
	suite.configPath = filepath.Join(tempDir, "config.json")
}

func (suite *ExtendedConfigTestSuite) TearDownTest() {
	if suite.tempDir != "" {
		os.RemoveAll(suite.tempDir)
	}
}

func TestExtendedConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ExtendedConfigTestSuite))
}

// Test basic configuration creation and validation
func (suite *ExtendedConfigTestSuite) TestBasicConfiguration() {
	config := &ExtendedBuildConfig{
		Name:    "test-config",
		Version: "1.0.0",
		Build: BuildOptions{
			Dockerfile: "Dockerfile",
			Context:    suite.tempDir,
			Pull:       true,
		},
		Runtime: RuntimeOptions{
			Runtime:   "runc",
			Isolation: define.IsolationDefault,
		},
		Storage: StorageOptions{
			Driver: "overlay",
			Root:   filepath.Join(suite.tempDir, "storage"),
		},
		Network: NetworkOptions{
			Mode: "bridge",
		},
		Security: SecurityOptions{
			Privileged: false,
			ReadOnly:   false,
		},
		Cache: CacheOptions{
			Enabled: true,
			Dir:     filepath.Join(suite.tempDir, "cache"),
			TTL:     time.Hour,
		},
		Features: FeatureFlags{
			LayerCaching: true,
		},
	}

	err := suite.manager.ValidateConfig(config)
	suite.NoError(err)
	suite.True(config.Cache.Enabled)
	suite.Equal("test-config", config.Name)
}

// Test configuration loading and saving
func (suite *ExtendedConfigTestSuite) TestConfigurationSaveAndLoad() {
	originalConfig := &ExtendedBuildConfig{
		Name:        "save-load-test",
		Version:     "2.0.0",
		Description: "Test configuration for save/load functionality",
		Build: BuildOptions{
			Dockerfile: "Dockerfile.test",
			Context:    suite.tempDir,
			Args: map[string]string{
				"BUILD_ARG": "test_value",
			},
			Tags: []string{"test:latest", "test:v2.0.0"},
		},
		Labels: map[string]string{
			"app":     "test-app",
			"version": "2.0.0",
		},
		Environment: map[string]string{
			"TEST_ENV": "test_value",
		},
	}

	// Save configuration
	err := suite.manager.SaveConfig(suite.configPath, originalConfig)
	suite.NoError(err)

	// Verify file exists
	suite.FileExists(suite.configPath)

	// Load configuration
	loadedConfig, err := suite.manager.LoadConfig(suite.configPath)
	suite.NoError(err)
	suite.NotNil(loadedConfig)

	// Verify loaded configuration matches original
	suite.Equal(originalConfig.Name, loadedConfig.Name)
	suite.Equal(originalConfig.Version, loadedConfig.Version)
	suite.Equal(originalConfig.Description, loadedConfig.Description)
	suite.Equal(originalConfig.Build.Dockerfile, loadedConfig.Build.Dockerfile)
	suite.Equal(originalConfig.Build.Args["BUILD_ARG"], loadedConfig.Build.Args["BUILD_ARG"])
	suite.Equal(len(originalConfig.Build.Tags), len(loadedConfig.Build.Tags))
	suite.Equal(originalConfig.Labels["app"], loadedConfig.Labels["app"])
}

// Test configuration inheritance
func (suite *ExtendedConfigTestSuite) TestConfigurationInheritance() {
	// Create parent configuration
	parentConfig := &ExtendedBuildConfig{
		Name:    "parent-config",
		Version: "1.0.0",
		Build: BuildOptions{
			Dockerfile: "Dockerfile.base",
			Pull:       true,
		},
		Runtime: RuntimeOptions{
			Runtime: "runc",
			Memory:  1024 * 1024 * 1024, // 1GB
		},
		Labels: map[string]string{
			"base": "true",
			"type": "parent",
		},
	}

	parentPath := filepath.Join(suite.tempDir, "parent.json")
	err := suite.manager.SaveConfig(parentPath, parentConfig)
	suite.NoError(err)

	// Create child configuration that inherits from parent
	childConfig := &ExtendedBuildConfig{
		Name:    "child-config",
		Version: "2.0.0",
		Inherits: []string{parentPath},
		Build: BuildOptions{
			Context: suite.tempDir, // Child-specific override
		},
		Labels: map[string]string{
			"type":  "child", // Override parent label
			"child": "true",  // Add new label
		},
	}

	childPath := filepath.Join(suite.tempDir, "child.json")
	err = suite.manager.SaveConfig(childPath, childConfig)
	suite.NoError(err)

	// Load child configuration (should resolve inheritance)
	loadedConfig, err := suite.manager.LoadConfig(childPath)
	suite.NoError(err)

	// Verify inheritance was resolved
	suite.Equal("child-config", loadedConfig.Name)
	suite.Equal("2.0.0", loadedConfig.Version)
	suite.Equal("Dockerfile.base", loadedConfig.Build.Dockerfile) // Inherited from parent
	suite.Equal(suite.tempDir, loadedConfig.Build.Context)        // Child override
	suite.True(loadedConfig.Build.Pull)                           // Inherited from parent
	suite.Equal(int64(1024*1024*1024), loadedConfig.Runtime.Memory) // Inherited from parent
	suite.Equal("true", loadedConfig.Labels["base"])              // Inherited from parent
	suite.Equal("child", loadedConfig.Labels["type"])             // Child override
	suite.Equal("true", loadedConfig.Labels["child"])             // Child addition
}

// Test profile application
func (suite *ExtendedConfigTestSuite) TestProfileApplication() {
	config := &ExtendedBuildConfig{
		Name:    "profile-test",
		Version: "1.0.0",
		Profile: "development",
		Profiles: map[string]Config{
			"development": {
				Build: &BuildOptions{
					NoCache: true,
					Pull:    true,
				},
				Runtime: &RuntimeOptions{
					Memory: 512 * 1024 * 1024, // 512MB
				},
				Features: &FeatureFlags{
					Experimental: true,
					LayerCaching: false,
				},
				Labels: map[string]string{
					"env":   "dev",
					"debug": "true",
				},
			},
			"production": {
				Build: &BuildOptions{
					NoCache: false,
					Pull:    false,
				},
				Runtime: &RuntimeOptions{
					Memory: 2 * 1024 * 1024 * 1024, // 2GB
				},
				Features: &FeatureFlags{
					Experimental: false,
					LayerCaching: true,
				},
				Labels: map[string]string{
					"env": "prod",
				},
			},
		},
	}

	err := suite.manager.applyProfile(config)
	suite.NoError(err)

	// Verify development profile was applied
	suite.True(config.Build.NoCache)
	suite.True(config.Build.Pull)
	suite.Equal(int64(512*1024*1024), config.Runtime.Memory)
	suite.True(config.Features.Experimental)
	suite.False(config.Features.LayerCaching)
	suite.Equal("dev", config.Labels["env"])
	suite.Equal("true", config.Labels["debug"])
}

// Test environment variable overrides
func (suite *ExtendedConfigTestSuite) TestEnvironmentOverrides() {
	// Set test environment variables
	os.Setenv("TEST_BUILD_ARG", "env_override_value")
	os.Setenv("TEST_MEMORY", "2147483648") // 2GB
	defer os.Unsetenv("TEST_BUILD_ARG")
	defer os.Unsetenv("TEST_MEMORY")

	config := &ExtendedBuildConfig{
		Name:    "env-test",
		Version: "1.0.0",
		Environment: map[string]string{
			"TEST_BUILD_ARG": "original_value",
			"TEST_MEMORY":    "1073741824", // 1GB
			"TEST_NO_OVERRIDE": "keep_this",
		},
	}

	err := suite.manager.applyEnvironmentOverrides(config)
	suite.NoError(err)

	// Verify environment overrides were applied
	suite.Equal("env_override_value", config.Environment["TEST_BUILD_ARG"])
	suite.Equal("2147483648", config.Environment["TEST_MEMORY"])
	suite.Equal("keep_this", config.Environment["TEST_NO_OVERRIDE"])
}

// Test validation rules
func (suite *ExtendedConfigTestSuite) TestValidationRules() {
	// Test invalid configuration name
	invalidConfig := &ExtendedBuildConfig{
		Name: "", // Invalid: empty name
	}

	err := suite.manager.ValidateConfig(invalidConfig)
	suite.Error(err)
	suite.Contains(err.Error(), "name is required")

	// Test invalid storage driver
	invalidConfig2 := &ExtendedBuildConfig{
		Name: "test",
		Storage: StorageOptions{
			Driver: "invalid_driver", // Invalid storage driver
		},
	}

	err = suite.manager.ValidateConfig(invalidConfig2)
	suite.Error(err)
	suite.Contains(err.Error(), "invalid storage driver")

	// Test invalid network mode
	invalidConfig3 := &ExtendedBuildConfig{
		Name: "test",
		Network: NetworkOptions{
			Mode: "invalid_mode", // Invalid network mode
		},
	}

	err = suite.manager.ValidateConfig(invalidConfig3)
	suite.Error(err)
	suite.Contains(err.Error(), "invalid network mode")

	// Test conflicting security options
	invalidConfig4 := &ExtendedBuildConfig{
		Name: "test",
		Security: SecurityOptions{
			Privileged: true,
			NoNewPrivs: true, // Conflict: privileged and no-new-privs
		},
	}

	err = suite.manager.ValidateConfig(invalidConfig4)
	suite.Error(err)
	suite.Contains(err.Error(), "mutually exclusive")
}

// Test configuration transformations
func (suite *ExtendedConfigTestSuite) TestConfigurationTransformations() {
	// Create configuration with relative paths
	relativeDir := "relative/path"
	config := &ExtendedBuildConfig{
		Name: "transform-test",
		Build: BuildOptions{
			Context: relativeDir,
		},
		Storage: StorageOptions{
			Root: relativeDir,
		},
		Cache: CacheOptions{
			Dir: relativeDir,
		},
	}

	// Apply transformations
	transformedConfig, err := suite.manager.applyTransformations(config)
	suite.NoError(err)

	// Verify paths were expanded to absolute paths
	suite.True(filepath.IsAbs(transformedConfig.Build.Context))
	suite.True(filepath.IsAbs(transformedConfig.Storage.Root))
	suite.True(filepath.IsAbs(transformedConfig.Cache.Dir))

	// Verify defaults were set
	suite.Equal("Dockerfile", transformedConfig.Build.Dockerfile)
	suite.Equal("runc", transformedConfig.Runtime.Runtime)
	suite.Equal("overlay", transformedConfig.Storage.Driver)
	suite.Equal("bridge", transformedConfig.Network.Mode)
}

// Test configuration merging
func (suite *ExtendedConfigTestSuite) TestConfigurationMerging() {
	parent := &ExtendedBuildConfig{
		Name: "parent",
		Build: BuildOptions{
			Dockerfile: "Dockerfile.base",
			Pull:       true,
			Args: map[string]string{
				"BASE_ARG":   "base_value",
				"SHARED_ARG": "parent_value",
			},
		},
		Labels: map[string]string{
			"base":   "true",
			"shared": "parent",
		},
	}

	child := &ExtendedBuildConfig{
		Name: "child",
		Build: BuildOptions{
			Context: "/child/context",
			Args: map[string]string{
				"CHILD_ARG":  "child_value",
				"SHARED_ARG": "child_value", // Override parent
			},
		},
		Labels: map[string]string{
			"child":  "true",
			"shared": "child", // Override parent
		},
	}

	err := suite.manager.mergeConfigs(child, parent)
	suite.NoError(err)

	// Verify merge results
	suite.Equal("child", child.Name) // Child name preserved
	suite.Equal("Dockerfile.base", child.Build.Dockerfile) // Parent value inherited
	suite.Equal("/child/context", child.Build.Context) // Child value preserved
	suite.True(child.Build.Pull) // Parent value inherited

	// Verify map merging
	suite.Equal("base_value", child.Build.Args["BASE_ARG"]) // Parent value inherited
	suite.Equal("child_value", child.Build.Args["CHILD_ARG"]) // Child value preserved
	suite.Equal("child_value", child.Build.Args["SHARED_ARG"]) // Child overrides parent

	suite.Equal("true", child.Labels["base"]) // Parent value inherited
	suite.Equal("true", child.Labels["child"]) // Child value preserved
	suite.Equal("child", child.Labels["shared"]) // Child overrides parent
}

// Test configuration caching
func (suite *ExtendedConfigTestSuite) TestConfigurationCaching() {
	config := &ExtendedBuildConfig{
		Name:    "cache-test",
		Version: "1.0.0",
	}

	// Save configuration
	err := suite.manager.SaveConfig(suite.configPath, config)
	suite.NoError(err)

	// Load configuration first time
	startTime := time.Now()
	loadedConfig1, err := suite.manager.LoadConfig(suite.configPath)
	firstLoadTime := time.Since(startTime)
	suite.NoError(err)
	suite.NotNil(loadedConfig1)

	// Load configuration second time (should be from cache)
	startTime = time.Now()
	loadedConfig2, err := suite.manager.LoadConfig(suite.configPath)
	secondLoadTime := time.Since(startTime)
	suite.NoError(err)
	suite.NotNil(loadedConfig2)

	// Second load should be significantly faster (cached)
	suite.True(secondLoadTime < firstLoadTime)

	// Clear cache and test cache is actually cleared
	suite.manager.ClearCache()
	startTime = time.Now()
	loadedConfig3, err := suite.manager.LoadConfig(suite.configPath)
	thirdLoadTime := time.Since(startTime)
	suite.NoError(err)
	suite.NotNil(loadedConfig3)

	// Third load should be similar to first load (cache cleared)
	suite.True(thirdLoadTime > secondLoadTime)
}

// Test configuration summary
func (suite *ExtendedConfigTestSuite) TestConfigurationSummary() {
	config := &ExtendedBuildConfig{
		Name:      "summary-test",
		Version:   "1.0.0",
		Profile:   "development",
		validated: true,
		merged:    true,
		Inherits:  []string{"parent1.json", "parent2.json"},
		Profiles: map[string]Config{
			"dev":  {},
			"prod": {},
		},
		Features: FeatureFlags{
			Experimental: true,
			LayerCaching: true,
		},
		Cache: CacheOptions{
			Enabled: true,
		},
		Labels: map[string]string{
			"app": "test",
			"env": "dev",
		},
		Annotations: map[string]string{
			"description": "test config",
		},
	}

	summary := suite.manager.GetConfigurationSummary(config)

	suite.Equal("summary-test", summary["name"])
	suite.Equal("1.0.0", summary["version"])
	suite.Equal("development", summary["profile"])
	suite.Equal(true, summary["validated"])
	suite.Equal(true, summary["merged"])
	suite.Equal(2, summary["inherits"])
	suite.Equal(2, summary["profiles"])
	suite.Equal(config.Features, summary["features"])
	suite.Equal(true, summary["cache"])
	suite.Equal(2, summary["labels"])
	suite.Equal(1, summary["annotations"])
}

// Test cache TTL functionality
func (suite *ExtendedConfigTestSuite) TestCacheTTL() {
	config := &ExtendedBuildConfig{
		Name:    "ttl-test",
		Version: "1.0.0",
	}

	// Set very short TTL for testing
	suite.manager.SetCacheTTL(time.Millisecond * 100)

	// Save and load config
	err := suite.manager.SaveConfig(suite.configPath, config)
	suite.NoError(err)

	loadedConfig1, err := suite.manager.LoadConfig(suite.configPath)
	suite.NoError(err)
	suite.NotNil(loadedConfig1)

	// Wait for TTL to expire
	time.Sleep(time.Millisecond * 150)

	// Load again - should reload from file due to TTL expiration
	loadedConfig2, err := suite.manager.LoadConfig(suite.configPath)
	suite.NoError(err)
	suite.NotNil(loadedConfig2)

	// Both configs should be equivalent but not the same instance
	suite.Equal(loadedConfig1.Name, loadedConfig2.Name)
	suite.Equal(loadedConfig1.Version, loadedConfig2.Version)
}

// Test cache enable/disable functionality
func (suite *ExtendedConfigTestSuite) TestCacheEnableDisable() {
	config := &ExtendedBuildConfig{
		Name:    "cache-toggle-test",
		Version: "1.0.0",
	}

	// Save configuration
	err := suite.manager.SaveConfig(suite.configPath, config)
	suite.NoError(err)

	// Load with caching enabled
	_, err = suite.manager.LoadConfig(suite.configPath)
	suite.NoError(err)

	// Verify config is cached
	suite.Equal(1, len(suite.manager.configs))

	// Disable caching
	suite.manager.EnableCache(false)
	suite.Equal(0, len(suite.manager.configs)) // Cache should be cleared

	// Load again - should not be cached
	_, err = suite.manager.LoadConfig(suite.configPath)
	suite.NoError(err)
	suite.Equal(0, len(suite.manager.configs)) // Still not cached

	// Re-enable caching
	suite.manager.EnableCache(true)
	_, err = suite.manager.LoadConfig(suite.configPath)
	suite.NoError(err)
	suite.Equal(1, len(suite.manager.configs)) // Now cached again
}

// Test unit tests for individual functions
func TestValidationRules(t *testing.T) {
	tests := []struct {
		name        string
		config      *ExtendedBuildConfig
		rule        ValidationRule
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid name",
			config: &ExtendedBuildConfig{
				Name: "valid-config-name",
			},
			rule:        validateName,
			expectError: false,
		},
		{
			name: "empty name",
			config: &ExtendedBuildConfig{
				Name: "",
			},
			rule:        validateName,
			expectError: true,
			errorMsg:    "name is required",
		},
		{
			name: "name too long",
			config: &ExtendedBuildConfig{
				Name: strings.Repeat("a", 254),
			},
			rule:        validateName,
			expectError: true,
			errorMsg:    "name too long",
		},
		{
			name: "valid storage driver",
			config: &ExtendedBuildConfig{
				Name: "test",
				Storage: StorageOptions{
					Driver: "overlay",
				},
			},
			rule:        validateStorageDriver,
			expectError: false,
		},
		{
			name: "invalid storage driver",
			config: &ExtendedBuildConfig{
				Name: "test",
				Storage: StorageOptions{
					Driver: "invalid",
				},
			},
			rule:        validateStorageDriver,
			expectError: true,
			errorMsg:    "invalid storage driver",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rule(tt.config)
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

func TestTransformers(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "transformer_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name        string
		config      *ExtendedBuildConfig
		transformer ConfigTransformer
		expectError bool
		validate    func(*testing.T, *ExtendedBuildConfig)
	}{
		{
			name: "expand paths",
			config: &ExtendedBuildConfig{
				Name: "test",
				Build: BuildOptions{
					Context: "relative/path",
				},
				Storage: StorageOptions{
					Root: "storage/path",
				},
				Cache: CacheOptions{
					Dir: "cache/path",
				},
			},
			transformer: expandPaths,
			expectError: false,
			validate: func(t *testing.T, config *ExtendedBuildConfig) {
				assert.True(t, filepath.IsAbs(config.Build.Context))
				assert.True(t, filepath.IsAbs(config.Storage.Root))
				assert.True(t, filepath.IsAbs(config.Cache.Dir))
			},
		},
		{
			name: "set defaults",
			config: &ExtendedBuildConfig{
				Name: "test",
			},
			transformer: setDefaults,
			expectError: false,
			validate: func(t *testing.T, config *ExtendedBuildConfig) {
				assert.Equal(t, "Dockerfile", config.Build.Dockerfile)
				assert.Equal(t, "runc", config.Runtime.Runtime)
				assert.Equal(t, "overlay", config.Storage.Driver)
				assert.Equal(t, "bridge", config.Network.Mode)
			},
		},
		{
			name: "optimize settings - parallel builds",
			config: &ExtendedBuildConfig{
				Name: "test",
				Features: FeatureFlags{
					ParallelBuilds: true,
				},
			},
			transformer: optimizeSettings,
			expectError: false,
			validate: func(t *testing.T, config *ExtendedBuildConfig) {
				assert.Equal(t, uint64(1024), config.Runtime.CPUShares)
			},
		},
		{
			name: "optimize settings - layer caching",
			config: &ExtendedBuildConfig{
				Name: "test",
				Features: FeatureFlags{
					LayerCaching: true,
				},
			},
			transformer: optimizeSettings,
			expectError: false,
			validate: func(t *testing.T, config *ExtendedBuildConfig) {
				assert.True(t, config.Cache.Enabled)
				assert.Contains(t, config.Cache.Dir, "buildah-cache")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.transformer(tt.config)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
		})
	}
}

// Benchmark tests
func BenchmarkConfigurationLoading(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "benchmark_test")
	require.NoError(b, err)
	defer os.RemoveAll(tempDir)

	store := &MockStore{}
	store.On("GraphDriverName").Return("overlay")
	store.On("GraphRoot").Return(filepath.Join(tempDir, "graph"))
	store.On("RunRoot").Return(filepath.Join(tempDir, "run"))

	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel) // Reduce logging noise in benchmarks

	manager := NewExtendedConfigManager(store, logger)
	configPath := filepath.Join(tempDir, "benchmark_config.json")

	config := &ExtendedBuildConfig{
		Name:    "benchmark-config",
		Version: "1.0.0",
		Build: BuildOptions{
			Dockerfile: "Dockerfile",
			Context:    tempDir,
		},
	}

	err = manager.SaveConfig(configPath, config)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := manager.LoadConfig(configPath)
		require.NoError(b, err)
	}
}