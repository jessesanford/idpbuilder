package build

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cnoe-io/idpbuilder/pkg/cmd/helpers"
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
		ctx := cmd.Context()
		return executeBuild(ctx, args[0], buildOptions{
			dockerfile: dockerfile,
			tag:        tag,
			platform:   platform,
			verbose:    helpers.LogLevel == "debug",
		})
	},
}

type buildOptions struct {
	dockerfile string
	tag        string
	platform   string
	verbose    bool
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

func executeBuild(ctx context.Context, contextPath string, opts buildOptions) error {
	// TODO: Will be implemented in Step 3
	fmt.Printf("Building image from %s\n", contextPath)
	fmt.Printf("  Dockerfile: %s\n", opts.dockerfile)
	fmt.Printf("  Tag: %s\n", opts.tag)
	if opts.platform != "" {
		fmt.Printf("  Platform: %s\n", opts.platform)
	}
	
	// Validate context path exists
	if contextPath != "-" {
		absPath, err := filepath.Abs(contextPath)
		if err != nil {
			return fmt.Errorf("invalid context path: %w", err)
		}
		
		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			return fmt.Errorf("build context does not exist: %s", absPath)
		}
	}
	
	fmt.Println("Build functionality will be implemented in Step 3")
	return nil
}