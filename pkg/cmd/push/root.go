package push

import (
	"context"
	"fmt"
	"strings"

	"github.com/cnoe-io/idpbuilder/pkg/cmd/helpers"
	"github.com/spf13/cobra"
)

var PushCmd = &cobra.Command{
	Use:   "push [OPTIONS] IMAGE",
	Short: "Push OCI images to registry with certificate auto-configuration",
	Long: `Push OCI images to registry with automatic certificate configuration.

The push command automatically extracts and configures certificates from the
local Kind cluster to enable secure pushes to the Gitea registry. Authentication
can be provided via flags or environment variables.

Examples:
  # Push image to default registry
  idpbuilder push myapp:v1.0

  # Push with authentication
  idpbuilder push --username admin --password secret myapp:v1.0

  # Push bypassing certificate verification
  idpbuilder push --insecure myapp:v1.0
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return executePush(ctx, args[0], pushOptions{
			insecure: insecure,
			username: username,
			password: password,
			verbose:  helpers.LogLevel == "debug",
		})
	},
}

type pushOptions struct {
	insecure bool
	username string
	password string
	verbose  bool
}

var (
	insecure bool
	username string
	password string
)

func init() {
	PushCmd.Flags().BoolVar(&insecure, "insecure", false, "Skip certificate verification")
	PushCmd.Flags().StringVar(&username, "username", "", "Registry username")
	PushCmd.Flags().StringVar(&password, "password", "", "Registry password")
}

func executePush(ctx context.Context, imageRef string, opts pushOptions) error {
	// TODO: Will be implemented in Step 4
	fmt.Printf("Pushing image %s\n", imageRef)
	if opts.username != "" {
		fmt.Printf("  Username: %s\n", opts.username)
		fmt.Println("  Password: [REDACTED]")
	}
	if opts.insecure {
		fmt.Println("  Insecure: true (skipping certificate verification)")
	}
	
	// Basic validation of image reference
	if !strings.Contains(imageRef, ":") {
		return fmt.Errorf("image reference must include tag: %s", imageRef)
	}
	
	parts := strings.Split(imageRef, "/")
	if len(parts) < 2 {
		return fmt.Errorf("invalid image reference format: %s", imageRef)
	}
	
	fmt.Println("Push functionality will be implemented in Step 4")
	return nil
}