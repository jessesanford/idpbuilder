package main

import (
	"github.com/cnoe-io/idpbuilder-push/client-interface-tests-split-004/pkg/cmd/push"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	config *PushConfig
)

// NewRootCommand creates the root command for the push CLI
func NewRootCommand() *cobra.Command {
	config = NewPushConfig()

	rootCmd := &cobra.Command{
		Use:   "push [PATH]",
		Short: "Build and push container images",
		Long: `Build and push container images to a registry or load them into a Kind cluster.

The push command builds a container image from the specified path and either:
- Pushes it to a container registry
- Loads it into a Kind cluster for local development`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config.BuildPath = args[0]
			return push.RunPush(cmd, args)
		},
	}

	// Registry flags
	rootCmd.Flags().StringVar(&config.RegistryURL, "registry", "", "Registry URL (can also use REGISTRY_URL env var)")
	rootCmd.Flags().StringVar(&config.RegistryUsername, "username", "", "Registry username")
	rootCmd.Flags().StringVar(&config.RegistryPassword, "password", "", "Registry password")
	rootCmd.Flags().BoolVar(&config.RegistryInsecure, "insecure", false, "Allow insecure registry connections")

	// Build flags
	rootCmd.Flags().StringVar(&config.Dockerfile, "dockerfile", "Dockerfile", "Path to Dockerfile")
	rootCmd.Flags().StringVar(&config.Context, "context", ".", "Build context path")
	rootCmd.Flags().StringVar(&config.Target, "target", "", "Target stage in multi-stage build")
	rootCmd.Flags().StringVar(&config.Platform, "platform", "linux/amd64", "Target platform")
	rootCmd.Flags().StringToStringVar(&config.BuildArgs, "build-arg", nil, "Build arguments")

	// Push flags
	rootCmd.Flags().StringVar(&config.ImageName, "name", "", "Image name (required)")
	rootCmd.Flags().StringVar(&config.ImageTag, "tag", "latest", "Image tag")
	rootCmd.Flags().BoolVar(&config.PushToKind, "kind", false, "Load image into Kind cluster instead of pushing to registry")
	rootCmd.Flags().StringVar(&config.KindCluster, "kind-cluster", "kind", "Kind cluster name")
	rootCmd.Flags().BoolVar(&config.Force, "force", false, "Force push even if image exists")
	rootCmd.Flags().BoolVar(&config.DryRun, "dry-run", false, "Show what would be done without executing")

	// Output flags
	rootCmd.Flags().BoolVarP(&config.Verbose, "verbose", "v", false, "Verbose output")
	rootCmd.Flags().BoolVarP(&config.Quiet, "quiet", "q", false, "Quiet output")

	// Mark required flags
	rootCmd.MarkFlagRequired("name")

	return rootCmd
}

// Execute executes the root command
func Execute() error {
	return NewRootCommand().Execute()
}

func init() {
	// Initialize viper for configuration
	viper.AutomaticEnv()
	viper.SetEnvPrefix("IDPBUILDER")
}