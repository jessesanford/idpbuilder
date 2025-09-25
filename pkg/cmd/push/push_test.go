package push

import (
	"context"
	"strings"
	"testing"

	"github.com/cnoe-io/idpbuilder/pkg/cmd/helpers"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPushCmd(t *testing.T) {
	t.Run("command exists", func(t *testing.T) {
		assert.NotNil(t, PushCmd)
		assert.Equal(t, "push IMAGE", PushCmd.Use)
		assert.Equal(t, "Push a container image to an OCI registry", PushCmd.Short)
		assert.True(t, len(PushCmd.Long) > 0)
	})

	t.Run("requires exactly one argument", func(t *testing.T) {
		err := PushCmd.Args(PushCmd, []string{})
		assert.Error(t, err)

		err = PushCmd.Args(PushCmd, []string{"image1", "image2"})
		assert.Error(t, err)

		err = PushCmd.Args(PushCmd, []string{"localhost:5000/test:latest"})
		assert.NoError(t, err)
	})
}

func TestPushFlags(t *testing.T) {
	t.Run("has all expected flags", func(t *testing.T) {
		flags := PushCmd.Flags()

		// Check username flag
		usernameFlag := flags.Lookup("username")
		assert.NotNil(t, usernameFlag)
		assert.Equal(t, "u", usernameFlag.Shorthand)

		// Check password flag
		passwordFlag := flags.Lookup("password")
		assert.NotNil(t, passwordFlag)
		assert.Equal(t, "p", passwordFlag.Shorthand)

		// Check insecure flag
		insecureFlag := flags.Lookup("insecure")
		assert.NotNil(t, insecureFlag)
		assert.Equal(t, "false", insecureFlag.DefValue)

		// Check auth-mode flag
		authModeFlag := flags.Lookup("auth-mode")
		assert.NotNil(t, authModeFlag)
		assert.Equal(t, "auto", authModeFlag.DefValue)

		// Check tag flag
		tagFlag := flags.Lookup("tag")
		assert.NotNil(t, tagFlag)
		assert.Equal(t, "t", tagFlag.Shorthand)

		// Check quiet flag
		quietFlag := flags.Lookup("quiet")
		assert.NotNil(t, quietFlag)
		assert.Equal(t, "q", quietFlag.Shorthand)

		// Check force flag
		forceFlag := flags.Lookup("force")
		assert.NotNil(t, forceFlag)
	})
}

func TestExtractRegistry(t *testing.T) {
	tests := []struct {
		name        string
		imageRef    string
		expected    string
		expectError bool
	}{
		{
			name:     "docker hub simple name",
			imageRef: "ubuntu",
			expected: "docker.io",
		},
		{
			name:     "docker hub with tag",
			imageRef: "ubuntu:latest",
			expected: "docker.io",
		},
		{
			name:     "localhost registry",
			imageRef: "localhost:5000/myapp:v1.0",
			expected: "localhost:5000",
		},
		{
			name:     "private registry",
			imageRef: "registry.example.com/myapp:v1.0",
			expected: "registry.example.com",
		},
		{
			name:     "gcr registry",
			imageRef: "gcr.io/project/image:tag",
			expected: "gcr.io",
		},
		{
			name:     "docker hub with namespace",
			imageRef: "library/ubuntu:latest",
			expected: "docker.io",
		},
		{
			name:        "empty image reference",
			imageRef:    "",
			expected:    "docker.io",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := extractRegistry(tt.imageRef)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestReplaceTag(t *testing.T) {
	tests := []struct {
		name     string
		imageRef string
		newTag   string
		expected string
	}{
		{
			name:     "replace existing tag",
			imageRef: "localhost:5000/myapp:v1.0",
			newTag:   "latest",
			expected: "localhost:5000/myapp:latest",
		},
		{
			name:     "add tag to image without tag",
			imageRef: "localhost:5000/myapp",
			newTag:   "v2.0",
			expected: "localhost:5000/myapp:v2.0",
		},
		{
			name:     "replace tag in complex registry path",
			imageRef: "registry.example.com:443/namespace/app:old",
			newTag:   "new",
			expected: "registry.example.com:443/namespace/app:new",
		},
		{
			name:     "simple image name",
			imageRef: "ubuntu",
			newTag:   "22.04",
			expected: "ubuntu:22.04",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceTag(tt.imageRef, tt.newTag)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSetupAuthentication(t *testing.T) {
	ctx := context.Background()

	t.Run("flags mode without credentials", func(t *testing.T) {
		// Reset global variables
		authMode = "flags"
		username = ""
		password = ""

		auth, err := setupAuthentication(ctx, "localhost:5000")
		assert.Error(t, err)
		assert.Nil(t, auth)
		assert.Contains(t, err.Error(), "username and password required")
	})

	t.Run("flags mode with credentials", func(t *testing.T) {
		// Reset global variables
		authMode = "flags"
		username = "testuser"
		password = "testpass"

		auth, err := setupAuthentication(ctx, "localhost:5000")
		assert.NoError(t, err)
		assert.NotNil(t, auth)
	})

	t.Run("docker-config mode", func(t *testing.T) {
		authMode = "docker-config"

		auth, err := setupAuthentication(ctx, "localhost:5000")
		assert.NoError(t, err)
		assert.NotNil(t, auth)
	})

	t.Run("auto mode with flags", func(t *testing.T) {
		authMode = "auto"
		username = "testuser"
		password = "testpass"

		auth, err := setupAuthentication(ctx, "localhost:5000")
		assert.NoError(t, err)
		assert.NotNil(t, auth)
	})

	t.Run("auto mode without flags", func(t *testing.T) {
		authMode = "auto"
		username = ""
		password = ""

		auth, err := setupAuthentication(ctx, "localhost:5000")
		assert.NoError(t, err)
		assert.NotNil(t, auth)
	})

	t.Run("invalid auth mode", func(t *testing.T) {
		authMode = "invalid"

		auth, err := setupAuthentication(ctx, "localhost:5000")
		assert.Error(t, err)
		assert.Nil(t, auth)
		assert.Contains(t, err.Error(), "unsupported auth mode")
	})
}

func TestCreateAuthConfig(t *testing.T) {
	t.Run("with valid credentials", func(t *testing.T) {
		result := createAuthConfig("testuser", "testpass")
		assert.NotEmpty(t, result)

		// The result should be base64 encoded JSON
		assert.True(t, len(result) > 0)
	})

	t.Run("with empty credentials", func(t *testing.T) {
		result := createAuthConfig("", "")
		assert.Empty(t, result)
	})

	t.Run("with partial credentials", func(t *testing.T) {
		result1 := createAuthConfig("user", "")
		assert.NotEmpty(t, result1)

		result2 := createAuthConfig("", "pass")
		assert.NotEmpty(t, result2)
	})
}

func TestPrePushE(t *testing.T) {
	// Set a valid log level to avoid the SetLogger error
	helpers.LogLevel = "info"

	cmd := &cobra.Command{}
	args := []string{"localhost:5000/test:latest"}

	err := prePushE(cmd, args)
	assert.NoError(t, err)
	assert.Equal(t, "localhost:5000/test:latest", imageRef)
}

// Integration test helpers and mocks would go here for testing the full push flow
// For now, we focus on unit testing the individual functions

func TestPushImageValidation(t *testing.T) {
	// Test the main pushImage function with various invalid inputs
	// These would typically require mocking the Docker client

	t.Run("validates required image argument", func(t *testing.T) {
		// Reset global variables
		imageRef = ""

		cmd := &cobra.Command{}
		cmd.SetContext(context.Background())

		// This would normally fail at the image validation step
		// For a full test, we'd need to mock the Docker client
	})
}

func TestCommandUsageMessages(t *testing.T) {
	t.Run("usage messages are present", func(t *testing.T) {
		assert.NotEmpty(t, imageUsage)
		assert.NotEmpty(t, usernameUsage)
		assert.NotEmpty(t, passwordUsage)
		assert.NotEmpty(t, insecureUsage)
		assert.NotEmpty(t, authModeUsage)
		assert.NotEmpty(t, tagUsage)
		assert.NotEmpty(t, quietUsage)
		assert.NotEmpty(t, forceUsage)
	})

	t.Run("usage messages contain expected content", func(t *testing.T) {
		assert.Contains(t, strings.ToLower(imageUsage), "image")
		assert.Contains(t, strings.ToLower(usernameUsage), "username")
		assert.Contains(t, strings.ToLower(passwordUsage), "password")
		assert.Contains(t, strings.ToLower(authModeUsage), "authentication")
		assert.Contains(t, strings.ToLower(tagUsage), "tag")
	})
}

func TestCommandExamples(t *testing.T) {
	t.Run("long description contains examples", func(t *testing.T) {
		longDesc := PushCmd.Long
		assert.Contains(t, longDesc, "Examples:")
		assert.Contains(t, longDesc, "idpbuilder push")
		assert.Contains(t, longDesc, "--username")
		assert.Contains(t, longDesc, "--password")
		assert.Contains(t, longDesc, "localhost:5000")
	})
}

// Benchmark tests for performance-critical functions
func BenchmarkExtractRegistry(b *testing.B) {
	imageRef := "registry.example.com:443/namespace/app:v1.0"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = extractRegistry(imageRef)
	}
}

func BenchmarkReplaceTag(b *testing.B) {
	imageRef := "registry.example.com:443/namespace/app:v1.0"
	newTag := "latest"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = replaceTag(imageRef, newTag)
	}
}

func BenchmarkCreateAuthConfig(b *testing.B) {
	username := "testuser"
	password := "testpassword"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = createAuthConfig(username, password)
	}
}

// Table-driven test for command flag validation
func TestCommandFlagDefaults(t *testing.T) {
	tests := []struct {
		flagName     string
		expectedType string
		hasShorthand bool
		shorthand    string
	}{
		{"username", "string", true, "u"},
		{"password", "string", true, "p"},
		{"insecure", "bool", false, ""},
		{"auth-mode", "string", false, ""},
		{"tag", "stringSlice", true, "t"},
		{"quiet", "bool", true, "q"},
		{"force", "bool", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.flagName, func(t *testing.T) {
			flag := PushCmd.Flags().Lookup(tt.flagName)
			require.NotNil(t, flag, "Flag %s should exist", tt.flagName)

			if tt.hasShorthand {
				assert.Equal(t, tt.shorthand, flag.Shorthand, "Flag %s shorthand mismatch", tt.flagName)
			} else {
				assert.Empty(t, flag.Shorthand, "Flag %s should not have shorthand", tt.flagName)
			}
		})
	}
}
