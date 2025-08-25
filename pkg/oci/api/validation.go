// Package api provides validation functions for OCI configuration types.
// This file implements comprehensive validation logic for build configurations,
// stack configurations, and OCI-specific formats.
package api

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	// imageTagRegex validates OCI image tag format
	imageTagRegex = regexp.MustCompile(`^[a-z0-9]([a-z0-9._-]*[a-z0-9])?$`)

	// semverRegex validates semantic versioning format
	semverRegex = regexp.MustCompile(`^v?(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)

	// platformRegex validates platform format (os/arch)
	platformRegex = regexp.MustCompile(`^[a-z]+/[a-z0-9]+$`)

	// Validator instance with custom validators registered
	validate *validator.Validate
)

// init registers custom validators
func init() {
	validate = validator.New()
	
	// Register custom tag validators
	validate.RegisterValidation("image_tag", validateImageTag)
	validate.RegisterValidation("semver", validateSemver)
	validate.RegisterValidation("platform", validatePlatform)
}

// validateImageTag validates OCI image tag format
func validateImageTag(fl validator.FieldLevel) bool {
	tag := fl.Field().String()
	if tag == "" {
		return false
	}
	
	// Tag must be 1-128 characters
	if len(tag) > 128 {
		return false
	}
	
	// Cannot start or end with separator
	if strings.HasPrefix(tag, ".") || strings.HasPrefix(tag, "_") || strings.HasPrefix(tag, "-") ||
		strings.HasSuffix(tag, ".") || strings.HasSuffix(tag, "_") || strings.HasSuffix(tag, "-") {
		return false
	}
	
	return imageTagRegex.MatchString(tag)
}

// validateSemver validates semantic versioning format
func validateSemver(fl validator.FieldLevel) bool {
	version := fl.Field().String()
	if version == "" {
		return false
	}
	
	return semverRegex.MatchString(version)
}

// validatePlatform validates platform format (os/arch)
func validatePlatform(fl validator.FieldLevel) bool {
	platform := fl.Field().String()
	if platform == "" {
		return true // Platform is optional in some contexts
	}
	
	return platformRegex.MatchString(platform)
}

// ValidateBuildConfig validates a BuildConfig struct with business logic
func ValidateBuildConfig(config *BuildConfig) error {
	if config == nil {
		return fmt.Errorf("build config cannot be nil")
	}
	
	// Validate struct tags
	if err := validate.Struct(config); err != nil {
		return fmt.Errorf("build config validation failed: %w", err)
	}
	
	// Business logic validation
	
	// Rootless mode requires VFS storage driver
	if config.Rootless && config.StorageDriver != "vfs" {
		return fmt.Errorf("rootless mode requires vfs storage driver, got %s", config.StorageDriver)
	}
	
	// Validate storage options for specific drivers
	if config.StorageDriver == "overlay" {
		if options := config.StorageOptions; options != nil {
			if val, exists := options["mountopt"]; exists {
				validMountOpts := []string{"nodev", "metacopy=on", "volatile"}
				found := false
				for _, validOpt := range validMountOpts {
					if strings.Contains(val, validOpt) {
						found = true
						break
					}
				}
				if !found {
					return fmt.Errorf("invalid overlay mountopt: %s", val)
				}
			}
		}
	}
	
	// Validate timeout ranges
	if config.BuildTimeout.Minutes() < 1 {
		return fmt.Errorf("build timeout must be at least 1 minute")
	}
	if config.BuildTimeout.Hours() > 24 {
		return fmt.Errorf("build timeout cannot exceed 24 hours")
	}
	
	return nil
}

// ValidateRegistryConfig validates a RegistryConfig struct
func ValidateRegistryConfig(config *RegistryConfig) error {
	if config == nil {
		return fmt.Errorf("registry config cannot be nil")
	}
	
	// Validate struct tags
	if err := validate.Struct(config); err != nil {
		return fmt.Errorf("registry config validation failed: %w", err)
	}
	
	// Business logic validation
	
	// Require authentication for non-localhost registries
	if !strings.Contains(config.URL, "localhost") && !strings.Contains(config.URL, "127.0.0.1") {
		if config.Username == "" && config.Token == "" {
			return fmt.Errorf("authentication required for non-localhost registry")
		}
	}
	
	// Validate mutual exclusion of authentication methods
	if config.Username != "" && config.Token != "" {
		return fmt.Errorf("cannot specify both username/password and token authentication")
	}
	
	// Username requires password
	if config.Username != "" && config.Password == "" {
		return fmt.Errorf("username authentication requires password")
	}
	
	// Validate timeout ranges
	if config.Timeout.Seconds() < 10 {
		return fmt.Errorf("registry timeout must be at least 10 seconds")
	}
	if config.Timeout.Minutes() > 30 {
		return fmt.Errorf("registry timeout cannot exceed 30 minutes")
	}
	
	// Validate retry configuration
	if config.RetryDelay.Seconds() < 1 {
		return fmt.Errorf("retry delay must be at least 1 second")
	}
	
	return nil
}

// ValidateStackConfig validates a StackOCIConfig struct
func ValidateStackConfig(config *StackOCIConfig) error {
	if config == nil {
		return fmt.Errorf("stack config cannot be nil")
	}
	
	// Validate struct tags
	if err := validate.Struct(config); err != nil {
		return fmt.Errorf("stack config validation failed: %w", err)
	}
	
	// Business logic validation
	
	// Validate stack name format (DNS label format)
	if !isValidDNSLabel(config.StackName) {
		return fmt.Errorf("stack name must be a valid DNS label: %s", config.StackName)
	}
	
	// Validate repository format
	if !isValidRepository(config.Repository) {
		return fmt.Errorf("invalid repository format: %s", config.Repository)
	}
	
	// Validate platform if specified
	if config.Platform != "" {
		validPlatforms := []string{
			"linux/amd64", "linux/arm64", "linux/arm/v7", "linux/arm/v6",
			"windows/amd64", "darwin/amd64", "darwin/arm64",
		}
		if !contains(validPlatforms, config.Platform) {
			return fmt.Errorf("unsupported platform: %s", config.Platform)
		}
	}
	
	// Validate labels
	for key, value := range config.Labels {
		if !isValidLabelKey(key) {
			return fmt.Errorf("invalid label key: %s", key)
		}
		if len(value) > 1024 {
			return fmt.Errorf("label value too long: %s", key)
		}
	}
	
	// Validate annotations
	for key, value := range config.Annotations {
		if !isValidAnnotationKey(key) {
			return fmt.Errorf("invalid annotation key: %s", key)
		}
		if len(value) > 2048 {
			return fmt.Errorf("annotation value too long: %s", key)
		}
	}
	
	return nil
}

// ValidateBuildRequest validates a BuildRequest struct
func ValidateBuildRequest(req *BuildRequest) error {
	if req == nil {
		return fmt.Errorf("build request cannot be nil")
	}
	
	// Validate struct tags
	if err := validate.Struct(req); err != nil {
		return fmt.Errorf("build request validation failed: %w", err)
	}
	
	// Business logic validation
	
	// Validate build ID format
	if len(req.ID) < 8 || len(req.ID) > 64 {
		return fmt.Errorf("build ID must be 8-64 characters")
	}
	
	// Validate Dockerfile path
	if !strings.HasSuffix(req.Dockerfile, "Dockerfile") && !strings.HasSuffix(req.Dockerfile, ".dockerfile") {
		return fmt.Errorf("dockerfile must end with 'Dockerfile' or '.dockerfile'")
	}
	
	return nil
}

// isValidDNSLabel validates DNS label format (RFC 1123)
func isValidDNSLabel(label string) bool {
	if len(label) == 0 || len(label) > 63 {
		return false
	}
	
	dnsLabelRegex := regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)
	return dnsLabelRegex.MatchString(label)
}

// isValidRepository validates OCI repository format
func isValidRepository(repo string) bool {
	if len(repo) == 0 || len(repo) > 255 {
		return false
	}
	
	// Repository format: [hostname[:port]/]name[:tag]
	parts := strings.Split(repo, "/")
	for _, part := range parts {
		if part == "" {
			return false
		}
		if !isValidDNSLabel(strings.Split(part, ":")[0]) {
			return false
		}
	}
	
	return true
}

// isValidLabelKey validates Docker label key format
func isValidLabelKey(key string) bool {
	if len(key) == 0 || len(key) > 63 {
		return false
	}
	
	labelKeyRegex := regexp.MustCompile(`^[a-z0-9]([a-z0-9._-]*[a-z0-9])?$`)
	return labelKeyRegex.MatchString(key)
}

// isValidAnnotationKey validates OCI annotation key format
func isValidAnnotationKey(key string) bool {
	if len(key) == 0 || len(key) > 255 {
		return false
	}
	
	// Annotation keys should be reverse DNS notation
	annotationKeyRegex := regexp.MustCompile(`^[a-z0-9]([a-z0-9.-]*[a-z0-9])?(/[a-z0-9]([a-z0-9._-]*[a-z0-9])?)*$`)
	return annotationKeyRegex.MatchString(key)
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}