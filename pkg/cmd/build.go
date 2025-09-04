package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cnoe-io/idpbuilder/pkg/builder"
	"github.com/cnoe-io/idpbuilder/pkg/logger"
	"github.com/spf13/cobra"
)

// BuildOptions contains options for the build command
type BuildOptions struct {
	Context    string
	Dockerfile string
	Tags       []string
	Platform   string
	Output     string
	Labels     map[string]string
	NoCache    bool
	Pull       bool
	Quiet      bool
}

var buildOpts = &BuildOptions{
	Labels: make(map[string]string),
}

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build [CONTEXT]",
	Short: "Build container images from source",
	Long: `Build container images from source code using the go-containerregistry library.

This command builds container images from a Dockerfile and build context, similar to 
'docker build' but using pure Go implementation without requiring Docker daemon.

Examples:
  # Build from current directory
  idpbuilder build .

  # Build with specific tag
  idpbuilder build -t myapp:latest .

  # Build with multiple tags
  idpbuilder build -t myapp:latest -t myapp:v1.0 .

  # Build with labels
  idpbuilder build --label version=1.0 --label env=prod .

  # Build with custom Dockerfile
  idpbuilder build -f custom.Dockerfile .`,
	Args: cobra.MaximumNArgs(1),
	RunE: runBuild,
}

func init() {
	// Feature flag check - only add command if enabled
	if os.Getenv("ENABLE_CLI_TOOLS") == "true" {
		rootCmd.AddCommand(buildCmd)
	}

	buildCmd.Flags().StringVarP(&buildOpts.Dockerfile, "file", "f", "Dockerfile", "Path to Dockerfile")
	buildCmd.Flags().StringArrayVarP(&buildOpts.Tags, "tag", "t", nil, "Repository name and optionally a tag (format: 'name:tag')")
	buildCmd.Flags().StringVar(&buildOpts.Platform, "platform", "linux/amd64", "Target platform for build")
	buildCmd.Flags().StringVarP(&buildOpts.Output, "output", "o", "", "Output format (tar, docker-archive)")
	buildCmd.Flags().StringToStringVar(&buildOpts.Labels, "label", nil, "Set metadata labels for image")
	buildCmd.Flags().BoolVar(&buildOpts.NoCache, "no-cache", false, "Do not use cache when building image")
	buildCmd.Flags().BoolVar(&buildOpts.Pull, "pull", false, "Always attempt to pull base image")
	buildCmd.Flags().BoolVarP(&buildOpts.Quiet, "quiet", "q", false, "Suppress build output")
}

func runBuild(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// Check feature flag
	if !isFeatureEnabled() {
		return fmt.Errorf("build command requires ENABLE_CLI_TOOLS=true")
	}

	// Set default context
	buildOpts.Context = "."
	if len(args) > 0 {
		buildOpts.Context = args[0]
	}

	// Validate options
	if err := validateBuildOptions(); err != nil {
		return fmt.Errorf("invalid build options: %w", err)
	}

	// Create logger
	log := logger.NewDefault()
	if buildOpts.Quiet {
		// Create a logger that outputs to a discard writer
		log = logger.New(io.Discard)
	}

	// Create build options
	opts := builder.DefaultBuildOptions().
		WithContext(buildOpts.Context).
		WithDockerfile(buildOpts.Dockerfile).
		WithPlatform(buildOpts.Platform).
		WithLogger(log).
		WithTimeout(30 * time.Minute)

	// Add tags
	if len(buildOpts.Tags) > 0 {
		opts.WithTags(buildOpts.Tags...)
	} else {
		// Default tag if none specified
		opts.WithTags("latest")
	}

	// Add labels
	for key, value := range buildOpts.Labels {
		opts.WithLabel(key, value)
	}

	// Create builder
	b, err := builder.New(opts)
	if err != nil {
		return fmt.Errorf("failed to create builder: %w", err)
	}

	// Build image
	fmt.Printf("Building image from context: %s\n", buildOpts.Context)
	if !buildOpts.Quiet {
		fmt.Printf("Using Dockerfile: %s\n", buildOpts.Dockerfile)
		fmt.Printf("Target platform: %s\n", buildOpts.Platform)
		if len(buildOpts.Tags) > 0 {
			fmt.Printf("Tags: %s\n", strings.Join(buildOpts.Tags, ", "))
		}
	}

	result, err := b.Build()
	if err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	// Handle output
	if buildOpts.Output != "" {
		if err := handleBuildOutput(ctx, result, buildOpts.Output); err != nil {
			return fmt.Errorf("failed to save output: %w", err)
		}
	}

	if !buildOpts.Quiet {
		fmt.Printf("Successfully built image: %s\n", result.ImageID)
		if len(result.Tags) > 0 {
			fmt.Printf("Tagged as: %s\n", strings.Join(result.Tags, ", "))
		}
	}

	return nil
}

func validateBuildOptions() error {
	// Check context exists
	if _, err := os.Stat(buildOpts.Context); os.IsNotExist(err) {
		return fmt.Errorf("build context does not exist: %s", buildOpts.Context)
	}

	// Check Dockerfile exists (if absolute path) or exists in context
	dockerfilePath := buildOpts.Dockerfile
	if !filepath.IsAbs(dockerfilePath) {
		dockerfilePath = filepath.Join(buildOpts.Context, dockerfilePath)
	}

	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		return fmt.Errorf("dockerfile does not exist: %s", dockerfilePath)
	}

	// Validate platform format
	parts := strings.Split(buildOpts.Platform, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid platform format: %s (expected: os/arch)", buildOpts.Platform)
	}

	// Validate tags
	for _, tag := range buildOpts.Tags {
		if strings.Contains(tag, " ") || tag == "" {
			return fmt.Errorf("invalid tag format: %s", tag)
		}
	}

	return nil
}

func handleBuildOutput(ctx context.Context, result *builder.BuildResult, output string) error {
	switch output {
	case "tar":
		return saveAsTarball(ctx, result, "image.tar")
	case "docker-archive":
		return saveAsTarball(ctx, result, "image.tar")
	default:
		return fmt.Errorf("unsupported output format: %s", output)
	}
}

func saveAsTarball(ctx context.Context, result *builder.BuildResult, filename string) error {
	// Create tarball manager
	opts := builder.DefaultBuildOptions().WithLogger(&quietLogger{})
	tm := builder.NewTarballManager(opts)

	// Save image to tarball
	tag := "latest"
	if len(result.Tags) > 0 {
		tag = result.Tags[0]
	}

	err := tm.SaveImageToTarball(ctx, result.Image, filename, tag)
	if err != nil {
		return fmt.Errorf("failed to save image as tarball: %w", err)
	}

	fmt.Printf("Image saved to: %s\n", filename)
	return nil
}

func isFeatureEnabled() bool {
	return os.Getenv("ENABLE_CLI_TOOLS") == "true"
}

// quietLogger is a logger that suppresses output (for testing)
type quietLogger struct{}

func (q *quietLogger) Debug(msg string, args ...interface{}) {}
func (q *quietLogger) Info(msg string, args ...interface{})  {}
func (q *quietLogger) Warn(msg string, args ...interface{})  {}
func (q *quietLogger) Error(msg string, args ...interface{}) {}