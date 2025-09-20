package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/cnoe-io/idpbuilder/pkg/cmd/create"
	"github.com/cnoe-io/idpbuilder/pkg/cmd/delete"
	"github.com/cnoe-io/idpbuilder/pkg/cmd/get"
	"github.com/cnoe-io/idpbuilder/pkg/cmd/helpers"
	"github.com/cnoe-io/idpbuilder/pkg/cmd/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "idpbuilder",
	Short: "Manage reference IDPs",
	Long:  "",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&helpers.LogLevel, "log-level", "l", "info", helpers.LogLevelMsg)
	rootCmd.PersistentFlags().BoolVar(&helpers.ColoredOutput, "color", false, helpers.ColoredOutputMsg)
	rootCmd.AddCommand(create.CreateCmd)
	rootCmd.AddCommand(get.GetCmd)
	rootCmd.AddCommand(delete.DeleteCmd)
	rootCmd.AddCommand(version.VersionCmd)

	// Add OCI commands - can be controlled by build tags or env vars
	if isOCIEnabled() {
		rootCmd.AddCommand(BuildCmd)
		rootCmd.AddCommand(PushCmd)
	}
}

// isOCIEnabled checks if OCI features should be enabled
// This can be controlled by environment variable or build tags
func isOCIEnabled() bool {
	// Always enable for now, can add env var check later
	// e.g., return os.Getenv("IDPBUILDER_ENABLE_OCI") == "true"
	return true
}

func Execute(ctx context.Context) {
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}