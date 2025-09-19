package gitea

import (
	"os"
	"testing"

	"github.com/cnoe-io/idpbuilder/pkg/certs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	registryURL := "https://gitea.example.com"
	certManager := certs.NewTrustStore()

	client, err := NewClient(registryURL, certManager)

	// Since we're using placeholder credentials, we expect this to succeed in initialization
	// The actual registry connection would fail, but client creation should work
	require.NotNil(t, client)
	require.NoError(t, err)

	// Verify client configuration
	assert.Equal(t, registryURL, client.config.URL)
	// With new credential management, defaults are empty when no credentials configured
	assert.Equal(t, "", client.config.Username) // From getRegistryUsername()
	assert.Equal(t, "", client.config.Token)    // From getRegistryPassword()
	assert.False(t, client.config.Insecure)
}

func TestNewInsecureClient(t *testing.T) {
	registryURL := "http://gitea.example.com"

	client, err := NewInsecureClient(registryURL)

	// Client creation should succeed
	require.NotNil(t, client)
	require.NoError(t, err)

	// Verify client configuration
	assert.Equal(t, registryURL, client.config.URL)
	// With new credential management, defaults are empty when no credentials configured
	assert.Equal(t, "", client.config.Username)
	assert.Equal(t, "", client.config.Token)
	assert.True(t, client.config.Insecure)
}

func TestGetRegistryUsername(t *testing.T) {
	// Test default username (empty when no credentials configured)
	username := getRegistryUsername()
	assert.Equal(t, "", username)
}

func TestGetRegistryPassword(t *testing.T) {
	// Test default password (empty when no credentials configured)
	password := getRegistryPassword()
	assert.Equal(t, "", password)
}

// Note: TestGetImageContentForReference removed as the method was replaced
// with real Docker daemon integration in the new implementation

func TestPushProgressStruct(t *testing.T) {
	// Test PushProgress struct
	progress := PushProgress{
		CurrentLayer: 3,
		TotalLayers:  10,
		Percentage:   30,
	}

	assert.Equal(t, 3, progress.CurrentLayer)
	assert.Equal(t, 10, progress.TotalLayers)
	assert.Equal(t, 30, progress.Percentage)
}

func TestClientConfiguration(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		insecure bool
	}{
		{
			name:     "secure client",
			url:      "https://gitea.secure.com",
			insecure: false,
		},
		{
			name:     "insecure client",
			url:      "http://gitea.insecure.com",
			insecure: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var client *Client
			var err error

			if tt.insecure {
				client, err = NewInsecureClient(tt.url)
			} else {
				certManager := certs.NewTrustStore()
				client, err = NewClient(tt.url, certManager)
			}

			require.NoError(t, err)
			require.NotNil(t, client)
			assert.Equal(t, tt.url, client.config.URL)
			assert.Equal(t, tt.insecure, client.config.Insecure)
		})
	}
}

// Test environment variable integration
func TestClientWithEnvironmentVariables(t *testing.T) {
	// Save original environment - use correct variable names
	originalUsername := os.Getenv("GITEA_USERNAME")
	originalPassword := os.Getenv("GITEA_PASSWORD")

	// Clean up after test
	defer func() {
		if originalUsername != "" {
			os.Setenv("GITEA_USERNAME", originalUsername)
		} else {
			os.Unsetenv("GITEA_USERNAME")
		}
		if originalPassword != "" {
			os.Setenv("GITEA_PASSWORD", originalPassword)
		} else {
			os.Unsetenv("GITEA_PASSWORD")
		}
	}()

	// Set environment variables with correct names
	os.Setenv("GITEA_USERNAME", "envuser")
	os.Setenv("GITEA_PASSWORD", "envpass")

	// Test that environment variables are now properly read
	username := getRegistryUsername()
	password := getRegistryPassword()

	assert.Equal(t, "envuser", username)
	assert.Equal(t, "envpass", password)
}
