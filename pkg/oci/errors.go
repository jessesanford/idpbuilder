package oci

import (
	"errors"
	"fmt"
)

// Common authentication errors
var (
	ErrNoCredentialsFound    = errors.New("no credentials found")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrTokenExpired          = errors.New("token expired")
	ErrRegistryUnauthorized  = errors.New("registry unauthorized")
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
	if ae.Registry != "" && ae.Source != 0 {
		return fmt.Sprintf("auth error for registry %s from source %s: %v",
			ae.Registry, ae.Source, ae.Err)
	}
	if ae.Registry != "" {
		return fmt.Sprintf("auth error for registry %s: %v", ae.Registry, ae.Err)
	}
	return fmt.Sprintf("auth error: %v", ae.Err)
}

// Unwrap returns the underlying error
func (ae *AuthError) Unwrap() error {
	return ae.Err
}

// IsAuthError checks if an error is an AuthError
func IsAuthError(err error) bool {
	var authErr *AuthError
	return errors.As(err, &authErr)
}

// NewAuthError creates a new AuthError
func NewAuthError(registry string, source CredentialSource, err error) *AuthError {
	return &AuthError{
		Registry: registry,
		Source:   source,
		Err:      err,
	}
}

// IsCredentialNotFound checks if error indicates missing credentials
func IsCredentialNotFound(err error) bool {
	return errors.Is(err, ErrNoCredentialsFound)
}

// IsTokenExpired checks if error indicates expired token
func IsTokenExpired(err error) bool {
	return errors.Is(err, ErrTokenExpired)
}

// IsUnauthorized checks if error indicates unauthorized access
func IsUnauthorized(err error) bool {
	return errors.Is(err, ErrRegistryUnauthorized)
}