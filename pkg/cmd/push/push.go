package push

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewPushCommand creates the push command with Viper integration for environment variable support
// This is a Phase 2 Wave 2 focus on configuration management only
func NewPushCommand(v *viper.Viper) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push IMAGE",
		Short: "Push a Docker image to an OCI registry",
		Long: `Push a local Docker image to an OCI-compliant container registry.

Configuration precedence: Flags > Environment Variables > Defaults

Environment Variables:
  IDPBUILDER_REGISTRY   Override registry URL (default: gitea.cnoe.localtest.me:8443)
  IDPBUILDER_USERNAME   Registry username (required if not provided via flag)
  IDPBUILDER_PASSWORD   Registry password (required if not provided via flag)
  IDPBUILDER_INSECURE   Skip TLS verification (true/false, 1/0, yes/no)
  IDPBUILDER_VERBOSE    Enable verbose output (true/false, 1/0, yes/no)

Examples:
  # Push using flags only
  idpbuilder push alpine:latest --username admin --password password

  # Push using environment variables
  export IDPBUILDER_USERNAME=admin
  export IDPBUILDER_PASSWORD=password
  idpbuilder push alpine:latest

  # Mix flags and environment variables (flags take precedence)
  export IDPBUILDER_REGISTRY=docker.io
  idpbuilder push alpine:latest --username admin --password password`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load configuration from flags, environment, and defaults
			config, err := LoadConfig(cmd, args, v)
			if err != nil {
				return fmt.Errorf("configuration error: %w", err)
			}

			// Validate configuration
			if err := config.Validate(); err != nil {
				return fmt.Errorf("validation error: %w", err)
			}

			// Display configuration sources in verbose mode
			if config.Verbose.Value == "true" {
				config.DisplaySources()
				fmt.Println()
			}

			// Convert to PushOptions for compatibility
			opts := config.ToPushOptions()

			// Execute push with error wrapping and exit code handling
			if err := runPush(cmd.Context(), opts); err != nil {
				// Return error to Cobra for proper handling
				// In CLI context, Cobra will exit with appropriate code
				// In test context, error is returned for assertion
				return err
			}

			return nil
		},
	}

	// Define flags with environment variable hints
	cmd.Flags().String("registry", DefaultRegistry,
		fmt.Sprintf("Registry URL (env: %s)", EnvRegistry))
	cmd.Flags().String("username", "",
		fmt.Sprintf("Registry username (env: %s, required)", EnvUsername))
	cmd.Flags().String("password", "",
		fmt.Sprintf("Registry password (env: %s, required)", EnvPassword))
	cmd.Flags().BoolP("insecure", "k", false,
		fmt.Sprintf("Skip TLS certificate verification (env: %s)", EnvInsecure))
	cmd.Flags().Bool("verbose", false,
		fmt.Sprintf("Enable verbose progress output (env: %s)", EnvVerbose))

	return cmd
}

// runPush executes the push operation with comprehensive error handling.
// Errors are wrapped with appropriate types (ValidationError, AuthenticationError,
// NetworkError, ImageNotFoundError) to enable proper exit code mapping.
func runPush(ctx context.Context, opts *PushOptions) error {
	// Validate options first (Wave 2.1 requirement)
	if err := opts.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// TODO: Phase 3 will integrate actual push implementation using:
	// - Docker client (pkg/docker) - wrap errors with WrapDockerError()
	// - Registry client (pkg/registry) - wrap errors with WrapRegistryError()
	// - Auth provider (pkg/auth)
	// - TLS config (pkg/tls)
	// - Progress reporting (pkg/progress)
	//
	// Example integration:
	//   dockerClient, err := docker.NewClient()
	//   if err != nil {
	//       return WrapDockerError(err, opts.ImageName)
	//   }
	//   defer dockerClient.Close()
	//
	//   image, err := dockerClient.GetImage(ctx, opts.ImageName)
	//   if err != nil {
	//       return WrapDockerError(err, opts.ImageName)
	//   }
	//
	//   registryClient, err := registry.NewClient(authProvider, tlsProvider)
	//   if err != nil {
	//       return WrapRegistryError(err, opts.Registry)
	//   }
	//
	//   err = registryClient.Push(ctx, image, targetRef, progressCallback)
	//   if err != nil {
	//       return WrapRegistryError(err, opts.Registry)
	//   }

	return fmt.Errorf("push implementation pending Phase 3 integration")
}

// RunPushForTesting exposes runPush for testing purposes
func RunPushForTesting(ctx context.Context, opts *PushOptions) error {
	return runPush(ctx, opts)
}
