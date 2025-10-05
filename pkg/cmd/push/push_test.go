package push

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPushCommand(t *testing.T) {
	// Test that the command is properly configured
	assert.Equal(t, "push [IMAGE] [REGISTRY_URL]", PushCmd.Use)
	assert.Equal(t, "Push OCI artifacts to a registry", PushCmd.Short)
	assert.NotEmpty(t, PushCmd.Long)
	assert.NotNil(t, PushCmd.RunE)
}

func TestBuildPushOptions(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		envVars     map[string]string
		expectedImg string
		expectedReg string
		wantErr     bool
	}{
		{
			name:        "basic args",
			args:        []string{"myimage:latest", "registry.example.com"},
			expectedImg: "myimage:latest",
			expectedReg: "registry.example.com",
			wantErr:     false,
		},
		{
			name:        "single arg only",
			args:        []string{"myimage:latest"},
			expectedImg: "myimage:latest",
			expectedReg: "",
			wantErr:     false,
		},
		{
			name:        "environment variables",
			args:        []string{"myimage:latest", "registry.example.com"},
			envVars: map[string]string{
				"REGISTRY_USERNAME": "envuser",
				"REGISTRY_PASSWORD": "envpass",
				"REGISTRY_INSECURE": "true",
			},
			expectedImg: "myimage:latest",
			expectedReg: "registry.example.com",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear and set environment
			os.Clearenv()
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			// Reset package-level variables
			username = ""
			password = ""
			insecure = false
			dryRun = false
			verbose = false

			// Build options
			cmd := &cobra.Command{}
			opts, err := buildPushOptions(cmd, tt.args)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedImg, opts.ImageRef)
			assert.Equal(t, tt.expectedReg, opts.RegistryURL)

			// Test environment variable handling
			if tt.envVars != nil {
				if user, ok := tt.envVars["REGISTRY_USERNAME"]; ok {
					assert.Equal(t, user, opts.Username)
				}
				if pass, ok := tt.envVars["REGISTRY_PASSWORD"]; ok {
					assert.Equal(t, pass, opts.Password)
				}
				if insec := tt.envVars["REGISTRY_INSECURE"]; insec == "true" {
					assert.True(t, opts.Insecure)
				}
			}
		})
	}
}

func TestValidateRegistryURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid https URL",
			url:     "https://registry.example.com",
			wantErr: false,
		},
		{
			name:    "valid http URL",
			url:     "http://localhost:5000",
			wantErr: false,
		},
		{
			name:    "URL without scheme",
			url:     "registry.example.com",
			wantErr: false, // Should add https://
		},
		{
			name:    "invalid scheme",
			url:     "ftp://registry.example.com",
			wantErr: true,
		},
		{
			name:    "empty URL",
			url:     "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRegistryURL(tt.url)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePushOptions(t *testing.T) {
	tests := []struct {
		name    string
		opts    *PushOptions
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid options with auth",
			opts: &PushOptions{
				ImageRef:    "myimage:latest",
				RegistryURL: "registry.example.com",
				Username:    "user",
				Password:    "pass",
			},
			wantErr: false,
		},
		{
			name: "missing image reference",
			opts: &PushOptions{
				RegistryURL: "registry.example.com",
			},
			wantErr: true,
			errMsg:  "image reference is required",
		},
		{
			name: "missing registry URL",
			opts: &PushOptions{
				ImageRef: "myimage:latest",
			},
			wantErr: true,
			errMsg:  "registry URL is required",
		},
		{
			name: "localhost doesn't require auth",
			opts: &PushOptions{
				ImageRef:    "myimage:latest",
				RegistryURL: "localhost:5000",
			},
			wantErr: false,
		},
		{
			name: "remote registry requires auth",
			opts: &PushOptions{
				ImageRef:    "myimage:latest",
				RegistryURL: "registry.example.com",
			},
			wantErr: true,
			errMsg:  "authentication required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePushOptions(tt.opts)
			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetStringValue(t *testing.T) {
	tests := []struct {
		name      string
		flagValue string
		envValue  string
		expected  string
	}{
		{
			name:      "flag takes precedence",
			flagValue: "flag_value",
			envValue:  "env_value",
			expected:  "flag_value",
		},
		{
			name:      "env used when flag empty",
			flagValue: "",
			envValue:  "env_value",
			expected:  "env_value",
		},
		{
			name:      "both empty",
			flagValue: "",
			envValue:  "",
			expected:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getStringValue(tt.flagValue, tt.envValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRequiresAuth(t *testing.T) {
	tests := []struct {
		name        string
		registryURL string
		expected    bool
	}{
		{
			name:        "localhost does not require auth",
			registryURL: "localhost:5000",
			expected:    false,
		},
		{
			name:        "127.0.0.1 does not require auth",
			registryURL: "127.0.0.1:5000",
			expected:    false,
		},
		{
			name:        "remote registry requires auth",
			registryURL: "registry.example.com",
			expected:    true,
		},
		{
			name:        "docker hub requires auth",
			registryURL: "docker.io",
			expected:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := requiresAuth(tt.registryURL)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateImageRef(t *testing.T) {
	tests := []struct {
		name     string
		imageRef string
		wantErr  bool
	}{
		{
			name:     "valid image ref",
			imageRef: "myimage:latest",
			wantErr:  false,
		},
		{
			name:     "image ref with namespace",
			imageRef: "registry.com/namespace/image:tag",
			wantErr:  false,
		},
		{
			name:     "empty image ref",
			imageRef: "",
			wantErr:  true,
		},
		{
			name:     "image ref with spaces",
			imageRef: "my image:latest",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateImageRef(tt.imageRef)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}