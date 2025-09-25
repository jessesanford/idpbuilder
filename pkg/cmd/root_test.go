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
	tests := []struct {
		name        string
		verbose     bool
		quiet       bool
		expectError bool
		errorMsg    string
		expectedEnv string
	}{
		{
			name:        "verbose_flag",
			verbose:     true,
			quiet:       false,
			expectError: false,
			expectedEnv: "debug",
		},
		{
			name:        "quiet_flag",
			verbose:     false,
			quiet:       true,
			expectError: false,
			expectedEnv: "error",
		},
		{
			name:        "default_flags",
			verbose:     false,
			quiet:       false,
			expectError: false,
			expectedEnv: "info",
		},
		{
			name:        "conflicting_flags",
			verbose:     true,
			quiet:       true,
			expectError: true,
			errorMsg:    "cannot use both --verbose and --quiet flags",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a root command to get flags
			rootCmd := NewRootCommand()
			cmd := rootCmd.GetCommand()

			// Set flag values
			cmd.PersistentFlags().Set("verbose", boolToString(tt.verbose))
			cmd.PersistentFlags().Set("quiet", boolToString(tt.quiet))

			// Clear environment variable first
			os.Unsetenv("LOG_LEVEL")

			err := ConfigureGlobalFlags(cmd)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedEnv, os.Getenv("LOG_LEVEL"))
			}

			// Clean up
			os.Unsetenv("LOG_LEVEL")
		})
	}
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