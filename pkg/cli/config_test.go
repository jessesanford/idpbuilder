package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	t.Run("load with defaults when no config file", func(t *testing.T) {
		config, err := LoadConfig("")
		require.NoError(t, err)
		assert.NotNil(t, config)

		// Check defaults
		assert.Equal(t, "https://gitea.cnoe.localtest.me:443", config.Registry.URL)
		assert.Equal(t, "gitea_admin", config.Registry.Username)
		assert.False(t, config.Registry.Insecure)
		assert.Equal(t, 30, config.Registry.Timeout)

		assert.Equal(t, ".", config.Build.Context)
		assert.Equal(t, "linux/amd64", config.Build.Platform)
		assert.False(t, config.Build.NoCache)

		assert.Equal(t, "info", config.Logging.Level)
		assert.Equal(t, "text", config.Logging.Format)
	})

	t.Run("load from valid YAML file", func(t *testing.T) {
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "config.yaml")

		yamlContent := `
registry:
  url: "https://custom-registry.example.com"
  username: "testuser"
  password: "testpass"
  insecure: true
  timeout: 60
build:
  context: "./custom"
  platform: "linux/arm64"
  no_cache: true
  parallelism: 4
logging:
  level: "debug"
  format: "json"
  file: "/var/log/idpbuilder.log"
`
		err := os.WriteFile(configFile, []byte(yamlContent), 0644)
		require.NoError(t, err)

		config, err := LoadConfig(configFile)
		require.NoError(t, err)
		assert.NotNil(t, config)

		// Check loaded values
		assert.Equal(t, "https://custom-registry.example.com", config.Registry.URL)
		assert.Equal(t, "testuser", config.Registry.Username)
		assert.Equal(t, "testpass", config.Registry.Password)
		assert.True(t, config.Registry.Insecure)
		assert.Equal(t, 60, config.Registry.Timeout)

		assert.Equal(t, "./custom", config.Build.Context)
		assert.Equal(t, "linux/arm64", config.Build.Platform)
		assert.True(t, config.Build.NoCache)
		assert.Equal(t, 4, config.Build.Parallelism)

		assert.Equal(t, "debug", config.Logging.Level)
		assert.Equal(t, "json", config.Logging.Format)
		assert.Equal(t, "/var/log/idpbuilder.log", config.Logging.File)
	})

	t.Run("load from nonexistent file returns error", func(t *testing.T) {
		config, err := LoadConfig("/nonexistent/config.yaml")
		assert.Error(t, err)
		assert.Nil(t, config)
		assert.Contains(t, err.Error(), "failed to read config")
	})

	t.Run("load from invalid YAML returns error", func(t *testing.T) {
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "invalid.yaml")

		invalidYAML := `
registry:
  url: "https://example.com"
  invalid_yaml_here: [unclosed bracket
`
		err := os.WriteFile(configFile, []byte(invalidYAML), 0644)
		require.NoError(t, err)

		config, err := LoadConfig(configFile)
		assert.Error(t, err)
		assert.Nil(t, config)
	})
}

func TestSaveConfig(t *testing.T) {
	t.Run("save config to file", func(t *testing.T) {
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "saved_config.yaml")

		config := &Config{
			Registry: RegistryConfig{
				URL:      "https://test-registry.example.com",
				Username: "testuser",
				Password: "testpass",
				Insecure: true,
				Timeout:  45,
			},
			Build: BuildConfig{
				Context:     "./test",
				Platform:    "linux/arm64",
				NoCache:     true,
				Parallelism: 2,
			},
			Logging: LoggingConfig{
				Level:  "warn",
				Format: "json",
			},
		}

		err := SaveConfig(config, configFile)
		require.NoError(t, err)

		// Verify file was created
		assert.FileExists(t, configFile)

		// Load it back and verify
		loadedConfig, err := LoadConfig(configFile)
		require.NoError(t, err)

		assert.Equal(t, config.Registry.URL, loadedConfig.Registry.URL)
		assert.Equal(t, config.Registry.Username, loadedConfig.Registry.Username)
		assert.Equal(t, config.Build.Platform, loadedConfig.Build.Platform)
		assert.Equal(t, config.Logging.Level, loadedConfig.Logging.Level)
	})

	t.Run("save config with empty path returns error", func(t *testing.T) {
		config := &Config{}
		err := SaveConfig(config, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config path cannot be empty")
	})
}

func TestExpandPath(t *testing.T) {
	t.Run("expand environment variables", func(t *testing.T) {
		os.Setenv("TEST_VAR", "test_value")
		defer os.Unsetenv("TEST_VAR")

		path, err := expandPath("$TEST_VAR/subdir")
		require.NoError(t, err)
		assert.Equal(t, "test_value/subdir", path)
	})

	t.Run("expand tilde", func(t *testing.T) {
		path, err := expandPath("~/test")
		require.NoError(t, err)

		home, _ := os.UserHomeDir()
		expected := filepath.Join(home, "test")
		assert.Equal(t, expected, path)
	})

	t.Run("regular path unchanged", func(t *testing.T) {
		originalPath := "/regular/path"
		path, err := expandPath(originalPath)
		require.NoError(t, err)
		assert.Equal(t, originalPath, path)
	})
}

func TestGetDefaultConfigPath(t *testing.T) {
	path := GetDefaultConfigPath()
	assert.NotEmpty(t, path)
	assert.Contains(t, path, ".idpbuilder")
	assert.Contains(t, path, "config.yaml")
}
