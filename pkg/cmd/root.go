// pkg/cmd/root.go
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/cnoe-io/idpbuilder/pkg/cmd/get"
)

// NewRootCmd creates the root command for idpbuilder
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "idpbuilder",
		Short: "IDP Builder - Internal Developer Platform Builder",
		Long: `IDP Builder is a tool for building and managing internal developer platforms.
It provides commands for creating, managing, and deploying platform components.`,
	}

	// Add subcommands
	cmd.AddCommand(get.NewGetCmd())

	return cmd
}