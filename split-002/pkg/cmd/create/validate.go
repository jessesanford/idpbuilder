package create

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cnoe-io/idpbuilder/pkg/cmd/helpers"
)

// ValidateCreateSpec validates the creation specification
func ValidateCreateSpec(spec CreateSpec) error {
	// Validate packages
	for i, pkg := range spec.Packages {
		if err := validatePackageSpec(pkg, i); err != nil {
			return fmt.Errorf("invalid package at index %d: %w", i, err)
		}
	}

	// Validate secrets
	for i, secret := range spec.Secrets {
		if err := validateSecretSpec(secret, i); err != nil {
			return fmt.Errorf("invalid secret at index %d: %w", i, err)
		}
	}

	// Validate configurations
	for i, config := range spec.Configs {
		if err := validateConfigSpec(config, i); err != nil {
			return fmt.Errorf("invalid config at index %d: %w", i, err)
		}
	}

	return nil
}

// validatePackageSpec validates a package specification
func validatePackageSpec(pkg PackageSpec, index int) error {
	if err := helpers.ValidateName(pkg.Name); err != nil {
		return fmt.Errorf("invalid package name: %w", err)
	}

	if pkg.Version == "" {
		return fmt.Errorf("package version cannot be empty")
	}

	if err := helpers.ValidateVersion(pkg.Version); err != nil {
		return fmt.Errorf("invalid package version: %w", err)
	}

	if pkg.Source == "" {
		return fmt.Errorf("package source cannot be empty")
	}

	return validatePackageSource(pkg.Source)
}

// validateSecretSpec validates a secret specification
func validateSecretSpec(secret SecretSpec, index int) error {
	if err := helpers.ValidateName(secret.Name); err != nil {
		return fmt.Errorf("invalid secret name: %w", err)
	}

	validTypes := []string{"Opaque", "kubernetes.io/tls", "kubernetes.io/dockerconfigjson", "kubernetes.io/service-account-token"}
	if secret.Type == "" {
		secret.Type = "Opaque" // Default type
	}

	isValidType := false
	for _, validType := range validTypes {
		if secret.Type == validType {
			isValidType = true
			break
		}
	}
	
	if !isValidType {
		return fmt.Errorf("invalid secret type: %s (must be one of: %v)", secret.Type, validTypes)
	}

	if len(secret.Data) == 0 && len(secret.StringData) == 0 {
		return fmt.Errorf("secret must have either data or stringData")
	}

	return nil
}

// validateConfigSpec validates a configuration specification
func validateConfigSpec(config ConfigSpec, index int) error {
	if err := helpers.ValidateName(config.Name); err != nil {
		return fmt.Errorf("invalid config name: %w", err)
	}

	if len(config.Data) == 0 {
		return fmt.Errorf("config data cannot be empty")
	}

	return nil
}

// validatePackageSource validates a package source URL or path
func validatePackageSource(source string) error {
	if source == "" {
		return fmt.Errorf("package source cannot be empty")
	}

	// Check if it's a URL
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		return validatePackageURL(source)
	}

	// Check if it's a file path
	if strings.Contains(source, "/") || strings.Contains(source, "\\") {
		return validatePackagePath(source)
	}

	// Check if it's a registry reference (e.g., "nginx:latest")
	if strings.Contains(source, ":") {
		return validatePackageRegistry(source)
	}

	return fmt.Errorf("invalid package source format: %s", source)
}

// validatePackageURL validates a package URL
func validatePackageURL(url string) error {
	if !strings.HasSuffix(url, ".tar.gz") && !strings.HasSuffix(url, ".zip") {
		helpers.LogWarning("Package URL does not have expected archive extension (.tar.gz or .zip): %s", url)
	}
	return nil
}

// validatePackagePath validates a package file path
func validatePackagePath(path string) error {
	ext := filepath.Ext(path)
	validExts := []string{".tar.gz", ".tgz", ".zip", ".yaml", ".yml"}
	
	isValid := false
	for _, validExt := range validExts {
		if strings.HasSuffix(path, validExt) {
			isValid = true
			break
		}
	}
	
	if !isValid {
		return fmt.Errorf("package path must have valid extension (%v), got: %s", validExts, ext)
	}

	return nil
}

// validatePackageRegistry validates a package registry reference
func validatePackageRegistry(ref string) error {
	parts := strings.Split(ref, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid registry reference format (expected name:tag): %s", ref)
	}

	name, tag := parts[0], parts[1]
	
	if err := helpers.ValidateName(name); err != nil {
		return fmt.Errorf("invalid package name in registry reference: %w", err)
	}

	if tag == "" {
		return fmt.Errorf("package tag cannot be empty")
	}

	return nil
}