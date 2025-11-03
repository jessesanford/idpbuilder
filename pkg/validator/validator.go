package validator

import (
	"github.com/cnoe-io/idpbuilder/pkg/errors"
)

// ValidateImageName validates the image name format.
// This is a stub implementation - full implementation in Effort 2.3.1.
func ValidateImageName(imageName string) error {
	if imageName == "" {
		return errors.NewValidationError("imageName", "image name is required", "provide a valid image name")
	}
	return nil
}

// ValidateRegistryURL validates the registry URL format.
// This is a stub implementation - full implementation in Effort 2.3.1.
func ValidateRegistryURL(registryURL string) error {
	if registryURL == "" {
		return errors.NewValidationError("registry", "registry URL is required", "provide a valid registry URL")
	}
	return nil
}

// ValidateCredentials validates username and password.
// This is a stub implementation - full implementation in Effort 2.3.1.
func ValidateCredentials(username, password string) error {
	if username == "" {
		return errors.NewValidationError("username", "username is required", "provide a username")
	}
	if password == "" {
		return errors.NewValidationError("password", "password is required", "provide a password")
	}
	return nil
}
