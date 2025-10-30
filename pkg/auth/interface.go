// Package auth provides interfaces and types for registry authentication.
package auth

import (
	"github.com/google/go-containerregistry/pkg/authn"
)

// Provider defines operations for providing authentication credentials to registries.
type Provider interface {
	// GetAuthenticator returns an authn.Authenticator compatible with go-containerregistry.
	//
	// Returns:
	//   - authn.Authenticator: Authenticator instance
	//   - error: ValidationError if credentials are malformed
	GetAuthenticator() (authn.Authenticator, error)

	// ValidateCredentials performs pre-flight validation of credentials.
	//
	// Returns:
	//   - error: ValidationError with details if invalid, nil if valid
	ValidateCredentials() error
}

// Credentials holds authentication information for basic auth.
type Credentials struct {
	// Username for registry authentication.
	Username string

	// Password for registry authentication.
	// Supports ALL special characters including quotes, spaces, unicode.
	Password string
}

