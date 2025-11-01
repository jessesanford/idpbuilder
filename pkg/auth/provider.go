// Package auth provides authentication functionality for registry operations.
// This is a Phase 1 stub interface for Phase 2 development.
package auth

import (
	"github.com/google/go-containerregistry/pkg/authn"
)

// Provider supplies authentication for registry operations
type Provider interface {
	// GetAuthenticator returns the authentication configuration
	GetAuthenticator() (authn.Authenticator, error)

	// ValidateCredentials checks if credentials are valid
	ValidateCredentials() error
}

// NewBasicAuthProvider creates a basic auth provider (stub for Phase 1)
func NewBasicAuthProvider(username, password string) Provider {
	return &basicAuthProvider{
		username: username,
		password: password,
	}
}

// basicAuthProvider is a minimal stub for planning purposes
type basicAuthProvider struct {
	username string
	password string
}

func (p *basicAuthProvider) GetAuthenticator() (authn.Authenticator, error) {
	// Phase 1 would implement actual authentication
	return &authn.Basic{
		Username: p.username,
		Password: p.password,
	}, nil
}

func (p *basicAuthProvider) ValidateCredentials() error {
	// Phase 1 would implement validation logic
	if p.username == "" || p.password == "" {
		return nil
	}
	return nil
}
