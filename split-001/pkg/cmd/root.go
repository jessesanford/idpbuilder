package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/cnoe-io/idpbuilder/pkg/cmd/helpers"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "idpbuilder",
	Short: "Manage reference IDPs",
	Long: `idpbuilder is a tool for building and managing 
reference Internal Developer Platforms (IDPs).

This CLI provides commands to create, delete, get, and version
IDP configurations and resources.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command
func Execute(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx)
}

// GetRootCmd returns the root command for testing
func GetRootCmd() *cobra.Command {
	return rootCmd
}

func init() {
	// Set up persistent flags
	rootCmd.PersistentFlags().StringVarP(&helpers.LogLevel, "log-level", "l", "info", helpers.LogLevelMsg)
	rootCmd.PersistentFlags().BoolVar(&helpers.ColoredOutput, "color", false, helpers.ColoredOutputMsg)
	rootCmd.PersistentFlags().BoolVarP(&helpers.Verbose, "verbose", "v", false, helpers.VerboseMsg)
	
	// Initialize logging and configuration
	helpers.InitializeLogging()
	helpers.LoadConfig()
}

// AddSubCommand adds a subcommand to the root command
func AddSubCommand(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}

// SetVersion sets the version for the root command
func SetVersion(version string) {
	rootCmd.Version = version
}