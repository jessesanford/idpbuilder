// Package push implements the push command for idpbuilder-push CLI.
// This command enables users to push OCI images to container registries
// with support for authentication and insecure connections.
package push

import (
	"context"
	"fmt"
	"time"

	"github.com/cnoe-io/idpbuilder-push/client-interface-tests-split-003/pkg/oci"
	"github.com/spf13/cobra"
)

// PushCommand represents the push command configuration and execution logic.
type PushCommand struct {
	// CLI flag values
	username string
	password string
	insecure bool
	timeout  string

	// Parsed options
	timeoutDuration time.Duration
}

// NewPushCommand creates and configures the push command with all necessary flags and validation.
func NewPushCommand() *cobra.Command {
	pushCmd := &PushCommand{}

	cmd := &cobra.Command{
		Use:   "push [image] [registry]",
		Short: "Push an OCI image to a container registry",
		Long: `Push an OCI image to a container registry with authentication support.

This command uploads container images to OCI-compliant registries such as
Docker Hub, GitLab Container Registry, or Gitea. It supports both username/password
authentication and can handle insecure (HTTP) connections for local registries.

Examples:
  # Push image with authentication
  idpbuilder push myapp:latest registry.example.com --username user --password pass

  # Push to insecure registry (HTTP)
  idpbuilder push myapp:latest localhost:5000 --insecure

  # Push with timeout
  idpbuilder push myapp:latest registry.example.com --timeout 5m
`,
		Args:    cobra.ExactArgs(2),
		PreRunE: pushCmd.validateArgs,
		RunE:    pushCmd.execute,
	}

	// Add command flags
	cmd.Flags().StringVarP(&pushCmd.username, "username", "u", "",
		"Username for registry authentication")
	cmd.Flags().StringVarP(&pushCmd.password, "password", "p", "",
		"Password for registry authentication")
	cmd.Flags().BoolVar(&pushCmd.insecure, "insecure", false,
		"Allow connections to insecure (HTTP) registries")
	cmd.Flags().StringVar(&pushCmd.timeout, "timeout", "5m",
		"Timeout for push operation (e.g., 5m, 30s)")

	return cmd
}

// validateArgs performs validation of command arguments and flags before execution.
func (p *PushCommand) validateArgs(cmd *cobra.Command, args []string) error {
	// Validate timeout format
	if p.timeout != "" {
		duration, err := time.ParseDuration(p.timeout)
		if err != nil {
			return fmt.Errorf("invalid timeout format '%s': %w (use formats like 5m, 30s)", p.timeout, err)
		}
		if duration <= 0 {
			return fmt.Errorf("timeout must be positive, got %v", duration)
		}
		p.timeoutDuration = duration
	} else {
		p.timeoutDuration = 5 * time.Minute // Default timeout
	}

	// Validate image reference (basic format check)
	imageRef := args[0]
	if imageRef == "" {
		return fmt.Errorf("image reference cannot be empty")
	}

	// Validate registry address
	registry := args[1]
	if registry == "" {
		return fmt.Errorf("registry address cannot be empty")
	}

	// Check for authentication consistency
	if (p.username != "" && p.password == "") || (p.username == "" && p.password != "") {
		return fmt.Errorf("both username and password must be provided for authentication")
	}

	return nil
}

// execute performs the actual push operation using the OCI library.
func (p *PushCommand) execute(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeoutDuration)
	defer cancel()

	imageRef := args[0]
	registry := args[1]

	// Build the full image reference
	fullImageRef := fmt.Sprintf("%s/%s", registry, imageRef)

	// Create push options
	pushOpts := &oci.PushOptions{
		ImageRef:  fullImageRef,
		Registry:  registry,
		Insecure:  p.insecure,
		Context:   ctx,
		Timeout:   p.timeoutDuration,
	}

	// Add authentication if provided
	if p.username != "" && p.password != "" {
		pushOpts.Auth = &oci.RegistryAuth{
			Username:      p.username,
			Password:      p.password,
			ServerAddress: registry,
		}
	}

	// Validate the push options
	if err := pushOpts.Validate(); err != nil {
		return fmt.Errorf("invalid push configuration: %w", err)
	}

	// Create authenticator if credentials are provided
	var authenticator oci.Authenticator
	if pushOpts.Auth != nil && !pushOpts.Auth.IsEmpty() {
		auth, err := oci.CreateAuthenticatorFromConfig(pushOpts.Auth)
		if err != nil {
			return fmt.Errorf("failed to create authenticator: %w", err)
		}

		if err := oci.ValidateAuthenticator(auth); err != nil {
			return fmt.Errorf("authenticator validation failed: %w", err)
		}

		authenticator = auth
	}

	// Display operation info
	fmt.Printf("Pushing image: %s\n", fullImageRef)
	fmt.Printf("Registry: %s\n", registry)
	if authenticator != nil {
		fmt.Printf("Authentication: %s\n", authenticator.GetType())
	} else {
		fmt.Printf("Authentication: none (anonymous)\n")
	}
	if p.insecure {
		fmt.Printf("Security: insecure (HTTP) connection\n")
	} else {
		fmt.Printf("Security: secure (HTTPS) connection\n")
	}
	fmt.Printf("Timeout: %v\n", p.timeoutDuration)

	// For now, this is a client interface implementation
	// The actual push logic will be implemented in a future split
	// This demonstrates the command interface and validation
	fmt.Printf("\n✓ Command interface validation successful\n")
	fmt.Printf("✓ Authentication setup complete\n")
	fmt.Printf("✓ Push options validated\n")
	fmt.Printf("\nNote: Actual push implementation will be added in the next phase\n")

	return nil
}

// GetPushOptions creates PushOptions from command line arguments and flags.
// This method is exposed for testing purposes.
func (p *PushCommand) GetPushOptions(imageRef, registry string) *oci.PushOptions {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeoutDuration)
	defer cancel()

	opts := &oci.PushOptions{
		ImageRef: fmt.Sprintf("%s/%s", registry, imageRef),
		Registry: registry,
		Insecure: p.insecure,
		Context:  ctx,
		Timeout:  p.timeoutDuration,
	}

	if p.username != "" && p.password != "" {
		opts.Auth = &oci.RegistryAuth{
			Username:      p.username,
			Password:      p.password,
			ServerAddress: registry,
		}
	}

	return opts
}