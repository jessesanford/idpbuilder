package push

import (
	"context"
	"fmt"

	"github.com/cnoe-io/idpbuilder/pkg/auth"
	"github.com/cnoe-io/idpbuilder/pkg/cmd/helpers"
	"github.com/cnoe-io/idpbuilder/pkg/push"
	"github.com/spf13/cobra"
)

// PushCmd represents the push command
var PushCmd = &cobra.Command{
	Use:   "push [IMAGE]",
	Short: "Push container images to a registry",
	Long: `Push container images to a registry with authentication support.

Examples:
  # Push an image without authentication
  idpbuilder push myimage:latest

  # Push an image with username and password
  idpbuilder push myimage:latest --username myuser --password mypass

  # Push an image with short flags
  idpbuilder push myimage:latest -u myuser -p mypass`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPush(cmd, cmd.Context(), args[0])
	},
}

func init() {
	// Add authentication flags to the push command
	auth.AddAuthenticationFlags(PushCmd)

	// Add common flags
	PushCmd.Flags().BoolP("verbose", "v", false, "Enable verbose logging")
	PushCmd.Flags().Bool("insecure", false, "Allow insecure registry connections")
}

// runPush executes the push command with the provided image name
func runPush(cmd *cobra.Command, ctx context.Context, imageName string) error {
	// Create logger
	logger := helpers.CmdLogger

	// Create push operation from command flags
	// NewPushOperationFromCommand extracts all flags (auth, TLS, registry, etc.)
	operation, err := push.NewPushOperationFromCommand(cmd, logger)
	if err != nil {
		return fmt.Errorf("failed to create push operation: %w", err)
	}

	// Execute the actual push operation
	// This calls the complete implementation in pkg/push/operations.go
	result, err := operation.Execute(ctx)
	if err != nil {
		return fmt.Errorf("push operation failed: %w", err)
	}

	// Report success with summary
	if result.ImagesPushed > 0 {
		fmt.Printf("\nSuccessfully pushed %d image(s) to registry\n", result.ImagesPushed)
		fmt.Printf("Total bytes transferred: %d\n", result.TotalBytes)
		fmt.Printf("Duration: %s\n", result.TotalDuration)
	} else {
		fmt.Println("\nNo images were pushed")
	}

	// Show failures if any
	if result.ImagesFailed > 0 {
		fmt.Printf("\nWarning: %d image(s) failed to push\n", result.ImagesFailed)
		for _, err := range result.Errors {
			fmt.Printf("  - %v\n", err)
		}
	}

	return nil
}