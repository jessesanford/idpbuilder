// Package cmd provides the CLI command structure for idpbuilder-push.
// This package defines the root command and integrates all subcommands
// including push, with proper configuration and error handling.
package cmd

import (
	"fmt"
	"os"

	"github.com/cnoe-io/idpbuilder-push/client-interface-tests-split-003/pkg/cmd/push"
	"github.com/spf13/cobra"
)

// Version information for the CLI tool
var (
	Version   = "dev"      // Set by build process
	BuildDate = "unknown"  // Set by build process
	GitCommit = "unknown"  // Set by build process
)

// RootCommand represents the base command when called without any subcommands.
type RootCommand struct {
	cmd *cobra.Command
}

// NewRootCommand creates and configures the root command for idpbuilder-push.
func NewRootCommand() *RootCommand {
	rootCmd := &RootCommand{}

	cmd := &cobra.Command{
		Use:   "idpbuilder-push",
		Short: "OCI image push utility for idpbuilder",
		Long: `idpbuilder-push is a command-line tool for pushing OCI container images
to registries. It provides authentication support, insecure connection handling,
and integration with Kubernetes cluster management.

This tool is designed to work with idpbuilder's development workflow,
enabling developers to easily push container images to local or remote
registries during development and testing.

Key Features:
  • Push OCI images to any registry
  • Username/password authentication
  • Support for insecure (HTTP) registries
  • Integration with Kind clusters
  • Configurable timeouts and retry logic
`,
		Version:      getVersionString(),
		SilenceUsage: true,
	}

	// Add global flags
	cmd.PersistentFlags().Bool("verbose", false, "Enable verbose logging")
	cmd.PersistentFlags().Bool("quiet", false, "Suppress non-error output")

	// Add subcommands
	cmd.AddCommand(push.NewPushCommand())

	// Add version information
	cmd.SetVersionTemplate(getVersionTemplate())

	rootCmd.cmd = cmd
	return rootCmd
}

// Execute runs the root command and handles any errors appropriately.
func (r *RootCommand) Execute() error {
	return r.cmd.Execute()
}

// GetCommand returns the underlying cobra command for testing purposes.
func (r *RootCommand) GetCommand() *cobra.Command {
	return r.cmd
}

// Execute is the main entry point for the CLI application.
// It creates the root command and executes it, handling any errors.
func Execute() {
	rootCmd := NewRootCommand()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// getVersionString returns a formatted version string with build information.
func getVersionString() string {
	if Version == "dev" {
		return "dev (development build)"
	}
	return Version
}

// getVersionTemplate returns a custom version template with detailed build info.
func getVersionTemplate() string {
	return fmt.Sprintf(`idpbuilder-push version: %s
Build date: %s
Git commit: %s
{{if .Version}}Version: {{.Version}}{{end}}
`, Version, BuildDate, GitCommit)
}

// ConfigureGlobalFlags sets up global configuration based on persistent flags.
func ConfigureGlobalFlags(cmd *cobra.Command) error {
	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		return fmt.Errorf("failed to get verbose flag: %w", err)
	}

	quiet, err := cmd.Flags().GetBool("quiet")
	if err != nil {
		return fmt.Errorf("failed to get quiet flag: %w", err)
	}

	// Configure logging based on flags
	if verbose && quiet {
		return fmt.Errorf("cannot use both --verbose and --quiet flags")
	}

	// Set up logging configuration
	if verbose {
		// Enable debug logging
		os.Setenv("LOG_LEVEL", "debug")
	} else if quiet {
		// Suppress non-error output
		os.Setenv("LOG_LEVEL", "error")
	} else {
		// Default logging level
		os.Setenv("LOG_LEVEL", "info")
	}

	return nil
}

// ValidateEnvironment checks that the runtime environment is suitable for operation.
func ValidateEnvironment() error {
	// Check that we can access the Docker daemon or containerd (for future image operations)
	// For now, this is a placeholder for environment validation

	return nil
}

// ShowHelp displays the help message for the root command.
func ShowHelp() {
	rootCmd := NewRootCommand()
	rootCmd.cmd.Help()
}