package gitea

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigFileProvider(t *testing.T) {
	t.Run("WithValidConfigFile", func(t *testing.T) {
		// Create temporary config file
		tmpDir := t.TempDir()
		configDir := filepath.Join(tmpDir, ".idpbuilder")
		configPath := filepath.Join(configDir, "config")

		err := os.MkdirAll(configDir, 0755)
		require.NoError(t, err)

		config := Config{
			Registries: map[string]RegistryCredentials{
				"gitea": {
					Username: "configuser",
					Password: "configpass",
					URL:      "https://gitea.example.com",
				},
			},
		}

		data, err := json.Marshal(config)
		require.NoError(t, err)

		err = os.WriteFile(configPath, data, 0600)
		require.NoError(t, err)

		// Test provider
		provider := &ConfigFileProvider{
			configPath: configPath,
		}

		assert.True(t, provider.IsAvailable())
		assert.Equal(t, 3, provider.Priority())

		username, err := provider.GetUsername()
		require.NoError(t, err)
		assert.Equal(t, "configuser", username)

		password, err := provider.GetPassword()
		require.NoError(t, err)
		assert.Equal(t, "configpass", password)
	})

	t.Run("WithDefaultRegistryConfig", func(t *testing.T) {
		// Create temporary config file with default registry
		tmpDir := t.TempDir()
		configDir := filepath.Join(tmpDir, ".idpbuilder")
		configPath := filepath.Join(configDir, "config")

		err := os.MkdirAll(configDir, 0755)
		require.NoError(t, err)

		config := Config{
			Registries: map[string]RegistryCredentials{
				"default": {
					Username: "defaultuser",
					Password: "defaultpass",
					URL:      "https://default.example.com",
				},
			},
		}

		data, err := json.Marshal(config)
		require.NoError(t, err)

		err = os.WriteFile(configPath, data, 0600)
		require.NoError(t, err)

		// Test provider
		provider := &ConfigFileProvider{
			configPath: configPath,
		}

		assert.True(t, provider.IsAvailable())

		username, err := provider.GetUsername()
		require.NoError(t, err)
		assert.Equal(t, "defaultuser", username)

		password, err := provider.GetPassword()
		require.NoError(t, err)
		assert.Equal(t, "defaultpass", password)
	})

	t.Run("WithNoGiteaCredentials", func(t *testing.T) {
		// Create temporary config file without gitea credentials
		tmpDir := t.TempDir()
		configDir := filepath.Join(tmpDir, ".idpbuilder")
		configPath := filepath.Join(configDir, "config")

		err := os.MkdirAll(configDir, 0755)
		require.NoError(t, err)

		config := Config{
			Registries: map[string]RegistryCredentials{
				"othercreds": {
					Username: "otheruser",
					Password: "otherpass",
					URL:      "https://other.example.com",
				},
			},
		}

		data, err := json.Marshal(config)
		require.NoError(t, err)

		err = os.WriteFile(configPath, data, 0600)
		require.NoError(t, err)

		// Test provider
		provider := &ConfigFileProvider{
			configPath: configPath,
		}

		assert.True(t, provider.IsAvailable()) // File exists

		_, err = provider.GetUsername()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no gitea credentials in config")

		_, err = provider.GetPassword()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no gitea credentials in config")
	})

	t.Run("WithNonExistentConfigFile", func(t *testing.T) {
		provider := &ConfigFileProvider{
			configPath: "/nonexistent/path/config",
		}

		assert.False(t, provider.IsAvailable())

		_, err := provider.GetUsername()
		assert.Error(t, err)

		_, err = provider.GetPassword()
		assert.Error(t, err)
	})

	t.Run("WithInvalidJSONConfigFile", func(t *testing.T) {
		// Create temporary config file with invalid JSON
		tmpDir := t.TempDir()
		configDir := filepath.Join(tmpDir, ".idpbuilder")
		configPath := filepath.Join(configDir, "config")

		err := os.MkdirAll(configDir, 0755)
		require.NoError(t, err)

		err = os.WriteFile(configPath, []byte("invalid json content"), 0600)
		require.NoError(t, err)

		// Test provider
		provider := &ConfigFileProvider{
			configPath: configPath,
		}

		assert.True(t, provider.IsAvailable()) // File exists

		_, err = provider.GetUsername()
		assert.Error(t, err)

		_, err = provider.GetPassword()
		assert.Error(t, err)
	})

	t.Run("ConfigCaching", func(t *testing.T) {
		// Create temporary config file
		tmpDir := t.TempDir()
		configDir := filepath.Join(tmpDir, ".idpbuilder")
		configPath := filepath.Join(configDir, "config")

		err := os.MkdirAll(configDir, 0755)
		require.NoError(t, err)

		config := Config{
			Registries: map[string]RegistryCredentials{
				"gitea": {
					Username: "cacheduser",
					Password: "cachedpass",
					URL:      "https://gitea.example.com",
				},
			},
		}

		data, err := json.Marshal(config)
		require.NoError(t, err)

		err = os.WriteFile(configPath, data, 0600)
		require.NoError(t, err)

		// Test provider
		provider := &ConfigFileProvider{
			configPath: configPath,
		}

		// First call should load and cache config
		username1, err := provider.GetUsername()
		require.NoError(t, err)
		assert.Equal(t, "cacheduser", username1)

		// Modify the file to ensure caching is working
		registry := config.Registries["gitea"]
		registry.Username = "modifieduser"
		config.Registries["gitea"] = registry
		data, err = json.Marshal(config)
		require.NoError(t, err)
		err = os.WriteFile(configPath, data, 0600)
		require.NoError(t, err)

		// Second call should return cached value
		username2, err := provider.GetUsername()
		require.NoError(t, err)
		assert.Equal(t, "cacheduser", username2) // Should still be cached value

		// Create new provider to test fresh load
		provider2 := &ConfigFileProvider{
			configPath: configPath,
		}
		username3, err := provider2.GetUsername()
		require.NoError(t, err)
		assert.Equal(t, "modifieduser", username3) // Should be new value
	})
}

func TestNewConfigFileProvider(t *testing.T) {
	provider := NewConfigFileProvider()

	// Should set up config path in user's home directory
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	expectedPath := filepath.Join(homeDir, ".idpbuilder", "config")
	assert.Equal(t, expectedPath, provider.configPath)
}