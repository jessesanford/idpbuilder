package auth

import (
	"context"
	"fmt"

	"github.com/cnoe-io/idpbuilder/pkg/push/errors"
	"github.com/cnoe-io/idpbuilder/pkg/push/retry"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// Authenticator handles registry authentication and credential management
type Authenticator struct {
	username string
	password string
	insecure bool
	keychain authn.Keychain
}

// Option is a functional option for configuring the Authenticator
type Option func(*Authenticator) error

// WithCredentials sets the username and password credentials for authentication
func WithCredentials(username, password string) Option {
	return func(a *Authenticator) error {
		creds, err := GetCredentials(username, password)
		if err != nil {
			return fmt.Errorf("failed to get credentials: %w", err)
		}
		if creds != nil {
			a.username = creds.Username
			a.password = creds.Password
		}
		return nil
	}
}

// WithInsecure enables or disables insecure mode for self-signed certificates
func WithInsecure(insecure bool) Option {
	return func(a *Authenticator) error {
		a.insecure = insecure
		return nil
	}
}

// WithKeychain sets a custom keychain for authentication
func WithKeychain(keychain authn.Keychain) Option {
	return func(a *Authenticator) error {
		if keychain == nil {
			return fmt.Errorf("keychain cannot be nil")
		}
		a.keychain = keychain
		return nil
	}
}

// NewAuthenticator creates a new Authenticator with the given options
func NewAuthenticator(opts ...Option) (*Authenticator, error) {
	auth := &Authenticator{
		keychain: authn.DefaultKeychain,
	}

	for _, opt := range opts {
		if err := opt(auth); err != nil {
			return nil, fmt.Errorf("failed to apply authenticator option: %w", err)
		}
	}

	return auth, nil
}

// GetAuthOptions returns the remote options for authentication
func (a *Authenticator) GetAuthOptions() []remote.Option {
	var opts []remote.Option

	// Add authentication
	if a.username != "" && a.password != "" {
		auth := &authn.Basic{
			Username: a.username,
			Password: a.password,
		}
		opts = append(opts, remote.WithAuth(auth))
	} else {
		opts = append(opts, remote.WithAuthFromKeychain(a.keychain))
	}

	// Add insecure transport if needed
	if a.insecure {
		opts = append(opts, GetInsecureOption())
	}

	return opts
}

// Validate checks if the authenticator configuration is valid
func (a *Authenticator) Validate() error {
	if a.keychain == nil {
		return fmt.Errorf("keychain is required")
	}

	// If credentials are provided, both username and password must be set
	if (a.username != "" && a.password == "") || (a.username == "" && a.password != "") {
		return fmt.Errorf("both username and password must be provided together")
	}

	return nil
}

// AuthenticateWithRetry performs authentication with retry logic for transient failures
func (a *Authenticator) AuthenticateWithRetry(ctx context.Context, registry string) error {
	if err := a.Validate(); err != nil {
		return fmt.Errorf("invalid authenticator configuration: %w", err)
	}

	strategy := retry.DefaultBackoff()

	return retry.WithRetry(ctx, func() error {
		// Test authentication by attempting to authenticate with the registry
		// This is a placeholder for the actual authentication test
		// In a real implementation, this might ping the registry or attempt to list catalog

		// The actual authentication happens when GetAuthOptions() is used
		// in conjunction with remote registry operations

		// For now, we validate that we have the necessary authentication components
		opts := a.GetAuthOptions()
		if len(opts) == 0 {
			return &errors.AuthenticationError{
				Registry: registry,
				Cause:    fmt.Errorf("no authentication options configured"),
			}
		}

		return nil
	}, strategy)
}

// HasCredentials returns true if explicit credentials are configured
func (a *Authenticator) HasCredentials() bool {
	return a.username != "" && a.password != ""
}

// IsInsecure returns true if insecure mode is enabled
func (a *Authenticator) IsInsecure() bool {
	return a.insecure
}

// GetKeychain returns the configured keychain
func (a *Authenticator) GetKeychain() authn.Keychain {
	return a.keychain
}