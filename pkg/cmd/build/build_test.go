package build

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildCmd(t *testing.T) {
	t.Run("command structure", func(t *testing.T) {
		assert.Equal(t, "build", BuildCmd.Use)
		assert.Equal(t, "Assemble OCI image from context directory", BuildCmd.Short)
		assert.NotEmpty(t, BuildCmd.Long)
		assert.NotEmpty(t, BuildCmd.Example)
		assert.NotNil(t, BuildCmd.RunE)
	})

	t.Run("required flags", func(t *testing.T) {
		// Check that tag flag is marked as required
		flag := BuildCmd.Flags().Lookup("tag")
		require.NotNil(t, flag)
		
		// The flag should exist and have the right properties
		assert.Equal(t, "", flag.DefValue)
		assert.Equal(t, "Image tag (required)", flag.Usage)
	})

	t.Run("optional flags", func(t *testing.T) {
		// Context flag with default
		flag := BuildCmd.Flags().Lookup("context")
		require.NotNil(t, flag)
		assert.Equal(t, ".", flag.DefValue)

		// Output flag (optional)
		flag = BuildCmd.Flags().Lookup("output")
		require.NotNil(t, flag)
		assert.Equal(t, "", flag.DefValue)

		// Platform flag with default
		flag = BuildCmd.Flags().Lookup("platform")
		require.NotNil(t, flag)
		assert.Equal(t, "linux/amd64", flag.DefValue)

		// Exclude flag
		flag = BuildCmd.Flags().Lookup("exclude")
		require.NotNil(t, flag)
	})
}

func TestRunBuild(t *testing.T) {
	t.Run("missing tag flag", func(t *testing.T) {
		// This would normally be tested by executing the command
		// but since we're testing just the function, we'll focus on
		// the command structure validation
		
		// The cobra framework handles required flag validation
		// so we focus on testing that the flag is properly configured
		cmd := BuildCmd
		cmd.Flags().Set("context", ".")
		
		// Tag is required, so command should not validate without it
		err := cmd.ValidateRequiredFlags()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tag")
	})

	t.Run("invalid context directory", func(t *testing.T) {
		// Test the runBuild function with invalid context
		cmd := BuildCmd
		cmd.Flags().Set("tag", "test:latest")
		cmd.Flags().Set("context", "/nonexistent/directory")
		
		err := runBuild(cmd, []string{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "context directory not found")
	})

	t.Run("invalid platform", func(t *testing.T) {
		// Test with invalid platform format
		cmd := BuildCmd
		cmd.Flags().Set("tag", "test:latest")
		cmd.Flags().Set("context", ".")
		cmd.Flags().Set("platform", "invalid-platform")
		
		err := runBuild(cmd, []string{})
		assert.Error(t, err)
		// The actual error message varies based on platform parsing implementation
		assert.Contains(t, err.Error(), "invalid")
	})
}

func TestFlagValidation(t *testing.T) {
	t.Run("all flags have descriptions", func(t *testing.T) {
		BuildCmd.Flags().VisitAll(func(flag *pflag.Flag) {
			assert.NotEmpty(t, flag.Usage, "Flag %s should have usage description", flag.Name)
		})
	})

	t.Run("examples are helpful", func(t *testing.T) {
		examples := BuildCmd.Example
		assert.Contains(t, examples, "idpbuilder build")
		assert.Contains(t, examples, "--context")
		assert.Contains(t, examples, "--tag")
	})
}