package oci

import (
	"context"
	"net/http"
	"time"
)

// Authenticator handles authentication for OCI registries
type Authenticator interface {
	// Authenticate returns credentials for the given registry
	Authenticate(ctx context.Context, registry string) (*Credentials, error)

	// RefreshToken refreshes an expired token
	RefreshToken(ctx context.Context, registry string) (*Credentials, error)

	// ValidateCredentials checks if credentials are still valid
	ValidateCredentials(ctx context.Context, creds *Credentials) (bool, error)
}

// Credentials represents authentication credentials
type Credentials struct {
	Username  string
	Password  string
	Token     string
	Registry  string
	ExpiresAt time.Time
}

// CredentialSource defines where credentials come from
type CredentialSource int

const (
	SourceDockerConfig CredentialSource = iota
	SourceEnvironment
	SourceKeychain
	SourceKubernetes
)

// String returns string representation of CredentialSource
func (cs CredentialSource) String() string {
	switch cs {
	case SourceDockerConfig:
		return "docker-config"
	case SourceEnvironment:
		return "environment"
	case SourceKeychain:
		return "keychain"
	case SourceKubernetes:
		return "kubernetes"
	default:
		return "unknown"
	}
}

// IsExpired checks if credentials have expired
func (c *Credentials) IsExpired() bool {
	if c.ExpiresAt.IsZero() {
		return false // No expiry set
	}
	return time.Now().After(c.ExpiresAt)
}

// IsValid checks if credentials are minimally valid
func (c *Credentials) IsValid() bool {
	return c != nil && (c.Username != "" || c.Token != "")
}

// RegistryClient defines the expected interface for OCI registry client
// This interface is implemented by OCIClient in client.go
type RegistryClient interface {
	Connect(ctx context.Context, registry string) error
	Authenticate(credentials *ClientCredentials) error
	SetInsecure(insecure bool)
	GetTransport() http.RoundTripper
	Close() error
}

// ClientCredentials represents the authentication information for registry access
// Integrates with Phase 2 authentication system
// Note: Renamed to avoid conflict with existing Credentials type
type ClientCredentials struct {
	Username string
	Password string
	Token    string
	Registry string
}