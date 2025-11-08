package push

import (
	"fmt"
	"strings"

	"github.com/cnoe-io/idpbuilder/pkg/errors"
	"github.com/cnoe-io/idpbuilder/pkg/validator"
)

// validatePushOptions validates all push options with security checks.
//
// This extends Wave 2.2's PushConfig.Validate() with additional
// security validation from Wave 2.3.
func validatePushOptions(opts *PushOptions) error {
	// Validate image name
	if err := validator.ValidateImageName(opts.ImageName); err != nil {
		// validator returns typed errors, just pass through
		return err
	}

	// Validate registry URL
	if err := validator.ValidateRegistryURL(opts.Registry); err != nil {
		// Check if it's a warning
		if errors.IsWarning(err) {
			// Log warning but continue
			fmt.Println(errors.FormatError(err))
		} else {
			return err
		}
	}

	// Validate credentials
	if err := validator.ValidateCredentials(opts.Username, opts.Password); err != nil {
		if errors.IsWarning(err) {
			fmt.Println(errors.FormatError(err))
		} else {
			return err
		}
	}

	return nil
}

// WrapDockerError wraps Docker client errors with appropriate types.
// It examines the error message to categorize the error and provide
// actionable suggestions.
//
// Recognized error patterns:
//   - "No such image": ImageNotFoundError (exit code 4)
//   - "connection refused" or "Cannot connect": NetworkError (exit code 3)
//   - Other errors: Generic error with context
//
// Example:
//
//	err := dockerClient.GetImage(ctx, "alpine:latest")
//	if err != nil {
//	    return WrapDockerError(err, "alpine:latest")
//	}
func WrapDockerError(err error, imageName string) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()

	// Check for "image not found" errors
	if strings.Contains(errMsg, "No such image") {
		return errors.NewImageNotFoundError(
			imageName,
			fmt.Sprintf("image '%s' not found in local Docker daemon", imageName),
			"pull the image first with: docker pull "+imageName,
		)
	}

	// Check for Docker daemon connection errors
	if strings.Contains(errMsg, "connection refused") || strings.Contains(errMsg, "Cannot connect") {
		return errors.NewNetworkError(
			"docker daemon",
			"cannot connect to Docker daemon",
			"ensure Docker daemon is running: systemctl start docker or start Docker Desktop",
		)
	}

	// Generic Docker error
	return fmt.Errorf("Docker error: %w", err)
}

// WrapRegistryError wraps registry client errors with appropriate types.
// It examines the error message to categorize the error and provide
// actionable suggestions.
//
// Recognized error patterns:
//   - "401" or "unauthorized": AuthenticationError (exit code 2)
//   - "connection refused" or "timeout": NetworkError (exit code 3)
//   - "x509" or "certificate": NetworkError (TLS issues, exit code 3)
//   - Other errors: Generic error with context
//
// Example:
//
//	err := registryClient.Push(ctx, image, ref, callback)
//	if err != nil {
//	    return WrapRegistryError(err, "docker.io")
//	}
func WrapRegistryError(err error, registry string) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()

	// Check for authentication errors
	if strings.Contains(errMsg, "401") || strings.Contains(errMsg, "unauthorized") {
		return errors.NewAuthenticationError(
			registry,
			fmt.Sprintf("authentication failed for registry %s", registry),
			"check your username and password, or verify registry credentials",
		)
	}

	// Check for network errors
	if strings.Contains(errMsg, "connection refused") || strings.Contains(errMsg, "timeout") {
		return errors.NewNetworkError(
			registry,
			fmt.Sprintf("cannot connect to registry %s", registry),
			"verify registry URL and network connectivity, or try with --insecure if using self-signed certificates",
		)
	}

	// TLS certificate errors
	if strings.Contains(errMsg, "x509") || strings.Contains(errMsg, "certificate") {
		return errors.NewNetworkError(
			registry,
			fmt.Sprintf("TLS certificate verification failed for %s", registry),
			"use --insecure flag to skip certificate verification (not recommended for production)",
		)
	}

	// Generic registry error
	return fmt.Errorf("registry error: %w", err)
}
