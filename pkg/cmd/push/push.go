package push

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cnoe-io/idpbuilder-push/client-interface-tests-split-004/pkg/kind"
	"github.com/cnoe-io/idpbuilder-push/client-interface-tests-split-004/pkg/oci"
	"github.com/spf13/cobra"
)

// PushOptions holds the options for the push command
type PushOptions struct {
	// Build options
	BuildPath   string
	Dockerfile  string
	Context     string
	Target      string
	Platform    string
	BuildArgs   map[string]string

	// Image options
	ImageName string
	ImageTag  string

	// Registry options
	RegistryURL      string
	RegistryUsername string
	RegistryPassword string
	RegistryInsecure bool

	// Kind options
	PushToKind  bool
	KindCluster string

	// General options
	Force   bool
	DryRun  bool
	Verbose bool
	Quiet   bool
}

// NewPushCommand creates a new push command
func NewPushCommand() *cobra.Command {
	opts := &PushOptions{
		Dockerfile:   "Dockerfile",
		Context:      ".",
		Platform:     "linux/amd64",
		ImageTag:     "latest",
		KindCluster:  "kind",
		BuildArgs:    make(map[string]string),
	}

	cmd := &cobra.Command{
		Use:   "push [PATH]",
		Short: "Build and push container images",
		Long: `Build and push container images to a registry or load them into a Kind cluster.

The push command builds a container image from the specified path and either:
- Pushes it to a container registry
- Loads it into a Kind cluster for local development

Examples:
  # Push to registry
  idpbuilder push ./app --name myapp --registry localhost:5000

  # Load into Kind cluster
  idpbuilder push ./app --name myapp --kind

  # Build with custom Dockerfile
  idpbuilder push ./app --name myapp --dockerfile custom.Dockerfile`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.BuildPath = args[0]
			return RunPush(cmd, args)
		},
	}

	// Build flags
	cmd.Flags().StringVar(&opts.Dockerfile, "dockerfile", opts.Dockerfile, "Path to Dockerfile")
	cmd.Flags().StringVar(&opts.Context, "context", opts.Context, "Build context path")
	cmd.Flags().StringVar(&opts.Target, "target", "", "Target stage in multi-stage build")
	cmd.Flags().StringVar(&opts.Platform, "platform", opts.Platform, "Target platform")
	cmd.Flags().StringToStringVar(&opts.BuildArgs, "build-arg", nil, "Build arguments")

	// Image flags
	cmd.Flags().StringVar(&opts.ImageName, "name", "", "Image name (required)")
	cmd.Flags().StringVar(&opts.ImageTag, "tag", opts.ImageTag, "Image tag")

	// Registry flags
	cmd.Flags().StringVar(&opts.RegistryURL, "registry", "", "Registry URL")
	cmd.Flags().StringVar(&opts.RegistryUsername, "username", "", "Registry username")
	cmd.Flags().StringVar(&opts.RegistryPassword, "password", "", "Registry password")
	cmd.Flags().BoolVar(&opts.RegistryInsecure, "insecure", false, "Allow insecure registry connections")

	// Kind flags
	cmd.Flags().BoolVar(&opts.PushToKind, "kind", false, "Load image into Kind cluster")
	cmd.Flags().StringVar(&opts.KindCluster, "kind-cluster", opts.KindCluster, "Kind cluster name")

	// General flags
	cmd.Flags().BoolVar(&opts.Force, "force", false, "Force push even if image exists")
	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "Show what would be done")
	cmd.Flags().BoolVarP(&opts.Verbose, "verbose", "v", false, "Verbose output")
	cmd.Flags().BoolVarP(&opts.Quiet, "quiet", "q", false, "Quiet output")

	// Mark required flags
	cmd.MarkFlagRequired("name")

	return cmd
}

