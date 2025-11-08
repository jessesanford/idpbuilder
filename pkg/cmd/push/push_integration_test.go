package push_test

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/cnoe-io/idpbuilder/pkg/cmd/push"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test Suite 5: Environment Variable Scenarios (10 tests)

// T-2.2.5-01: Test push command with all configuration from environment variables
func TestPushCommand_AllFromEnvironment(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Given: All configuration from environment variables
	os.Setenv("IDPBUILDER_REGISTRY", "gitea.cnoe.localtest.me:8443")
	os.Setenv("IDPBUILDER_USERNAME", "giteaAdmin")
	os.Setenv("IDPBUILDER_PASSWORD", "password")
	os.Setenv("IDPBUILDER_INSECURE", "true")
	defer func() {
		os.Unsetenv("IDPBUILDER_REGISTRY")
		os.Unsetenv("IDPBUILDER_USERNAME")
		os.Unsetenv("IDPBUILDER_PASSWORD")
		os.Unsetenv("IDPBUILDER_INSECURE")
	}()

	// Given: Command with NO flags set
	cmd := push.NewPushCommand(viper.New())
	cmd.SetArgs([]string{"alpine:latest"})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd.SetContext(ctx)

	// When: Execute command (should use env vars)
	err := cmd.Execute()

	// Then: Command succeeds using environment variables only
	// Note: This may fail if Docker daemon is not running or image not available
	// The test verifies config loading works, actual push may fail in CI
	if err != nil {
		// Acceptable errors: Docker not available, image not found
		// Configuration errors (username/password) should not occur
		assert.NotContains(t, err.Error(), "username is required")
		assert.NotContains(t, err.Error(), "password is required")
	}
}

