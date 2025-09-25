package cmd

import (
	"fmt"
	"os"

	"github.com/cnoe-io/idpbuilder-push/client-interface-tests-split-004/pkg/cmd/push"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewRootCommand creates the root command for the idpbuilder CLI
func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "idpbuilder",
		Short: "IDPBuilder - Build and manage IDP components",
		Long: `IDPBuilder is a tool for building and managing Internal Developer Platform components.

It provides commands to build, push, and manage container images and Kubernetes resources
for your IDP infrastructure.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add subcommands
	rootCmd.AddCommand(push.NewPushCommand())

	// Global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Enable quiet output")

	// Bind flags to viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet"))

	return rootCmd
}

// Execute executes the root command and handles errors
func Execute() {
	if err := NewRootCommand().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}