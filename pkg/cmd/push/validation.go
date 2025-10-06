package push

import (
	"fmt"
	"net/url"
	"strings"
)

// validatePushOptions validates all push command options
func validatePushOptions(opts *PushOptions) error {
	// Validate required fields
	if opts.ImageRef == "" {
		return fmt.Errorf("image reference is required")
	}

	if opts.RegistryURL == "" {
		return fmt.Errorf("registry URL is required")
	}

	// Validate registry URL format
	if err := validateRegistryURL(opts.RegistryURL); err != nil {
		return fmt.Errorf("invalid registry URL: %w", err)
	}

	// Validate image reference format
	if err := validateImageRef(opts.ImageRef); err != nil {
		return fmt.Errorf("invalid image reference: %w", err)
	}

	// Validate authentication if registry requires it
	if requiresAuth(opts.RegistryURL) {
		if opts.Username == "" || opts.Password == "" {
			return fmt.Errorf("authentication required: provide --username and --password or set REGISTRY_USERNAME and REGISTRY_PASSWORD")
		}
	}

	return nil
}

// validateRegistryURL checks if the registry URL is valid
func validateRegistryURL(registryURL string) error {
	// Add scheme if missing
	if !strings.Contains(registryURL, "://") {
		registryURL = "https://" + registryURL
	}

	u, err := url.Parse(registryURL)
	if err != nil {
		return err
	}

	if u.Host == "" {
		return fmt.Errorf("registry host cannot be empty")
	}

	// Check for valid schemes
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("unsupported scheme: %s (use http or https)", u.Scheme)
	}

	return nil
}

// validateImageRef validates the image reference format
func validateImageRef(imageRef string) error {
	if imageRef == "" {
		return fmt.Errorf("image reference cannot be empty")
	}

	// Basic validation - more comprehensive validation will be in E1.2.3
	// Check for invalid characters
	if strings.ContainsAny(imageRef, " \t\n") {
		return fmt.Errorf("image reference contains whitespace")
	}

	return nil
}

// requiresAuth checks if the registry requires authentication
func requiresAuth(registryURL string) bool {
	// Docker Hub and most registries require auth
	// Local registries (localhost, 127.0.0.1) typically don't
	if strings.Contains(registryURL, "localhost") ||
	   strings.Contains(registryURL, "127.0.0.1") {
		return false
	}
	return true
}