// RunPush executes the push command
func RunPush(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Get options from flags
	opts, err := getOptionsFromFlags(cmd, args)
	if err != nil {
		return fmt.Errorf("parsing options: %w", err)
	}

	// Validate arguments
	if err := validatePushArgs(args); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Validate options
	if err := opts.validate(); err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}

	if opts.DryRun {
		return printDryRun(opts)
	}

	// Build the image
	imageRef, err := buildImage(ctx, opts)
	if err != nil {
		return fmt.Errorf("building image: %w", err)
	}

	if !opts.Quiet {
		fmt.Printf("Successfully built image: %s\n", imageRef)
	}

	// Push or load the image
	if opts.PushToKind {
		if err := loadToKind(ctx, imageRef, opts); err != nil {
			return fmt.Errorf("loading to Kind: %w", err)
		}
		if !opts.Quiet {
			fmt.Printf("Successfully loaded image to Kind cluster: %s\n", opts.KindCluster)
		}
	} else {
		if err := pushToRegistry(ctx, imageRef, opts); err != nil {
			return fmt.Errorf("pushing to registry: %w", err)
		}
		if !opts.Quiet {
			fmt.Printf("Successfully pushed image to registry: %s\n", opts.RegistryURL)
		}
	}

	return nil
}

// validatePushArgs validates the command arguments
func validatePushArgs(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("exactly one build path is required")
	}

	buildPath := args[0]
	if buildPath == "" {
		return fmt.Errorf("build path cannot be empty")
	}

	// Check if path exists
	if _, err := os.Stat(buildPath); os.IsNotExist(err) {
		return fmt.Errorf("build path does not exist: %s", buildPath)
	}

	return nil
}

// buildImage builds the container image
func buildImage(ctx context.Context, opts *PushOptions) (string, error) {
	if opts.Verbose {
		fmt.Printf("Building image from: %s\n", opts.BuildPath)
	}

	// Construct full image reference
	imageRef := opts.ImageName
	if opts.ImageTag != "" {
		imageRef += ":" + opts.ImageTag
	}

	// For now, simulate the build process
	// In a real implementation, this would use docker/buildkit APIs
	if opts.Verbose {
		fmt.Printf("Build context: %s\n", filepath.Join(opts.BuildPath, opts.Context))
		fmt.Printf("Dockerfile: %s\n", filepath.Join(opts.BuildPath, opts.Dockerfile))
		fmt.Printf("Target: %s\n", opts.Target)
		fmt.Printf("Platform: %s\n", opts.Platform)
		if len(opts.BuildArgs) > 0 {
			fmt.Println("Build args:")
			for key, value := range opts.BuildArgs {
				fmt.Printf("  %s=%s\n", key, value)
			}
		}
	}

	if opts.Verbose {
		fmt.Printf("Build completed for image: %s\n", imageRef)
	}

	return imageRef, nil
}

// pushToRegistry pushes the image to a container registry
func pushToRegistry(ctx context.Context, image string, opts *PushOptions) error {
	if opts.Verbose {
		fmt.Printf("Pushing image to registry: %s\n", opts.RegistryURL)
	}

	// Create authentication config
	var auth *oci.RegistryAuth
	if opts.RegistryUsername != "" || opts.RegistryPassword != "" || opts.RegistryInsecure {
		auth = &oci.RegistryAuth{
			Username:      opts.RegistryUsername,
			Password:      opts.RegistryPassword,
			ServerAddress: opts.RegistryURL,
		}
	}

	// Create push options using the existing OCI types
	pushOpts := &oci.PushOptions{
		ImageRef: image,
		Auth:     auth,
		Insecure: opts.RegistryInsecure,
		Context:  ctx,
	}

	// Create and execute push flow
	flow, err := oci.NewPushFlow(pushOpts)
	if err != nil {
		return fmt.Errorf("creating push flow: %w", err)
	}

	// Execute the push
	if err := flow.Execute(); err != nil {
		return fmt.Errorf("push failed: %w", err)
	}

	return nil
}

// loadToKind loads the image into a Kind cluster
func loadToKind(ctx context.Context, image string, opts *PushOptions) error {
	if opts.Verbose {
		fmt.Printf("Loading image to Kind cluster: %s\n", opts.KindCluster)
	}

	// Create Kind cluster manager
	config := &kind.ClusterConfig{
		Name:        opts.KindCluster,
		KubeVersion: "v1.28.0", // Default version
	}

	manager, err := kind.NewClusterManager(config)
	if err != nil {
		return fmt.Errorf("creating cluster manager: %w", err)
	}

	// Load the image
	if err := manager.Load(ctx, []string{image}); err != nil {
		return fmt.Errorf("loading image to Kind: %w", err)
	}

	return nil
}

// Helper functions

