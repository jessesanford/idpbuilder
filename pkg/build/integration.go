package build

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Integration provides a wrapper around Phase 2 Wave 1 buildah functionality
type Integration struct {
	// Uses buildah command-line tool directly
}

// NewIntegration creates a new build integration instance
func NewIntegration() *Integration {
	return &Integration{}
}

// Build executes a container image build using buildah
func (i *Integration) Build(ctx context.Context, opts BuildOptions) error {
	// Validate options
	if err := i.validateOptions(opts); err != nil {
		return fmt.Errorf("invalid build options: %w", err)
	}

	// Check if buildah is available
	if _, err := exec.LookPath("buildah"); err != nil {
		return fmt.Errorf("buildah not found in PATH: %w", err)
	}

	// Build command arguments
	args := []string{"bud"}
	
	// Add tag if specified
	if opts.Tag != "" {
		args = append(args, "-t", opts.Tag)
	}
	
	// Add Dockerfile path
	args = append(args, "-f", opts.DockerfilePath)
	
	// Add insecure flag if needed
	if opts.Insecure {
		args = append(args, "--tls-verify=false")
	}
	
	// Add context directory
	args = append(args, opts.ContextDir)
	
	// Execute buildah command
	fmt.Printf("Building image with buildah...\n")
	fmt.Printf("  Context: %s\n", opts.ContextDir)
	fmt.Printf("  Dockerfile: %s\n", opts.DockerfilePath)
	fmt.Printf("  Tag: %s\n", opts.Tag)
	if opts.Insecure {
		fmt.Printf("  Mode: Insecure (skipping TLS verification)\n")
	}
	
	cmd := exec.CommandContext(ctx, "buildah", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("buildah build failed: %w", err)
	}
	
	fmt.Printf("✅ Build successful!\n")
	return nil
}

// validateOptions validates the build options
func (i *Integration) validateOptions(opts BuildOptions) error {
	if opts.Tag == "" {
		return fmt.Errorf("tag is required")
	}

	if opts.ContextDir == "" {
		opts.ContextDir = "."
	}

	if opts.DockerfilePath == "" {
		opts.DockerfilePath = "Dockerfile"
	}

	// Validate dockerfile exists
	dockerfilePath := filepath.Join(opts.ContextDir, opts.DockerfilePath)
	if !filepath.IsAbs(opts.DockerfilePath) {
		dockerfilePath = filepath.Join(opts.ContextDir, opts.DockerfilePath)
	} else {
		dockerfilePath = opts.DockerfilePath
	}

	// TODO: Add file existence check when implementing actual integration
	fmt.Printf("Would validate dockerfile exists at: %s\n", dockerfilePath)

	return nil
}

// GetCertificateConfig returns certificate configuration for secure builds
func (i *Integration) GetCertificateConfig(insecure bool) (map[string]interface{}, error) {
	// TODO: Integrate with Phase 2 Wave 1 certificate trust manager
	// This will provide certificate bundles for secure registry operations

	if insecure {
		return map[string]interface{}{
			"skip_verify": true,
		}, nil
	}

	// Placeholder for certificate configuration
	return map[string]interface{}{
		"certificate_bundle_path": "/etc/ssl/certs/ca-certificates.crt",
		"trust_manager_enabled":   true,
	}, nil
}
