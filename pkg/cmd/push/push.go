package push

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/cnoe-io/idpbuilder/pkg/docker"
	"github.com/cnoe-io/idpbuilder/pkg/registry"
	"github.com/cnoe-io/idpbuilder/pkg/auth"
	"github.com/cnoe-io/idpbuilder/pkg/tls"
)

// NewPushCommand creates the push command that integrates all Phase 1 packages
func NewPushCommand() *cobra.Command {
	opts := &PushOptions{}

	cmd := &cobra.Command{
		Use:   "push IMAGE",
		Short: "Push a Docker image to an OCI registry",
		Long: `Push a local Docker image to an OCI-compliant container registry.

The command retrieves the image from the local Docker daemon and pushes it to
the specified registry using credentials provided via flags or environment variables.

Examples:
  # Push to default Gitea registry
  idpbuilder push alpine:latest --username admin --password password

  # Push to custom registry
  idpbuilder push myapp:v1.0 --registry docker.io --username user --password pass

  # Push with verbose progress
  idpbuilder push alpine:latest --verbose --username admin --password password

  # Push with insecure TLS (development only)
  idpbuilder push alpine:latest --insecure --username admin --password password`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ImageName = args[0]
			return runPush(cmd.Context(), opts)
		},
	}

	// Define flags
	cmd.Flags().StringVar(&opts.Registry, "registry", "gitea.cnoe.localtest.me:8443",
		"Registry URL (default: Gitea registry)")
	cmd.Flags().StringVar(&opts.Username, "username", "",
		"Registry username (required)")
	cmd.Flags().StringVar(&opts.Password, "password", "",
		"Registry password (required)")
	cmd.Flags().BoolVarP(&opts.Insecure, "insecure", "k", false,
		"Skip TLS certificate verification (insecure)")
	cmd.Flags().BoolVar(&opts.Verbose, "verbose", false,
		"Enable verbose progress output")

	// Mark required flags
	cmd.MarkFlagRequired("username")
	cmd.MarkFlagRequired("password")

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

	// STAGE 7: Create progress callback (basic implementation for Wave 1)
	progressCallback := func(update registry.ProgressUpdate) {
		if opts.Verbose {
			fmt.Fprintf(os.Stdout, "Layer %s: %d/%d bytes (%s)\n",
				truncateDigest(update.LayerDigest, 12),
				update.BytesPushed,
				update.LayerSize,
				update.Status)
		}
	}

	// STAGE 8: Execute push
	if opts.Verbose {
		fmt.Fprintf(os.Stdout, "Pushing to %s...\n", targetRef)
	}
	if err := registryClient.Push(ctx, image, targetRef, progressCallback); err != nil {
		return fmt.Errorf("push failed: %w", err)
	}

	fmt.Fprintf(os.Stdout, "✓ Successfully pushed %s to %s\n", opts.ImageName, opts.Registry)
	return nil
}

// truncateDigest truncates digest to specified length for display
func truncateDigest(digest string, length int) string {
	if len(digest) <= length {
		return digest
	}
	return digest[:length]
}

// RunPushForTesting exposes runPush for testing purposes
func RunPushForTesting(ctx context.Context, opts *PushOptions) error {
	return runPush(ctx, opts)
}

// TruncateDigestForTesting exposes truncateDigest for testing purposes
func TruncateDigestForTesting(digest string, length int) string {
	return truncateDigest(digest, length)
}