// getOptionsFromFlags extracts options from command flags
func getOptionsFromFlags(cmd *cobra.Command, args []string) (*PushOptions, error) {
	opts := &PushOptions{
		BuildPath: args[0],
	}

	// Extract flag values
	var err error
	if opts.Dockerfile, err = cmd.Flags().GetString("dockerfile"); err != nil {
		return nil, err
	}
	if opts.Context, err = cmd.Flags().GetString("context"); err != nil {
		return nil, err
	}
	if opts.Target, err = cmd.Flags().GetString("target"); err != nil {
		return nil, err
	}
	if opts.Platform, err = cmd.Flags().GetString("platform"); err != nil {
		return nil, err
	}
	if opts.BuildArgs, err = cmd.Flags().GetStringToString("build-arg"); err != nil {
		return nil, err
	}
	if opts.ImageName, err = cmd.Flags().GetString("name"); err != nil {
		return nil, err
	}
	if opts.ImageTag, err = cmd.Flags().GetString("tag"); err != nil {
		return nil, err
	}
	if opts.RegistryURL, err = cmd.Flags().GetString("registry"); err != nil {
		return nil, err
	}
	if opts.RegistryUsername, err = cmd.Flags().GetString("username"); err != nil {
		return nil, err
	}
	if opts.RegistryPassword, err = cmd.Flags().GetString("password"); err != nil {
		return nil, err
	}
	if opts.RegistryInsecure, err = cmd.Flags().GetBool("insecure"); err != nil {
		return nil, err
	}
	if opts.PushToKind, err = cmd.Flags().GetBool("kind"); err != nil {
		return nil, err
	}
	if opts.KindCluster, err = cmd.Flags().GetString("kind-cluster"); err != nil {
		return nil, err
	}
	if opts.Force, err = cmd.Flags().GetBool("force"); err != nil {
		return nil, err
	}
	if opts.DryRun, err = cmd.Flags().GetBool("dry-run"); err != nil {
		return nil, err
	}
	if opts.Verbose, err = cmd.Flags().GetBool("verbose"); err != nil {
		return nil, err
	}
	if opts.Quiet, err = cmd.Flags().GetBool("quiet"); err != nil {
		return nil, err
	}

	// Use environment variables as fallbacks
	if opts.RegistryURL == "" {
		opts.RegistryURL = os.Getenv("REGISTRY_URL")
	}
	if opts.RegistryUsername == "" {
		opts.RegistryUsername = os.Getenv("REGISTRY_USERNAME")
	}
	if opts.RegistryPassword == "" {
		opts.RegistryPassword = os.Getenv("REGISTRY_PASSWORD")
	}

	return opts, nil
}

// validate validates the push options
func (opts *PushOptions) validate() error {
	if opts.ImageName == "" {
		return fmt.Errorf("image name is required")
	}

	// Validate registry URL if pushing to registry
	if !opts.PushToKind {
		if opts.RegistryURL == "" {
			return fmt.Errorf("registry URL is required when pushing to registry")
		}
	}

	// Validate Kind cluster if pushing to Kind
	if opts.PushToKind && opts.KindCluster == "" {
		return fmt.Errorf("Kind cluster name is required when pushing to Kind")
	}

	// Validate build path
	if opts.BuildPath == "" {
		return fmt.Errorf("build path is required")
	}

	if _, err := os.Stat(opts.BuildPath); os.IsNotExist(err) {
		return fmt.Errorf("build path does not exist: %s", opts.BuildPath)
	}

	return nil
}

// printDryRun prints what would be done in dry-run mode
func printDryRun(opts *PushOptions) error {
	fmt.Println("Dry run mode - showing what would be executed:")
	fmt.Printf("Build path: %s\n", opts.BuildPath)
	fmt.Printf("Dockerfile: %s\n", opts.Dockerfile)
	fmt.Printf("Context: %s\n", opts.Context)
	fmt.Printf("Image name: %s\n", opts.ImageName)
	fmt.Printf("Image tag: %s\n", opts.ImageTag)
	fmt.Printf("Platform: %s\n", opts.Platform)

	if len(opts.BuildArgs) > 0 {
		fmt.Println("Build args:")
		for key, value := range opts.BuildArgs {
			fmt.Printf("  %s=%s\n", key, value)
		}
	}

	if opts.PushToKind {
		fmt.Printf("Action: Load image to Kind cluster '%s'\n", opts.KindCluster)
	} else {
		fmt.Printf("Action: Push image to registry '%s'\n", opts.RegistryURL)
		if opts.RegistryUsername != "" {
			fmt.Printf("Registry username: %s\n", opts.RegistryUsername)
		}
	}

	return nil
}