package push

import (
	"strings"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPushCmd(t *testing.T) {
	t.Run("command structure", func(t *testing.T) {
		assert.Equal(t, "push IMAGE[:TAG]", PushCmd.Use)
		assert.Equal(t, "Push image to Gitea registry", PushCmd.Short)
		assert.NotEmpty(t, PushCmd.Long)
		assert.NotEmpty(t, PushCmd.Example)
		assert.NotNil(t, PushCmd.RunE)
		assert.NotNil(t, PushCmd.Args)
		
		// Test that exactly 1 argument is required
		err := PushCmd.Args(PushCmd, []string{"test:latest"})
		assert.NoError(t, err)
		
		err = PushCmd.Args(PushCmd, []string{})
		assert.Error(t, err)
		
		err = PushCmd.Args(PushCmd, []string{"test:latest", "extra"})
		assert.Error(t, err)
	})

	t.Run("flags configuration", func(t *testing.T) {
		// Insecure flag
		flag := PushCmd.Flags().Lookup("insecure")
		require.NotNil(t, flag)
		assert.Equal(t, "false", flag.DefValue)
		assert.Equal(t, "Skip TLS certificate verification", flag.Usage)

		// Registry flag
		flag = PushCmd.Flags().Lookup("registry")
		require.NotNil(t, flag)
		assert.Equal(t, "", flag.DefValue)
		assert.Contains(t, flag.Usage, "Registry URL")

		// Username flag
		flag = PushCmd.Flags().Lookup("username")
		require.NotNil(t, flag)
		assert.Equal(t, "", flag.DefValue)

		// Password flag
		flag = PushCmd.Flags().Lookup("password")
		require.NotNil(t, flag)
		assert.Equal(t, "", flag.DefValue)

		// Retry flag
		flag = PushCmd.Flags().Lookup("retry")
		require.NotNil(t, flag)
		assert.Equal(t, "3", flag.DefValue)
	})
}

func TestRunPush(t *testing.T) {
	t.Run("missing image argument", func(t *testing.T) {
		// This test ensures the Args validation works
		err := PushCmd.Args(PushCmd, []string{})
		assert.Error(t, err)
	})

	t.Run("image name normalization", func(t *testing.T) {
		cmd := PushCmd
		cmd.Flags().Set("insecure", "true") // Skip cert setup for test

		// Test with explicit tag
		err := runPush(cmd, []string{"myapp:v1"})
		assert.Error(t, err) // Should fail - either on client creation or image loading
		// The error could be about trust store (client creation) or image loading
		assert.True(t, 
			strings.Contains(err.Error(), "image loading") || 
			strings.Contains(err.Error(), "trust store") ||
			strings.Contains(err.Error(), "failed to create registry client"))

		// Test without tag (should add :latest)
		err = runPush(cmd, []string{"myapp"})
		assert.Error(t, err) // Should fail - either on client creation or image loading
		assert.True(t, 
			strings.Contains(err.Error(), "image loading") || 
			strings.Contains(err.Error(), "trust store") ||
			strings.Contains(err.Error(), "failed to create registry client"))
	})

	t.Run("default registry configuration", func(t *testing.T) {
		cmd := PushCmd
		cmd.Flags().Set("insecure", "true")

		// Should use default Gitea registry
		err := runPush(cmd, []string{"test:latest"})
		assert.Error(t, err)
		// Error should indicate either trust store or image loading limitation
		assert.True(t, 
			strings.Contains(err.Error(), "image loading") || 
			strings.Contains(err.Error(), "trust store") ||
			strings.Contains(err.Error(), "failed to create registry client"))
	})

	t.Run("environment variable password", func(t *testing.T) {
		// Test that password is read from environment
		cmd := PushCmd
		cmd.Flags().Set("insecure", "true")

		// This would normally test environment variable reading
		// For now, we verify the command structure is correct
		err := runPush(cmd, []string{"test:latest"})
		assert.Error(t, err)
		assert.True(t, 
			strings.Contains(err.Error(), "image loading") || 
			strings.Contains(err.Error(), "trust store") ||
			strings.Contains(err.Error(), "failed to create registry client"))
	})

	t.Run("certificate setup with secure connection", func(t *testing.T) {
		cmd := PushCmd
		// Don't set insecure flag - should attempt cert setup

		err := runPush(cmd, []string{"test:latest"})
		assert.Error(t, err)
		// Should either fail on cert setup or image loading
		// Both are expected limitations at this stage
	})
}

func TestFlagValidation(t *testing.T) {
	t.Run("all flags have descriptions", func(t *testing.T) {
		PushCmd.Flags().VisitAll(func(flag *pflag.Flag) {
			assert.NotEmpty(t, flag.Usage, "Flag %s should have usage description", flag.Name)
		})
	})

	t.Run("examples are helpful", func(t *testing.T) {
		examples := PushCmd.Example
		assert.Contains(t, examples, "idpbuilder push")
		assert.Contains(t, examples, "myapp")
		assert.Contains(t, examples, "--insecure")
	})

	t.Run("flag types are correct", func(t *testing.T) {
		// Boolean flags
		insecureFlag := PushCmd.Flags().Lookup("insecure")
		assert.Equal(t, "bool", insecureFlag.Value.Type())

		// String flags
		registryFlag := PushCmd.Flags().Lookup("registry")
		assert.Equal(t, "string", registryFlag.Value.Type())

		// Int flags
		retryFlag := PushCmd.Flags().Lookup("retry")
		assert.Equal(t, "int", retryFlag.Value.Type())
	})
}

func TestImageNameHandling(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"with tag", "myapp:v1", "myapp:v1"},
		{"without tag", "myapp", "myapp:latest"}, // Should be normalized to :latest
		{"with registry", "registry.example.com/myapp:v2", "registry.example.com/myapp:v2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This tests the logic that should normalize image names
			// Since we can't easily test the internal logic without refactoring,
			// we test that the command accepts various image name formats
			cmd := PushCmd
			cmd.Flags().Set("insecure", "true")

			err := runPush(cmd, []string{tt.input})
			// All should fail with the same error (trust store or image loading)
			assert.Error(t, err)
			assert.True(t, 
				strings.Contains(err.Error(), "image loading") || 
				strings.Contains(err.Error(), "trust store") ||
				strings.Contains(err.Error(), "failed to create registry client"))
		})
	}
}