package push

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/cnoe-io/idpbuilder/pkg/docker"
	"github.com/cnoe-io/idpbuilder/pkg/registry"
	"github.com/cnoe-io/idpbuilder/pkg/auth"
	"github.com/cnoe-io/idpbuilder/pkg/tls"
	"github.com/cnoe-io/idpbuilder/pkg/progress"
)

// NewPushCommand creates the push command that integrates all Phase 1 packages
// with environment variable support via Viper
func NewPushCommand(v *viper.Viper) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push IMAGE",
		Short: "Push a Docker image to an OCI registry",
		Long: `Push a local Docker image to an OCI-compliant container registry.

The command retrieves the image from the local Docker daemon and pushes it to
the specified registry using credentials provided via flags or environment variables.

Configuration precedence: Flags > Environment Variables > Defaults

Environment Variables:
  IDPBUILDER_REGISTRY   Override registry URL (default: gitea.cnoe.localtest.me:8443)
  IDPBUILDER_USERNAME   Registry username (required if not provided via flag)
  IDPBUILDER_PASSWORD   Registry password (required if not provided via flag)
  IDPBUILDER_INSECURE   Skip TLS verification (true/false, 1/0, yes/no)
  IDPBUILDER_VERBOSE    Enable verbose output (true/false, 1/0, yes/no)

Examples:
  # Push using flags only (Wave 2.1 compatibility)
  idpbuilder push alpine:latest --username admin --password password

  # Push using environment variables
  export IDPBUILDER_USERNAME=admin
  export IDPBUILDER_PASSWORD=password
  idpbuilder push alpine:latest

  # Mix flags and environment variables (flags take precedence)
  export IDPBUILDER_REGISTRY=docker.io
  idpbuilder push alpine:latest --username admin --password password

  # Push to custom registry with verbose output
  idpbuilder push myapp:v1.0 --registry docker.io --username user --password pass --verbose

  # Push with insecure TLS (development only)
  idpbuilder push alpine:latest --insecure --username admin --password password`,
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

			// Convert to PushOptions for Wave 2.1 compatibility
			opts := config.ToPushOptions()

			// Call runPush unchanged
			return runPush(cmd.Context(), opts)
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

// runPush orchestrates the 8-stage push pipeline using Phase 1 interfaces
func runPush(ctx context.Context, opts *PushOptions) error {
	// Validate options first
	if err := opts.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// STAGE 1: Initialize Docker client
	dockerClient, err := docker.NewClient()
	if err != nil {
		return fmt.Errorf("failed to connect to Docker daemon: %w", err)
	}
	defer dockerClient.Close()

	// STAGE 2: Retrieve image from Docker daemon
	if opts.Verbose {
		fmt.Fprintf(os.Stdout, "Retrieving image %s from Docker daemon...\n", opts.ImageName)
	}
	image, err := dockerClient.GetImage(ctx, opts.ImageName)
	if err != nil {
		return fmt.Errorf("failed to get image: %w", err)
	}

	// STAGE 3: Setup authentication
	authProvider := auth.NewBasicAuthProvider(opts.Username, opts.Password)
	if err := authProvider.ValidateCredentials(); err != nil {
		return fmt.Errorf("invalid credentials: %w", err)
	}

	// STAGE 4: Setup TLS configuration
	tlsProvider := tls.NewConfigProvider(opts.Insecure)

	// STAGE 5: Create registry client
	registryClient, err := registry.NewClient(authProvider, tlsProvider)
	if err != nil {
		return fmt.Errorf("failed to create registry client: %w", err)
	}

	// STAGE 6: Build target reference
	targetRef, err := registryClient.BuildImageReference(opts.Registry, opts.ImageName)
	if err != nil {
		return fmt.Errorf("invalid registry or image name: %w", err)
	}

	// STAGE 7: Create progress reporter (replaces basic callback)
	reporter := progress.NewReporter(opts.Verbose)

	// STAGE 8: Execute push with reporter
	fmt.Fprintf(os.Stdout, "Pushing to %s...\n", targetRef)
	if err := registryClient.Push(ctx, image, targetRef, reporter.GetCallback()); err != nil {
		return fmt.Errorf("push failed: %w", err)
	}

	// Display final summary
	reporter.DisplaySummary()

	fmt.Fprintf(os.Stdout, "✓ Successfully pushed %s to %s\n", opts.ImageName, opts.Registry)
	return nil
}

// RunPushForTesting exposes runPush for testing purposes
func RunPushForTesting(ctx context.Context, opts *PushOptions) error {
	return runPush(ctx, opts)
}
