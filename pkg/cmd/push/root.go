package push

import (
	"github.com/cnoe-io/idpbuilder/pkg/cmd/helpers"
	"github.com/cnoe-io/idpbuilder/pkg/oci/commands"
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
		// Check if buildah/podman is available
		if err := commands.CheckPodmanAvailable(); err != nil {
			return err
		}
		
		ctx := cmd.Context()
		return commands.ExecutePush(ctx, args[0], commands.PushOptions{
			Insecure: insecure,
			Username: username,
			Password: password,
			Verbose:  helpers.LogLevel == "debug",
		})
	},
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