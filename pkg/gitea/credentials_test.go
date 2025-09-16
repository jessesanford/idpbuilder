package gitea

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvCredentialProvider(t *testing.T) {
	// Save and restore environment
	oldUsername := os.Getenv("GITEA_USERNAME")
	oldPassword := os.Getenv("GITEA_PASSWORD")
	defer func() {
		if oldUsername != "" {
			os.Setenv("GITEA_USERNAME", oldUsername)
		} else {
			os.Unsetenv("GITEA_USERNAME")
		}
		if oldPassword != "" {
			os.Setenv("GITEA_PASSWORD", oldPassword)
		} else {
			os.Unsetenv("GITEA_PASSWORD")
		}
	}()

	provider := NewEnvCredentialProvider()

	t.Run("WithEnvironmentVariablesSet", func(t *testing.T) {
		// Test with environment variables set
		os.Setenv("GITEA_USERNAME", "testuser")
		os.Setenv("GITEA_PASSWORD", "testpass")

		assert.True(t, provider.IsAvailable())
		assert.Equal(t, 2, provider.Priority())

		username, err := provider.GetUsername()
		require.NoError(t, err)
		assert.Equal(t, "testuser", username)

		password, err := provider.GetPassword()
		require.NoError(t, err)
		assert.Equal(t, "testpass", password)
	})

	t.Run("WithEnvironmentVariablesNotSet", func(t *testing.T) {
		// Test with environment variables not set
		os.Unsetenv("GITEA_USERNAME")
		os.Unsetenv("GITEA_PASSWORD")

		assert.False(t, provider.IsAvailable())

		_, err := provider.GetUsername()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "GITEA_USERNAME not set")

		_, err = provider.GetPassword()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "GITEA_PASSWORD not set")
	})

	t.Run("WithPartialEnvironmentVariables", func(t *testing.T) {
		// Test with only username set
		os.Setenv("GITEA_USERNAME", "testuser")
		os.Unsetenv("GITEA_PASSWORD")

		assert.False(t, provider.IsAvailable())

		// Username should work
		username, err := provider.GetUsername()
		require.NoError(t, err)
		assert.Equal(t, "testuser", username)

		// Password should fail
		_, err = provider.GetPassword()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "GITEA_PASSWORD not set")
	})
}

func TestCLICredentialProvider(t *testing.T) {
	provider := NewCLICredentialProvider()

	t.Run("InitialState", func(t *testing.T) {
		assert.False(t, provider.IsAvailable())
		assert.Equal(t, 1, provider.Priority())

		_, err := provider.GetUsername()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no CLI username provided")

		_, err = provider.GetPassword()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no CLI password/token provided")
	})

	t.Run("WithCredentialsSet", func(t *testing.T) {
		provider.SetCredentials("cliuser", "clitoken")

		assert.True(t, provider.IsAvailable())

		username, err := provider.GetUsername()
		require.NoError(t, err)
		assert.Equal(t, "cliuser", username)

		password, err := provider.GetPassword()
		require.NoError(t, err)
		assert.Equal(t, "clitoken", password)
	})

	t.Run("WithPartialCredentials", func(t *testing.T) {
		provider.SetCredentials("onlyuser", "")
		assert.False(t, provider.IsAvailable())

		provider.SetCredentials("", "onlypass")
		assert.False(t, provider.IsAvailable())
	})
}

func TestCredentialManager(t *testing.T) {
	// Clean environment for testing
	oldUsername := os.Getenv("GITEA_USERNAME")
	oldPassword := os.Getenv("GITEA_PASSWORD")
	defer func() {
		if oldUsername != "" {
			os.Setenv("GITEA_USERNAME", oldUsername)
		} else {
			os.Unsetenv("GITEA_USERNAME")
		}
		if oldPassword != "" {
			os.Setenv("GITEA_PASSWORD", oldPassword)
		} else {
			os.Unsetenv("GITEA_PASSWORD")
		}
	}()

	t.Run("PriorityOrder", func(t *testing.T) {
		manager := NewCredentialManager()

		// Set environment variables as fallback
		os.Setenv("GITEA_USERNAME", "envuser")
		os.Setenv("GITEA_PASSWORD", "envpass")

		// Without CLI credentials, should use env
		username, password, err := manager.GetCredentials()
		require.NoError(t, err)
		assert.Equal(t, "envuser", username)
		assert.Equal(t, "envpass", password)

		// With CLI credentials, should use CLI (higher priority)
		manager.SetCLICredentials("cliuser", "clipass")
		username, password, err = manager.GetCredentials()
		require.NoError(t, err)
		assert.Equal(t, "cliuser", username)
		assert.Equal(t, "clipass", password)
	})

	t.Run("BackwardCompatibilityMethods", func(t *testing.T) {
		manager := NewCredentialManager()
		os.Setenv("GITEA_USERNAME", "testuser")
		os.Setenv("GITEA_PASSWORD", "testpass")

		// Test backward compatibility methods
		username := manager.GetUsername()
		assert.Equal(t, "testuser", username)

		password := manager.GetPassword()
		assert.Equal(t, "testpass", password)
	})

	t.Run("NoCredentialsAvailable", func(t *testing.T) {
		manager := NewCredentialManager()

		// Clear environment
		os.Unsetenv("GITEA_USERNAME")
		os.Unsetenv("GITEA_PASSWORD")

		// Should return empty strings for backward compatibility
		username := manager.GetUsername()
		assert.Equal(t, "", username)

		password := manager.GetPassword()
		assert.Equal(t, "", password)

		// GetCredentials should return error
		_, _, err := manager.GetCredentials()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no credentials available")
	})

	t.Run("SetCLICredentials", func(t *testing.T) {
		manager := NewCredentialManager()

		// Initially no credentials
		assert.False(t, manager.providers[0].IsAvailable())

		// Set CLI credentials
		manager.SetCLICredentials("testuser", "testtoken")

		// CLI provider should now be available
		assert.True(t, manager.providers[0].IsAvailable())

		username, password, err := manager.GetCredentials()
		require.NoError(t, err)
		assert.Equal(t, "testuser", username)
		assert.Equal(t, "testtoken", password)
	})
}