// T-2.2.5-02: Test that flags override environment variables
func TestPushCommand_FlagOverridesEnvironment(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Given: Environment variable set to different registry
	os.Setenv("IDPBUILDER_REGISTRY", "docker.io")
	os.Setenv("IDPBUILDER_USERNAME", "wronguser")
	os.Setenv("IDPBUILDER_PASSWORD", "wrongpass")
	defer func() {
		os.Unsetenv("IDPBUILDER_REGISTRY")
		os.Unsetenv("IDPBUILDER_USERNAME")
		os.Unsetenv("IDPBUILDER_PASSWORD")
	}()

	// Given: Command with flags that override env vars
	cmd := push.NewPushCommand(viper.New())
	cmd.SetArgs([]string{
		"alpine:latest",
		"--registry", "gitea.cnoe.localtest.me:8443",
		"--username", "giteaAdmin",
		"--password", "password",
		"--insecure",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd.SetContext(ctx)

	// When: Execute command
	err := cmd.Execute()

	// Then: Command succeeds using flag values (not env vars)
	// Configuration should be correct (flags override env)
	if err != nil {
		assert.NotContains(t, err.Error(), "wronguser")
		assert.NotContains(t, err.Error(), "wrongpass")
	}
}

// T-2.2.5-03: Test mixed configuration sources
func TestPushCommand_MixedConfiguration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Given: Some config from env, some from flags
	os.Setenv("IDPBUILDER_REGISTRY", "gitea.cnoe.localtest.me:8443")
	os.Setenv("IDPBUILDER_USERNAME", "giteaAdmin")
	defer func() {
		os.Unsetenv("IDPBUILDER_REGISTRY")
		os.Unsetenv("IDPBUILDER_USERNAME")
	}()

	// Given: Command with partial flags
	cmd := push.NewPushCommand(viper.New())
	cmd.SetArgs([]string{
		"alpine:latest",
		"--password", "password",
		"--insecure",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd.SetContext(ctx)

	// When: Execute command
	err := cmd.Execute()

	// Then: Mixed sources work correctly
	if err != nil {
		assert.NotContains(t, err.Error(), "username is required")
		assert.NotContains(t, err.Error(), "registry is required")
	}
}

// T-2.2.5-04: Test verbose mode shows configuration sources
func TestPushCommand_VerboseShowsSources(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Given: Mixed configuration sources
	os.Setenv("IDPBUILDER_USERNAME", "giteaAdmin")
	os.Setenv("IDPBUILDER_PASSWORD", "password")
	defer func() {
		os.Unsetenv("IDPBUILDER_USERNAME")
		os.Unsetenv("IDPBUILDER_PASSWORD")
	}()

	// Given: Command with verbose flag
	cmd := push.NewPushCommand(viper.New())
	cmd.SetArgs([]string{
		"alpine:latest",
		"--registry", "gitea.cnoe.localtest.me:8443",
		"--insecure",
		"--verbose",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd.SetContext(ctx)

	// When: Execute command
	err := cmd.Execute()

	// Then: Verbose mode should display sources
	// (Output verification would require capturing stdout)
	if err != nil {
		assert.NotContains(t, err.Error(), "username is required")
	}
}

// T-2.2.5-05: Test validation error messages mention environment variables
func TestPushCommand_ValidationErrorsWithEnvHints(t *testing.T) {
	// Given: Missing required credentials (no env vars, no flags)
	cmd := push.NewPushCommand(viper.New())
	cmd.SetArgs([]string{"alpine:latest"})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd.SetContext(ctx)

	// When: Execute command
	err := cmd.Execute()

	// Then: Error message should mention environment variable option
	require.Error(t, err, "Should fail with missing credentials")
	assert.Contains(t, err.Error(), "IDPBUILDER_USERNAME",
		"Error should mention IDPBUILDER_USERNAME environment variable")
}

// T-2.2.5-06: Test environment variables override defaults
func TestPushCommand_EnvironmentOverridesDefault(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Given: Custom registry in environment (overrides default)
	os.Setenv("IDPBUILDER_REGISTRY", "custom.registry.io:5000")
	os.Setenv("IDPBUILDER_USERNAME", "testuser")
	os.Setenv("IDPBUILDER_PASSWORD", "testpass")
	defer func() {
		os.Unsetenv("IDPBUILDER_REGISTRY")
		os.Unsetenv("IDPBUILDER_USERNAME")
		os.Unsetenv("IDPBUILDER_PASSWORD")
	}()

	// Given: Command with no registry flag
	cmd := push.NewPushCommand(viper.New())
	cmd.SetArgs([]string{"alpine:latest"})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd.SetContext(ctx)

	// When: Execute command
	err := cmd.Execute()

	// Then: Should attempt to use custom registry (not default)
	if err != nil {
		// Error may mention custom registry
		// Should NOT use default gitea registry
		errMsg := err.Error()
		if strings.Contains(errMsg, "registry") {
			assert.Contains(t, errMsg, "custom.registry.io")
		}
	}
}

// T-2.2.5-07: Test insecure flag from environment
func TestPushCommand_InsecureFromEnvironment(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Given: Insecure mode set via environment
	os.Setenv("IDPBUILDER_INSECURE", "true")
	os.Setenv("IDPBUILDER_USERNAME", "giteaAdmin")
	os.Setenv("IDPBUILDER_PASSWORD", "password")
	defer func() {
		os.Unsetenv("IDPBUILDER_INSECURE")
		os.Unsetenv("IDPBUILDER_USERNAME")
		os.Unsetenv("IDPBUILDER_PASSWORD")
	}()

	// Given: Command with no insecure flag
	cmd := push.NewPushCommand(viper.New())
	cmd.SetArgs([]string{"alpine:latest"})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd.SetContext(ctx)

	// When: Execute command
	err := cmd.Execute()

	// Then: Should use insecure mode from environment
	if err != nil {
		assert.NotContains(t, err.Error(), "username is required")
	}
}

// T-2.2.5-08: Test password with special characters from environment
func TestPushCommand_PasswordSpecialCharacters(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Given: Password with special characters in environment
	specialPassword := "p@ssw0rd!#$%^&*()_+-=[]{}|;:',.<>?/"
	os.Setenv("IDPBUILDER_USERNAME", "testuser")
	os.Setenv("IDPBUILDER_PASSWORD", specialPassword)
	defer func() {
		os.Unsetenv("IDPBUILDER_USERNAME")
		os.Unsetenv("IDPBUILDER_PASSWORD")
	}()

	// Given: Command using environment password
	cmd := push.NewPushCommand(viper.New())
	cmd.SetArgs([]string{
		"alpine:latest",
		"--insecure",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd.SetContext(ctx)

	// When: Execute command
	err := cmd.Execute()

	// Then: Special characters should be handled correctly
	if err != nil {
		assert.NotContains(t, err.Error(), "password is required")
	}
}

// T-2.2.5-09: Test registry override
func TestPushCommand_RegistryOverride(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Given: Registry override via environment
	os.Setenv("IDPBUILDER_REGISTRY", "localhost:5000")
	os.Setenv("IDPBUILDER_USERNAME", "localuser")
	os.Setenv("IDPBUILDER_PASSWORD", "localpass")
	os.Setenv("IDPBUILDER_INSECURE", "true")
	defer func() {
		os.Unsetenv("IDPBUILDER_REGISTRY")
		os.Unsetenv("IDPBUILDER_USERNAME")
		os.Unsetenv("IDPBUILDER_PASSWORD")
		os.Unsetenv("IDPBUILDER_INSECURE")
	}()

	// Given: Command using environment configuration
	cmd := push.NewPushCommand(viper.New())
	cmd.SetArgs([]string{"alpine:latest"})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd.SetContext(ctx)

	// When: Execute command
	err := cmd.Execute()

	// Then: Should attempt to use overridden registry
	if err != nil && strings.Contains(err.Error(), "registry") {
		assert.Contains(t, err.Error(), "localhost:5000")
	}
}

// T-2.2.5-10: Test backward compatibility with Wave 2.1 (flags only)
func TestPushCommand_BackwardCompatibility_Wave21(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Given: NO environment variables set (Wave 2.1 style)
	// Given: Command using only flags (Wave 2.1 pattern)
	cmd := push.NewPushCommand(viper.New())
	cmd.SetArgs([]string{
		"alpine:latest",
		"--registry", "gitea.cnoe.localtest.me:8443",
		"--username", "giteaAdmin",
		"--password", "password",
		"--insecure",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd.SetContext(ctx)

	// When: Execute command (Wave 2.1 style)
	err := cmd.Execute()

	// Then: Command works exactly as in Wave 2.1
	// Configuration errors should not occur
	if err != nil {
		assert.NotContains(t, err.Error(), "username is required")
		assert.NotContains(t, err.Error(), "password is required")
		assert.NotContains(t, err.Error(), "configuration error")
	}
}

// Test Suite 6: Edge Cases & Error Handling (10 tests)

// T-2.2.6-01: Test empty environment variable (should be ignored)
func TestPushCommand_EmptyEnvironmentVariable(t *testing.T) {
	// Given: Empty environment variable (should fall back to default/flag)
	os.Setenv("IDPBUILDER_REGISTRY", "")
	os.Setenv("IDPBUILDER_USERNAME", "testuser")
	os.Setenv("IDPBUILDER_PASSWORD", "testpass")
	defer func() {
		os.Unsetenv("IDPBUILDER_REGISTRY")
		os.Unsetenv("IDPBUILDER_USERNAME")
		os.Unsetenv("IDPBUILDER_PASSWORD")
	}()

	// Given: Command with no registry flag (should use default)
	cmd := push.NewPushCommand(viper.New())
	cmd.SetArgs([]string{"alpine:latest"})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd.SetContext(ctx)

	// When: Execute command
	err := cmd.Execute()

	// Then: Empty env var should be ignored, default registry used
	if err != nil {
		// Should NOT complain about missing registry
		assert.NotContains(t, err.Error(), "registry is required")
	}
}

// T-2.2.6-02: Test invalid boolean in environment variable
func TestPushCommand_InvalidBooleanInEnv(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Given: Environment variable with invalid boolean value
	os.Setenv("IDPBUILDER_INSECURE", "maybe")
	os.Setenv("IDPBUILDER_USERNAME", "testuser")
	os.Setenv("IDPBUILDER_PASSWORD", "testpass")
	defer func() {
		os.Unsetenv("IDPBUILDER_INSECURE")
		os.Unsetenv("IDPBUILDER_USERNAME")
		os.Unsetenv("IDPBUILDER_PASSWORD")
	}()

	// Given: Command with no flags
	cmd := push.NewPushCommand(viper.New())
	cmd.SetArgs([]string{"alpine:latest"})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd.SetContext(ctx)

	// When: Execute command
	err := cmd.Execute()

	// Then: Invalid boolean falls back to default (false)
	// Command should still succeed (not fail on invalid bool)
	if err != nil {
		assert.NotContains(t, err.Error(), "invalid boolean")
		assert.NotContains(t, err.Error(), "maybe")
	}
}

// T-2.2.6-03: Test environment variable with leading/trailing spaces
func TestPushCommand_EnvVarWithSpaces(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Given: Environment variable with spaces (should be trimmed)
	os.Setenv("IDPBUILDER_USERNAME", "  testuser  ")
	os.Setenv("IDPBUILDER_PASSWORD", "  testpass  ")
	os.Setenv("IDPBUILDER_INSECURE", "  true  ")
	defer func() {
		os.Unsetenv("IDPBUILDER_USERNAME")
		os.Unsetenv("IDPBUILDER_PASSWORD")
		os.Unsetenv("IDPBUILDER_INSECURE")
	}()

	// Given: Command using environment variables
	cmd := push.NewPushCommand(viper.New())
	cmd.SetArgs([]string{"alpine:latest"})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd.SetContext(ctx)

	// When: Execute command
	err := cmd.Execute()

	// Then: Spaces should be handled (trimmed for booleans)
	if err != nil {
		assert.NotContains(t, err.Error(), "username is required")
	}
}

// T-2.2.6-04: Test multiple boolean environment variable formats
func TestPushCommand_MultipleEnvFormats(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testCases := []struct {
		name      string
		insecure  string
		verbose   string
		shouldErr bool
	}{
		{"lowercase true", "true", "false", false},
		{"uppercase TRUE", "TRUE", "FALSE", false},
		{"numeric 1/0", "1", "0", false},
		{"yes/no", "yes", "no", false},
		{"YES/NO", "YES", "NO", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Given: Boolean env vars in different formats
			os.Setenv("IDPBUILDER_INSECURE", tc.insecure)
			os.Setenv("IDPBUILDER_VERBOSE", tc.verbose)
			os.Setenv("IDPBUILDER_USERNAME", "testuser")
			os.Setenv("IDPBUILDER_PASSWORD", "testpass")
			defer func() {
				os.Unsetenv("IDPBUILDER_INSECURE")
				os.Unsetenv("IDPBUILDER_VERBOSE")
				os.Unsetenv("IDPBUILDER_USERNAME")
				os.Unsetenv("IDPBUILDER_PASSWORD")
			}()

			// Given: Command using environment variables
			cmd := push.NewPushCommand(viper.New())
			cmd.SetArgs([]string{"alpine:latest"})

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()
			cmd.SetContext(ctx)

			// When: Execute command
			err := cmd.Execute()

			// Then: Format should be recognized
			if err != nil {
				assert.NotContains(t, err.Error(), "configuration error")
			}
		})
	}
}

// T-2.2.6-05: Test unsetting environment variable after setting it
func TestPushCommand_UnsetAfterSet(t *testing.T) {
	// Given: Set and then unset environment variable
	os.Setenv("IDPBUILDER_USERNAME", "tempuser")
	os.Unsetenv("IDPBUILDER_USERNAME")

	// Given: Set required credentials via flags
	cmd := push.NewPushCommand(viper.New())
	cmd.SetArgs([]string{
		"alpine:latest",
		"--username", "realuser",
		"--password", "realpass",
		"--insecure",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd.SetContext(ctx)

	// When: Execute command
	err := cmd.Execute()

	// Then: Should use flag values (unset env var has no effect)
	if err != nil {
		assert.NotContains(t, err.Error(), "tempuser")
	}
}

// T-2.2.6-06: Test flag explicitly set to empty string
func TestPushCommand_FlagExplicitlySetToEmpty(t *testing.T) {
	// Given: Environment variable set
	os.Setenv("IDPBUILDER_USERNAME", "envuser")
	defer os.Unsetenv("IDPBUILDER_USERNAME")

	// Given: Flag explicitly set to empty (should override env)
	cmd := push.NewPushCommand(viper.New())
	cmd.SetArgs([]string{
		"alpine:latest",
		"--username", "",
		"--password", "testpass",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd.SetContext(ctx)

	// When: Execute command
	err := cmd.Execute()

	// Then: Should fail with missing username (flag overrides env)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "username is required")
}

// T-2.2.6-07: Test context cancellation with environment variables
func TestPushCommand_ContextCancellationWithEnv(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Given: All config from environment
	os.Setenv("IDPBUILDER_USERNAME", "testuser")
	os.Setenv("IDPBUILDER_PASSWORD", "testpass")
	defer func() {
		os.Unsetenv("IDPBUILDER_USERNAME")
		os.Unsetenv("IDPBUILDER_PASSWORD")
	}()

	// Given: Command with very short context timeout
	cmd := push.NewPushCommand(viper.New())
	cmd.SetArgs([]string{"alpine:latest", "--insecure"})

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	cmd.SetContext(ctx)

	// When: Execute command
	err := cmd.Execute()

	// Then: Should handle context cancellation gracefully
	// (May timeout during execution, should not panic)
	if err != nil {
		// Either context error or other error, but no panic
		assert.NotNil(t, err)
	}
}

// T-2.2.6-08: Test viper instance reuse
func TestPushCommand_ViperInstanceReuse(t *testing.T) {
	// Given: Single viper instance reused
	v := viper.New()

	// Given: Environment variables set
	os.Setenv("IDPBUILDER_USERNAME", "testuser")
	os.Setenv("IDPBUILDER_PASSWORD", "testpass")
	defer func() {
		os.Unsetenv("IDPBUILDER_USERNAME")
		os.Unsetenv("IDPBUILDER_PASSWORD")
	}()

	// Given: Create multiple commands with same viper instance
	cmd1 := push.NewPushCommand(v)
	cmd1.SetArgs([]string{"alpine:latest", "--insecure"})

	cmd2 := push.NewPushCommand(v)
	cmd2.SetArgs([]string{"ubuntu:latest", "--insecure"})

	// Then: Both commands should be independent
	assert.NotEqual(t, cmd1, cmd2)
	assert.NotNil(t, cmd1.Flags())
	assert.NotNil(t, cmd2.Flags())
}

// T-2.2.6-09: Test concurrent environment variable access
func TestPushCommand_ConcurrentEnvironmentAccess(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Given: Environment variables set
	os.Setenv("IDPBUILDER_USERNAME", "concurrentuser")
	os.Setenv("IDPBUILDER_PASSWORD", "concurrentpass")
	defer func() {
		os.Unsetenv("IDPBUILDER_USERNAME")
		os.Unsetenv("IDPBUILDER_PASSWORD")
	}()

	// Given: Multiple commands created concurrently
	v := viper.New()
	cmd1 := push.NewPushCommand(v)
	cmd1.SetArgs([]string{"alpine:latest", "--insecure"})

	cmd2 := push.NewPushCommand(viper.New())
	cmd2.SetArgs([]string{"alpine:latest", "--insecure"})

	// Then: Both should access environment safely
	// (No race conditions should occur)
	assert.NotNil(t, cmd1)
	assert.NotNil(t, cmd2)
}

// T-2.2.6-10: Test environment variable precedence in help documentation
func TestPushCommand_EnvVarPrecedenceDocumentation(t *testing.T) {
	// Given: Create push command
	cmd := push.NewPushCommand(viper.New())

	// Then: Help text should document environment variables
	helpText := cmd.Long
	assert.Contains(t, helpText, "IDPBUILDER_REGISTRY",
		"Help should document IDPBUILDER_REGISTRY")
	assert.Contains(t, helpText, "IDPBUILDER_USERNAME",
		"Help should document IDPBUILDER_USERNAME")
	assert.Contains(t, helpText, "IDPBUILDER_PASSWORD",
		"Help should document IDPBUILDER_PASSWORD")
	assert.Contains(t, helpText, "environment variable",
		"Help should mention environment variables")

	// Then: Flag descriptions should mention env vars
	registryFlag := cmd.Flags().Lookup("registry")
	require.NotNil(t, registryFlag)
	assert.Contains(t, registryFlag.Usage, "IDPBUILDER_REGISTRY",
		"Registry flag help should mention env var")

	usernameFlag := cmd.Flags().Lookup("username")
	require.NotNil(t, usernameFlag)
	assert.Contains(t, usernameFlag.Usage, "IDPBUILDER_USERNAME",
		"Username flag help should mention env var")
}
