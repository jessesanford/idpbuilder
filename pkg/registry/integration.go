package registry

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Integration provides a wrapper around Phase 2 Wave 1 registry functionality
type Integration struct {
	// Uses buildah command-line tool for pushing
}

// NewIntegration creates a new registry integration instance
func NewIntegration() *Integration {
	return &Integration{}
}

// Push executes a container image push to a registry
func (i *Integration) Push(ctx context.Context, opts PushOptions) error {
	// Validate that we have an image and registry
	if opts.ImageID == "" {
		return fmt.Errorf("image ID is required")
	}
	if opts.Repository == "" {
		return fmt.Errorf("repository/registry is required")
	}

	// Check if buildah is available (we use it for pushing too)
	if _, err := exec.LookPath("buildah"); err != nil {
		return fmt.Errorf("buildah not found in PATH: %w", err)
	}

	// Build command arguments
	args := []string{"push"}
	
	// Add insecure flag if needed
	if opts.Insecure {
		args = append(args, "--tls-verify=false")
	}
	
	// Add credentials if provided
	if opts.Username != "" && opts.Password != "" {
		args = append(args, "--creds", fmt.Sprintf("%s:%s", opts.Username, opts.Password))
	}
	
	// Source image (local)
	args = append(args, opts.ImageID)
	
	// Destination (registry)
	// Format: docker://registry/namespace/image:tag
	// Extract just the image name and tag from the source
	imageParts := strings.Split(opts.ImageID, "/")
	shortName := imageParts[len(imageParts)-1] // Get last part (image:tag)
	
	destRef := fmt.Sprintf("docker://%s/%s", opts.Repository, shortName)
	args = append(args, destRef)
	
	// Execute buildah push command
	fmt.Printf("Pushing image to registry...\n")
	fmt.Printf("  Source: %s\n", opts.ImageID)
	fmt.Printf("  Destination: %s\n", destRef)
	if opts.Insecure {
		fmt.Printf("  Mode: Insecure (skipping TLS verification)\n")
	}
	
	cmd := exec.CommandContext(ctx, "buildah", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("buildah push failed: %w", err)
	}
	
	fmt.Printf("✅ Push successful!\n")
	return nil
}


// GetRegistryConfig returns registry configuration for secure operations
func (i *Integration) GetRegistryConfig(insecure bool) (map[string]interface{}, error) {
	// TODO: Integrate with Phase 2 Wave 1 certificate trust manager
	// This will provide certificate bundles for secure registry operations

	if insecure {
		return map[string]interface{}{
			"skip_tls_verify": true,
		}, nil
	}

	// Placeholder for registry configuration
	return map[string]interface{}{
		"certificate_bundle_path": "/etc/ssl/certs/ca-certificates.crt",
		"trust_manager_enabled":   true,
		"gitea_integration":       true,
	}, nil
}

// IsGiteaRegistry checks if the registry URL is a Gitea registry
func (i *Integration) IsGiteaRegistry(registryURL string) bool {
	// TODO: Implement Gitea registry detection logic
	// This might check for specific URL patterns or registry metadata
	return false
}

// GetGiteaAuth returns Gitea-specific authentication if available
func (i *Integration) GetGiteaAuth(ctx context.Context, registryURL string) (string, string, error) {
	// TODO: Integrate with Gitea authentication mechanisms
	// This might use service accounts, tokens, or other Gitea-specific auth
	return "", "", fmt.Errorf("gitea authentication not yet implemented")
}
