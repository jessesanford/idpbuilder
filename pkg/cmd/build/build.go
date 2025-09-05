package build

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/cnoe-io/idpbuilder/pkg/builder"
	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// BuildCmd is the build command
var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Assemble OCI image from context directory",
	Long:  `Assemble a single-layer OCI image from a directory using go-containerregistry.
The image is stored locally as an OCI tarball.`,
	Example: `  idpbuilder build --context ./app --tag myapp:v1
  idpbuilder build --context . --tag myimage:latest
  idpbuilder build --context ./src --tag registry.example.com/myapp:v2`,
	RunE: runBuild,
}

func init() {
	BuildCmd.Flags().String("context", ".", "Build context directory")
	BuildCmd.Flags().String("tag", "", "Image tag (required)")
	BuildCmd.Flags().String("output", "", "Output tarball path (optional)")
	BuildCmd.Flags().String("platform", "linux/amd64", "Target platform")
	BuildCmd.Flags().StringSlice("exclude", []string{}, "Exclude patterns")
	BuildCmd.MarkFlagRequired("tag")
}

func runBuild(cmd *cobra.Command, args []string) error {
	contextPath, _ := cmd.Flags().GetString("context")
	tag, _ := cmd.Flags().GetString("tag")
	output, _ := cmd.Flags().GetString("output")
	platformStr, _ := cmd.Flags().GetString("platform")
	excludePatterns, _ := cmd.Flags().GetStringSlice("exclude")

	// Validate context directory exists
	if _, err := os.Stat(contextPath); err != nil {
		return fmt.Errorf("context directory not found: %w", err)
	}

	// Create progress reporter
	// TODO: Integrate cli progress bar
	// progress := cli.NewProgressBar("Building image")
	// defer progress.Finish()

	// Parse platform
	platform, err := v1.ParsePlatform(platformStr)
	if err != nil {
		return fmt.Errorf("invalid platform: %w", err)
	}

	// Create build options
	opts := builder.BuildOptions{
		Platform:     *platform,
		Labels:       map[string]string{"org.opencontainers.image.ref.name": tag},
		FeatureFlags: make(map[string]bool),
	}

	// Note: Exclude patterns would need to be handled at the layer creation level
	// For now, we'll store them in labels for future reference
	if len(excludePatterns) > 0 {
		if opts.Labels == nil {
			opts.Labels = make(map[string]string)
		}
		opts.Labels["idpbuilder.exclude"] = fmt.Sprintf("%v", excludePatterns)
	}

	// Create builder
	b, err := builder.NewBuilder(opts)
	if err != nil {
		return fmt.Errorf("failed to create builder: %w", err)
	}

	// progress.UpdateMessage("Assembling layers from " + contextPath)
	fmt.Printf("Assembling layers from %s\n", contextPath)

	// If output is specified, build as tarball
	if output != "" {
		if err := b.BuildTarball(context.Background(), contextPath, output, opts); err != nil {
			return fmt.Errorf("build failed: %w", err)
		}
		// progress.UpdateMessage("Image built and saved to: " + output)
		fmt.Printf("Image built and saved to: %s\n", output)
		return nil
	}

	// Build the image
	image, err := b.Build(context.Background(), contextPath, opts)
	if err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	// Get image digest for output
	digest, err := image.Digest()
	if err != nil {
		return fmt.Errorf("failed to get image digest: %w", err)
	}

	// progress.UpdateMessage(fmt.Sprintf("Image built successfully: %s@%s", tag, digest.String()))
	fmt.Printf("Image built successfully: %s@%s\n", tag, digest.String())
	return nil
}