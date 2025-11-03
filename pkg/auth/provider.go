// Package auth provides authentication functionality for registry operations.
package auth

import (
	"fmt"

	"github.com/google/go-containerregistry/pkg/authn"
)

// Provider supplies authentication for registry operations
type Provider interface {
	// GetAuthenticator returns the authentication configuration
	GetAuthenticator() (authn.Authenticator, error)

	// ValidateCredentials checks if credentials are valid
	ValidateCredentials() error
}

// NewBasicAuthProvider creates a basic auth provider
func NewBasicAuthProvider(username, password string) Provider {
	return &basicAuthProvider{
		username: username,
		password: password,
	}
}

// basicAuthProvider implements basic authentication
type basicAuthProvider struct {
	username string
	password string
}

func (p *basicAuthProvider) GetAuthenticator() (authn.Authenticator, error) {
	return &authn.Basic{
		Username: p.username,
		Password: p.password,
	}, nil
}

func (p *basicAuthProvider) ValidateCredentials() error {
	if p.username == "" || p.password == "" {
		return fmt.Errorf("username and password are required")
	}
	return nil
}
