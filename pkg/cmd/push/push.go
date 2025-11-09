package push

import (
	"context"
	"fmt"
	"os"

	"github.com/cnoe-io/idpbuilder/pkg/auth"
	"github.com/cnoe-io/idpbuilder/pkg/docker"
	"github.com/cnoe-io/idpbuilder/pkg/errors"
	"github.com/cnoe-io/idpbuilder/pkg/registry"
	"github.com/cnoe-io/idpbuilder/pkg/tls"
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
				// Format error with appropriate styling
				fmt.Fprintln(os.Stderr, errors.FormatError(err))

				// In test mode, return the error instead of exiting
				// This allows tests to verify error handling behavior
				if os.Getenv("IDPBUILDER_TEST_MODE") == "true" || os.Getenv("IDPBUILDER_TEST_MODE") == "1" {
					return err
				}

				// Get exit code for this error type and exit
				// (production behavior only, not in tests)
				exitCode := errors.GetExitCode(err)
				os.Exit(exitCode)
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
	// Validate options with security checks (Phase 2 Wave 3)
	if err := validatePushOptions(opts); err != nil {
		// validatePushOptions returns typed errors, pass through
		return err
	}

	// Create Docker client and retrieve image
	dockerClient, err := docker.NewClient()
	if err != nil {
		return WrapDockerError(err, opts.ImageName)
	}
	defer dockerClient.Close()

	// Get image from local Docker daemon
	image, err := dockerClient.GetImage(ctx, opts.ImageName)
	if err != nil {
		return WrapDockerError(err, opts.ImageName)
	}

	// Create authentication provider
	authProvider := auth.NewBasicAuthProvider(opts.Username, opts.Password)

	// Create TLS configuration provider
	tlsProvider := tls.NewConfigProvider(opts.Insecure)

	// Create registry client
	registryClient, err := registry.NewClient(authProvider, tlsProvider)
	if err != nil {
		return WrapRegistryError(err, opts.Registry)
	}

	// Build fully qualified image reference
	targetRef, err := registryClient.BuildImageReference(opts.Registry, opts.ImageName)
	if err != nil {
		return WrapRegistryError(err, opts.Registry)
	}

	// Display verbose output if enabled
	if opts.Verbose {
		fmt.Printf("Pushing image %s to %s\n", opts.ImageName, targetRef)
	}

	// Push image to registry with optional progress callback
	var progressCallback registry.ProgressCallback
	if opts.Verbose {
		progressCallback = func(update registry.ProgressUpdate) {
			fmt.Printf("  %s: %d/%d bytes (%s)\n",
				update.Status, update.BytesPushed, update.LayerSize, update.LayerDigest)
		}
	}

	err = registryClient.Push(ctx, image, targetRef, progressCallback)
	if err != nil {
		return WrapRegistryError(err, opts.Registry)
	}

	// Success message
	if opts.Verbose {
		fmt.Printf("Successfully pushed %s to %s\n", opts.ImageName, targetRef)
	}

	return nil
}

// RunPushForTesting exposes runPush for testing purposes
func RunPushForTesting(ctx context.Context, opts *PushOptions) error {
	return runPush(ctx, opts)
}
