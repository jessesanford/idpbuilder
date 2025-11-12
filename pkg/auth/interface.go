// Package auth provides authentication credential management for OCI registries.
package auth

import (
	"github.com/google/go-containerregistry/pkg/authn"
)

// AuthProvider supplies authentication credentials for registry operations.
// Implementations may read credentials from flags, environment variables, or configuration files.
type AuthProvider interface {
	// GetAuthenticator returns a go-containerregistry compatible authenticator.
	// The authenticator is used by RegistryClient for HTTP Basic Authentication.
	//
	// Example:
	//   authenticator, err := authProvider.GetAuthenticator()
	//   if err != nil {
	//       return fmt.Errorf("getting authenticator: %w", err)
	//   }
	//   // Use authenticator with go-containerregistry's remote.Write()
	GetAuthenticator() (authn.Authenticator, error)

	// ValidateCredentials checks that the credentials meet format and security requirements.
	// Returns an error if username or password is empty, too long, or contains invalid characters.
	//
	// Example:
	//   if err := authProvider.ValidateCredentials(); err != nil {
	//       return fmt.Errorf("invalid credentials: %w", err)
	//   }
	ValidateCredentials() error
}

// NewAuthProvider creates a new authentication provider.
// Credentials are read from the provided username and password parameters first.
// If empty, falls back to environment variables: IDPBUILDER_REGISTRY_USERNAME, IDPBUILDER_REGISTRY_PASSWORD.
//
// Returns an error if no valid credentials are available.
//
// Example:
//   // From flags:
//   provider, err := auth.NewAuthProvider("admin", "secretpass")
//
//   // From environment variables:
//   os.Setenv("IDPBUILDER_REGISTRY_USERNAME", "admin")
//   os.Setenv("IDPBUILDER_REGISTRY_PASSWORD", "secretpass")
//   provider, err := auth.NewAuthProvider("", "")
func NewAuthProvider(username, password string) (AuthProvider, error) {
	// Implementation will be provided in Wave 2
	panic("not implemented")
}

// InvalidCredentialsError indicates that credentials do not meet requirements.
type InvalidCredentialsError struct {
	Reason string
}

func (e *InvalidCredentialsError) Error() string {
	return "invalid credentials: " + e.Reason
}

// MissingCredentialsError indicates that required credentials are not provided.
type MissingCredentialsError struct {
	Field string // "username" or "password"
}

func (e *MissingCredentialsError) Error() string {
	return "missing required credential: " + e.Field
}
