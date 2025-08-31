package build

import (
	"context"
	"fmt"

	buildintegration "github.com/cnoe-io/idpbuilder/pkg/build"
	"github.com/cnoe-io/idpbuilder/pkg/cmd/helpers"
	"github.com/spf13/cobra"
)

const (
	dockerfilePathUsage = "Path to the Dockerfile to build"
	contextDirUsage     = "Build context directory"
	tagUsage            = "Name and optionally a tag in the 'name:tag' format"
	insecureUsage       = "Skip certificate validation for registries"
)

var (
	dockerfilePath string
	contextDir     string
	tag            string
	insecure       bool
)

var BuildCmd = &cobra.Command{
	Use:          "build [context]",
	Short:        "Build container images using Buildah",
	Long:         `Build OCI container images from Dockerfiles with certificate support for secure registry operations`,
	RunE:         build,
	PreRunE:      preBuildE,
	SilenceUsage: true,
}

func init() {
	BuildCmd.Flags().StringVarP(&dockerfilePath, "file", "f", "Dockerfile", dockerfilePathUsage)
	BuildCmd.Flags().StringVar(&contextDir, "context", ".", contextDirUsage)
	BuildCmd.Flags().StringVarP(&tag, "tag", "t", "", tagUsage)
	BuildCmd.Flags().BoolVar(&insecure, "insecure", false, insecureUsage)
}

func preBuildE(cmd *cobra.Command, args []string) error {
	return helpers.SetLogger()
}

func build(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// Use context from args if provided, otherwise use flag
	buildContext := contextDir
	if len(args) > 0 {
		buildContext = args[0]
	}

	if tag == "" {
		return fmt.Errorf("tag is required for build")
	}

	return runBuild(ctx, buildContext, dockerfilePath, tag, insecure)
}

func runBuild(ctx context.Context, contextDir, dockerfilePath, tag string, insecure bool) error {
	integration := buildintegration.NewIntegration()

	opts := buildintegration.BuildOptions{
		ContextDir:     contextDir,
		DockerfilePath: dockerfilePath,
		Tag:            tag,
		Insecure:       insecure,
	}

	return integration.Build(ctx, opts)
}
