package gitea

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockKeyringProvider for testing keyring functionality without system keyring
type MockKeyringProvider struct {
	username  string
	password  string
	available bool
	shouldErr bool
}

func NewMockKeyringProvider(username, password string, available bool) *MockKeyringProvider {
	return &MockKeyringProvider{
		username:  username,
		password:  password,
		available: available,
		shouldErr: false,
	}
}

func (m *MockKeyringProvider) GetUsername() (string, error) {
	if m.shouldErr {
		return "", fmt.Errorf("mock keyring error")
	}
	if m.username == "" {
		return "", fmt.Errorf("no username in keyring")
	}
	return m.username, nil
}

func (m *MockKeyringProvider) GetPassword() (string, error) {
	if m.shouldErr {
		return "", fmt.Errorf("mock keyring error")
	}
	if m.password == "" {
		return "", fmt.Errorf("no password in keyring")
	}
	return m.password, nil
}

func (m *MockKeyringProvider) IsAvailable() bool {
	return m.available
}

func (m *MockKeyringProvider) Priority() int {
	return 4
}

func (m *MockKeyringProvider) SetShouldError(shouldErr bool) {
	m.shouldErr = shouldErr
}

func TestMockKeyringProvider(t *testing.T) {
	t.Run("WithCredentials", func(t *testing.T) {
		provider := NewMockKeyringProvider("keyringuser", "keyringpass", true)

		assert.True(t, provider.IsAvailable())
		assert.Equal(t, 4, provider.Priority())

		username, err := provider.GetUsername()
		assert.NoError(t, err)
		assert.Equal(t, "keyringuser", username)

		password, err := provider.GetPassword()
		assert.NoError(t, err)
		assert.Equal(t, "keyringpass", password)
	})

	t.Run("WithoutCredentials", func(t *testing.T) {
		provider := NewMockKeyringProvider("", "", false)

		assert.False(t, provider.IsAvailable())

		_, err := provider.GetUsername()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no username in keyring")

		_, err = provider.GetPassword()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no password in keyring")
	})

	t.Run("WithErrors", func(t *testing.T) {
		provider := NewMockKeyringProvider("user", "pass", true)
		provider.SetShouldError(true)

		_, err := provider.GetUsername()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "mock keyring error")

		_, err = provider.GetPassword()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "mock keyring error")
	})
}

func TestKeyringProvider(t *testing.T) {
	// These tests focus on the structure and constants of the real KeyringProvider
	// We can't test actual keyring operations in a unit test environment

	t.Run("ProviderCreation", func(t *testing.T) {
		provider := NewKeyringProvider()

		assert.NotNil(t, provider)
		assert.Equal(t, keyringService, provider.service)
		assert.Equal(t, keyringUser, provider.user)
		assert.Equal(t, 4, provider.Priority())
	})

	t.Run("Constants", func(t *testing.T) {
		assert.Equal(t, "idpbuilder", keyringService)
		assert.Equal(t, "gitea", keyringUser)
	})
}

// TestCredentialManagerWithMockKeyring tests the credential manager using mock keyring
func TestCredentialManagerWithMockKeyring(t *testing.T) {
	t.Run("KeyringAsLowestPriority", func(t *testing.T) {
		// Create a credential manager with mock keyring as the last provider
		mockKeyring := NewMockKeyringProvider("keyringuser", "keyringpass", true)

		manager := &CredentialManager{
			providers: []CredentialProvider{
				NewCLICredentialProvider(),
				NewEnvCredentialProvider(),
				NewConfigFileProvider(),
				mockKeyring,
			},
		}

		// With no other credentials, should fall back to keyring
		username, password, err := manager.GetCredentials()
		assert.NoError(t, err)
		assert.Equal(t, "keyringuser", username)
		assert.Equal(t, "keyringpass", password)
	})

	t.Run("KeyringNotUsedWhenHigherPriorityAvailable", func(t *testing.T) {
		// Set up environment variables
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

		os.Setenv("GITEA_USERNAME", "envuser")
		os.Setenv("GITEA_PASSWORD", "envpass")

		mockKeyring := NewMockKeyringProvider("keyringuser", "keyringpass", true)

		manager := &CredentialManager{
			providers: []CredentialProvider{
				NewCLICredentialProvider(),
				NewEnvCredentialProvider(),
				NewConfigFileProvider(),
				mockKeyring,
			},
		}

		// Should use environment variables, not keyring
		username, password, err := manager.GetCredentials()
		assert.NoError(t, err)
		assert.Equal(t, "envuser", username)
		assert.Equal(t, "envpass", password)
	})

	t.Run("KeyringErrorHandling", func(t *testing.T) {
		// Clear environment to force keyring usage
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

		os.Unsetenv("GITEA_USERNAME")
		os.Unsetenv("GITEA_PASSWORD")

		mockKeyring := NewMockKeyringProvider("keyringuser", "keyringpass", true)
		mockKeyring.SetShouldError(true)

		manager := &CredentialManager{
			providers: []CredentialProvider{
				NewCLICredentialProvider(),
				NewEnvCredentialProvider(),
				NewConfigFileProvider(),
				mockKeyring,
			},
		}

		// Should return error when keyring fails
		_, _, err := manager.GetCredentials()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no credentials available")
	})
}
