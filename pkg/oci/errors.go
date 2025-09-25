package oci

import (
	"errors"
	"fmt"
)

// Common authentication errors
var (
	ErrNoCredentialsFound     = errors.New("no credentials found")
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrTokenExpired           = errors.New("token expired")
	ErrRegistryUnauthorized   = errors.New("registry unauthorized")
	ErrCredentialSourceFailed = errors.New("credential source failed")
)

// AuthError represents an authentication error with context
type AuthError struct {
	Registry string
	Source   CredentialSource
	Err      error
}

// Error implements the error interface
func (ae *AuthError) Error() string {
	return fmt.Sprintf("auth error for registry %s: %v", ae.Registry, ae.Err)
}

// NewAuthError creates a new AuthError
func NewAuthError(registry string, source CredentialSource, err error) *AuthError {
	return &AuthError{
		Registry: registry,
		Source:   source,
		Err:      err,
	}
}
