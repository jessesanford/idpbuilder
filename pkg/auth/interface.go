// Package auth provides authentication interfaces and types for OCI registry access.
//
// This package defines the core authentication contract used by registry clients
// to obtain authenticators compatible with go-containerregistry.
package auth

import (
	"github.com/google/go-containerregistry/pkg/authn"
)

// Provider defines the interface for authentication providers.
//
// Implementations must provide go-containerregistry compatible authenticators
// and validate credentials before use.
type Provider interface {
	// GetAuthenticator returns an authn.Authenticator for use with go-containerregistry.
	//
	// This method validates credentials and returns an authenticator that can be
	// used with remote.Push(), remote.Pull(), and other go-containerregistry functions.
	//
	// Returns:
	//   - authn.Authenticator: Authenticator instance compatible with go-containerregistry
	//   - error: ValidationError if credentials are invalid or malformed
	GetAuthenticator() (authn.Authenticator, error)

	// ValidateCredentials performs pre-flight validation of credentials.
	//
	// This method checks credential format and basic requirements without
	// contacting the registry. It should be called before attempting authentication
	// to fail fast on invalid credentials.
	//
	// Returns:
	//   - error: ValidationError if credentials are invalid, nil if valid
	ValidateCredentials() error
}

// Credentials holds authentication credentials.
//
// This struct stores the username and password for basic authentication.
// Credentials are validated by Provider implementations before use.
type Credentials struct {
	// Username is the registry username
	Username string

	// Password is the registry password
	Password string
}
