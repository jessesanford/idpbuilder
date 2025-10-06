package push

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// PushOptions holds all configuration for the push command
type PushOptions struct {
	// Registry configuration
	RegistryURL  string
	Repository   string
	Tag          string

	// Authentication
	Username     string
	Password     string

	// TLS configuration
	Insecure     bool

	// Image source
	ImagePath    string
	ImageRef     string

	// Behavior
	DryRun       bool
	Verbose      bool
}

var PushCmd = &cobra.Command{
	Use:   "push [IMAGE] [REGISTRY_URL]",
	Short: "Push OCI artifacts to a registry",
	Long: `Push OCI artifacts to any OCI-compliant registry.

The push command supports authentication via flags or environment variables:
  - Authentication: --username/--password or REGISTRY_USERNAME/REGISTRY_PASSWORD
  - TLS: --insecure flag for self-signed certificates

Examples:
  # Push with authentication via flags
  idpbuilder push myimage:latest registry.example.com/repo --username user --password pass

  # Push with environment variables
  export REGISTRY_USERNAME=user
  export REGISTRY_PASSWORD=pass
  idpbuilder push myimage:latest registry.example.com/repo

  # Push to insecure registry (self-signed cert)
  idpbuilder push myimage:latest registry.example.com/repo --insecure`,
	Args:         cobra.RangeArgs(1, 2),
	RunE:         runPush,
	SilenceUsage: true,
}

func runPush(cmd *cobra.Command, args []string) error {
	opts, err := buildPushOptions(cmd, args)
	if err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}

	// Validate options
	if err := validatePushOptions(opts); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Log configuration
	log.Printf("INFO: Starting push operation - image=%s registry=%s insecure=%t dry-run=%t",
		opts.ImageRef, opts.RegistryURL, opts.Insecure, opts.DryRun)

	// Warning for insecure mode
	if opts.Insecure {
		fmt.Fprintf(os.Stderr, "⚠️  WARNING: TLS certificate verification disabled\n")
		fmt.Fprintf(os.Stderr, "   Only use with self-signed certificates in development\n\n")
	}

	// TODO: In E1.2.3, implement actual push logic here
	// For now, just log the configuration
	log.Printf("INFO: Push configuration validated successfully")

	if opts.DryRun {
		fmt.Printf("DRY RUN: Would push %s to %s\n", opts.ImageRef, opts.RegistryURL)
		return nil
	}

	fmt.Printf("✅ Push command structure ready (implementation pending E1.2.3)\n")
	return nil
}