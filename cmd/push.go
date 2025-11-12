// Package cmd provides the CLI commands for idpbuilder.
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/cnoe-io/idpbuilder/pkg/auth"
	"github.com/cnoe-io/idpbuilder/pkg/docker"
	"github.com/cnoe-io/idpbuilder/pkg/registry"
	"github.com/cnoe-io/idpbuilder/pkg/tls"
	"github.com/spf13/cobra"
)

// PushCommand orchestrates the image push workflow.
// It coordinates Docker, Registry, Auth, and TLS packages to push an image to a registry.
type PushCommand struct {
	// dockerClient provides access to the local Docker daemon.
	dockerClient docker.DockerClient

	// registryClient handles OCI registry operations.
	registryClient registry.RegistryClient

	// authProvider supplies authentication credentials.
	authProvider auth.AuthProvider

	// tlsProvider generates TLS configurations.
	tlsProvider tls.TLSProvider
}

// PushFlags holds command-line flag values for the push command.
type PushFlags struct {
	// Registry is the target registry URL (e.g., "https://gitea.cnoe.localtest.me:8443")
	Registry string

	// Username for registry authentication
	Username string

	// Password for registry authentication
	Password string

	// Insecure enables TLS certificate verification bypass (for self-signed certs)
	Insecure bool

	// Verbose enables detailed progress output
	Verbose bool
}

// NewPushCommand creates and configures the push command.
//
// Example registration with Cobra:
//
//	rootCmd.AddCommand(cmd.NewPushCommand())
func NewPushCommand() *cobra.Command {
	flags := &PushFlags{}

	cmd := &cobra.Command{
		Use:   "push IMAGE_NAME",
		Short: "Push a container image to the IDPBuilder registry",
		Long: `Push a container image to the IDPBuilder Gitea registry.

The image must exist in your local Docker daemon. Use 'docker images' to list available images.

Examples:
  # Push with default registry
  idpbuilder push myapp:latest --password mypassword

  # Push to custom registry
  idpbuilder push myapp:v1.0 --registry https://registry.example.com --username admin --password secret

  # Push with insecure TLS (for self-signed certificates)
  idpbuilder push myapp:dev --insecure --password mypassword

  # Push with verbose output
  idpbuilder push myapp:latest --verbose --password mypassword
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			imageName := args[0]
			return runPush(cmd.Context(), imageName, flags)
		},
	}

	// Define flags
	cmd.Flags().StringVar(&flags.Registry, "registry", "https://gitea.cnoe.localtest.me:8443",
		"Target registry URL")
	cmd.Flags().StringVar(&flags.Username, "username", "giteaadmin",
		"Registry username (default: giteaadmin)")
	cmd.Flags().StringVar(&flags.Password, "password", "",
		"Registry password (required)")
	cmd.Flags().BoolVar(&flags.Insecure, "insecure", false,
		"Skip TLS certificate verification (use only for self-signed certs)")
	cmd.Flags().BoolVar(&flags.Verbose, "verbose", false,
		"Enable verbose output")

	// Mark password as required
	cmd.MarkFlagRequired("password")

	return cmd
}

// runPush executes the push workflow.
// This is the main orchestration function that coordinates all packages.
func runPush(ctx context.Context, imageName string, flags *PushFlags) error {
	// Implementation will be provided in Phase 2 Wave 1
	// This skeleton shows the structure and interface usage
	fmt.Fprintf(os.Stderr, "Push command not yet implemented\n")
	fmt.Fprintf(os.Stderr, "Will push %s to %s\n", imageName, flags.Registry)
	return nil
}

// Exit codes for the push command
const (
	ExitSuccess      = 0 // Push successful
	ExitGeneralError = 1 // Invalid arguments, unexpected failures
	ExitAuthError    = 2 // Authentication failure
	ExitNetworkError = 3 // Registry unreachable, TLS failure
	ExitImageNotFound = 4 // Image not in Docker daemon
)
