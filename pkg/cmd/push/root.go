package push

import (
	"context"
	"fmt"

	"github.com/cnoe-io/idpbuilder/pkg/cmd/helpers"
	"github.com/cnoe-io/idpbuilder/pkg/registry"
	"github.com/spf13/cobra"
)

const (
	insecureUsage = "Skip TLS certificate validation for registries"
	usernameUsage = "Username for registry authentication"
	passwordUsage = "Password for registry authentication"
	authFileUsage = "Path to authentication file"
)

var (
	insecure bool
	username string
	password string
	authFile string
)

var PushCmd = &cobra.Command{
	Use:          "push [image] [registry]",
	Short:        "Push container images to registry",
	Long:         `Push OCI container images to Gitea or other registries with certificate support`,
	RunE:         push,
	PreRunE:      prePushE,
	SilenceUsage: true,
	Args:         cobra.RangeArgs(1, 2),
}

func init() {
	PushCmd.Flags().BoolVar(&insecure, "insecure", false, insecureUsage)
	PushCmd.Flags().StringVarP(&username, "username", "u", "", usernameUsage)
	PushCmd.Flags().StringVarP(&password, "password", "p", "", passwordUsage)
	PushCmd.Flags().StringVar(&authFile, "authfile", "", authFileUsage)
}

func prePushE(cmd *cobra.Command, args []string) error {
	return helpers.SetLogger()
}

func push(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	if len(args) < 1 {
		return fmt.Errorf("image name is required")
	}

	imageName := args[0]
	var registryURL string

	if len(args) > 1 {
		registryURL = args[1]
	}

	return runPush(ctx, imageName, registryURL, username, password, authFile, insecure)
}

func runPush(ctx context.Context, imageName, registryURL, username, password, authFile string, insecure bool) error {
	integration := registry.NewIntegration()

	// Pass credentials via environment for buildah
	if username != "" && password != "" {
		ctx = context.WithValue(ctx, "username", username)
		ctx = context.WithValue(ctx, "password", password)
	}

	opts := registry.PushOptions{
		ImageID:    imageName,  // The local image to push
		Repository: registryURL, // The registry destination
		Tag:        "",         // Tag is already in imageName
		Insecure:   insecure,
		Username:   username,
		Password:   password,
	}

	return integration.Push(ctx, opts)
}
