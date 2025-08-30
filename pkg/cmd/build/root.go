package build

import (
	"github.com/cnoe-io/idpbuilder/pkg/cmd/helpers"
	"github.com/cnoe-io/idpbuilder/pkg/oci/commands"
	"github.com/spf13/cobra"
)

var BuildCmd = &cobra.Command{
	Use:   "build [OPTIONS] PATH | URL | -",
	Short: "Build OCI images with certificate auto-configuration",
	Long: `Build OCI images using buildah with automatic certificate configuration.

The build command automatically extracts and configures certificates from the
local Kind cluster to enable secure builds. Images are built using buildah
with proper certificate trust configuration.

Examples:
  # Build image from current directory
  idpbuilder build .

  # Build with custom Dockerfile and tag
  idpbuilder build --file custom.Dockerfile --tag myapp:v1.0 .

  # Build for multiple platforms
  idpbuilder build --platform linux/amd64,linux/arm64 .
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if buildah is available
		if err := commands.CheckBuildahAvailable(); err != nil {
			return err
		}
		
		ctx := cmd.Context()
		return commands.ExecuteBuild(ctx, args[0], commands.BuildOptions{
			Dockerfile: dockerfile,
			Tag:        tag,
			Platform:   platform,
			Verbose:    helpers.LogLevel == "debug",
			Context:    args[0],
		})
	},
}

var (
	dockerfile string
	tag        string
	platform   string
)

func init() {
	BuildCmd.Flags().StringVarP(&dockerfile, "file", "f", "Dockerfile", "Name of the Dockerfile")
	BuildCmd.Flags().StringVarP(&tag, "tag", "t", "", "Name and optionally tag in 'name:tag' format")
	BuildCmd.Flags().StringVar(&platform, "platform", "", "Set platform if server is multi-platform capable")
	
	// Mark tag as required
	BuildCmd.MarkFlagRequired("tag")
}