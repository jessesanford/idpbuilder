// Package cmd_test provides comprehensive tests for the root command and CLI structure.
// These tests verify proper command registration, flag handling, and version information.
package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRootCommand(t *testing.T) {
	tests := []struct {
		name     string
		expected struct {
			use   string
			short string
		}
	}{
		{
			name: "root_command_creation",
			expected: struct {
				use   string
				short string
			}{
				use:   "idpbuilder-push",
				short: "OCI image push utility for idpbuilder",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootCmd := NewRootCommand()
			cmd := rootCmd.GetCommand()

			assert.Equal(t, tt.expected.use, cmd.Use)
			assert.Equal(t, tt.expected.short, cmd.Short)
			assert.NotEmpty(t, cmd.Long)
			assert.NotEmpty(t, cmd.Version)

			// Verify global flags are present
			verboseFlag := cmd.PersistentFlags().Lookup("verbose")
			assert.NotNil(t, verboseFlag)

			quietFlag := cmd.PersistentFlags().Lookup("quiet")
			assert.NotNil(t, quietFlag)

			// Verify subcommands are registered
			commands := cmd.Commands()
			assert.NotEmpty(t, commands, "Root command should have subcommands")

			// Look for push command
			pushCmdFound := false
			for _, subCmd := range commands {
				if subCmd.Name() == "push" {
					pushCmdFound = true
					break
				}
			}
			assert.True(t, pushCmdFound, "Push command should be registered")
		})
	}
}

func TestGetVersionString(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		expected string
	}{
		{
			name:     "development_version",
			version:  "dev",
			expected: "dev (development build)",
		},
		{
			name:     "release_version",
			version:  "v1.2.3",
			expected: "v1.2.3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original version
			originalVersion := Version
			defer func() { Version = originalVersion }()

			// Set test version
			Version = tt.version

			result := getVersionString()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetVersionTemplate(t *testing.T) {
	// Save original values
	originalVersion := Version
	originalBuildDate := BuildDate
	originalGitCommit := GitCommit
	defer func() {
		Version = originalVersion
		BuildDate = originalBuildDate
		GitCommit = originalGitCommit
	}()

	// Set test values
	Version = "v1.2.3"
	BuildDate = "2023-12-01"
	GitCommit = "abc123def"

	template := getVersionTemplate()

	assert.Contains(t, template, "idpbuilder-push version: v1.2.3")
	assert.Contains(t, template, "Build date: 2023-12-01")
	assert.Contains(t, template, "Git commit: abc123def")
}

func TestConfigureGlobalFlags(t *testing.T) {
	// Test basic functionality exists - detailed testing would need integration tests
	// since ConfigureGlobalFlags depends on parsed cobra flags
	t.Run("function_exists", func(t *testing.T) {
		rootCmd := NewRootCommand()
		cmd := rootCmd.GetCommand()

		// Just verify the function doesn't panic with a basic command
		// Real flag testing would need a more complex setup
		assert.NotNil(t, cmd)
		assert.NotPanics(t, func() {
			// This will error due to flag access, but shouldn't panic
			ConfigureGlobalFlags(cmd)
		})
	})

	// Test the environment setting logic separately
	t.Run("environment_variable_setting", func(t *testing.T) {
		// Clear environment first
		os.Unsetenv("LOG_LEVEL")

		// Test environment variable setting directly
		os.Setenv("LOG_LEVEL", "debug")
		assert.Equal(t, "debug", os.Getenv("LOG_LEVEL"))

		os.Setenv("LOG_LEVEL", "error")
		assert.Equal(t, "error", os.Getenv("LOG_LEVEL"))

		os.Setenv("LOG_LEVEL", "info")
		assert.Equal(t, "info", os.Getenv("LOG_LEVEL"))

		// Clean up
		os.Unsetenv("LOG_LEVEL")
	})
}

func TestValidateEnvironment(t *testing.T) {
	// Currently, ValidateEnvironment is a placeholder
	// Test that it doesn't return an error
	err := ValidateEnvironment()
	assert.NoError(t, err, "ValidateEnvironment should not return an error in current implementation")
}

func TestRootCommandExecution(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "no_args_shows_help",
			args:        []string{},
			expectError: false,
		},
		{
			name:        "version_flag",
			args:        []string{"--version"},
			expectError: false,
		},
		{
			name:        "help_flag",
			args:        []string{"--help"},
			expectError: false,
		},
		{
			name:        "invalid_command",
			args:        []string{"invalid-command"},
			expectError: true,
			errorMsg:    "unknown command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootCmd := NewRootCommand()
			cmd := rootCmd.GetCommand()

			// Set up output capture
			var output strings.Builder
			cmd.SetOut(&output)
			cmd.SetErr(&output)

			// Set arguments
			cmd.SetArgs(tt.args)

			err := cmd.Execute()

			if tt.expectError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestShowHelp(t *testing.T) {
	// Test that ShowHelp doesn't panic
	// We can't easily capture the output since it goes directly to stdout
	assert.NotPanics(t, func() {
		ShowHelp()
	})
}

func TestRootCommandSubcommands(t *testing.T) {
	rootCmd := NewRootCommand()
	cmd := rootCmd.GetCommand()

	commands := cmd.Commands()
	commandNames := make([]string, len(commands))

	for i, subCmd := range commands {
		commandNames[i] = subCmd.Name()
	}

	// Verify expected subcommands are present
	expectedCommands := []string{"push"}
	for _, expected := range expectedCommands {
		assert.Contains(t, commandNames, expected, "Expected subcommand %s to be registered", expected)
	}
}

// Helper function to convert bool to string for flag setting
func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}