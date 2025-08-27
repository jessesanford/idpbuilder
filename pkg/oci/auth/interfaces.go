// Package auth provides authentication and certificate management interfaces
// for OCI registry operations. This package defines the contracts that will
// be implemented by various authentication providers in later phases.
package auth

import (
	"context"
)

// AuthProvider defines the interface for managing authentication credentials
// across different registry providers. Implementations handle credential
// retrieval, validation, refresh, and storage.
type AuthProvider interface {
	// GetCredentials retrieves credentials for the specified registry.
	// Returns ErrCredentialsNotFound if no credentials exist for the registry.
	GetCredentials(ctx context.Context, registry string) (*Credentials, error)

	// ValidateCredentials checks if the provided credentials are still valid
	// and haven't expired. Implementations may perform network calls to verify.
	ValidateCredentials(ctx context.Context, creds *Credentials) error

	// RefreshToken attempts to refresh expired tokens using refresh tokens
	// or other credential refresh mechanisms. Only applicable for token-based auth.
	RefreshToken(ctx context.Context, creds *Credentials) (*Token, error)

	// StoreCredentials securely persists credentials for the specified registry.
	// Implementations should encrypt sensitive data before storage.
	StoreCredentials(ctx context.Context, registry string, creds *Credentials) error

	// RemoveCredentials deletes stored credentials for the specified registry.
	// This is used when credentials are revoked or no longer needed.
	RemoveCredentials(ctx context.Context, registry string) error
}

// CredentialStore defines the interface for persistent storage of credentials.
// Different implementations can provide file-based, keyring, or cloud-based storage.
type CredentialStore interface {
	// Load retrieves credentials from persistent storage for the specified registry.
	// Returns ErrCredentialsNotFound if no credentials exist.
	Load(ctx context.Context, registry string) (*Credentials, error)

	// Save persists credentials to storage for the specified registry.
	// Implementations should handle encryption and secure storage.
	Save(ctx context.Context, registry string, creds *Credentials) error

	// Delete removes credentials from storage for the specified registry.
	// No error is returned if credentials don't exist.
	Delete(ctx context.Context, registry string) error

	// ListRegistries returns a list of all registries that have stored credentials.
	// This is useful for credential management and cleanup operations.
	ListRegistries(ctx context.Context) ([]string, error)
}

// CertificateProvider defines the interface for managing TLS certificates
// and trust bundles for secure registry connections.
type CertificateProvider interface {
	// GetCertBundle retrieves the certificate bundle for secure connections
	// to the specified registry. Returns system certs if none specified.
	GetCertBundle(ctx context.Context, registry string) (*CertificateBundle, error)

	// LoadFromSystem loads system certificate authorities from the OS.
	// This provides the default trust store for most scenarios.
	LoadFromSystem() (*CertificateBundle, error)

	// LoadFromFile loads certificates from a PEM file.
	// Used for custom CA certificates or client certificates.
	LoadFromFile(ctx context.Context, path string) (*CertificateBundle, error)

	// LoadFromKind loads certificates from a Kind cluster configuration.
	// This is specific to development environments using Kind clusters.
	LoadFromKind(ctx context.Context, clusterName string) (*CertificateBundle, error)

	// ValidateCertificate validates a certificate against the trust chain
	// and checks expiration, revocation, and other validity criteria.
	ValidateCertificate(ctx context.Context, certData []byte) (*CertificateInfo, error)

	// GetTLSConfig returns a complete TLS configuration for the registry,
	// including certificates, cipher suites, and security settings.
	GetTLSConfig(ctx context.Context, registry string) (*TLSConfig, error)
}

// TokenManager defines the interface for managing authentication tokens,
// including validation, refresh, and caching operations.
type TokenManager interface {
	// IsTokenValid checks if a token is still valid and hasn't expired.
	// May perform network calls for server-side validation.
	IsTokenValid(ctx context.Context, token *Token) bool

	// RefreshToken attempts to refresh an expired token using refresh tokens
	// or other token refresh mechanisms.
	RefreshToken(ctx context.Context, token *Token) (*Token, error)

	// ParseToken parses a token string into a structured Token object,
	// extracting metadata like expiration, type, and scope.
	ParseToken(ctx context.Context, tokenString string) (*Token, error)

	// CacheToken stores a token in memory cache with appropriate TTL.
	// This reduces the need for repeated authentication calls.
	CacheToken(ctx context.Context, registry string, token *Token) error

	// GetCachedToken retrieves a cached token for the registry if available
	// and not expired. Returns nil if no valid cached token exists.
	GetCachedToken(ctx context.Context, registry string) (*Token, error)

	// ClearTokenCache removes cached tokens for the specified registry.
	// This is used when tokens are revoked or authentication changes.
	ClearTokenCache(ctx context.Context, registry string) error
}